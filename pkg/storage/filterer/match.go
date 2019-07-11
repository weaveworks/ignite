package filterer

import (
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
)

// Match describes the result of filtering an Object
// If the Object to be filtered didn't match, return nil
type Match interface {
	// Get the matched Object
	Object() meta.Object
	// Check if the match was exact
	Exact() bool
}

// GenericMatch is the simplest implementation
// of Match, carrying no additional data
type GenericMatch struct {
	object meta.Object
	exact  bool
}

var _ Match = &GenericMatch{}

func NewMatch(object meta.Object, exact bool) *GenericMatch {
	return &GenericMatch{
		object: object,
		exact:  exact,
	}
}

func (m *GenericMatch) Object() meta.Object {
	return m.object
}

func (m *GenericMatch) Exact() bool {
	return m.exact
}
