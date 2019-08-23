package containerd

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	"github.com/containerd/console"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/containers"
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
	"github.com/weaveworks/ignite/pkg/runtime"
	"github.com/weaveworks/ignite/pkg/util"
	"golang.org/x/sys/unix"
)

const (
	ctdSocket    = "/run/containerd/containerd.sock"
	ctdNamespace = "firecracker"
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

	t, err := cont.Task(cc.ctx, nil)
	if err != nil {
		return nil, err
	}

	info, err := cont.Info(cc.ctx)
	if err != nil {
		return nil, err
	}

	pids, err := t.Pids(cc.ctx)
	if err != nil {
		return nil, err
	}

	if len(pids) == 0 {
		return nil, fmt.Errorf("no running tasks found for container %q", container)
	}

	return &runtime.ContainerInspectResult{
		ID:        info.ID,
		Image:     info.Image,  // TODO: This may be incorrect
		Status:    "",          // TODO: This
		IPAddress: nil,         // TODO: This, containerd only supports CNI
		PID:       pids[0].Pid, // TODO: This should respect multiple tasks, we need a way to identify ignite-spawn
	}, nil
}

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
		task    containerd.Task
		statusC <-chan containerd.ExitStatus
	)

	if task, err = cont.Task(cc.ctx, cio.NewAttach(cio.WithStdio)); err != nil {
		return
	}
	defer util.DeferErr(&err, func() error { _, err := task.Delete(cc.ctx); return err })

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

	ec := <-statusC
	code, _, err := ec.Result()
	if err != nil {
		return
	}

	if code != 0 {
		err = fmt.Errorf("attach exited with code %d", code)
	}

	return
}

func (cc *ctdClient) RunContainer(image meta.OCIImageRef, config *runtime.ContainerConfig, name string) (s string, err error) {
	img, err := cc.client.GetImage(cc.ctx, image.Normalized())
	if err != nil {
		return
	}

	// TODO: Fix this, simulates the Docker "entrypoint"
	config.Cmd = append([]string{"ignite-spawn"}, config.Cmd...)

	// Add the /etc/resolv.conf mount, this isn't done automatically by containerd
	config.Binds = append(config.Binds, runtime.BindBoth("/etc/resolv.conf"))

	// Build the OCI specification
	// TODO: Refine this, add missing options
	opts := []oci.SpecOpts{
		oci.WithDefaultSpec(),
		oci.WithDefaultUnixDevices,
		oci.WithImageConfig(img),
		oci.WithProcessArgs(config.Cmd...),
		withAddedCaps(config.CapAdds),
		withHostname(config.Hostname),
		withMounts(config.Binds),
		withDevices(config.Devices),
	}

	tty := true
	if tty {
		opts = append(opts, oci.WithTTY)
	}

	// TODO: Handle CapAdd & Hostname
	// Known limitations, containerd doesn't support the following config fields
	// StopTimeout
	// AutoRemove
	// NetworkMode (only CNI supported)
	// PortBindings

	cOpts := []containerd.NewContainerOpts{
		containerd.WithImage(img),
		// Even when "readonly" is set, we don't use KindView snapshot here. (#1495)
		// We pass writable snapshot to the OCI runtime, and the runtime remounts it as read-only,
		// after creating some mount points on demand.
		//containerd.WithSnapshot(name),
		containerd.WithNewSnapshot(name, img),
		containerd.WithImageStopSignal(img, "SIGTERM"),
		containerd.WithNewSpec(opts...),
		// TODO: Upgrade to v2
		containerd.WithRuntime(plugin.RuntimeRuncV1, nil),
		containerd.WithContainerLabels(config.Labels),
	}

	cont, err := cc.client.NewContainer(cc.ctx, name, cOpts...)
	if err != nil {
		return
	}

	var con console.Console
	if tty {
		con = console.Current()
		defer con.Reset()
		if err = con.SetRaw(); err != nil {
			return
		}
	}

	/*stdinC := &stdinCloser{
		stdin: os.Stdin,
	}*/
	ioOpts := []cio.Opt{cio.WithFIFODir("/run/containerd/fifo")}
	//stdio := cio.NewCreator(append([]cio.Opt{cio.WithStreams(stdinC, os.Stdout, os.Stderr)}, ioOpts...)...)
	ioCreator := cio.NewCreator(append([]cio.Opt{cio.WithStreams(con, con, nil), cio.WithTerminal}, ioOpts...)...)

	/*task, err := tasks.NewTask(ctx, client, container, context.String("checkpoint"), con, context.Bool("null-io"), context.String("log-uri"), ioOpts, opts...)
	if err != nil {
		return err
	}*/

	// TODO: Set up TTY and stdio streams
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
	mounts := make([]specs.Mount, 0)
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

	task, err := cont.Task(cc.ctx, cio.Load)
	if err != nil {
		return
	}

	_, err = task.Delete(cc.ctx) // TODO: This probably force-kills, also handle the exit status
	return
}

func (cc *ctdClient) KillContainer(container, signal string) error {
	return cc.StopContainer(container, nil) // TODO: Handle this separately
}

func (cc *ctdClient) RemoveContainer(container string) error {
	cont, err := cc.client.LoadContainer(cc.ctx, container)
	if err != nil {
		return err
	}

	return cont.Delete(cc.ctx)
}

func (cc *ctdClient) ContainerLogs(container string) (io.ReadCloser, error) {
	// TODO: Implement logs for containerd
	return nil, unsupported("logs")
}

func (cc *ctdClient) RawClient() interface{} {
	return cc.client
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

func unsupported(feature string) error {
	return fmt.Errorf("containerd: %q is currently unsupported", feature)
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
