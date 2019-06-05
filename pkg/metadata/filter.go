package metadata

import (
	"fmt"
	"github.com/luxas/ignite/pkg/util"
	"strings"
)

type Match struct {
	Object  AnyMetadata
	Strings []string
}

// Match needs to be printable
var _ fmt.Stringer = &Match{}

func (m *Match) String() string {
	return strings.Join(m.Strings, " ")
}

type Filter interface {
	Filter(AnyMetadata) []string
	ErrNonexistent() error
	ErrAmbiguous([]*Match) error
}

// Compile-time assert to verify interface compatibility
var _ Filter = &IDNameFilter{}

type IDNameFilter struct {
	prefix     string
	filterType ObjectType
}

func NewIDNameFilter(p string, t ObjectType) *IDNameFilter {
	return &IDNameFilter{
		prefix:     p,
		filterType: t,
	}
}

func (f *IDNameFilter) Filter(any AnyMetadata) []string {
	md := any.GetMD()
	return util.MatchPrefix(f.prefix, md.ID, md.Name.String())
}

func (f *IDNameFilter) ErrNonexistent() error {
	return fmt.Errorf("can't find %s: no ID/name matches for %q", f.filterType, f.prefix)
}

func (f *IDNameFilter) ErrAmbiguous(matches []*Match) error {
	return fmt.Errorf("ambiguous %s query: %q matched the following IDs/names: %s", f.filterType, f.prefix, formatMatches(matches))
}

func formatMatches(matches []*Match) string {
	var sb strings.Builder

	for i, match := range matches {
		sb.WriteString(match.String())

		if i+1 < len(matches) {
			sb.WriteString(", ")
		}
	}

	return sb.String()
}
