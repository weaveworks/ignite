package network

import (
	"net"

	"github.com/weaveworks/ignite/pkg/runtime"
)

// Plugin describes a generic network plugin
type Plugin interface {
	// Name returns the network plugin's name.
	Name() string

	// PrepareContainerSpec sets any needed options on the container spec before starting the container
	PrepareContainerSpec(container *runtime.ContainerConfig) error

	// SetupContainerNetwork sets up the networking for a container
	// This is ran _after_ the container has been started
	SetupContainerNetwork(containerID string) (*Result, error)

	// RemoveContainerNetwork is the method called before a container using the network plugin can be deleted
	RemoveContainerNetwork(containerID string) error

	// Status returns error if the network plugin is in error state
	Status() error
}

type Result struct {
	Addresses []Address
}

type Address struct {
	net.IPNet
	Gateway net.IP
}
