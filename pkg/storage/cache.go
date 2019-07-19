package storage

import (
	log "github.com/sirupsen/logrus"
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"

	"k8s.io/apimachinery/pkg/runtime/schema"
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
		objectCache: newObjectCache(backingStorage.Get),
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

type loadFunc func(schema.GroupVersionKind, meta.UID) (meta.Object, error)

type objectCache struct {
	objects  map[schema.GroupVersionKind]map[meta.UID]cacheObject
	loadFunc loadFunc
}

func newObjectCache(l loadFunc) *objectCache {
	return &objectCache{
		objects:  make(map[schema.GroupVersionKind]map[meta.UID]cacheObject),
		loadFunc: l,
	}
}

func (c *cache) New(gvk schema.GroupVersionKind) (meta.Object, error) {
	// passthrough to storage
	return c.storage.New(gvk)
}

func (c *cache) Get(gvk schema.GroupVersionKind, uid meta.UID) (meta.Object, error) {
	// If the requested Object resides in the cache, return it
	if obj, err := c.objectCache.loadByID(gvk, uid); err != nil {
		return nil, err
	} else if obj != nil {
		return obj, nil
	}

	obj, err := c.storage.Get(gvk, uid)

	// If no errors occurred while loading, store the Object in the cache
	if err != nil {
		c.objectCache.store(obj)
	}

	return obj, err
}

func (c *cache) Set(gvk schema.GroupVersionKind, obj meta.Object) error {
	// Store the changed Object in the cache
	c.objectCache.store(obj)

	// TODO: For now the cache always flushes, we might add a Cache.Flush() later
	return c.storage.Set(gvk, obj)
}

// Patch performs a strategic merge patch on the object with the given UID, using the byte-encoded patch given
func (s *cache) Patch(gvk schema.GroupVersionKind, uid meta.UID, patch []byte) error {
	// TODO: Should we do something here to cache the change? I don't think so, but...
	// Just passthrough to the storage here
	return s.storage.Patch(gvk, uid, patch)
}

func (c *cache) Delete(gvk schema.GroupVersionKind, uid meta.UID) error {
	// Delete the given Object from the cache and storage
	c.objectCache.delete(gvk, uid)
	return c.storage.Delete(gvk, uid)
}

type listFunc func(gvk schema.GroupVersionKind) ([]meta.Object, error)
type cacheStoreFunc func([]meta.Object)

// list is a common handler for List and ListMeta
func (c *cache) list(gvk schema.GroupVersionKind, slf, clf listFunc, csf cacheStoreFunc) ([]meta.Object, error) {
	var objs []meta.Object
	var storageCount uint64
	var err error

	if storageCount, err = c.storage.Count(gvk); err != nil {
		return nil, err
	}

	if c.objectCache.count(gvk) != storageCount {
		// If the cache doesn't track all of the Objects, request them from the storage
		if objs, err = slf(gvk); err != nil {
			// If no errors occurred, store the Objects in the cache
			csf(objs)
		}
	} else {
		// If the cache tracks everything, return the cache's contents
		objs, err = clf(gvk)
	}

	return objs, err
}

func (c *cache) List(gvk schema.GroupVersionKind) ([]meta.Object, error) {
	return c.list(gvk, c.storage.List, c.objectCache.list, c.objectCache.storeAll)
}

func (c *cache) ListMeta(gvk schema.GroupVersionKind) ([]meta.Object, error) {
	return c.list(gvk, c.storage.ListMeta, c.objectCache.listMeta, c.objectCache.storeAllMeta)
}

func (c *cache) Count(gvk schema.GroupVersionKind) (uint64, error) {
	// The cache is transparent about how many items it has cached
	return c.storage.Count(gvk)
}

func (c *cache) Flush() error {
	// Load the entire cache
	allObjects, err := c.objectCache.loadAll()
	if err != nil {
		return err
	}

	for _, obj := range allObjects {
		// Request the storage to save each Object
		gvk := obj.GroupVersionKind()
		if err := c.storage.Set(gvk, obj); err != nil {
			return err
		}
	}

	return nil
}

// loadFull checks if the Object is an APIType, and loads the full Object in that case
func (c *objectCache) loadFull(obj cacheObject) (meta.Object, error) {
	if !obj.apiType {
		log.Debugf("cache: full %s object cached\n", obj.object.GetKind())
		return obj.object, nil
	}

	log.Debugf("cache: loading full %s object\n", obj.object.GetKind())
	result, err := c.loadFunc(obj.object.GroupVersionKind(), obj.object.GetUID())
	if err != nil {
		return nil, err
	}

	c.store(result)
	return result, nil
}

func (c *objectCache) load(obj meta.Object) (meta.Object, error) {
	return c.loadByID(obj.GroupVersionKind(), obj.GetUID())
}

func (c *objectCache) loadByID(gvk schema.GroupVersionKind, uid meta.UID) (meta.Object, error) {
	if uids, ok := c.objects[gvk]; ok {
		if obj, ok := uids[uid]; ok {
			return c.loadFull(obj)
		}
	}

	return nil, nil
}

func (c *objectCache) loadAll() ([]meta.Object, error) {
	var size uint64

	for gvk := range c.objects {
		size += c.count(gvk)
	}

	all := make([]meta.Object, 0, size)

	for gvk := range c.objects {
		if objects, err := c.list(gvk); err == nil {
			all = append(all, objects...)
		} else {
			return nil, err
		}
	}

	return all, nil
}

func (c *objectCache) storeCO(obj cacheObject) {
	gvk := obj.object.GroupVersionKind()

	if _, ok := c.objects[gvk]; !ok {
		c.objects[gvk] = make(map[meta.UID]cacheObject)
	}

	c.objects[gvk][obj.object.GetUID()] = obj
}

func (c *objectCache) store(obj meta.Object) {
	log.Debugf("cache: storing %s object\n", obj.GetKind())
	c.storeCO(cacheObject{object: obj})
}

func (c *objectCache) storeAll(objs []meta.Object) {
	for _, obj := range objs {
		c.store(obj)
	}
}

func (c *objectCache) storeMeta(obj meta.Object) {
	log.Debugf("cache: storing %s meta object\n", obj.GetKind())
	c.storeCO(cacheObject{object: obj, apiType: true})
}

func (c *objectCache) storeAllMeta(objs []meta.Object) {
	for _, obj := range objs {
		if uids, ok := c.objects[obj.GroupVersionKind()]; ok {
			if _, ok := uids[obj.GetUID()]; ok {
				continue
			}
		}

		c.storeMeta(obj)
	}
}

func (c *objectCache) delete(gvk schema.GroupVersionKind, uid meta.UID) {
	if uids, ok := c.objects[gvk]; ok {
		delete(uids, uid)
	}
}

func (c *objectCache) count(gvk schema.GroupVersionKind) uint64 {
	count := uint64(len(c.objects[gvk]))
	log.Debugf("cache: counted %d %s objects\n", count, gvk.Kind)
	return count
}

func (c *objectCache) list(gvk schema.GroupVersionKind) ([]meta.Object, error) {
	uids := c.objects[gvk]
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

func (c *objectCache) listMeta(gvk schema.GroupVersionKind) ([]meta.Object, error) {
	uids := c.objects[gvk]
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
