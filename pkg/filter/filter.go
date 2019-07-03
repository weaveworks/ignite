package filter

import (
	"github.com/weaveworks/ignite/pkg/metadata"
)

type filterer struct {
	filter  metadata.Filter
	sources []metadata.Metadata
}

func NewFilterer(filter metadata.Filter, sources []metadata.Metadata) (*filterer, error) {
	return &filterer{
		filter:  filter,
		sources: sources,
	}, nil
}

func (f *filterer) match() []*metadata.Match {
	var matches []*metadata.Match

	for _, object := range f.sources {
		strings := f.filter.Filter(object)

		if len(strings) > 0 {
			matches = append(matches, &metadata.Match{
				Object:  object,
				Strings: strings,
			})
		}
	}

	return matches
}

func (f *filterer) Single() (metadata.Metadata, error) {
	matches := f.match()

	if len(matches) == 0 {
		return nil, f.filter.ErrNonexistent()
	}

	if len(matches) > 1 {
		return nil, f.filter.ErrAmbiguous(matches)
	}

	return matches[0].Object, nil
}

func (f *filterer) All() []metadata.Metadata {
	matches := f.match()

	objects := make([]metadata.Metadata, 0, len(matches))
	for _, match := range matches {
		objects = append(objects, match.Object)
	}

	return objects
}
