package containerd

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/docker/cli/cli/config/credentials"
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/preflight"
	"github.com/weaveworks/ignite/pkg/resolvconf"
	"github.com/weaveworks/ignite/pkg/runtime"
	"github.com/weaveworks/ignite/pkg/runtime/auth"
	"github.com/weaveworks/ignite/pkg/util"

	"github.com/containerd/console"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/containers"
	"github.com/containerd/containerd/defaults"
	"github.com/containerd/containerd/errdefs"
	"github.com/containerd/containerd/mount"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/oci"
	"github.com/containerd/containerd/plugin"
	refdocker "github.com/containerd/containerd/reference/docker"
	"github.com/containerd/containerd/remotes"
	"github.com/containerd/containerd/remotes/docker"
	v2shim "github.com/containerd/containerd/runtime/v2/shim"
	"github.com/containerd/containerd/snapshots"
	"github.com/opencontainers/go-digest"
	"github.com/opencontainers/image-spec/identity"
	imagespec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/opencontainers/runtime-spec/specs-go"
	log "github.com/sirupsen/logrus"
	"github.com/weaveworks/ignite/pkg/providers"
	"golang.org/x/sys/unix"
)

const (
	ctdNamespace     = "firecracker"
	stopTimeoutLabel = "IgniteStopTimeout"
	logPathTemplate  = "/tmp/%s.log"
	resolvConfName   = "runtime.containerd.resolv.conf"

	// InsecureRegistriesEnvVar helps set insecure registries.
	InsecureRegistriesEnvVar = "IGNITE_CONTAINERD_INSECURE_REGISTRIES"
)

var (
	// containerdSocketLocations is a list of socket locations to stat for
	containerdSocketLocations = []string{
		defaults.DefaultAddress, // "/run/containerd/containerd.sock"
		"/run/docker/containerd/containerd.sock",
	}
	// v2ShimRuntimes is a list of containerd runtimes we support.
	// Note that we only list `runc` runtimes -- containerd plugin runtimes are omitted.
	// This package also supports a fallback to the legacy runtime: `plugin.RuntimeLinuxV1`.
	v2ShimRuntimes = []string{
		plugin.RuntimeRuncV2,
		plugin.RuntimeRuncV1,
	}
)

// ctdClient is a runtime.Interface
// implementation serving the containerd client
type ctdClient struct {
	ctx    context.Context
	client *containerd.Client
}

var _ runtime.Interface = &ctdClient{}

// statContainerdSocket returns the first existing file in the containerdSocketLocations list
func statContainerdSocket() (string, error) {
	for _, socket := range containerdSocketLocations {
		if _, err := os.Stat(socket); err == nil {
			return socket, nil
		}
	}
	return "", fmt.Errorf("Could not stat a containerd socket: %v", containerdSocketLocations)
}

// getNewestAvailableContainerdRuntime returns the newest possible runtime for the shims available in the PATH.
// If no shim is found, the legacy Linux V1 runtime is returned along with an error.
// Use of this function couples ignite to the PATH and mount namespace of containerd which is undesireable.
//
// TODO(stealthybox): PR CheckRuntime() to containerd libraries instead of using exec.LookPath()
func getNewestAvailableContainerdRuntime() (string, error) {
	for _, rt := range v2ShimRuntimes {
		binary := v2shim.BinaryName(rt)
		if binary == "" {
			// this shouldn't happen if the matching test is passing, but it's not fatal -- just log and continue
			log.Errorf("shim binary could not be found -- %q is an invalid runtime/v2/shim", rt)
		} else if _, err := exec.LookPath(binary); err == nil {
			return rt, nil
		}
	}

	// legacy fallback needs hard-coded binary name -- it's not exported by containerd/runtime/v1/shim
	if _, err := exec.LookPath("containerd-shim"); err == nil {
		return plugin.RuntimeLinuxV1, nil
	}

	// legacy fallback needs hard-coded binary name -- it's not exported by containerd/runtime/v1/shim
	// this is for debian's packaging of docker.io
	if _, err := exec.LookPath("docker-containerd-shim"); err == nil {
		return plugin.RuntimeLinuxV1, nil
	}

	// default to the legacy runtime and return an error so the caller can decide what to do
	return plugin.RuntimeLinuxV1, fmt.Errorf("a containerd-shim could not be found for runtimes %v, %s", v2ShimRuntimes, plugin.RuntimeLinuxV1)
}

