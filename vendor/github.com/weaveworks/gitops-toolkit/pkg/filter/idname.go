package filter

import (
	"strings"

	"github.com/weaveworks/gitops-toolkit/pkg/runtime"
	"github.com/weaveworks/gitops-toolkit/pkg/storage/filterer"
	"github.com/weaveworks/gitops-toolkit/pkg/util"
)

type IDNameMatch struct {
	*filterer.GenericMatch
	matches []string
}

var _ filterer.Match = &IDNameMatch{}

// The IDNameFilter is the basic filter matching objects by their ID/name
type IDNameFilter struct {
	prefix string
	kind   runtime.Kind
}

var _ filterer.MetaFilter = &IDNameFilter{}

func NewIDNameFilter(p string) *IDNameFilter {
	return &IDNameFilter{
		prefix: p,
	}
}

func (f *IDNameFilter) FilterMeta(object runtime.Object) (filterer.Match, error) {
	if len(f.kind) == 0 {
		f.kind = object.GetKind() // reflect.Indirect(reflect.ValueOf(object)).Type().Name()
	}

	if matches, exact := util.MatchPrefix(f.prefix, string(object.GetUID()), object.GetName()); len(matches) > 0 {
		return &IDNameMatch{
			filterer.NewMatch(object, exact),
			matches,
		}, nil
	}

	return nil, nil
}

func (f *IDNameFilter) SetKind(k runtime.Kind) {
	f.kind = k
}

func (f *IDNameFilter) AmbiguousError(matches []filterer.Match) *filterer.AmbiguousError {
	return filterer.NewAmbiguousError("ambiguous %s query: %q matched the following IDs/names: %s", f.kind, f.prefix, formatMatches(matches))
}

func (f *IDNameFilter) NonexistentError() *filterer.NonexistentError {
	return filterer.NewNonexistentError("can't find %s: no ID/name matches for %q", f.kind, f.prefix)
}

func formatMatches(input []filterer.Match) string {
	var sb strings.Builder
	var matches []string

	for _, match := range input {
		matches = append(matches, match.(*IDNameMatch).matches...)
	}

	for i, str := range matches {
		sb.WriteString(str)

		if i+1 < len(matches) {
			sb.WriteString(", ")
		}
	}

	return sb.String()
}
