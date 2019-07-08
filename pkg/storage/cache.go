package storage

import (
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
)

// Cache is an intermediate caching layer, which conforms to Storage
// Typically you back the cache with an actual storage
type Cache interface {
	Storage
	// Flush is used to write the state of the entire cache to storage
	// Warning: this is a very expensive operation
	Flush() error
}

func NewCache(backingStorage Storage) Cache {
	return &cache{
		storage:   backingStorage,
		cache:     newObjectCache(),
		metaCache: newObjectCache(),
	}
}

type cache struct {
	// storage is the backing storage for the cache
	// used to look up non-cached Objects
	storage Storage

	// cache caches the Objects by Kind and UID
	// This guarantees uniqueness when looking up a specific Object
	cache *objectCache

	// metaCache caches the Objects as meta.APITypes by Kind and UID
	metaCache *objectCache
}

var _ Cache = &cache{}

type objectCache struct {
	objects map[meta.Kind]map[meta.UID]meta.Object
}

func newObjectCache() *objectCache {
	return &objectCache{
		make(map[meta.Kind]map[meta.UID]meta.Object),
	}
}

func (c *cache) Get(obj meta.Object) error {
	// Try to load from cache first
	if obj = c.cache.load(obj); obj != nil {
		return nil
	}

	// If not cached, request loading from storage
	return c.storage.Get(obj)
}

func (c *cache) Set(obj meta.Object) error {
	// Store the changed Object in the cache
	c.cache.store(obj)

	// TODO: For now the cache always flushes, we might add a Cache.Flush() later
	return c.storage.Set(obj)
}

func (c *cache) GetByID(kind meta.Kind, uid meta.UID) (meta.Object, error) {
	// If the requested Object resides in the cache, return it
	if obj := c.cache.loadByID(kind, uid); obj != nil {
		return obj, nil
	}

	obj, err := c.storage.GetByID(kind, uid)

	// If no errors occurred while loading, store the Object in the cache
	if err != nil {
		c.cache.store(obj)
	}

	return obj, err
}

func (c *cache) Delete(kind meta.Kind, uid meta.UID) error {
	// Delete the given Object from the cache and storage
	c.cache.delete(kind, uid)
	return c.storage.Delete(kind, uid)
}

func (c *cache) List(kind meta.Kind) ([]meta.Object, error) {
	var objs []meta.Object
	var storageCount uint64
	var err error

	if storageCount, err = c.storage.Count(kind); err != nil {
		return nil, err
	}

	if c.cache.count(kind) != storageCount {
		// If the cache doesn't track all of the Objects, request them from the storage
		if objs, err = c.storage.List(kind); err != nil {
			// If no errors occurred, store the Objects in the cache
			c.cache.storeAll(objs)
		}
	} else {
		// If the cache tracks everything, return the cache's contents
		objs = c.cache.list(kind)
	}

	return objs, err
}

// TODO: The metaCache falls out of sync on any updates, the cache should support saving
// headers and objects coupled together loading the object to replace the header when needed
func (c *cache) ListMeta(kind meta.Kind) ([]meta.Object, error) {
	// TODO: Support extracting meta.APIType Objects from fully loaded objects to back the meta cache
	var objs []meta.Object
	var storageCount uint64
	var err error

	if storageCount, err = c.storage.Count(kind); err != nil {
		return nil, err
	}

	if c.metaCache.count(kind) != storageCount {
		// If the cache doesn't track all of the Objects, request them from the storage
		if objs, err = c.storage.ListMeta(kind); err != nil {
			// If no errors occurred, store the Objects in the cache
			c.metaCache.storeAll(objs)
		}
	} else {
		// If the cache tracks everything, return the cache's contents
		objs = c.metaCache.list(kind)
	}

	return objs, err
}

func (c *cache) Count(kind meta.Kind) (uint64, error) {
	// The cache is transparent about how many items it has cached
	return c.storage.Count(kind)
}

func (c *cache) Flush() error {
	// Load the entire cache
	for _, obj := range c.cache.loadAll() {
		// Request the storage to save each Object
		if err := c.storage.Set(obj); err != nil {
			return err
		}
	}

	return nil
}

func (c *objectCache) load(obj meta.Object) meta.Object {
	return c.loadByID(obj.GetKind(), obj.GetUID())
}

func (c *objectCache) loadByID(kind meta.Kind, uid meta.UID) meta.Object {
	if uids, ok := c.objects[kind]; ok {
		if obj, ok := uids[uid]; ok {
			return obj
		}
	}

	return nil
}

func (c *objectCache) loadAll() []meta.Object {
	var size uint64

	for kind := range c.objects {
		size += c.count(kind)
	}

	all := make([]meta.Object, 0, size)

	for _, uids := range c.objects {
		for _, obj := range uids {
			all = append(all, obj)
		}
	}

	return all
}

func (c *objectCache) store(obj meta.Object) {
	kind := obj.GetKind()

	if _, ok := c.objects[kind]; !ok {
		c.objects[kind] = make(map[meta.UID]meta.Object)
	}

	c.objects[kind][obj.GetUID()] = obj
}

func (c *objectCache) storeAll(objs []meta.Object) {
	for _, obj := range objs {
		c.store(obj)
	}
}

func (c *objectCache) delete(kind meta.Kind, uid meta.UID) {
	if uids, ok := c.objects[kind]; ok {
		delete(uids, uid)
	}
}

func (c *objectCache) count(kind meta.Kind) uint64 {
	return uint64(len(c.objects[kind]))
}

func (c *objectCache) list(kind meta.Kind) []meta.Object {
	uids := c.objects[kind]
	list := make([]meta.Object, 0, len(uids))

	for _, obj := range uids {
		list = append(list, obj)
	}

	return list
}