// GetContainerdClient builds a client for talking to containerd
func GetContainerdClient() (*ctdClient, error) {
	ctdSocket, err := statContainerdSocket()
	if err != nil {
		return nil, err
	}

	runtime, err := getNewestAvailableContainerdRuntime()
	if err != nil {
		// proceed with the default runtime -- our PATH can't see a shim binary, but containerd might be able to
		log.Warningf("Proceeding with default runtime %q: %v", runtime, err)
	}

	cli, err := containerd.New(
		ctdSocket,
		containerd.WithDefaultRuntime(runtime),
	)
	if err != nil {
		return nil, err
	}

	return &ctdClient{
		client: cli,
		ctx:    namespaces.WithNamespace(context.Background(), ctdNamespace),
	}, nil
}

// newRemoteResolver returns a remote resolver with auth info for a given
// host name.
func newRemoteResolver(refHostname string, configPath string) (remotes.Resolver, error) {
	var authzOpts []docker.AuthorizerOpt
	regOpts := []docker.RegistryOpt{}
	insecureAllowed := false
	client := &http.Client{}

	// Allow setting insecure_registries through a client-side ENV variable.
	// dockerconfig.json does not have a place to set this.
	// We would have to override the parser to add a field otherwise.
	for _, reg := range strings.Split(os.Getenv(InsecureRegistriesEnvVar), ",") {
		// image hostnames don't have protocols, this is the most forgiving parsing logic.
		if credentials.ConvertToHostname(reg) == refHostname {
			insecureAllowed = true
		}
	}

	if authCreds, serverAddress, err := auth.NewAuthCreds(refHostname, configPath); err != nil {
		return nil, err
	} else {
		authzOpts = append(authzOpts, docker.WithAuthCreds(authCreds))
		// Allow the dockerconfig.json to specify HTTP as a specific protocol override, defaults to HTTPS
		if strings.HasPrefix(serverAddress, "http://") {
			if !insecureAllowed {
				return nil, fmt.Errorf("Registry %q uses plain HTTP, but is not in the %s env var", serverAddress, InsecureRegistriesEnvVar)
			}
			regOpts = append(regOpts, docker.WithPlainHTTP(docker.MatchAllHosts))
		} else {
			if insecureAllowed {
				log.Warnf("Disabling TLS Verification for %q via %s env var", serverAddress, InsecureRegistriesEnvVar)
				client.Transport = &http.Transport{
					TLSClientConfig: &tls.Config{
						InsecureSkipVerify: true,
					},
				}
			}
		}
	}
	authz := docker.NewDockerAuthorizer(authzOpts...)

	regOpts = append(regOpts, docker.WithAuthorizer(authz))
	regOpts = append(regOpts, docker.WithClient(client))

	resolverOpts := docker.ResolverOptions{
		Hosts: docker.ConfigureDefaultRegistries(regOpts...),
	}

	resolver := docker.NewResolver(resolverOpts)
	return resolver, nil
}

func (cc *ctdClient) PullImage(image meta.OCIImageRef) error {
	log.Debugf("containerd: Pulling image %q", image)

	// Get the domain name from the image.
	named, err := refdocker.ParseDockerRef(image.String())
	if err != nil {
		return err
	}
	refDomain := refdocker.Domain(named)

	// Create a remote resolver for the domain.
	resolver, err := newRemoteResolver(refDomain, providers.RegistryConfigDir)
	if err != nil {
		return err
	}

	opts := []containerd.RemoteOpt{
		containerd.WithResolver(resolver),
		containerd.WithPullUnpack,
	}

	_, err = cc.client.Pull(cc.ctx, image.Normalized(), opts...)
	return err
}

