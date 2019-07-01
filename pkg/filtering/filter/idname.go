package filter

import (
	"fmt"
	"github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"
	"github.com/weaveworks/ignite/pkg/filtering/filterer"
	"strings"

	"github.com/weaveworks/ignite/pkg/util"
)

// Compile-time assert to verify interface compatibility
var _ filterer.Filter = &IDNameFilter{}

type IDNameFilter struct {
	prefix     string
	filterType v1alpha1.PoolDeviceType
}

func NewIDNameFilter(p string, t v1alpha1.PoolDeviceType) *IDNameFilter {
	return &IDNameFilter{
		prefix:     p,
		filterType: t,
	}
}

func (f *IDNameFilter) Filter(object v1alpha1.Object) []string {
	return util.MatchPrefix(f.prefix, string(object.GetUID()), object.GetName())
}

func (f *IDNameFilter) ErrNonexistent() filterer.NonexistentError {
	return fmt.Errorf("can't find %s: no ID/name matches for %q", f.filterType, f.prefix)
}

func (f *IDNameFilter) ErrAmbiguous(matches []*filterer.Match) filterer.AmbiguousError {
	return fmt.Errorf("ambiguous %s query: %q matched the following IDs/names: %s", f.filterType, f.prefix, formatMatches(matches))
}

func formatMatches(matches []*filterer.Match) string {
	var sb strings.Builder

	for i, match := range matches {
		sb.WriteString(match.String())

		if i+1 < len(matches) {
			sb.WriteString(", ")
		}
	}

	return sb.String()
}
