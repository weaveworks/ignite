package update

type Event int

const (
	EventCreate Event = iota + 1 // 1
	EventDelete                  // 2
	EventModify                  // 3
)

type FileUpdate struct {
	Event Event
	Path  string
}