func (cc *ctdClient) InspectImage(image meta.OCIImageRef) (result *runtime.ImageInspectResult, err error) {
	var img containerd.Image
	var config imagespec.Descriptor
	var id *meta.OCIContentID
	var usage snapshots.Usage

	log.Debugf("containerd: Inspecting image %q", image)
	if img, err = cc.client.GetImage(cc.ctx, image.Normalized()); err != nil {
		return
	}

	if config, err = img.Config(cc.ctx); err != nil {
		return
	}

	if usage, err = cc.imageUsage(img); err != nil {
		return
	}

	// img.Name() -> "docker.io/weaveworks/ignite-ubuntu:latest"
	// config.Digest.String() -> "sha256:9552fe790974f7232205bf8219934d49af38fd4b47aaeb61f539ea735e93f26e"
	if id, err = meta.ParseOCIContentID(fmt.Sprintf("%s@%s", img.Name(), config.Digest.String())); err != nil {
		return
	}

	result = &runtime.ImageInspectResult{
		ID:   id,
		Size: usage.Size,
	}

	return
}

// ExportImage exports the root filesystem of the given image. It does so by fetching the snapshots
// created from the image's layers, creating a read-only, mountable view snapshot on top with a random
// key, mounting that snapshot into a temporary directory, and starting a tar streamer capturing the
// assembled snapshot stack's root filesystem.
func (cc *ctdClient) ExportImage(image meta.OCIImageRef) (r io.ReadCloser, cleanup func() error, err error) {
	var (
		viewKey string
		img     containerd.Image
		diffIDs []digest.Digest
		mounts  []mount.Mount
		dir     string
	)

	// Fetch the image based on the given ID
	if img, err = cc.client.GetImage(cc.ctx, image.Normalized()); err != nil {
		return
	}

	// Load the default snapshotter and the diff IDs for the image
	snapshotter := cc.client.SnapshotService(containerd.DefaultSnapshotter)
	if diffIDs, err = img.RootFS(cc.ctx); err != nil {
		return
	}

	for {
		// Generate a new key for the mountable view snapshot
		if viewKey, err = util.NewUID(); err != nil {
			return
		}

		// Verify that no snapshot exists matching the viewKey
		if _, err = snapshotter.Stat(cc.ctx, viewKey); err != nil {
			break
		}
	}

	// Fetch the key for the top-most snapshot of the image
	// and create the mountable view snapshot on top of it using viewKey
	key := identity.ChainID(diffIDs).String()
	if mounts, err = snapshotter.View(cc.ctx, viewKey, key); err != nil {
		return
	}

	// Create a temporary directory to mount the view snapshot
	if dir, err = ioutil.TempDir("", ""); err != nil {
		return
	}

	// Perform the mount using syscalls
	if err = mount.All(mounts, dir); err != nil {
		return
	}

	// Get the info of each entry in the mount
	var infos []os.FileInfo
	if infos, err = ioutil.ReadDir(dir); err != nil {
		return
	}

	// Construct the arguments, append each entry's name
	// Each entry is appended separately to avoid the paths
	// becoming "./boot/" instead of "boot/", that messes
	// with the selective extraction of kernels
	args := append(make([]string, 0, len(infos)+3), "-c", "-C", dir)
	for _, info := range infos {
		args = append(args, info.Name())
	}

	// Create the tar streaming command and assign the io.ReadCloser to be returned
	tarCmd := exec.Command("tar", args...)
	if r, err = tarCmd.StdoutPipe(); err != nil {
		return
	}

	// Start the tar streamer
	if err = tarCmd.Start(); err != nil {
		return
	}

	// Construct the cleanup function
	cleanup = func() (err error) {
		defer util.DeferErr(&err, snapshotter.Close)
		defer util.DeferErr(&err, func() error { return snapshotter.Remove(cc.ctx, viewKey) })
		defer util.DeferErr(&err, func() error { return mount.UnmountAll(dir, 0) })
		defer util.DeferErr(&err, tarCmd.Wait)
		return
	}

	return
}

