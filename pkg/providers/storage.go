package providers

import (
	"github.com/weaveworks/ignite/pkg/apis/ignite/scheme"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/storage"
	"github.com/weaveworks/ignite/pkg/storage/cache"
	"github.com/weaveworks/ignite/pkg/storage/manifest"
	"github.com/weaveworks/ignite/pkg/storage/sync"
	"github.com/weaveworks/ignite/pkg/storage/watch"
)

// Storage is the default storage implementation
var Storage storage.Storage
var SS *sync.SyncStorage

func SetCachedStorage() error {
	Storage = cache.NewCache(
		storage.NewGenericStorage(
			storage.NewDefaultRawStorage(constants.DATA_DIR), scheme.Serializer))
	return nil
}

// TODO: Special constructor for this setup
func SetTestManifestStorage() error {
	ws, err := watch.NewGenericWatchStorage(storage.NewGenericStorage(manifest.NewManifestRawStorage("/etc/firecracker/manifests"), scheme.Serializer))
	if err != nil {
		return err
	}

	ss := sync.NewSyncStorage(
		storage.NewGenericStorage(storage.NewDefaultRawStorage(constants.DATA_DIR), scheme.Serializer),
		ws)

	SS = ss.(*sync.SyncStorage)

	Storage = cache.NewCache(ss)

	return nil
}
