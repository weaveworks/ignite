package filterer

import (
	"fmt"

	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/storage"
)

type Filterer struct {
	storage storage.Storage
}

func NewFilterer(storage storage.Storage) *Filterer {
	return &Filterer{
		storage: storage,
	}
}

type filterFunc func(meta.Object) (meta.Object, error)

// Find a single meta.Object of the given kind using the given filter
func (f *Filterer) Find(kind meta.Kind, filter BaseFilter) (meta.Object, error) {
	var result meta.Object

	// Fetch the sources, correct filtering method and if we're dealing with meta.APIType objects
	sources, filterFunc, metaObjects, err := f.parseFilter(kind, filter)
	if err != nil {
		return nil, err
	}

	// Perform the filtering
	for _, object := range sources {
		if match, err := filterFunc(object); err != nil { // The filter returns meta.Object if it matches, otherwise nil
			return nil, err
		} else if match != nil {
			if result != nil {
				return nil, filter.AmbiguousError()
			} else {
				result = match
			}
		}
	}

	if result == nil {
		return nil, filter.NonexistentError()
	}

	// If we're filtering meta.APIType objects, load the full Object to be returned
	if metaObjects {
		return f.storage.GetByID(result.GetKind(), result.GetUID())
	}

	return result, nil
}

// Find all meta.Objects of the given kind using the given filter
func (f *Filterer) FindAll(kind meta.Kind, filter BaseFilter) ([]meta.Object, error) {
	var results []meta.Object

	// Fetch the sources, correct filtering method and if we're dealing with meta.APIType objects
	sources, filterFunc, metaObjects, err := f.parseFilter(kind, filter)
	if err != nil {
		return nil, err
	}

	// Perform the filtering
	for _, object := range sources {
		if match, err := filterFunc(object); err != nil { // The filter returns meta.Object if it matches, otherwise nil
			return nil, err
		} else if match != nil {
			results = append(results, match)
		}
	}

	// If we're filtering meta.APIType objects, load the full Objects to be returned
	if metaObjects {
		objects := make([]meta.Object, len(results))
		for i, result := range results {
			if objects[i], err = f.storage.GetByID(result.GetKind(), result.GetUID()); err != nil {
				return nil, err
			}
		}

		return objects, nil
	}

	return results, nil
}

func (f *Filterer) parseFilter(kind meta.Kind, filter BaseFilter) (sources []meta.Object, filterFunc filterFunc, metaObjects bool, err error) {
	// Parse ObjectFilters before MetaFilters, so ObjectFilters can embed MetaFilters
	if objectFilter, ok := filter.(ObjectFilter); ok {
		filterFunc = objectFilter.Filter
		sources, err = f.storage.List(kind)
	} else if metaFilter, ok := filter.(MetaFilter); ok {
		filterFunc = metaFilter.FilterMeta
		sources, err = f.storage.ListMeta(kind)
		metaObjects = true
	} else {
		err = fmt.Errorf("invalid filter type: %T", filter)
	}
	// Make sure the desired kind propagates down to the filter
	filter.SetKind(kind)

	return
}
