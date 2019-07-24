package sync

import (
	"fmt"

	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/storage"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// SyncStorage is a Storage implementation taking in multiple Storages and
// keeping them in sync. Any write operation executed on the SyncStorage
// is propagated to all of the Storages it manages (including the embedded
// one). For any retrieval or generation operation, the embedded Storage
// will be used (it is treated as read-write). As all other Storages only
// receive write operations, they can be thought of as write-only.
type SyncStorage struct {
	storage.Storage
	storages []storage.Storage
	watchers []*watcher
}

var _ storage.Storage = &SyncStorage{}

// NewSyncStorage constructs a new SyncStorage
func NewSyncStorage(rwStorage storage.Storage) storage.Storage {
	return &SyncStorage{Storage: rwStorage}
}

func (s *SyncStorage) Add(wStorages ...storage.Storage) storage.Storage {
	s.storages = append(s.storages, wStorages...)
	return s
}

func (s *SyncStorage) AddWatched(wStorages ...storage.Storage) (storage.Storage, error) {
	for _, st := range wStorages {
		w, err := newWatcher(st, s.update)
		if err != nil {
			return nil, err
		}

		s.watchers = append(s.watchers, w)
	}

	return s.Add(wStorages...), nil
}

// Set is propagated to all Storages
func (s *SyncStorage) Set(gvk schema.GroupVersionKind, obj meta.Object) error {
	return s.runAll(func(s storage.Storage) error {
		return s.Set(gvk, obj)
	})
}

// Patch is propagated to all Storages
func (s *SyncStorage) Patch(gvk schema.GroupVersionKind, uid meta.UID, patch []byte) error {
	return s.runAll(func(s storage.Storage) error {
		return s.Patch(gvk, uid, patch)
	})
}

// Delete is propagated to all Storages
func (s *SyncStorage) Delete(gvk schema.GroupVersionKind, uid meta.UID) error {
	return s.runAll(func(s storage.Storage) error {
		return s.Delete(gvk, uid)
	})
}

type callFunc func(storage.Storage) error

// runAll runs the given callFunc for all Storages in parallel and aggregates all errors
func (s *SyncStorage) runAll(f callFunc) (err error) {
	type result struct {
		int
		error
	}
	errC := make(chan result)

	for i, s := range s.storages {
		go func() {
			errC <- result{i, f(s)}
		}()
	}

	for i := 0; i < len(s.storages); i++ {
		if result := <-errC; result.error != nil {
			if err == nil {
				err = fmt.Errorf("SyncStorage: error in Storage %d: %v", result.int, result.error)
			} else {
				err = fmt.Errorf("%v\n%29s %d: %v", err, "and error in Storage", result.int, result.error)
			}
		}
	}

	return
}

type event int

const (
	eventCreate event = iota + 1
	eventDelete
	eventModify
)

type fileUpdate struct {
	storage storage.Storage
	event   event
	path    string
}

func (s *SyncStorage) update(fu *fileUpdate) error {
	fmt.Println("Update:", fu)

	return nil
}
