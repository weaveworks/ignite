package update

import "fmt"

// ObjectEvent is an enum describing a change in an Object's state.
type ObjectEvent byte

var _ fmt.Stringer = ObjectEvent(0)

const (
	ObjectEventNone   ObjectEvent = iota // 0
	ObjectEventCreate                    // 1
	ObjectEventModify                    // 2
	ObjectEventDelete                    // 3
)

func (o ObjectEvent) String() string {
	switch o {
	case 0:
		return "NONE"
	case 1:
		return "CREATE"
	case 2:
		return "MODIFY"
	case 3:
		return "DELETE"
	}

	// Should never happen
	return "UNKNOWN"
}
