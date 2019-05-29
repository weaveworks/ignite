package kernmd

import (
	"fmt"
	"github.com/luxas/ignite/pkg/filter"
	"strings"
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

func (n *KernelFilter) Filter(f filter.Filterable) (bool, error) {
	md, ok := f.(*KernelMetadata)
	if !ok {
		return false, fmt.Errorf("failed to assert Filterable %v to KernelMetadata", f)
	}

	return strings.HasPrefix(md.ID, n.prefix) || strings.HasPrefix(md.Name, n.prefix), nil
}

func LoadKernelMetadata(id string) (*KernelMetadata, error) {
	md := NewKernelMetadata(id, "-") // A blank name triggers an unnecessary name generation
	err := md.Load()
	return md, err
}

func LoadKernelMetadataFilterable(id string) (filter.Filterable, error) {
	return LoadKernelMetadata(id)
}
