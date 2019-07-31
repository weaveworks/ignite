package spawn

import (
	"github.com/weaveworks/ignite/pkg/providers"
	clientprovider "github.com/weaveworks/ignite/pkg/providers/client"
	storageprovider "github.com/weaveworks/ignite/pkg/providers/storage"
)

// NOTE: Provider initialization is order-dependent!
// E.g. the network plugin depends on the runtime.
var Providers = []providers.ProviderInitFunc{
	storageprovider.SetGenericStorage, // Use a generic storage implementation backed by a cache
	clientprovider.SetClient,          // Set the globally available client
}
