package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/weaveworks/ignite/pkg/containerruntime"
)

const dockerNetNSFmt = "/proc/%v/ns/net"

// GetDockerClient builds a client for talking to docker
func GetDockerClient() (containerruntime.Interface, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}
	return &dockerClient{
		client: cli,
	}, nil
}

type dockerClient struct {
	client *client.Client
}

func (dc *dockerClient) RawClient() interface{} {
	return dc.client
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