func (cc *ctdClient) InspectContainer(container string) (*runtime.ContainerInspectResult, error) {
	var cont containerd.Container

	cont, err := cc.client.LoadContainer(cc.ctx, container)
	if err != nil {
		return nil, err
	}

	task, err := cont.Task(cc.ctx, nil)
	if err != nil {
		return nil, err
	}

	tStatus, err := task.Status(cc.ctx)
	if err != nil {
		return nil, err
	}

	info, err := cont.Info(cc.ctx)
	if err != nil {
		return nil, err
	}

	return &runtime.ContainerInspectResult{
		ID:        info.ID,
		Image:     info.Image, // TODO: This may be incorrect
		Status:    string(tStatus.Status),
		IPAddress: nil, // TODO: This, containerd only supports CNI
		PID:       task.Pid(),
	}, nil
}

/*
FIFO handling with attach and logs:
- Attach will read the stdout FIFO and in addition to writing to screen, copy the output to a file
- Logs will first read and print the file, then read the FIFO and write it's output to the file and the screen
*/

func (cc *ctdClient) AttachContainer(container string) (err error) {
	var (
		cont containerd.Container
		spec *oci.Spec
	)

	if cont, err = cc.client.LoadContainer(cc.ctx, container); err != nil {
		return
	}

	if spec, err = cont.Spec(cc.ctx); err != nil {
		return
	}

	var (
		con console.Console
		tty = spec.Process.Terminal
	)

	if tty {
		con = console.Current()
		defer util.DeferErr(&err, con.Reset)
		if err = con.SetRaw(); err != nil {
			return
		}
	}

	var (
		task     containerd.Task
		statusC  <-chan containerd.ExitStatus
		igniteIO *igniteIO
	)

	if igniteIO, err = newIgniteIO(fmt.Sprintf(logPathTemplate, container)); err != nil {
		return
	}
	defer util.DeferErr(&err, igniteIO.Close)

	if task, err = cont.Task(cc.ctx, cio.NewAttach(igniteIO.Opt())); err != nil {
		return
	}

	if statusC, err = task.Wait(cc.ctx); err != nil {
		return
	}

	if tty {
		if err := HandleConsoleResize(cc.ctx, task, con); err != nil {
			log.Errorf("console resize failed: %v", err)
		}
	} else {
		sigc := ForwardAllSignals(cc.ctx, task)
		defer StopCatch(sigc)
	}

	var code uint32
	select {
	case ec := <-statusC:
		code, _, err = ec.Result()
	case <-igniteIO.Detach():
		fmt.Println() // Use a new line for the log entry
		log.Println("Detached")
	}

	if code != 0 && err == nil {
		err = fmt.Errorf("attach exited with code %d", code)
	}

	return
}

