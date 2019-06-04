package metadata

import (
	"encoding/json"
	"fmt"
	"github.com/luxas/ignite/pkg/util"
)

type Name struct {
	string
}

// Compile-time assert to verify interface compatibility
var _ fmt.Stringer = &Name{}

func (n *Name) randomize() {
	if n.string == "" {
		n.string = util.RandomName()
	}
}

func NewName(input string) (*Name, error) {
	// TODO: Check if input is valid
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
	return json.Marshal(n.String())
}

func (n *Name) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	name, err := NewName(s)
	if err != nil {
		return err
	}

	*n = *name

	return nil
}
