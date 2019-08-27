package network

import (
	"fmt"

	"github.com/weaveworks/ignite/pkg/network"
	"github.com/weaveworks/ignite/pkg/providers"
	cniprovider "github.com/weaveworks/ignite/pkg/providers/cni"
	dockerprovider "github.com/weaveworks/ignite/pkg/providers/docker"
)

func SetNetworkPlugin() error {
	switch providers.NetworkPluginName {
	case network.PluginDockerBridge:
		return dockerprovider.SetDockerNetwork() // Use the Docker bridge network
	case network.PluginCNI:
		return cniprovider.SetCNINetworkPlugin() // Use the CNI Network plugin
	}

	return fmt.Errorf("unknown network plugin %q", providers.NetworkPluginName)
}
