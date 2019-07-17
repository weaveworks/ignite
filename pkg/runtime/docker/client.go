package docker

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/weaveworks/ignite/pkg/runtime"
)

const dockerNetNSFmt = "/proc/%v/ns/net"

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

func (dc *dockerClient) RemoveContainer(container string) error {
	return dc.client.ContainerRemove(context.Background(), container, types.ContainerRemoveOptions{})
}

func (dc *dockerClient) StopContainer(container string, timeout *time.Duration) error {
	return dc.client.ContainerStop(context.Background(), container, timeout)
}

func (dc *dockerClient) KillContainer(container, signal string) error {
	return dc.client.ContainerKill(context.Background(), container, signal)
}

// GetNetNS returns the network namespace of the given containerID. The ID
// supplied is typically the ID of a pod sandbox. This getter doesn't try
// to map non-sandbox IDs to their respective sandboxes.
func (dc *dockerClient) GetNetNS(podSandboxID string) (string, error) {
	r, err := dc.client.ContainerInspect(context.TODO(), podSandboxID)
	if err != nil {
		return "", err
	}

	return getNetworkNamespace(&r)
}

func getNetworkNamespace(c *types.ContainerJSON) (string, error) {
	if c.State.Pid == 0 {
		// Docker reports pid 0 for an exited container.
		return "", fmt.Errorf("cannot find network namespace for the terminated container %q", c.ID)
	}

	return fmt.Sprintf(dockerNetNSFmt, c.State.Pid), nil
}