func (cc *ctdClient) RunContainer(image meta.OCIImageRef, config *runtime.ContainerConfig, name, id string) (s string, err error) {
	img, err := cc.client.GetImage(cc.ctx, image.Normalized())
	if err != nil {
		return
	}

	// Remove the container if it exists
	if err = cc.RemoveContainer(name); err != nil {
		return
	}

	// Load the default snapshotter
	snapshotter := cc.client.SnapshotService(containerd.DefaultSnapshotter)

	// Add the /etc/resolv.conf mount, this isn't done automatically by containerd
	// Ensure a resolv.conf exists in the vmDir. Calculate path using the vm id
	// TODO(stealthybox):
	//  - create snapshot ahead of time?
	//    - is there a performance penalty for creating snapshot ahead of time?
	//    - maybe we can use containerd.NewContainerOpts{} to do it just-in-time
	//  - write this file to snapshot mount instead of vmDir
	//  - commit snapshot?
	//  - deprecate vm id from this function's signature
	resolvConfPath := filepath.Join(constants.VM_DIR, id, resolvConfName)
	err = resolvconf.EnsureResolvConf(resolvConfPath, constants.DATA_DIR_FILE_PERM)
	if err != nil {
		return
	}
	config.Binds = append(
		config.Binds,
		&runtime.Bind{
			HostPath:      resolvConfPath,
			ContainerPath: "/etc/resolv.conf",
		},
	)

	// Add the stop timeout as a label, as containerd doesn't natively support it
	config.Labels[stopTimeoutLabel] = strconv.FormatUint(uint64(config.StopTimeout), 10)

	// Build the OCI specification
	opts := []oci.SpecOpts{
		oci.WithDefaultSpec(),
		oci.WithDefaultUnixDevices,
		oci.WithTTY,
		oci.WithImageConfigArgs(img, config.Cmd),
		oci.WithEnv(config.EnvVars),
		withAddedCaps(config.CapAdds),
		withHostname(config.Hostname),
		withMounts(config.Binds),
		withDevices(config.Devices),
	}

	// Known limitations, containerd doesn't support the following config fields:
	// - StopTimeout
	// - AutoRemove
	// - NetworkMode (only CNI supported)
	// - PortBindings

	snapshotOpt := containerd.WithSnapshot(name)
	if _, err = snapshotter.Stat(cc.ctx, name); errdefs.IsNotFound(err) {
		// Even if "read only" is set, we don't use a KindView snapshot here (#1495).
		// We pass the writable snapshot to the OCI runtime, and the runtime remounts
		// it as read-only after creating some mount points on-demand.
		snapshotOpt = containerd.WithNewSnapshot(name, img)
	} else if err != nil {
		return
	}

	cOpts := []containerd.NewContainerOpts{
		containerd.WithImage(img),
		snapshotOpt,
		//containerd.WithImageStopSignal(img, "SIGTERM"),
		containerd.WithNewSpec(opts...),
		containerd.WithContainerLabels(config.Labels),
	}

	cont, err := cc.client.NewContainer(cc.ctx, name, cOpts...)
	if err != nil {
		return
	}

	// This is a dummy PTY to silence output
	// when starting without attach breaking
	con, _, err := console.NewPty()
	if err != nil {
		return
	}
	defer util.DeferErr(&err, con.Close)

	// We need a temporary dummy stdin reader that
	// actually works, can't use nullReader here
	dummyReader, _, err := os.Pipe()
	if err != nil {
		return
	}
	defer util.DeferErr(&err, dummyReader.Close)

	// Spawn the Creator with the dummy streams
	ioCreator := cio.NewCreator(cio.WithTerminal, cio.WithStreams(dummyReader, con, con))

	task, err := cont.NewTask(cc.ctx, ioCreator)
	if err != nil {
		return
	}

	if err = task.Start(cc.ctx); err != nil {
		return
	}

	// TODO: Save task.Pid() somewhere for attaching?
	s = task.ID()
	return
}

func withAddedCaps(caps []string) oci.SpecOpts {
	prefixed := make([]string, 0, len(caps))

	for _, c := range caps {
		// TODO: Make the CAPs have an unified format between Docker and containerd
		prefixed = append(prefixed, fmt.Sprintf("CAP_%s", c))
	}

	return oci.WithAddedCapabilities(prefixed)
}

func withHostname(hostname string) oci.SpecOpts {
	return func(_ context.Context, _ oci.Client, _ *containers.Container, s *specs.Spec) error {
		s.Hostname = hostname
		return nil
	}
}

func withMounts(binds []*runtime.Bind) oci.SpecOpts {
	mounts := make([]specs.Mount, 0, len(binds))
	for _, bind := range binds {
		mounts = append(mounts, specs.Mount{
			Source:      bind.HostPath,
			Destination: bind.ContainerPath,
			Type:        "bind",
			Options: []string{
				"rbind",
				"rw",
			},
		})
	}

	return oci.WithMounts(mounts)
}

