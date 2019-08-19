package cache

import (
	log "github.com/sirupsen/logrus"
	"github.com/weaveworks/gitops-toolkit/pkg/runtime"
	"github.com/weaveworks/gitops-toolkit/pkg/storage"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type index struct {
	storage storage.Storage
	objects map[schema.GroupVersionKind]map[runtime.UID]*cacheObject
}

func newIndex(storage storage.Storage) *index {
	return &index{
		storage: storage,
		objects: make(map[schema.GroupVersionKind]map[runtime.UID]*cacheObject),
	}
}

func (i *index) loadByID(gvk schema.GroupVersionKind, uid runtime.UID) (runtime.Object, error) {
	if uids, ok := i.objects[gvk]; ok {
		if obj, ok := uids[uid]; ok {
			log.Tracef("index: cache hit for %s with UID %q", gvk.Kind, uid)
			return obj.loadFull()
		}
	}

	log.Tracef("index: cache miss for %s with UID %q", gvk.Kind, uid)
	return nil, nil
}

func (i *index) load(obj runtime.Object) (runtime.Object, error) {
	return i.loadByID(obj.GroupVersionKind(), obj.GetUID())
}

func (i *index) loadAll() ([]runtime.Object, error) {
	var size uint64

	for gvk := range i.objects {
		size += i.count(gvk)
	}

	all := make([]runtime.Object, 0, size)

	for gvk := range i.objects {
		if objects, err := i.list(gvk); err == nil {
			all = append(all, objects...)
		} else {
			return nil, err
		}
	}

	return all, nil
}

func store(i *index, obj runtime.Object, apiType bool) error {
	// If store is called for an invalid Object lacking an UID,
	// panic and print the stack trace. This should never happen.
	if obj.GetUID() == "" {
		panic("Attempt to cache invalid Object: missing UID")
	}

	co, err := newCacheObject(i.storage, obj, apiType)
	if err != nil {
		return err
	}

	gvk := co.object.GroupVersionKind()

	if _, ok := i.objects[gvk]; !ok {
		i.objects[gvk] = make(map[runtime.UID]*cacheObject)
	}

	log.Tracef("index: storing %s object with UID %q, meta: %t", gvk.Kind, obj.GetName(), apiType)
	i.objects[gvk][co.object.GetUID()] = co

	return nil
}

func (i *index) store(obj runtime.Object) error {
	return store(i, obj, false)
}

func (i *index) storeAll(objs []runtime.Object) (err error) {
	for _, obj := range objs {
		if err = i.store(obj); err != nil {
			break
		}
	}

	return
}

func (i *index) storeMeta(obj runtime.Object) error {
	return store(i, obj, true)
}

func (i *index) storeAllMeta(objs []runtime.Object) (err error) {
	for _, obj := range objs {
		if uids, ok := i.objects[obj.GroupVersionKind()]; ok {
			if _, ok := uids[obj.GetUID()]; ok {
				continue
			}
		}

		if err = i.storeMeta(obj); err != nil {
			break
		}
	}

	return
}

func (i *index) delete(gvk schema.GroupVersionKind, uid runtime.UID) {
	if uids, ok := i.objects[gvk]; ok {
		delete(uids, uid)
	}
}

func (i *index) count(gvk schema.GroupVersionKind) (count uint64) {
	count = uint64(len(i.objects[gvk]))
	log.Tracef("index: counted %d %s object(s)", count, gvk.Kind)
	return
}

func list(i *index, gvk schema.GroupVersionKind, apiTypes bool) ([]runtime.Object, error) {
	uids := i.objects[gvk]
	list := make([]runtime.Object, 0, len(uids))

	log.Tracef("index: listing %s objects, meta: %t", gvk, apiTypes)
	for _, obj := range uids {
		loadFunc := obj.loadFull
		if apiTypes {
			loadFunc = obj.loadAPI
		}

		if result, err := loadFunc(); err != nil {
			return nil, err
		} else {
			list = append(list, result)
		}
	}

	return list, nil
}

func (i *index) list(gvk schema.GroupVersionKind) ([]runtime.Object, error) {
	return list(i, gvk, false)
}

func (i *index) listMeta(gvk schema.GroupVersionKind) ([]runtime.Object, error) {
	return list(i, gvk, true)
}
