package filter

import "fmt"

// Compile-time asserts to verify interface compatibility
var _ error = &ErrNonexistent{}
var _ error = &ErrAmbiguous{}

type ErrNonexistent struct{}

func NewErrNonexistent() *ErrNonexistent {
	return &ErrNonexistent{}
}

type ErrAmbiguous struct {
	matches []string
}

func NewErrAmbiguous(matches []*match) *ErrAmbiguous {
	e := &ErrAmbiguous{}
	e.matches = make([]string, 0, len(matches))

	for _, match := range matches {
		e.matches = append(e.matches, match.strings...)
	}

	return e
}

func (e *ErrNonexistent) Error() string {
	return "can't find %s matching %q"
}

func (e *ErrAmbiguous) Error() string {
	return fmt.Sprintf("ambiguous %%s search, %%q matched %v", e.matches)
}
