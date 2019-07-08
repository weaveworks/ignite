package filter

import (
	"fmt"

	api "github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/storage/filterer"
)

// The VMFilter filters only VMs, but has special functionality for matching
// If wanting to match all VMs, input a blank string as the prefix
// This ObjectFilter embeds a MetaFilter, which is OK, as ObjectFilter
// interface compatibility is checked before the MetaFilter interface
type VMFilter struct {
	*IDNameFilter
	all bool
}

var _ filterer.ObjectFilter = &VMFilter{}

func NewVMFilter(p string) *VMFilter {
	return NewVMFilterAll(p, true)
}

func NewVMFilterAll(p string, all bool) *VMFilter {
	return &VMFilter{
		IDNameFilter: NewIDNameFilter(p),
		all:          all,
	}
}

func (f *VMFilter) Filter(object meta.Object) (meta.Object, error) {
	// Option to list just running VMs
	if !f.all {
		vm, ok := object.(*api.VM)
		if !ok {
			return nil, fmt.Errorf("invalid Object type for VMFilter: %T", object)
		}

		if vm.Status.State != api.VMStateRunning {
			return nil, nil
		}
	}

	return f.IDNameFilter.FilterMeta(object)
}
