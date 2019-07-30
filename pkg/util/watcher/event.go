package watcher

import (
	"fmt"
	"strings"
)

// Event is an enum describing a change in a file's/Object's state.
// Unknown state changes can be signaled with a zero value.
type Event byte

const (
	EventNone   Event = iota // 0
	EventCreate              // 1
	EventDelete              // 2
	EventModify              // 3
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

// Events is a slice of Events
type Events []Event

var _ fmt.Stringer = Events{}

func (e Events) String() string {
	strs := make([]string, 0, len(e))
	for _, ev := range e {
		strs = append(strs, ev.String())
	}

	return strings.Join(strs, ",")
}

func (e Events) Bytes() []byte {
	b := make([]byte, 0, len(e))
	for _, event := range e {
		b = append(b, byte(event))
	}

	return b
}

// FileUpdate is used by watchers to
// signal the state change of a file.
type FileUpdate struct {
	Event Event
	Path  string
}