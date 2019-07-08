package filter

import (
	"strings"

	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/storage/filterer"
	"github.com/weaveworks/ignite/pkg/util"
)

// The IDNameFilter is the basic filter matching objects by their ID/name
type IDNameFilter struct {
	prefix  string
	matches []string
	kind    meta.Kind
}

var _ filterer.MetaFilter = &IDNameFilter{}

func NewIDNameFilter(p string) *IDNameFilter {
	return &IDNameFilter{
		prefix: p,
	}
}

func (f *IDNameFilter) FilterMeta(object meta.Object) (meta.Object, error) {
	if len(f.kind) == 0 {
		f.kind = object.GetKind() // reflect.Indirect(reflect.ValueOf(object)).Type().Name()
	}

	if matches := util.MatchPrefix(f.prefix, string(object.GetUID()), object.GetName()); len(matches) > 0 {
		f.matches = append(f.matches, matches...)
		return object, nil
	}

	return nil, nil
}

func (f *IDNameFilter) SetKind(k meta.Kind) {
	f.kind = k
}

func (f *IDNameFilter) AmbiguousError() *filterer.AmbiguousError {
	return filterer.NewAmbiguousError("ambiguous %s query: %q matched the following IDs/names: %s", f.kind, f.prefix, formatMatches(f.matches))
}

func (f *IDNameFilter) NonexistentError() *filterer.NonexistentError {
	return filterer.NewNonexistentError("can't find %s: no ID/name matches for %q", f.kind, f.prefix)
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
