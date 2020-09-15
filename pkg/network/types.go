package network

import (
	"fmt"
	"net"

	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/runtime"
)

// Plugin describes a generic network plugin
type Plugin interface {
	// Name returns the network plugin's name.
	Name() PluginName

	// PrepareContainerSpec sets any needed options on the container spec before starting the container
	PrepareContainerSpec(container *runtime.ContainerConfig) error

	// SetupContainerNetwork sets up the networking for a container
	// This is ran _after_ the container has been started
	SetupContainerNetwork(containerID string, portmappings ...meta.PortMapping) (*Result, error)

	// RemoveContainerNetwork is the method called before a container using the network plugin can be deleted
	RemoveContainerNetwork(containerID string, portmappings ...meta.PortMapping) error
}

type Result struct {
	Addresses []Address
}

type Address struct {
	IP      net.IP
	Gateway net.IP
}

// PluginName defines a name for a network plugin
type PluginName string

var _ fmt.Stringer = PluginName("")

func (pn PluginName) String() string {
	return string(pn)
}

const (
	// PluginCNI specifies the network mode where CNI is used
	PluginCNI PluginName = "cni"
	// PluginDockerBridge specifies the default docker bridge network is used
	PluginDockerBridge PluginName = "docker-bridge"
)

// ListPlugins gets the list of available network plugins
func ListPlugins() []PluginName {
	return []PluginName{
		PluginCNI,
		PluginDockerBridge,
	}
}
