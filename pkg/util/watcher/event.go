package watcher

import (
	"fmt"
	"strings"
)

// FileEvent is an enum describing a change in a file's state
type FileEvent byte

const (
	FileEventNone   FileEvent = iota // 0
	FileEventModify                  // 1
	FileEventDelete                  // 2
	FileEventMove                    // 3
)

func (e FileEvent) String() string {
	switch e {
	case 0:
		return "NONE"
	case 1:
		return "MODIFY"
	case 2:
		return "DELETE"
	case 3:
		return "MOVE"
	}

	return "UNKNOWN"
}

// FileEvents is a slice of FileEvents
type FileEvents []FileEvent

var _ fmt.Stringer = FileEvents{}

func (e FileEvents) String() string {
	strs := make([]string, 0, len(e))
	for _, ev := range e {
		strs = append(strs, ev.String())
	}

	return strings.Join(strs, ",")
}

func (e FileEvents) Bytes() []byte {
	b := make([]byte, 0, len(e))
	for _, event := range e {
		b = append(b, byte(event))
	}

	return b
}

// FileUpdates is a slice of FileUpdate pointers
type FileUpdates []*FileUpdate

// FileUpdate is used by watchers to
// signal the state change of a file.
type FileUpdate struct {
	Event FileEvent
	Path  string
}
