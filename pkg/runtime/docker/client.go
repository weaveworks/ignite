package docker

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	cont "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/runtime"
	"github.com/weaveworks/ignite/pkg/util"
)

const (
	dockerNetNSFmt = "/proc/%v/ns/net"
	portFormat     = "%d/tcp" // TODO: Support protocols other than TCP
)

// dockerClient is a runtime.Interface
// implementation serving the Docker client
type dockerClient struct {
	client *client.Client
}

var _ runtime.Interface = &dockerClient{}

// GetDockerClient builds a client for talking to docker
func GetDockerClient() (runtime.Interface, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion("1.35"))
	if err != nil {
		return nil, err
	}

	return &dockerClient{
		client: cli,
	}, nil
}

func (dc *dockerClient) RawClient() interface{} {
	return dc.client
}

func (dc *dockerClient) InspectImage(image string) (*runtime.ImageInspectResult, error) {
	res, _, err := dc.client.ImageInspectWithRaw(context.Background(), image)
	if err != nil {
		return nil, err
	}

	return &runtime.ImageInspectResult{
		ID:          res.ID,
		RepoDigests: res.RepoDigests,
		Size:        res.Size,
	}, nil
}

func (dc *dockerClient) PullImage(image string) (io.ReadCloser, error) {
	return dc.client.ImagePull(context.Background(), image, types.ImagePullOptions{})
}

func (dc *dockerClient) ExportImage(image string) (io.ReadCloser, string, error) {
	config, err := dc.client.ContainerCreate(context.Background(), &container.Config{
		Cmd:   []string{"sh"}, // We need a temporary command, this doesn't need to exist in the image
		Image: image,
	}, nil, nil, "")
	if err != nil {
		return nil, "", err
	}

	rc, err := dc.client.ContainerExport(context.Background(), config.ID)
	return rc, config.ID, err
}

func (dc *dockerClient) InspectContainer(container string) (*runtime.ContainerInspectResult, error) {
	res, _, err := dc.client.ContainerInspectWithRaw(context.Background(), container, false)
	if err != nil {
		return nil, err
	}

	return &runtime.ContainerInspectResult{
		ID:     res.ID,
		Image:  res.Image,
		Status: res.State.Status,
	}, nil
}

func (dc *dockerClient) AttachContainer(container string) (err error) {
	// TODO: Rework to perform the attach via the Docker client,
	// this will require manual TTY and signal emulation/handling.
	// Implement the pseudo-TTY and remove this call, see
	// https://github.com/weaveworks/ignite/pull/211#issuecomment-512809841
	var ec int
	if ec, err = util.ExecForeground("docker", "attach", container); err != nil {
		if ec == 1 { // Docker's detach sequence (^P^Q) has an exit code of -1
			err = nil // Don't treat it as an error
		}
	}

	return
}

func (dc *dockerClient) RunContainer(image string, config *runtime.ContainerConfig, name string) (string, error) {
	portBindings := make(nat.PortMap)
	for _, portMapping := range config.PortBindings {
		var hostIP string
		if portMapping.BindAddress != nil {
			hostIP = portMapping.BindAddress.String()
		}

		protocol := portMapping.Protocol
		if len(protocol) == 0 {
			// Docker uses TCP by default
			protocol = meta.ProtocolTCP
		}

		portBindings[nat.Port(fmt.Sprintf("%d/%s", portMapping.VMPort, protocol.String()))] = []nat.PortBinding{
			{
				HostIP:   hostIP,
				HostPort: fmt.Sprintf(portFormat, portMapping.HostPort),
			},
		}
	}

	binds := make([]string, 0, len(config.Binds))
	for _, bind := range config.Binds {
		binds = append(binds, fmt.Sprintf("%s:%s", bind.HostPath, bind.ContainerPath))
	}

	devices := make([]container.DeviceMapping, 0, len(config.Devices))
	for _, device := range config.Devices {
		devices = append(devices, container.DeviceMapping{
			PathOnHost:        device.HostPath,
			PathInContainer:   device.ContainerPath,
			CgroupPermissions: "rwm",
		})
	}

	stopTimeout := int(config.StopTimeout)

	c, err := dc.client.ContainerCreate(context.Background(), &container.Config{
		Hostname:    config.Hostname,
		Tty:         true, // --tty
		OpenStdin:   true, // --interactive
		Cmd:         config.Cmd,
		Image:       image,
		Labels:      config.Labels,
		StopTimeout: &stopTimeout,
	}, &container.HostConfig{
		Binds:        binds,
		NetworkMode:  container.NetworkMode(config.NetworkMode),
		PortBindings: portBindings,
		AutoRemove:   config.AutoRemove,
		CapAdd:       config.CapAdds,
		Resources: container.Resources{
			Devices: devices,
		},
	}, nil, name)
	if err != nil {
		return "", err
	}

	return c.ID, dc.client.ContainerStart(context.Background(), c.ID, types.ContainerStartOptions{})
}

func (dc *dockerClient) StopContainer(container string, timeout *time.Duration) error {
	if err := dc.client.ContainerStop(context.Background(), container, timeout); err != nil {
		return err
	}

	// Wait for the container to be stopped
	return dc.waitForContainer(container, cont.WaitConditionNotRunning, nil)
}

func (dc *dockerClient) KillContainer(container, signal string) error {
	if err := dc.client.ContainerKill(context.Background(), container, signal); err != nil {
		return err
	}

	// Wait for the container to be killed
	return dc.waitForContainer(container, cont.WaitConditionNotRunning, nil)
}

func (dc *dockerClient) RemoveContainer(container string) error {
	// Container waiting can fail if the container gets removed before
	// we reach the waiting fence. Start the waiter in a goroutine before
	// performing the removal to ensure proper removal detection.
	errC := make(chan error)
	readyC := make(chan bool)
	go func() {
		errC <- dc.waitForContainer(container, cont.WaitConditionRemoved, &readyC)
	}()

	<-readyC // The ready channel is used to wait until removal detection has started
	if err := dc.client.ContainerRemove(context.Background(), container, types.ContainerRemoveOptions{}); err != nil {
		return err
	}

	// Wait for the container to be removed
	return <-errC
}

func (dc *dockerClient) ContainerLogs(container string) (io.ReadCloser, error) {
	return dc.client.ContainerLogs(context.Background(), container, types.ContainerLogsOptions{
		ShowStdout: true, // We only need stdout, as TTY mode merges stderr into it
	})
}

// ContainerNetNS returns the network namespace of the given container.
func (dc *dockerClient) ContainerNetNS(container string) (string, error) {
	r, err := dc.client.ContainerInspect(context.TODO(), container)
	if err != nil {
		return "", err
	}

	return getNetworkNamespace(&r)
}

func (dc *dockerClient) waitForContainer(container string, condition cont.WaitCondition, readyC *chan bool) error {
	resultC, errC := dc.client.ContainerWait(context.Background(), container, condition)

	// The ready channel can be used to wait until
	// the container wait request has been sent to
	// Docker to ensure proper ordering
	if readyC != nil {
		*readyC <- true
	}

	select {
	case result := <-resultC:
		if result.Error != nil {
			return fmt.Errorf("failed to wait for container %q: %s", container, result.Error.Message)
		}
	case err := <-errC:
		return err
	}

	return nil
}

func getNetworkNamespace(c *types.ContainerJSON) (string, error) {
	if c.State.Pid == 0 {
		// Docker reports pid 0 for an exited container.
		return "", fmt.Errorf("cannot find network namespace for the terminated container %q", c.ID)
	}

	return fmt.Sprintf(dockerNetNSFmt, c.State.Pid), nil
}
