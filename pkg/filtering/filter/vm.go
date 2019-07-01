package filter

import (
	"github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"
	"github.com/weaveworks/ignite/pkg/filtering/filterer"
	"github.com/weaveworks/ignite/pkg/snapshotter"
)

// Compile-time assert to verify interface compatibility
var _ filterer.Filter = &VMFilter{}

type VMFilter struct {
	*IDNameFilter
	all bool
}

func NewVMFilter(p string) *VMFilter {
	return NewVMFilterAll(p, true)
}

func NewVMFilterAll(p string, all bool) *VMFilter {
	return &VMFilter{
		IDNameFilter: NewIDNameFilter(p, v1alpha1.PoolDeviceTypeVM),
		all:          all,
	}
}

func (n *VMFilter) Filter(object v1alpha1.Object) []string {
	// Option to list just running VMs
	if !n.all {
		if snapshotter.ObjectToVM(object).Status.State != v1alpha1.VMStateRunning {
			return nil
		}
	}

	return n.IDNameFilter.Filter(object)
}
