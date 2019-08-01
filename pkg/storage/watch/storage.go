package watch

import (
	"fmt"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	"github.com/weaveworks/ignite/pkg/apis/ignite/scheme"
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/storage"
	"github.com/weaveworks/ignite/pkg/storage/watch/update"
	"github.com/weaveworks/ignite/pkg/util/sync"
	"github.com/weaveworks/ignite/pkg/util/watcher"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/yaml"
)

// WatchStorage is an extended Storage implementation, which provides a watcher
// for watching changes in the directory managed by the embedded Storage's RawStorage.
// If the RawStorage is a MappedRawStorage instance, it's mappings will automatically
// be updated by the WatchStorage. Update events are sent to the given event stream.
type WatchStorage interface {
	// WatchStorage extends the Storage interface
	storage.Storage
	// GetTrigger returns a hook that can be used to detect a watch event
	SetEventStream(AssociatedEventStream)
}

type AssociatedEventStream chan update.AssociatedUpdate

// NewGenericWatchStorage constructs a new WatchStorage
func NewGenericWatchStorage(s storage.Storage) (WatchStorage, error) {
	ws := &GenericWatchStorage{
		Storage: s,
	}

	var err error
	var files []string
	if ws.watcher, files, err = watcher.NewFileWatcher(s.RawStorage().Dir()); err != nil {
		return nil, err
	}

	if mapped, ok := ws.RawStorage().(storage.MappedRawStorage); ok {
		ws.monitor = sync.RunMonitor(func() {
			ws.monitorFunc(mapped, files) // Offload the file registration to the goroutine
		})
	}

	return ws, nil
}

// GenericWatchStorage implements the WatchStorage interface
type GenericWatchStorage struct {
	storage.Storage
	watcher *watcher.FileWatcher
	events  *AssociatedEventStream
	monitor *sync.Monitor
}

var _ WatchStorage = &GenericWatchStorage{}

// Suspend modify events during Set
func (s *GenericWatchStorage) Set(gvk schema.GroupVersionKind, obj meta.Object) error {
	s.watcher.Suspend(watcher.EventModify)
	return s.Storage.Set(gvk, obj)
}

// Suspend modify events during Patch
func (s *GenericWatchStorage) Patch(gvk schema.GroupVersionKind, uid meta.UID, patch []byte) error {
	s.watcher.Suspend(watcher.EventModify)
	return s.Storage.Patch(gvk, uid, patch)
}

// Suspend delete events during Delete
func (s *GenericWatchStorage) Delete(gvk schema.GroupVersionKind, uid meta.UID) error {
	s.watcher.Suspend(watcher.EventDelete)
	return s.Storage.Delete(gvk, uid)
}

func (s *GenericWatchStorage) SetEventStream(eventStream AssociatedEventStream) {
	s.events = &eventStream
}

func (s *GenericWatchStorage) Close() error {
	s.watcher.Close()
	s.monitor.Wait()
	return nil
}

func (s *GenericWatchStorage) monitorFunc(mapped storage.MappedRawStorage, files []string) {
	defer log.Debug("GenericWatchStorage: monitoring thread stopped")

	// Fill the mappings of the MappedRawStorage before starting to monitor changes
	for _, file := range files {
		if obj, err := resolveAPIType(file); err != nil {
			log.Warnf("Ignoring %q: %v", file, err)
		} else {
			mapped.AddMapping(storage.NewKey(obj.GetKind(), obj.GetUID()), file)
			// Send the event to the events channel
			s.sendEvent(watcher.EventModify, obj)
		}
	}

	for {
		if event, ok := <-s.watcher.GetFileUpdateStream(); ok {
			var obj meta.Object
			var err error

			if event.Event == watcher.EventDelete {
				var key storage.Key
				if key, err = mapped.GetMapping(event.Path); err != nil {
					log.Warnf("Failed to retrieve data for %q: %v", event.Path, err)
					continue
				}

				// This creates a "fake" Object from the key to be used for
				// deletion, as the original has already been removed from disk
				obj = meta.NewAPIType()
				obj.SetUID(key.UID)
				obj.SetGroupVersionKind(schema.GroupVersionKind{
					Group:   api.GroupName,
					Version: runtime.APIVersionInternal,
					Kind:    key.Kind.Title(),
				})
			} else {
				if obj, err = resolveAPIType(event.Path); err != nil {
					log.Warnf("Ignoring %q: %v", event.Path, err)
					continue
				}

				// This is based on the key's existence instead of watcher.EventCreate,
				// as Objects can get updated (via watcher.EventModify) to be conformant
				if _, err = mapped.GetMapping(event.Path); err != nil {
					mapped.AddMapping(storage.NewKey(obj.GetKind(), obj.GetUID()), event.Path)
				}
			}

			// Send the event to the events channel
			s.sendEvent(event.Event, obj)
		} else {
			return
		}
	}
}

func (s *GenericWatchStorage) sendEvent(event watcher.Event, obj meta.Object) {
	if s.events != nil {
		*s.events <- update.AssociatedUpdate{
			Update: update.Update{
				Event:   event,
				APIType: obj,
			},
			Storage: s,
		}
	}
}

func resolveAPIType(path string) (meta.Object, error) {
	obj := meta.NewAPIType()
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// The yaml package supports both YAML and JSON
	if err := yaml.Unmarshal(content, obj); err != nil {
		return nil, err
	}

	gvk := obj.GroupVersionKind()

	// Don't decode API objects unknown to Ignite (e.g. Kubernetes manifests)
	if !scheme.Scheme.Recognizes(gvk) {
		return nil, fmt.Errorf("unknown API version %q and/or kind %q", obj.APIVersion, obj.Kind)
	}

	// Require the UID field to be set
	if len(obj.GetUID()) == 0 {
		return nil, fmt.Errorf(".metadata.uid not set")
	}

	return obj, nil
}
