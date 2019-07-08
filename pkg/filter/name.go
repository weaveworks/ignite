package filter

import (
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/storage/filterer"
)

// The NameFilter matches Objects by their exact name
type NameFilter struct {
	name    string
	matches uint64
	kind    meta.Kind
}

var _ filterer.MetaFilter = &NameFilter{}

func NewNameFilter(n string) *NameFilter {
	return &NameFilter{
		name: n,
	}
}

func (f *NameFilter) FilterMeta(object meta.Object) (meta.Object, error) {
	if len(f.kind) == 0 {
		f.kind = object.GetKind() // reflect.Indirect(reflect.ValueOf(object)).Type().Name()
	}

	if object.GetName() == f.name {
		f.matches++
		return object, nil
	}

	return nil, nil
}

func (f *NameFilter) SetKind(k meta.Kind) {
	f.kind = k
}

func (f *NameFilter) AmbiguousError() *filterer.AmbiguousError {
	return filterer.NewAmbiguousError("ambiguous %s query: %q matched %d names", f.kind, f.name, f.matches)
}

func (f *NameFilter) NonexistentError() *filterer.NonexistentError {
	return filterer.NewNonexistentError("can't find %s: no name matches for %q", f.kind, f.name)
}
