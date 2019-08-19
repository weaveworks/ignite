package cache

import (
	log "github.com/sirupsen/logrus"
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/serializer"
	"github.com/weaveworks/ignite/pkg/storage"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// Cache is an intermediate caching layer, which conforms to Storage
// Typically you back the cache with an actual storage
type Cache interface {
	storage.Storage
	// Flush is used to write the state of the entire cache to storage
	// Warning: this is a very expensive operation
	Flush() error
}

type cache struct {
	// storage is the backing Storage for the cache
	// used to look up non-cached Objects
	storage storage.Storage

	// index caches the Objects by GroupVersionKind and UID
	// This guarantees uniqueness when looking up a specific Object
	index *index
}

var _ Cache = &cache{}

func NewCache(backingStorage storage.Storage) Cache {
	c := &cache{
		storage: backingStorage,
		index:   newIndex(backingStorage),
	}

	return c
}

func (s *cache) Serializer() serializer.Serializer {
	return s.storage.Serializer()
}

func (c *cache) New(gvk schema.GroupVersionKind) (meta.Object, error) {
	// Request the storage to create the Object. The
	// newly generated Object has not got an UID which
	// is required for indexing, so just return it
	// without storing it into the cache
	return c.storage.New(gvk)
}

func (c *cache) Get(gvk schema.GroupVersionKind, uid meta.UID) (obj meta.Object, err error) {
	log.Tracef("cache: Get %s with UID %q", gvk.Kind, uid)

	// If the requested Object resides in the cache, return it
	if obj, err = c.index.loadByID(gvk, uid); err != nil || obj != nil {
		return
	}

	// Request the Object from the storage
	obj, err = c.storage.Get(gvk, uid)

	// If no errors occurred, cache it
	if err == nil {
		err = c.index.store(obj)
	}

	return
}

func (c *cache) GetMeta(gvk schema.GroupVersionKind, uid meta.UID) (obj meta.Object, err error) {
	log.Tracef("cache: GetMeta %s with UID %q", gvk.Kind, uid)

	obj, err = c.storage.GetMeta(gvk, uid)

	// If no errors occurred while loading, store the Object in the cache
	if err == nil {
		err = c.index.storeMeta(obj)
	}

	return
}

func (c *cache) Set(gvk schema.GroupVersionKind, obj meta.Object) error {
	log.Tracef("cache: Set %s with UID %q", gvk.Kind, obj.GetUID())

	// Store the changed Object in the cache
	if err := c.index.store(obj); err != nil {
		return err
	}

	// TODO: For now the cache always flushes, we might add automatic flushing later
	return c.storage.Set(gvk, obj)
}

func (c *cache) Patch(gvk schema.GroupVersionKind, uid meta.UID, patch []byte) error {
	// TODO: For now patches are always flushed, the cache will load the updated Object on-demand on access
	return c.storage.Patch(gvk, uid, patch)
}

func (c *cache) Delete(gvk schema.GroupVersionKind, uid meta.UID) error {
	log.Tracef("cache: Delete %s with UID %q", gvk.Kind, uid)

	// Delete the given Object from the cache and storage
	c.index.delete(gvk, uid)
	return c.storage.Delete(gvk, uid)
}

type listFunc func(gvk schema.GroupVersionKind) ([]meta.Object, error)
type cacheStoreFunc func([]meta.Object) error

// list is a common handler for List and ListMeta
func (c *cache) list(gvk schema.GroupVersionKind, slf, clf listFunc, csf cacheStoreFunc) (objs []meta.Object, err error) {
	var storageCount uint64
	if storageCount, err = c.storage.Count(gvk); err != nil {
		return
	}

	if c.index.count(gvk) != storageCount {
		log.Tracef("cache: miss when listing: %s", gvk)
		// If the cache doesn't track all of the Objects, request them from the storage
		if objs, err = slf(gvk); err != nil {
			// If no errors occurred, store the Objects in the cache
			err = csf(objs)
		}
	} else {
		log.Tracef("cache: hit when listing: %s", gvk)
		// If the cache tracks everything, return the cache's contents
		objs, err = clf(gvk)
	}

	return
}

func (c *cache) List(gvk schema.GroupVersionKind) ([]meta.Object, error) {
	return c.list(gvk, c.storage.List, c.index.list, c.index.storeAll)
}

func (c *cache) ListMeta(gvk schema.GroupVersionKind) ([]meta.Object, error) {
	return c.list(gvk, c.storage.ListMeta, c.index.listMeta, c.index.storeAllMeta)
}

func (c *cache) Count(gvk schema.GroupVersionKind) (uint64, error) {
	// The cache is transparent about how many items it has cached
	return c.storage.Count(gvk)
}

func (c *cache) Checksum(gvk schema.GroupVersionKind, uid meta.UID) (string, error) {
	// The cache is transparent about the checksums
	return c.storage.Checksum(gvk, uid)
}

func (c *cache) RawStorage() storage.RawStorage {
	return c.storage.RawStorage()
}

func (c *cache) Close() error {
	return c.storage.Close()
}

func (c *cache) Flush() error {
	// Load the entire cache
	allObjects, err := c.index.loadAll()
	if err != nil {
		return err
	}

	for _, obj := range allObjects {
		// Request the storage to save each Object
		if err := c.storage.Set(obj.GroupVersionKind(), obj); err != nil {
			return err
		}
	}

	return nil
}
