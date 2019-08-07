package manifeststorage

import (
	log "github.com/sirupsen/logrus"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/providers"
	"github.com/weaveworks/ignite/pkg/storage/cache"
	"github.com/weaveworks/ignite/pkg/storage/manifest"
)

var ManifestStorage *manifest.ManifestStorage

func SetManifestStorage() (err error) {
	log.Trace("Initializing the ManifestStorage provider...")
	ManifestStorage, err = manifest.NewManifestStorage(constants.MANIFEST_DIR)
	if err != nil {
		return
	}

	providers.Storage = cache.NewCache(ManifestStorage)
	return
}
