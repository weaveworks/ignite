package ignite

import (
	"github.com/weaveworks/ignite/pkg/providers"
	clientprovider "github.com/weaveworks/ignite/pkg/providers/client"
	cniprovider "github.com/weaveworks/ignite/pkg/providers/cni"
	dockerprovider "github.com/weaveworks/ignite/pkg/providers/docker"
	storageprovider "github.com/weaveworks/ignite/pkg/providers/storage"
)

// NOTE: Provider initialization is order-dependent!
// E.g. the network plugin depends on the runtime.
var Providers = []providers.ProviderInitFunc{
	dockerprovider.SetDockerRuntime,   // Use the Docker runtime
	cniprovider.SetCNINetworkPlugin,   // Use the CNI Network plugin
	storageprovider.SetGenericStorage, // Use a generic storage implementation backed by a cache
	clientprovider.SetClient,          // Set the globally available client
}
