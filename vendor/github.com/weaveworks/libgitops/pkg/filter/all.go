package filter

import (
	"github.com/weaveworks/libgitops/pkg/runtime"
	"github.com/weaveworks/libgitops/pkg/storage/filterer"
)

// The AllFilter matches anything it's given
type AllFilter struct{}

// It's more efficient for this to be an ObjectFilter, as it loads everything anyways
var _ filterer.ObjectFilter = &AllFilter{}

func NewAllFilter() *AllFilter {
	return &AllFilter{}
}

func (f *AllFilter) Filter(object runtime.Object) (filterer.Match, error) {
	return filterer.NewMatch(object, false), nil
}

// The AllFilter shouldn't be used to match single Objects
func (f *AllFilter) AmbiguousError(_ []filterer.Match) *filterer.AmbiguousError {
	return filterer.NewAmbiguousError("ambiguous query: AllFilter used to match single Object")
}

func (f *AllFilter) NonexistentError() *filterer.NonexistentError {
	return filterer.NewNonexistentError("no results: AllFilter used to match single Object")
}

// SetKind is a no-op as this filter doesn't have a dedicated kind
func (f *AllFilter) SetKind(_ runtime.Kind) {}
