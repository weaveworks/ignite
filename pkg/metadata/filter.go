package metadata

import (
	"fmt"
	"github.com/luxas/ignite/pkg/filter"
	"github.com/luxas/ignite/pkg/util"
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

func (n *IDNameFilter) Filter(f filter.Filterable) ([]string, error) {
	md, ok := f.(*Metadata)
	if !ok {
		return nil, fmt.Errorf("failed to assert Filterable %v to Metadata", f)
	}

	return util.MatchPrefix(n.prefix, md.ID, md.Name), nil
}
