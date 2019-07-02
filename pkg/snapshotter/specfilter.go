package snapshotter

import (
	"fmt"
	"github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"
)

// Compile-time assert to verify interface compatibility
var _ Filter = &objectFilter{}

// The objectFilter is used to get specific objects from the snapshotter
type objectFilter struct {
	reference  *Object
	filterType v1alpha1.PoolDeviceType
}

func newObjectFilter(reference *Object) *objectFilter {
	return &objectFilter{
		reference: reference,
	}
}

func (f *objectFilter) SetType(t v1alpha1.PoolDeviceType) {
	f.filterType = t
}

func (f *objectFilter) Filter(o *Object) (*Object, error) {
	if f.reference.device != nil && f.reference.device != o.device {
		return nil, nil
	}

	if f.reference.object != nil && f.reference.object != o.object {
		return nil, nil
	}

	if f.reference.parent != nil && f.reference.parent != o.parent {
		return nil, nil
	}

	return o, nil
}

func (f *objectFilter) ErrAmbiguous() ErrAmbiguous {
	return fmt.Errorf("ambiguous")
}

func (f *objectFilter) ErrNonexistent() ErrNonexistent {
	return fmt.Errorf("nonexistent")
}
