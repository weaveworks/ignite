package update

import (
	"sort"
	"strings"
)

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

// Events is a sortable slice of Events
type Events []Event

// Events implements sort.Interface
var _ sort.Interface = Events{}

func (e Events) Len() int           { return len(e) }
func (e Events) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }
func (e Events) Less(i, j int) bool { return e[i] < e[j] }
func (e Events) String() string {
	strs := []string{}
	for _, ev := range e {
		strs = append(strs, ev.String())
	}
	return strings.Join(strs, ",")
}
