package update

type Event int

const (
	EventCreate Event = iota + 1
	EventDelete
	EventModify
)

type FileUpdate struct {
	Event Event
	Path  string
}
