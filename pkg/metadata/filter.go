package metadata

import (
	"fmt"
	"github.com/luxas/ignite/pkg/filter"
	"strings"
)

// Compile-time assert to verify interface compatibility
var _ filter.Filter = &IDNameFilter{}

type IDNameFilter struct {
	prefix string
}

func NewIDNameFilter(p string) *IDNameFilter {
	return &IDNameFilter{
		prefix: p,
	}
}

func (n *IDNameFilter) Filter(f filter.Filterable) (bool, error) {
	md, ok := f.(*Metadata)
	if !ok {
		return false, fmt.Errorf("failed to assert Filterable %v to Metadata", f)
	}

	return strings.HasPrefix(md.ID, n.prefix) || strings.HasPrefix(md.Name, n.prefix), nil
}

func loadMetadata(id string, t ObjectType) (interface{}, error) {
	md := &Metadata{
		ID:   id,
		Type: t,
	}

	err := md.Load()
	return md, err
}

// TODO: Temporary until own type
func LoadImageMetadata(id string) (filter.Filterable, error) {
	return loadMetadata(id, Image)
}

// TODO: Temporary until own type
func LoadKernelMetadata(id string) (filter.Filterable, error) {
	return loadMetadata(id, Kernel)
}
