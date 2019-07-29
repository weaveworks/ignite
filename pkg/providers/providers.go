package providers

// NOTE: Provider initialization is order-dependent!
// E.g. the network plugin depends on the runtime.
var Providers = []func() error{
	SetDockerRuntime,    // Use the Docker runtime
	SetCNINetworkPlugin, // Use the CNI Network plugin
	SetCachedStorage,    // Use a generic storage implementation backed by a cache
	SetClient,           // Set the globally available client
}

// `ignite daemon` overwrites/re-initializes the Storage and Client providers
var DaemonProviders = []func() error{
	SetManifestStorage, // Use the ManifestStorage implementation
	SetClient,          // Set the globally available client
}

// Populate initializes all providers
func Populate(providers []func() error) error {

	for _, init := range providers {
		if err := init(); err != nil {
			return err
		}
	}

	return nil
}
