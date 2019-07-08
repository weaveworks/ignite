package storage

import (
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"log"
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
	c := &cache{
		storage:     backingStorage,
		objectCache: newObjectCache(backingStorage.GetByID),
	}

	return c
}

type cache struct {
	// storage is the backing storage for the cache
	// used to look up non-cached Objects
	storage Storage

	// cache caches the Objects by Kind and UID
	// This guarantees uniqueness when looking up a specific Object
	objectCache *objectCache
}

var _ Cache = &cache{}

type cacheObject struct {
	object  meta.Object
	apiType bool
}

type loadFunc func(meta.Kind, meta.UID) (meta.Object, error)

type objectCache struct {
	objects  map[meta.Kind]map[meta.UID]cacheObject
	loadFunc loadFunc
}

func newObjectCache(l loadFunc) *objectCache {
	return &objectCache{
		objects:  make(map[meta.Kind]map[meta.UID]cacheObject),
		loadFunc: l,
	}
}

func (c *cache) Get(obj meta.Object) error {
	var err error

	// Try to load from cache first
	if obj, err = c.objectCache.load(obj); err != nil {
		return err
	} else if obj != nil {
		return nil
	}

	// If not cached, request loading from storage
	return c.storage.Get(obj)
}

func (c *cache) Set(obj meta.Object) error {
	// Store the changed Object in the cache
	c.objectCache.store(obj)

	// TODO: For now the cache always flushes, we might add a Cache.Flush() later
	return c.storage.Set(obj)
}

func (c *cache) GetByID(kind meta.Kind, uid meta.UID) (meta.Object, error) {
	// If the requested Object resides in the cache, return it
	if obj, err := c.objectCache.loadByID(kind, uid); err != nil {
		return nil, err
	} else if obj != nil {
		return obj, nil
	}

	obj, err := c.storage.GetByID(kind, uid)

	// If no errors occurred while loading, store the Object in the cache
	if err != nil {
		c.objectCache.store(obj)
	}

	return obj, err
}

func (c *cache) Delete(kind meta.Kind, uid meta.UID) error {
	// Delete the given Object from the cache and storage
	c.objectCache.delete(kind, uid)
	return c.storage.Delete(kind, uid)
}

type listFunc func(meta.Kind) ([]meta.Object, error)
type cacheStoreFunc func([]meta.Object)

// list is a common handler for List and ListMeta
func (c *cache) list(kind meta.Kind, slf, clf listFunc, csf cacheStoreFunc) ([]meta.Object, error) {
	var objs []meta.Object
	var storageCount uint64
	var err error

	if storageCount, err = c.storage.Count(kind); err != nil {
		return nil, err
	}

	if c.objectCache.count(kind) != storageCount {
		// If the cache doesn't track all of the Objects, request them from the storage
		if objs, err = slf(kind); err != nil {
			// If no errors occurred, store the Objects in the cache
			csf(objs)
		}
	} else {
		// If the cache tracks everything, return the cache's contents
		objs, err = clf(kind)
	}

	return objs, err
}

func (c *cache) List(kind meta.Kind) ([]meta.Object, error) {
	return c.list(kind, c.storage.List, c.objectCache.list, c.objectCache.storeAll)
}

func (c *cache) ListMeta(kind meta.Kind) ([]meta.Object, error) {
	return c.list(kind, c.storage.ListMeta, c.objectCache.listMeta, c.objectCache.storeAllMeta)
}

func (c *cache) Count(kind meta.Kind) (uint64, error) {
	// The cache is transparent about how many items it has cached
	return c.storage.Count(kind)
}

func (c *cache) Flush() error {
	// Load the entire cache
	allObjects, err := c.objectCache.loadAll()
	if err != nil {
		return err
	}

	for _, obj := range allObjects {
		// Request the storage to save each Object
		if err := c.storage.Set(obj); err != nil {
			return err
		}
	}

	return nil
}

// loadFull checks if the Object is an APIType, and loads the full Object in that case
func (c *objectCache) loadFull(obj cacheObject) (meta.Object, error) {
	if !obj.apiType {
		log.Printf("cache: full %s object cached\n", obj.object.GetKind())
		return obj.object, nil
	}

	log.Printf("cache: loading full %s object\n", obj.object.GetKind())
	result, err := c.loadFunc(obj.object.GetKind(), obj.object.GetUID())
	if err != nil {
		return nil, err
	}

	c.store(result)
	return result, nil
}

func (c *objectCache) load(obj meta.Object) (meta.Object, error) {
	return c.loadByID(obj.GetKind(), obj.GetUID())
}

func (c *objectCache) loadByID(kind meta.Kind, uid meta.UID) (meta.Object, error) {
	if uids, ok := c.objects[kind]; ok {
		if obj, ok := uids[uid]; ok {
			return c.loadFull(obj)
		}
	}

	return nil, nil
}

func (c *objectCache) loadAll() ([]meta.Object, error) {
	var size uint64

	for kind := range c.objects {
		size += c.count(kind)
	}

	all := make([]meta.Object, 0, size)

	for kind := range c.objects {
		if objects, err := c.list(kind); err == nil {
			all = append(all, objects...)
		} else {
			return nil, err
		}
	}

	return all, nil
}

func (c *objectCache) storeCO(obj cacheObject) {
	kind := obj.object.GetKind()

	if _, ok := c.objects[kind]; !ok {
		c.objects[kind] = make(map[meta.UID]cacheObject)
	}

	c.objects[kind][obj.object.GetUID()] = obj
}

func (c *objectCache) store(obj meta.Object) {
	log.Printf("cache: storing %s object\n", obj.GetKind())
	c.storeCO(cacheObject{object: obj})
}

func (c *objectCache) storeAll(objs []meta.Object) {
	for _, obj := range objs {
		c.store(obj)
	}
}

func (c *objectCache) storeMeta(obj meta.Object) {
	log.Printf("cache: storing %s meta object\n", obj.GetKind())
	c.storeCO(cacheObject{object: obj, apiType: true})
}

func (c *objectCache) storeAllMeta(objs []meta.Object) {
	for _, obj := range objs {
		if uids, ok := c.objects[obj.GetKind()]; ok {
			if _, ok := uids[obj.GetUID()]; ok {
				continue
			}
		}

		c.storeMeta(obj)
	}
}

func (c *objectCache) delete(kind meta.Kind, uid meta.UID) {
	if uids, ok := c.objects[kind]; ok {
		delete(uids, uid)
	}
}

func (c *objectCache) count(kind meta.Kind) uint64 {
	count := uint64(len(c.objects[kind]))
	log.Printf("cache: counted %d %s objects\n", count, kind)
	return count
}

func (c *objectCache) list(kind meta.Kind) ([]meta.Object, error) {
	uids := c.objects[kind]
	list := make([]meta.Object, 0, len(uids))

	for _, obj := range uids {
		if result, err := c.loadFull(obj); err == nil {
			list = append(list, result)
		} else {
			return nil, err
		}
	}

	return list, nil
}

func (c *objectCache) listMeta(kind meta.Kind) ([]meta.Object, error) {
	uids := c.objects[kind]
	list := make([]meta.Object, 0, len(uids))

	for _, obj := range uids {
		apiType := obj.object

		if !obj.apiType {
			apiType = meta.APITypeFrom(obj.object)
		}

		list = append(list, apiType)
	}

	return list, nil
}
