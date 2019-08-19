package filter

import (
	"github.com/weaveworks/gitops-toolkit/pkg/runtime"
	"github.com/weaveworks/gitops-toolkit/pkg/storage/filterer"
)

// The NameFilter matches Objects by their exact name
type NameFilter struct {
	name string
	kind runtime.Kind
}

var _ filterer.MetaFilter = &NameFilter{}

func NewNameFilter(n string) *NameFilter {
	return &NameFilter{
		name: n,
	}
}

func (f *NameFilter) FilterMeta(object runtime.Object) (filterer.Match, error) {
	if object.GetName() == f.name {
		return filterer.NewMatch(object, true), nil
	}

	return nil, nil
}

func (f *NameFilter) SetKind(k runtime.Kind) {
	f.kind = k
}

func (f *NameFilter) AmbiguousError(_ []filterer.Match) *filterer.AmbiguousError {
	return filterer.NewAmbiguousError("ambiguous %s query: %q matched multiple names", f.kind, f.name)
}

func (f *NameFilter) NonexistentError() *filterer.NonexistentError {
	return filterer.NewNonexistentError("can't find %s: no name matches for %q", f.kind, f.name)
}
