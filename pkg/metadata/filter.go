package metadata

import (
	"fmt"
	"strings"

	api "github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"

	"github.com/weaveworks/ignite/pkg/util"
)

type Match struct {
	Object  Metadata
	Strings []string
}

// Match needs to be printable
var _ fmt.Stringer = &Match{}

func (m *Match) String() string {
	return strings.Join(m.Strings, " ")
}

type NonexistentError struct {
	error
}

type AmbiguousError struct {
	error
}

type Filter interface {
	Filter(Metadata) []string
	ErrNonexistent() error
	ErrAmbiguous([]*Match) error
}

// Compile-time assert to verify interface compatibility
var _ Filter = &IDNameFilter{}

type IDNameFilter struct {
	prefix     string
	filterType api.PoolDeviceType
}

func NewIDNameFilter(p string, t api.PoolDeviceType) *IDNameFilter {
	return &IDNameFilter{
		prefix:     p,
		filterType: t,
	}
}

func (f *IDNameFilter) Filter(md Metadata) []string {
	return util.MatchPrefix(f.prefix, md.GetUID().String(), md.GetName())
}

func (f *IDNameFilter) ErrNonexistent() error {
	return &NonexistentError{fmt.Errorf("can't find %s: no ID/name matches for %q", f.filterType, f.prefix)}
}

func (f *IDNameFilter) ErrAmbiguous(matches []*Match) error {
	return &AmbiguousError{fmt.Errorf("ambiguous %s query: %q matched the following IDs/names: %s", f.filterType, f.prefix, formatMatches(matches))}
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
