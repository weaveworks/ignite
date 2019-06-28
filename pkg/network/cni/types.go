package cni

const (
	// CNIPluginName describes the name of the CNI network plugin
	CNIPluginName = "cni"
	// DefaultInterfaceName describes the interface name that the CNI network plugin will set up
	DefaultInterfaceName = "eth0"
	// CNIBinDir describes the directory where the CNI binaries are stored
	CNIBinDir = "/opt/cni/bin"
	// CNIConfDir describes the directory where the CNI plugin's configuration is stored
	// TODO: CNIBinDir and CNIConfDir should maybe be globally configurable?
	CNIConfDir = "/etc/cni/net.d"

	loopbackCNIConfig = `{
	"cniVersion": "0.2.0",
	"name": "cni-loopback",
	"plugins":[{
		"type": "loopback"
	}]
}`
)

// NetworkPlugin describes a network plugin for docker
type NetworkPlugin interface {
	// Name returns the network plugin's name.
	Name() string

	// SetupContainerNetwork sets up the networking for a Docker container using CNI
	SetupContainerNetwork(containerID string) error

	// RemoveContainerNetwork is the method called before a container using the CNI network can be deleted
	RemoveContainerNetwork(containerID string) error

	// Status returns error if the network plugin is in error state
	Status() error
}
