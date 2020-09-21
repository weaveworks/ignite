package docker

import (
	"fmt"
	"net"

	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/network"
	"github.com/weaveworks/ignite/pkg/runtime"
)

type dockerNetworkPlugin struct {
	runtime runtime.Interface
}

func GetDockerNetworkPlugin(r runtime.Interface) network.Plugin {
	return &dockerNetworkPlugin{r}
}

func (*dockerNetworkPlugin) Name() network.PluginName {
	return network.PluginDockerBridge
}

func (*dockerNetworkPlugin) PrepareContainerSpec(_ *runtime.ContainerConfig) error {
	// no-op, we don't need to set any special parameters on the container
	return nil
}

func (plugin *dockerNetworkPlugin) SetupContainerNetwork(containerID string, _ ...meta.PortMapping) (*network.Result, error) {
	// This is used to fetch the IP address the runtime gives to the VM container
	result, err := plugin.runtime.InspectContainer(containerID)
	if err != nil {
		return nil, fmt.Errorf("failed to inspect container %s: %v", containerID, err)
	}

	return &network.Result{
		Addresses: []network.Address{
			{
				IP: result.IPAddress,
				// TODO: Make this auto-detect if the gateway is not using the standard setup
				Gateway: net.IPv4(result.IPAddress[0], result.IPAddress[1], result.IPAddress[2], 1),
			},
		},
	}, nil
}

func (*dockerNetworkPlugin) RemoveContainerNetwork(_ string, _ ...meta.PortMapping) error {
	// no-op for docker, this is handled automatically
	return nil
}
