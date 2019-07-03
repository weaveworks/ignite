package vmmd

import (
	api "github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"
	"github.com/weaveworks/ignite/pkg/metadata"
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
		IDNameFilter: metadata.NewIDNameFilter(p, api.PoolDeviceTypeVM),
		all:          all,
	}
}

func (n *VMFilter) Filter(md metadata.Metadata) []string {
	// Option to list just running VMs
	if !n.all {
		if !ToVMMetadata(md).Running() {
			return nil
		}
	}

	return n.IDNameFilter.Filter(md)
}
