package providers

import (
	"github.com/weaveworks/ignite/pkg/apis/ignite/scheme"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/storage"
	"github.com/weaveworks/ignite/pkg/storage/cache"
)

// Storage is the default storage implementation
var Storage storage.Storage

func SetCachedStorage() error {
	Storage = cache.NewCache(
		storage.NewGenericStorage(
			storage.NewDefaultRawStorage(constants.DATA_DIR), scheme.Serializer))
	return nil
}