func withDevices(devices []*runtime.Bind) oci.SpecOpts {
	return func(_ context.Context, _ oci.Client, _ *containers.Container, s *specs.Spec) error {
		for _, dev := range devices {
			// Support only character and block devices
			var devType string
			if err := deviceType(dev.HostPath, &devType); err != nil {
				return err
			}

			if s.Linux == nil {
				s.Linux = &specs.Linux{}
			}
			if s.Linux.Resources == nil {
				s.Linux.Resources = &specs.LinuxResources{}
			}

			var stat unix.Stat_t
			if err := unix.Stat(dev.HostPath, &stat); err != nil {
				return err
			}

			major := int64(unix.Major(stat.Rdev))
			minor := int64(unix.Minor(stat.Rdev))

			s.Linux.Resources.Devices = append(s.Linux.Resources.Devices, specs.LinuxDeviceCgroup{
				Type:   devType,
				Major:  &major,
				Minor:  &minor,
				Access: "rwm",
				Allow:  true,
			})
			s.Linux.Devices = append(s.Linux.Devices, specs.LinuxDevice{
				Path:  dev.ContainerPath, // dev.HostPath is irrelevant for the container Spec -- major,minor is the primary key
				Type:  devType,
				Major: major,
				Minor: minor,
				// Maybe set FileMode too?
				UID: &stat.Uid,
				GID: &stat.Gid,
			})
		}

		return nil
	}
}

func (cc *ctdClient) StopContainer(container string, timeout *time.Duration) (err error) {
	cont, err := cc.client.LoadContainer(cc.ctx, container)
	if err != nil {
		// If the container is not found, return nil, no-op.
		if errdefs.IsNotFound(err) {
			log.Warn(err)
			err = nil
		}
		return
	}

	// Use the container-specific timeout if no timeout is given
	if timeout == nil {
		var labels map[string]string
		var duration uint64

		if labels, err = cont.Labels(cc.ctx); err != nil {
			return
		}

		if duration, err = strconv.ParseUint(labels[stopTimeoutLabel], 10, 32); err != nil {
			return
		}

		to := time.Duration(duration) * time.Second
		timeout = &to
	}

	task, err := cont.Task(cc.ctx, cio.Load)
	if err != nil {
		// If the task is not found, return nil, no-op.
		if errdefs.IsNotFound(err) {
			log.Warn(err)
			err = nil
		}
		return
	}

	// Initiate a wait
	waitC, err := task.Wait(cc.ctx)
	if err != nil {
		return
	}

	// Send a SIGTERM signal to request a clean shutdown
	if err = task.Kill(cc.ctx, syscall.SIGTERM); err != nil {
		return
	}

	// After sending the signal, start the timer to force-kill the task
	timeoutC := make(chan error)
	timer := time.AfterFunc(*timeout, func() {
		timeoutC <- task.Kill(cc.ctx, syscall.SIGQUIT)
	})

	// Wait for the task to stop or the timer to fire
	select {
	case exitStatus := <-waitC:
		timer.Stop()             // Cancel the force-kill timer
		err = exitStatus.Error() // TODO: Handle exit code
	case err = <-timeoutC: // The kill timer has fired
	}

	// Delete the task
	if _, e := task.Delete(cc.ctx); e != nil {
		if err != nil {
			err = fmt.Errorf("%v, task deletion failed: %v", err, e) // TODO: Multierror
		} else {
			err = e
		}
	}

	return
}

func (cc *ctdClient) KillContainer(container, signal string) (err error) {
	cont, err := cc.client.LoadContainer(cc.ctx, container)
	if err != nil {
		// If the container is not found, return nil, no-op.
		if errdefs.IsNotFound(err) {
			log.Warn(err)
			err = nil
		}
		return
	}

	task, err := cont.Task(cc.ctx, cio.Load)
	if err != nil {
		// If the task is not found, return nil, no-op.
		if errdefs.IsNotFound(err) {
			log.Warn(err)
			err = nil
		}
		return
	}

	// Initiate a wait
	waitC, err := task.Wait(cc.ctx)
	if err != nil {
		return
	}

	// Send a SIGQUIT signal to force stop
	if err = task.Kill(cc.ctx, syscall.SIGQUIT); err != nil {
		return
	}

	// Wait for the container to stop
	<-waitC

	// Delete the task
	_, err = task.Delete(cc.ctx)
	return
}

