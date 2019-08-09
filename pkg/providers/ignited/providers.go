package ignited

import (
	"github.com/weaveworks/ignite/pkg/providers"
	clientprovider "github.com/weaveworks/ignite/pkg/providers/client"
	cniprovider "github.com/weaveworks/ignite/pkg/providers/cni"
	dockerprovider "github.com/weaveworks/ignite/pkg/providers/docker"
	manifeststorageprovider "github.com/weaveworks/ignite/pkg/providers/manifeststorage"
)

// NOTE: Provider initialization is order-dependent!
// E.g. the network plugin depends on the runtime.
var Providers = []providers.ProviderInitFunc{
	dockerprovider.SetDockerRuntime,            // Use the Docker runtime
	dockerprovider.SetDockerNetwork,            // Use the Docker bridge network
	cniprovider.SetCNINetworkPlugin,            // Use the CNI Network plugin
	manifeststorageprovider.SetManifestStorage, // Use the ManifestStorage implementation, backed by a cache
	clientprovider.SetClient,                   // Set the globally available client
}
