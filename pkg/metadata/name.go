package metadata

import (
	"encoding/json"
	"fmt"

	"github.com/luxas/ignite/pkg/util"
	"regexp"
)

type Name struct {
	string
}

const (
	nameRegex = "^[a-z-_]*$"
)

// Compile-time assert to verify interface compatibility
var _ fmt.Stringer = &Name{}

func (n *Name) randomize() {
	if n.string == "" {
		n.string = util.RandomName()
	}
}

func NewName(input string, matches *[]*Name) (*Name, error) {
	matched, err := regexp.MatchString(nameRegex, input)
	if err != nil {
		return nil, fmt.Errorf("failed to validate name input %q: %v", input, err)
	}

	if !matched {
		return nil, fmt.Errorf("invalid name %q: does not match required format %s", input, nameRegex)
	}

	// Check the given matches for uniqueness
	if matches != nil {
		for _, match := range *matches {
			if input == match.string {
				return nil, fmt.Errorf("invalid name %q: not unique", input)
			}
		}
	}

	return &Name{
		input,
	}, nil
}

func newUnsetName() *Name {
	return &Name{
		"<unset>", // This should never be visible
	}
}

func (n *Name) String() string {
	return n.string
}

func (n Name) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.string)
}

func (n *Name) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	name, err := NewName(s, nil)
	if err != nil {
		return err
	}

	*n = *name

	return nil
}
