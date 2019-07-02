package filter

import (
	"fmt"
	"github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"
	"github.com/weaveworks/ignite/pkg/snapshotter"
	"strings"

	"github.com/weaveworks/ignite/pkg/util"
)

// Compile-time assert to verify interface compatibility
var _ snapshotter.Filter = &IDNameFilter{}

// The IDNameFilter is the basic filter matching objects by their ID/name
type IDNameFilter struct {
	prefix     string
	matches    []string
	filterType v1alpha1.PoolDeviceType
}

func NewIDNameFilter(p string) *IDNameFilter {
	return &IDNameFilter{
		prefix: p,
	}
}

func (f *IDNameFilter) SetType(t v1alpha1.PoolDeviceType) {
	f.filterType = t
}

func (f *IDNameFilter) Filter(object *snapshotter.Object) (*snapshotter.Object, error) {
	mo, err := object.GetMetaObject()
	if err != nil {
		return nil, err
	}

	matches := util.MatchPrefix(f.prefix, string(mo.GetUID()), mo.GetName())
	if len(matches) > 0 {
		f.matches = append(f.matches, matches...)
		return object, nil
	}

	return nil, nil
}

func (f *IDNameFilter) ErrAmbiguous() snapshotter.ErrAmbiguous {
	return fmt.Errorf("ambiguous %s query: %q matched the following IDs/names: %s", f.filterType, f.prefix, formatMatches(f.matches))
}

func (f *IDNameFilter) ErrNonexistent() snapshotter.ErrNonexistent {
	return fmt.Errorf("can't find %s: no ID/name matches for %q", f.filterType, f.prefix)
}

func formatMatches(matches []string) string {
	var sb strings.Builder

	for i, match := range matches {
		sb.WriteString(match)

		if i+1 < len(matches) {
			sb.WriteString(", ")
		}
	}

	return sb.String()
}
