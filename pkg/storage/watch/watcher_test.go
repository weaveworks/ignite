package watch

import (
	"reflect"
	"testing"

	"github.com/weaveworks/ignite/pkg/storage/watch/update"
)

var testEvents = []update.Events{
	{
		update.EventDelete,
		update.EventCreate,
		update.EventModify,
	},
	{
		update.EventCreate,
		update.EventModify,
		update.EventDelete,
	},
	{
		update.EventCreate,
		update.EventModify,
		update.EventDelete,
		update.EventCreate,
	},
}

var targets = []update.Events{
	{
		update.EventModify,
	},
	{
		update.EventNone,
	},
	{
		update.EventNone,
		update.EventCreate,
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
