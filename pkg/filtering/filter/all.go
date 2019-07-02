package filter

import (
	"fmt"
	"github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"
	"github.com/weaveworks/ignite/pkg/snapshotter"
)

// Compile-time assert to verify interface compatibility
var _ snapshotter.Filter = &AllFilter{}

// The AllFilter matches everything it's given
type AllFilter struct {
	filterType v1alpha1.PoolDeviceType
}

func NewAllFilter() *AllFilter {
	return &AllFilter{}
}

func (f *AllFilter) SetType(t v1alpha1.PoolDeviceType) {
	f.filterType = t
}

func (f *AllFilter) Filter(object *snapshotter.Object) (*snapshotter.Object, error) {
	return object, nil
}

func (f *AllFilter) ErrAmbiguous() snapshotter.ErrAmbiguous {
	return fmt.Errorf("ambiguous %s query: AllFilter used to match single object", f.filterType)
}

func (f *AllFilter) ErrNonexistent() snapshotter.ErrNonexistent {
	return fmt.Errorf("no %s objects to query", f.filterType)
}
