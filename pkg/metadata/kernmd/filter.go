package kernmd

import (
	"github.com/luxas/ignite/pkg/filter"
	"github.com/luxas/ignite/pkg/util"
)

// Compile-time assert to verify interface compatibility
var _ filter.Filter = &KernelFilter{}

type KernelFilter struct {
	prefix string
}

func NewKernelFilter(p string) *KernelFilter {
	return &KernelFilter{
		prefix: p,
	}
}

func (n *KernelFilter) Filter(f filter.Filterable) ([]string, error) {
	md, err := ToKernelMetadata(f)
	if err != nil {
		return nil, err
	}

	return util.MatchPrefix(n.prefix, md.ID, md.Name), nil
}

func LoadKernelMetadata(id string) (*KernelMetadata, error) {
	md := NewKernelMetadata(id, "-") // A blank name triggers an unnecessary name generation
	err := md.Load()
	return md, err
}

func LoadKernelMetadataFilterable(id string) (filter.Filterable, error) {
	return LoadKernelMetadata(id)
}
