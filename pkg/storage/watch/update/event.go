package update

// Event is an enum describing a change in a file's/Object's state.
// Unknown state changes can be signaled with a zero value.
type Event uint8

const (
	EventCreate Event = iota + 1 // 1
	EventDelete                  // 2
	EventModify                  // 3
)

func (e Event) String() string {
	switch e {
	case 1:
		return "CREATE"
	case 2:
		return "DELETE"
	case 3:
		return "MODIFY"
	}

	return "NONE"
}
