package filterer

import (
	"github.com/weaveworks/libgitops/pkg/runtime"
)

// Match describes the result of filtering an Object
// If the Object to be filtered didn't match, return nil
type Match interface {
	// Get the matched Object
	Object() runtime.Object
	// Check if the match was exact
	Exact() bool
}

// GenericMatch is the simplest implementation
// of Match, carrying no additional data
type GenericMatch struct {
	object runtime.Object
	exact  bool
}

var _ Match = &GenericMatch{}

func NewMatch(object runtime.Object, exact bool) *GenericMatch {
	return &GenericMatch{
		object: object,
		exact:  exact,
	}
}

func (m *GenericMatch) Object() runtime.Object {
	return m.object
}

func (m *GenericMatch) Exact() bool {
	return m.exact
}
