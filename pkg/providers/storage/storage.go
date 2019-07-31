package providers

import (
	log "github.com/sirupsen/logrus"
	"github.com/weaveworks/ignite/pkg/apis/ignite/scheme"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/providers"
	"github.com/weaveworks/ignite/pkg/storage"
	"github.com/weaveworks/ignite/pkg/storage/cache"
	"github.com/weaveworks/ignite/pkg/storage/manifest"
)

func SetGenericStorage() error {
	log.Trace("Initializing the GenericStorage provider...")
	providers.Storage = cache.NewCache(
		storage.NewGenericStorage(
			storage.NewDefaultRawStorage(constants.DATA_DIR), scheme.Serializer))
	return nil
}

func SetManifestStorage() error {
	log.Trace("Initializing the ManifestStorage provider...")
	ms, err := manifest.NewManifestStorage("/etc/firecracker/manifests")
	if err != nil {
		return err
	}
	providers.Storage = cache.NewCache(ms)
	return nil
}
