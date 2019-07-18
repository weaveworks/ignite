package providers

// Populate initializes all providers
func Populate() error {
	// NOTE: This initialization is order-dependent!
	// E.g. the network plugin depends on the runtime.
	providers := []func() error{
		SetDockerRuntime,    // Use the Docker runtime
		SetCNINetworkPlugin, // Use the CNI Network plugin
		SetCachedStorage,    // Use a generic storage implementation backed by a cache
		SetClient,           // Set the globally available client
	}

	for _, init := range providers {
		if err := init(); err != nil {
			return err
		}
	}

	return nil
}
