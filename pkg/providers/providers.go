package providers

// Populate initializes all providers
func Populate() error {
	// NOTE: This initialization is order-dependent!
	// E.g. the network plugin depends on the runtime.
	providers := []func() error{
		SetDockerRuntime,    // Use the Docker runtime
		SetCNINetworkPlugin, // Use the CNI Network plugin
		SetCachedStorage,    // Use a generic storage implementation backed by a cache
		//SetTestManifestStorage, // Use the test version of the Manifest Storage
		SetClient, // Set the globally available client
	}

	for _, init := range providers {
		if err := init(); err != nil {
			return err
		}
	}

	return nil
}
