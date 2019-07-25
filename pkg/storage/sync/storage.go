package sync

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/client"
	"github.com/weaveworks/ignite/pkg/storage"
	"github.com/weaveworks/ignite/pkg/storage/watch"
	"github.com/weaveworks/ignite/pkg/storage/watch/update"
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
	storages    []storage.Storage
	eventStream watch.AssociatedEventStream
}

var _ storage.Storage = &SyncStorage{}

// NewSyncStorage constructs a new SyncStorage
func NewSyncStorage(rwStorage storage.Storage, wStorages ...storage.Storage) storage.Storage {
	ss := &SyncStorage{
		Storage:  rwStorage,
		storages: append(wStorages, rwStorage),
	}

	fmt.Println("Storages: %v", ss.storages)

	for _, s := range ss.storages {
		if watchStorage, ok := s.(watch.WatchStorage); ok {
			watchStorage.SetEventStream(ss.getEventStream())
		}
	}

	if ss.eventStream != nil {
		fmt.Println("Started thread!")
		go ss.monitor()
	}

	return ss
}

func (ss *SyncStorage) getEventStream() watch.AssociatedEventStream {
	if ss.eventStream == nil {
		ss.eventStream = make(watch.AssociatedEventStream)
	}

	return ss.eventStream
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

// runAll runs the given function for all Storages in parallel and aggregates all errors
func (ss *SyncStorage) runAll(f func(storage.Storage) error) (err error) {
	type result struct {
		int
		error
	}

	errC := make(chan result)
	for i, s := range ss.storages {
		go func(i int, s storage.Storage) {
			errC <- result{i, f(s)}
		}(i, s) // NOTE: This requires i and s as arguments, otherwise they will be evaluated for one Storage only
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

func (ss *SyncStorage) monitor() {
	// This is the internal client for propagating updates
	c := client.NewClient(ss)

	// TODO: Support detecting changes done when Ignite isn't running
	// This is difficult to do though, as we have don't know which state is the latest
	// For now, only update the state on write when Ignite is running
	for {
		if event, ok := <-ss.eventStream; ok {
			switch event.Event {
			case update.EventModify, update.EventCreate:
				// First load the Object using the Storage given in the event,
				// then set it using the client constructed above
				if obj, err := client.NewClient(event.Storage).Dynamic(event.APIType.GetKind()).Get(event.APIType.GetUID()); err != nil {
					log.Errorf("Failed to get Object with UID %q: %v", event.APIType.GetUID(), err)
				} else if err = c.Dynamic(event.APIType.GetKind()).Set(obj); err != nil {
					log.Errorf("Failed to set Object with UID %q: %v", event.APIType.GetUID(), err)
				}
			case update.EventDelete:
				// For deletion we use the generated "fake" APIType object
				if err := c.Dynamic(event.APIType.GetKind()).Delete(event.APIType.GetUID()); err != nil {
					log.Errorf("Failed to delete Object with UID %q: %v", event.APIType.GetUID(), err)
				}
			}
		} else {
			return
		}
	}
}
