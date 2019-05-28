package vmmd

import (
	"fmt"
	"github.com/luxas/ignite/pkg/filter"
	"github.com/luxas/ignite/pkg/metadata"
	"strings"
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

func (n *VMFilter) Filter(f filter.Filterable) (bool, error) {
	md, ok := f.(*VMMetadata)
	if !ok {
		return false, fmt.Errorf("failed to assert Filterable %v to VMMetadata", f)
	}

	// Option to list just running VMs
	running := true
	if !n.all {
		running = md.VMOD().State == Running
	}

	return running && (strings.HasPrefix(md.ID, n.prefix) || strings.HasPrefix(md.Name, n.prefix)), nil
}

func LoadVMMetadata(id string) (filter.Filterable, error) {
	md := &VMMetadata{
		Metadata: &metadata.Metadata{
			ID:         id,
			Type:       metadata.VM,
			ObjectData: &VMObjectData{},
		},
	}

	err := md.Load()
	return md, err
}
