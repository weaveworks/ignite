package containerd

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"time"

	"github.com/containerd/console"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/containers"
	"github.com/containerd/containerd/errdefs"
	"github.com/containerd/containerd/mount"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/oci"
	"github.com/containerd/containerd/plugin"
	"github.com/containerd/containerd/snapshots"
	"github.com/opencontainers/go-digest"
	"github.com/opencontainers/image-spec/identity"
	imagespec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/opencontainers/runtime-spec/specs-go"
	log "github.com/sirupsen/logrus"
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/preflight"
	"github.com/weaveworks/ignite/pkg/preflight/checkers"
	"github.com/weaveworks/ignite/pkg/runtime"
	"github.com/weaveworks/ignite/pkg/util"
	"golang.org/x/sys/unix"
)

const (
	ctdSocket        = "/run/containerd/containerd.sock"
	ctdNamespace     = "firecracker"
	stopTimeoutLabel = "IgniteStopTimeout"
	logPathTemplate  = "/tmp/%s.log"
)

// ctdClient is a runtime.Interface
// implementation serving the containerd client
type ctdClient struct {
	ctx    context.Context
	client *containerd.Client
}

var _ runtime.Interface = &ctdClient{}

// GetContainerdClient builds a client for talking to containerd
func GetContainerdClient() (*ctdClient, error) {
	cli, err := containerd.New(ctdSocket)
	if err != nil {
		return nil, err
	}

	return &ctdClient{
		client: cli,
		ctx:    namespaces.WithNamespace(context.Background(), ctdNamespace),
	}, nil
}

func (cc *ctdClient) PullImage(image meta.OCIImageRef) error {
	log.Debugf("containerd: Pulling image %q", image)
	_, err := cc.client.Pull(cc.ctx, image.Normalized(), containerd.WithPullUnpack)
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

	info, err := cont.Info(cc.ctx)
	if err != nil {
		return nil, err
	}

	return &runtime.ContainerInspectResult{
		ID:        info.ID,
		Image:     info.Image, // TODO: This may be incorrect
		Status:    "",         // TODO: This
		IPAddress: nil,        // TODO: This, containerd only supports CNI
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

func (cc *ctdClient) RunContainer(image meta.OCIImageRef, config *runtime.ContainerConfig, name string) (s string, err error) {
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
	config.Binds = append(config.Binds, runtime.BindBoth("/etc/resolv.conf"))

	// Add the stop timeout as a label, as containerd doesn't natively support it
	config.Labels[stopTimeoutLabel] = strconv.FormatUint(uint64(config.StopTimeout), 10)

	// Build the OCI specification
	opts := []oci.SpecOpts{
		oci.WithDefaultSpec(),
		oci.WithDefaultUnixDevices,
		oci.WithTTY,
		oci.WithImageConfigArgs(img, config.Cmd),
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
		// TODO: Upgrade to v2
		containerd.WithRuntime(plugin.RuntimeRuncV1, nil),
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
				Path:  dev.HostPath,
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
		return
	}

	task, err := cont.Task(cc.ctx, cio.Load)
	if err != nil {
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

func (cc *ctdClient) RemoveContainer(container string) (err error) {
	// Remove the container if it exists
	var cont containerd.Container
	var task containerd.Task
	if cont, err = cc.client.LoadContainer(cc.ctx, container); ifFound(&err) {
		// Load the container's task without attaching
		if task, err = cont.Task(cc.ctx, nil); ifFound(&err) {
			_, err = task.Delete(cc.ctx)
		}

		// Delete the container
		if err == nil {
			err = cont.Delete(cc.ctx, containerd.WithSnapshotCleanup)
		}

		// Remove the log file if it exists
		logFile := fmt.Sprintf(logPathTemplate, container)
		if util.FileExists(logFile) && err == nil {
			err = os.RemoveAll(logFile)
		}
	}

	return
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

func (cc *ctdClient) PreflightChecker() preflight.Checker {
	return checkers.NewExistingFileChecker(ctdSocket)
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

// ifFound is a helper for functions returning errdefs.ErrNotFound in if-statements.
// If in points to that error, the value is changed to nil and false is returned.
// If in points to any other error, no change is performed and false is returned.
// If in points to nil, no change is performed and true is returned.
func ifFound(in *error) bool {
	if in == nil {
		panic("nil pointer given to ifFound")
	}

	// If the given error is an errdefs.ErrNotFound, clear it
	if errdefs.IsNotFound(*in) {
		*in = nil
		return false
	}

	return *in == nil
}
