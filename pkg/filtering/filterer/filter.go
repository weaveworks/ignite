package filterer

import (
	"fmt"
	"github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"
	"strings"
)

type Match struct {
	Object  v1alpha1.Object
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
	Filter(v1alpha1.Object) []string
	ErrNonexistent() NonexistentError
	ErrAmbiguous([]*Match) AmbiguousError
}
