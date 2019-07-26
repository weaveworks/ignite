package update

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
	}

	return "MODIFY"
}

type FileUpdate struct {
	Event Event
	Path  string
}
