package ignite

import (
	"github.com/weaveworks/ignite/pkg/providers"
	clientprovider "github.com/weaveworks/ignite/pkg/providers/client"
	"github.com/weaveworks/ignite/pkg/providers/network"
	"github.com/weaveworks/ignite/pkg/providers/runtime"
	storageprovider "github.com/weaveworks/ignite/pkg/providers/storage"
)

// Preload providers need to be loaded before flag parsing has finished
var Preload = []providers.ProviderInitFunc{
	storageprovider.SetGenericStorage, // Use a generic storage implementation backed by a cache
	clientprovider.SetClient,          // Set the globally available client
}

// NOTE: Provider initialization is order-dependent!
// E.g. the network plugin depends on the runtime.
var Providers = []providers.ProviderInitFunc{
	runtime.SetRuntime,       // Set the selected runtime
	network.SetNetworkPlugin, // Set the selected network plugin
}
