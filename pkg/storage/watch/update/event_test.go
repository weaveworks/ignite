package update

import (
	"reflect"
	"sort"
	"testing"
)

func TestEventSort(t *testing.T) {
	e1 := Events{EventDelete, EventModify, EventCreate}
	e2 := Events{EventCreate, EventDelete, EventModify}

	sort.Sort(e1)
	sort.Sort(e2)
	if !reflect.DeepEqual(e1, e2) {
		t.Errorf("events do not match: %v %v", e1, e2)
	}
}
