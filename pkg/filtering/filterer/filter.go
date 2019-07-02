package filterer

import (
	"fmt"
	"strings"

	ignitemeta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
)

type Match struct {
	Object  ignitemeta.Object
	Strings []string
}

// Match needs to be printable
var _ fmt.Stringer = &Match{}

func (m *Match) String() string {
	return strings.Join(m.Strings, " ")
}

type NonexistentError error
type AmbiguousError error

type Filter interface {
	Filter(ignitemeta.Object) []string
	ErrNonexistent() NonexistentError
	ErrAmbiguous([]*Match) AmbiguousError
}
