package watcher

import (
	"reflect"
	"testing"
)

var testEvents = []FileEvents{
	{
		FileEventDelete,
		EventCreate,
		FileEventModify,
	},
	{
		EventCreate,
		FileEventModify,
		FileEventDelete,
	},
	{
		EventCreate,
		FileEventModify,
		FileEventDelete,
		EventCreate,
	},
}

var targets = []FileEvents{
	{
		FileEventModify,
	},
	{
		FileEventNone,
	},
	{
		FileEventNone,
		EventCreate,
	},
}

func TestEventConcatenation(t *testing.T) {
	for i, e := range testEvents {
		result := concatenateEvents(e)
		if !reflect.DeepEqual(result, targets[i]) {
			t.Errorf("wrong concatenation result: %v != %v", result, targets[i])
		}
	}
}
