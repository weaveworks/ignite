package watcher

import (
	"reflect"
	"testing"
)

var testEvents = []Events{
	{
		EventDelete,
		EventCreate,
		EventModify,
	},
	{
		EventCreate,
		EventModify,
		EventDelete,
	},
	{
		EventCreate,
		EventModify,
		EventDelete,
		EventCreate,
	},
}

var targets = []Events{
	{
		EventModify,
	},
	{
		EventNone,
	},
	{
		EventNone,
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
