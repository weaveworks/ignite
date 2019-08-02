package update

// ObjectEvent is an enum describing a change in an Object's state.
type ObjectEvent byte

const (
	ObjectEventNone   ObjectEvent = iota // 0
	ObjectEventCreate                    // 1
	ObjectEventModify                    // 2
	ObjectEventDelete                    // 3
)
