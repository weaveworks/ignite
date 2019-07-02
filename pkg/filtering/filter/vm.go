package filter

import (
	"fmt"
	"github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"
	"github.com/weaveworks/ignite/pkg/snapshotter"
)

// Compile-time assert to verify interface compatibility
var _ snapshotter.Filter = &VMFilter{}

// The VMFilter filters only VMs, but has special functionality for matching
type VMFilter struct {
	*IDNameFilter
	all bool
}

func NewVMFilter(p string) *VMFilter {
	return NewVMFilterAll(p, true)
}

func NewVMFilterAll(p string, all bool) *VMFilter {
	return &VMFilter{
		IDNameFilter: NewIDNameFilter(p),
		all:          all,
	}
}

func (f *VMFilter) Filter(object *snapshotter.Object) (*snapshotter.Object, error) {
	// Option to list just running VMs
	if !f.all {
		mo, err := object.GetMetaObject()
		if err != nil {
			return nil, err
		}

		vm, ok := mo.(*v1alpha1.VM)
		if !ok {
			return nil, fmt.Errorf("invalid object type for VMFilter: %T", mo)
		}

		if vm.Status.State != v1alpha1.VMStateRunning {
			return nil, nil
		}
	}

	return f.IDNameFilter.Filter(object)
}
