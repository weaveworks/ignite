package filterer

import (
	"fmt"

	"github.com/weaveworks/libgitops/pkg/runtime"
)

// BaseFilter provides shared functionality for filter types
type BaseFilter interface {
	// AmbiguousError specifies what to error if
	// a single request returned multiple matches
	// The matches are given as an argument
	AmbiguousError([]Match) *AmbiguousError
	// NonexistentError specifies what to error if
	// a single request returned no matches
	NonexistentError() *NonexistentError
	// SetKind sets the kind for the filter
	SetKind(runtime.Kind)
}

// ObjectFilter implementations filter fully loaded runtime.Objects
type ObjectFilter interface {
	BaseFilter
	// Every Object to be filtered is passed though Filter, which should
	// return the Object on match, or nil if it doesn't match.
	// The boolean indicates an exact match.
	Filter(runtime.Object) (Match, error)
}

// MetaFilter implementations operate on runtime.APIType objects,
// which are more light weight, but provide only name/UID matching.
type MetaFilter interface {
	BaseFilter
	// Every Object to be filtered is passed though FilterMeta, which should
	// return the Object on match, or nil if it doesn't match. The Objects
	// given to FilterMeta are of type runtime.APIType, stripped of other contents.
	// The boolean indicates an exact match.
	FilterMeta(runtime.Object) (Match, error)
}

type AmbiguousError struct {
	error
}

func NewAmbiguousError(format string, data ...interface{}) *AmbiguousError {
	return &AmbiguousError{
		fmt.Errorf(format, data...),
	}
}

func IsAmbiguousError(err error) bool {
	_, ok := err.(*AmbiguousError)
	return ok
}

type NonexistentError struct {
	error
}

func NewNonexistentError(format string, data ...interface{}) *NonexistentError {
	return &NonexistentError{
		fmt.Errorf(format, data...),
	}
}

func IsNonexistentError(err error) bool {
	_, ok := err.(*NonexistentError)
	return ok
}
