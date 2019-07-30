package manifest

import (
	"github.com/weaveworks/ignite/pkg/apis/ignite/scheme"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/storage"
	"github.com/weaveworks/ignite/pkg/storage/manifest/raw"
	"github.com/weaveworks/ignite/pkg/storage/sync"
	"github.com/weaveworks/ignite/pkg/storage/watch"
)

func NewManifestStorage(dataDir string) (*ManifestStorage, error) {

	ws, err := watch.NewGenericWatchStorage(storage.NewGenericStorage(raw.NewGenericMappedRawStorage(dataDir), scheme.Serializer))
	if err != nil {
		return nil, err
	}

	ss := sync.NewSyncStorage(
		storage.NewGenericStorage(
			storage.NewDefaultRawStorage(constants.DATA_DIR), scheme.Serializer),
		ws)

	return &ManifestStorage{
		Storage: ss,
	}, nil
}

// ManifestStorage implements the storage interface for GitOps purposes
type ManifestStorage struct {
	storage.Storage
}

// GetUpdateStream gets the channel with updates
func (s *ManifestStorage) GetUpdateStream() sync.UpdateStream {
	return s.Storage.(*sync.SyncStorage).GetUpdateStream()
}
