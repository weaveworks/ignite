package network

// Plugin describes a generic network plugin
type Plugin interface {
	// Name returns the network plugin's name.
	Name() string

	// SetupContainerNetwork sets up the networking for a container
	SetupContainerNetwork(containerID string) error

	// RemoveContainerNetwork is the method called before a container using the network plugin can be deleted
	RemoveContainerNetwork(containerID string) error

	// Status returns error if the network plugin is in error state
	Status() error
}
