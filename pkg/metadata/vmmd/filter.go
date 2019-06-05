package vmmd

import (
	"github.com/luxas/ignite/pkg/metadata"
)

// Compile-time assert to verify interface compatibility
var _ metadata.Filter = &VMFilter{}

type VMFilter struct {
	*metadata.IDNameFilter
	all bool
}

func NewVMFilter(p string) *VMFilter {
	return NewVMFilterAll(p, true)
}

func NewVMFilterAll(p string, all bool) *VMFilter {
	return &VMFilter{
		IDNameFilter: metadata.NewIDNameFilter(p, metadata.VM),
		all:          all,
	}
}

func (n *VMFilter) Filter(any metadata.AnyMetadata) []string {
	// Option to list just running VMs
	if !n.all {
		if ToVMMetadata(any).VMOD().State != Running {
			return nil
		}
	}

	return n.IDNameFilter.Filter(any)
}