func (cc *ctdClient) RemoveContainer(container string) error {
	// Remove the container if it exists
	cont, contLoadErr := cc.client.LoadContainer(cc.ctx, container)
	if errdefs.IsNotFound(contLoadErr) {
		log.Debug(contLoadErr)
		return nil
	} else if contLoadErr != nil {
		return contLoadErr
	}

	// Load the container's task without attaching
	task, taskLoadErr := cont.Task(cc.ctx, nil)
	if errdefs.IsNotFound(taskLoadErr) {
		log.Debug(taskLoadErr)
	} else if taskLoadErr != nil {
		return taskLoadErr
	} else {
		_, taskDeleteErr := task.Delete(cc.ctx)
		if taskDeleteErr != nil {
			log.Debug(taskDeleteErr)
		}
	}

	// Delete the container
	deleteContErr := cont.Delete(cc.ctx, containerd.WithSnapshotCleanup)
	if errdefs.IsNotFound(contLoadErr) {
		log.Debug(contLoadErr)
	} else if deleteContErr != nil {
		return deleteContErr
	}

	// Remove the log file if it exists
	logFile := fmt.Sprintf(logPathTemplate, container)
	if util.FileExists(logFile) {
		logDeleteErr := os.RemoveAll(logFile)
		if logDeleteErr != nil {
			return logDeleteErr
		}
	}

	return nil
}

func (cc *ctdClient) ContainerLogs(container string) (r io.ReadCloser, err error) {
	var (
		cont containerd.Container
	)

	if cont, err = cc.client.LoadContainer(cc.ctx, container); err != nil {
		return
	}

	var retriever *logRetriever
	if retriever, err = newlogRetriever(fmt.Sprintf(logPathTemplate, container)); err != nil {
		return
	}

	if _, err = cont.Task(cc.ctx, cio.NewAttach(retriever.Opt())); err != nil {
		return
	}

	// Currently we have no way of detecting if the task's attach has filled the stdout and stderr
	// buffers without asynchronous I/O (syscall.Conn and syscall.Splice). If the read reaches
	// the end, the application hangs indefinitely waiting for new output from the container.
	// TODO: Get rid of this, implement asynchronous I/O and read until the streams have been exhausted
	time.Sleep(time.Second)

	// Close the writer to signal EOF
	if err = retriever.CloseWriter(); err != nil {
		return
	}

	return retriever, nil
}

func (cc *ctdClient) Name() runtime.Name {
	return runtime.RuntimeContainerd
}

func (cc *ctdClient) RawClient() interface{} {
	return cc.client
}

type containerdSocketChecker struct{}

func (ctdsc containerdSocketChecker) Check() error {
	_, err := statContainerdSocket()
	return err
}
func (ctdsc containerdSocketChecker) Name() string {
	return "containerdSocketChecker"
}
func (ctdsc containerdSocketChecker) Type() string {
	return "containerdSocketChecker"
}
func (cc *ctdClient) PreflightChecker() preflight.Checker {
	return containerdSocketChecker{}
}

// imageUsage returns the size/inode usage of the given image by
// summing up the resource usage of each of its snapshot layers
func (cc *ctdClient) imageUsage(image containerd.Image) (usage snapshots.Usage, err error) {
	var digestIDs []digest.Digest

	snapshotter := cc.client.SnapshotService(containerd.DefaultSnapshotter)
	defer util.DeferErr(&err, snapshotter.Close)
	if digestIDs, err = image.RootFS(cc.ctx); err != nil {
		return
	}

	chainIDs := identity.ChainIDs(digestIDs)
	for _, chainID := range chainIDs {
		var u snapshots.Usage
		if u, err = snapshotter.Usage(cc.ctx, chainID.String()); err != nil {
			return
		}

		usage.Add(u)
	}

	return
}

func deviceType(device string, devType *string) (err error) {
	if exists, info := util.PathExists(device); exists {
		if info.Mode()&os.ModeCharDevice != 0 {
			*devType = "c"
		} else if info.Mode()&os.ModeDevice != 0 {
			*devType = "b"
		} else {
			err = fmt.Errorf("not a device file: %q", device)
		}
	} else {
		err = fmt.Errorf("device path %q not found", device)
	}

	return
}
