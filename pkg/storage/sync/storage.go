package sync

import (
	"fmt"
	"github.com/weaveworks/ignite/pkg/storage/watch"

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
	watched  map[watch.EventStream]storage.Storage
}

var _ storage.Storage = &SyncStorage{}

// NewSyncStorage constructs a new SyncStorage
func NewSyncStorage(rwStorage storage.Storage, wStorages ...storage.Storage) storage.Storage {
	ss := &SyncStorage{
		Storage:  rwStorage,
		storages: append(wStorages, rwStorage),
	}

	for _, s := range ss.storages {
		if watchStorage, ok := s.(watch.WatchStorage); ok {
			ss.watched[watchStorage.EventStream()] = watchStorage
		}
	}

	if len(ss.watched) > 0 {

	}
}

func (ss *SyncStorage) Add(wStorages ...storage.Storage) storage.Storage {
	ss.storages = append(ss.storages, wStorages...)
	return ss
}

func (ss *SyncStorage) AddWatched(wStorages ...storage.Storage) (storage.Storage, error) {
	for _, st := range wStorages {
		w, err := newWatcher(st, ss.update)
		if err != nil {
			return nil, err
		}

		ss.watchers = append(ss.watchers, w)
	}

	return ss.Add(wStorages...), nil
}

// Set is propagated to all Storages
func (ss *SyncStorage) Set(gvk schema.GroupVersionKind, obj meta.Object) error {
	return ss.runAll(func(s storage.Storage) error {
		return s.Set(gvk, obj)
	})
}

// Patch is propagated to all Storages
func (ss *SyncStorage) Patch(gvk schema.GroupVersionKind, uid meta.UID, patch []byte) error {
	return ss.runAll(func(s storage.Storage) error {
		return s.Patch(gvk, uid, patch)
	})
}

// Delete is propagated to all Storages
func (ss *SyncStorage) Delete(gvk schema.GroupVersionKind, uid meta.UID) error {
	return ss.runAll(func(s storage.Storage) error {
		return s.Delete(gvk, uid)
	})
}

type callFunc func(storage.Storage) error

// runAll runs the given callFunc for all Storages in parallel and aggregates all errors
func (ss *SyncStorage) runAll(f callFunc) (err error) {
	type result struct {
		int
		error
	}
	errC := make(chan result)

	for i, s := range ss.storages {
		go func() {
			errC <- result{i, f(s)}
		}()
	}

	for i := 0; i < len(ss.storages); i++ {
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

func (ss *SyncStorage) update(fu *fileUpdate) error {
	fmt.Println("Update:", fu)

	return nil
}
