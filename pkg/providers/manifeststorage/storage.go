package manifeststorage

import (
	log "github.com/sirupsen/logrus"
	"github.com/weaveworks/ignite/pkg/providers"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/storage/cache"
	"github.com/weaveworks/ignite/pkg/storage/manifest"
)

func SetManifestStorage() error {
	log.Trace("Initializing the ManifestStorage provider...")
	ms, err := manifest.NewManifestStorage(constants.MANIFEST_DIR)
	if err != nil {
		return err
	}
	providers.Storage = cache.NewCache(ms)
	return nil
}
