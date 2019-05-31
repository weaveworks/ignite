package vmmd

import (
	"github.com/luxas/ignite/pkg/filter"
	"github.com/luxas/ignite/pkg/util"
)

// Compile-time assert to verify interface compatibility
var _ filter.Filter = &VMFilter{}

type VMFilter struct {
	prefix string
	all    bool
}

func NewVMFilter(p string) *VMFilter {
	return NewVMFilterAll(p, true)
}

func NewVMFilterAll(p string, all bool) *VMFilter {
	return &VMFilter{
		prefix: p,
		all:    all,
	}
}

func (n *VMFilter) Filter(f filter.Filterable) ([]string, error) {
	md, err := ToVMMetadata(f)
	if err != nil {
		return nil, err
	}

	// Option to list just running VMs
	if !n.all {
		if md.VMOD().State != Running {
			return nil, nil
		}
	}

	return util.MatchPrefix(n.prefix, md.ID, md.Name), nil
}

func LoadVMMetadata(id string) (*VMMetadata, error) {
	md := NewVMMetadata(id, "-", &VMObjectData{}) // A blank name triggers an unnecessary name generation
	err := md.Load()
	return md, err
}

func LoadVMMetadataFilterable(id string) (filter.Filterable, error) {
	return LoadVMMetadata(id)
}
