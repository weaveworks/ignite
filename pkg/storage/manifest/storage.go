package manifest

import (
	"github.com/weaveworks/ignite/pkg/apis/ignite/scheme"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/storage"
	"github.com/weaveworks/ignite/pkg/storage/manifest/raw"
	"github.com/weaveworks/ignite/pkg/storage/sync"
	"github.com/weaveworks/ignite/pkg/storage/watch"
	"github.com/weaveworks/ignite/pkg/storage/watch/update"
	"os"
)

// TODO: Re-implement this with SyncStorage and an update aggregator
func NewManifestStorage(dataDir string) (*ManifestStorage, error) {
	dataDir, err := resolveLink(dataDir)
	if err != nil {
		return nil, err
	}

	ws, err := watch.NewGenericWatchStorage(storage.NewGenericStorage(raw.NewManifestRawStorage(dataDir), scheme.Serializer))
	if err != nil {
		return nil, err
	}

	ss := sync.NewSyncStorage(
		storage.NewGenericStorage(
			storage.NewDefaultRawStorage(constants.DATA_DIR), scheme.Serializer),
		ws)

	s := &ManifestStorage{
		Storage: ss,
	}

	go s.aggregate()

	return s, nil
}

type UpdateCache []update.Update

// ManifestStorage implements the storage interface for GitOps purposes
type ManifestStorage struct {
	storage.Storage
	cache UpdateCache
}

func (s *ManifestStorage) Sync() (c UpdateCache) {
	c, s.cache = s.cache, nil
	return
}

func (s *ManifestStorage) aggregate() {
	updateStream := s.Storage.(*sync.SyncStorage).GetUpdateStream()

	for {
		if upd, ok := <-updateStream; ok {
			s.cache = append(s.cache, upd)
		} else {
			return
		}
	}
}

func resolveLink(path string) (string, error) {
	fi, err := os.Lstat(path)
	if err != nil {
		return "", err
	}

	if fi.Mode()&os.ModeSymlink == 0 {
		// The path is not a symlink, return it as is
		return path, nil
	}

	// The target is a symlink
	return os.Readlink(path)
}
