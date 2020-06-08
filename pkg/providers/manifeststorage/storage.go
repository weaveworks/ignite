package manifeststorage

import (
	log "github.com/sirupsen/logrus"
	"github.com/weaveworks/ignite/pkg/apis/ignite/scheme"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/providers"
	"github.com/weaveworks/libgitops/pkg/storage/cache"
	"github.com/weaveworks/libgitops/pkg/storage/manifest"
)

var ManifestStorage *manifest.ManifestStorage

func SetManifestStorage() (err error) {
	log.Trace("Initializing the ManifestStorage provider...")
	ManifestStorage, err = manifest.NewTwoWayManifestStorage(constants.MANIFEST_DIR, constants.DATA_DIR, scheme.Serializer)
	if err != nil {
		return
	}

	providers.Storage = cache.NewCache(ManifestStorage)
	return
}
