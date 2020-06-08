package filter

import (
	"fmt"

	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	"github.com/weaveworks/libgitops/pkg/filter"
	"github.com/weaveworks/libgitops/pkg/runtime"
	"github.com/weaveworks/libgitops/pkg/storage/filterer"
)

// The VMFilter filters only VMs, but has special functionality for matching
// If wanting to match all VMs, input a blank string as the prefix
// This ObjectFilter embeds a MetaFilter, which is OK, as ObjectFilter
// interface compatibility is checked before the MetaFilter interface
type VMFilter struct {
	*filter.IDNameFilter
	all bool
}

var _ filterer.ObjectFilter = &VMFilter{}

func NewVMFilter(p string) *VMFilter {
	return NewVMFilterAll(p, true)
}

func NewVMFilterAll(p string, all bool) *VMFilter {
	return &VMFilter{
		IDNameFilter: filter.NewIDNameFilter(p),
		all:          all,
	}
}

func (f *VMFilter) Filter(object runtime.Object) (filterer.Match, error) {
	// Option to list just running VMs
	if !f.all {
		vm, ok := object.(*api.VM)
		if !ok {
			return nil, fmt.Errorf("invalid Object type for VMFilter: %T", object)
		}

		if !vm.Running() {
			return nil, nil
		}
	}

	return f.IDNameFilter.FilterMeta(object)
}
