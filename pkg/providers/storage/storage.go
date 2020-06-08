package storage

import (
	log "github.com/sirupsen/logrus"
	"github.com/weaveworks/ignite/pkg/apis/ignite/scheme"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/providers"
	"github.com/weaveworks/libgitops/pkg/storage"
	"github.com/weaveworks/libgitops/pkg/storage/cache"
)

func SetGenericStorage() error {
	log.Trace("Initializing the GenericStorage provider...")
	providers.Storage = cache.NewCache(
		storage.NewGenericStorage(
			storage.NewGenericRawStorage(constants.DATA_DIR), scheme.Serializer))
	return nil
}
