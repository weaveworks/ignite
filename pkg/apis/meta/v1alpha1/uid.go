package v1alpha1

import "fmt"

// UID represents an unique ID for a type
type UID string

var _ fmt.Stringer = UID("")

// String returns the UID in string representation
func (u UID) String() string {
	return string(u)
}
