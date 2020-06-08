package runtime

import (
	"fmt"
	"strconv"
	"unicode/utf8"
)

// UID represents an unique ID for a type
type UID string

var _ fmt.Stringer = UID("")

// String returns the UID in string representation
func (u UID) String() string {
	return string(u)
}

// This unmarshaler enables the UID to be passed in as an
// unquoted string in JSON. Upon marshaling, quotes will
// be automatically added.
func (u *UID) UnmarshalJSON(b []byte) error {
	if !utf8.Valid(b) {
		return fmt.Errorf("invalid UID string: %s", b)
	}

	uid, err := strconv.Unquote(string(b))
	if err != nil {
		return err
	}

	*u = UID(uid)
	return nil
}
