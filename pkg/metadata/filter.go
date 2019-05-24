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
