package filterer

import (
	"github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"
)

type filterer struct {
	filter  Filter
	sources []v1alpha1.Object
}

func NewFilterer(filter Filter, sources []v1alpha1.Object) (*filterer, error) {
	return &filterer{
		filter:  filter,
		sources: sources,
	}, nil
}

func (f *filterer) match() []*Match {
	var matches []*Match

	for _, object := range f.sources {
		strings := f.filter.Filter(object)

		if len(strings) > 0 {
			matches = append(matches, &Match{
				Object:  object,
				Strings: strings,
			})
		}
	}

	return matches
}

func (f *filterer) Single() (v1alpha1.Object, error) {
	matches := f.match()

	if len(matches) == 0 {
		return nil, f.filter.ErrNonexistent()
	}

	if len(matches) > 1 {
		return nil, f.filter.ErrAmbiguous(matches)
	}

	return matches[0].Object, nil
}

func (f *filterer) All() []v1alpha1.Object {
	matches := f.match()

	objects := make([]v1alpha1.Object, 0, len(matches))
	for _, match := range matches {
		objects = append(objects, match.Object)
	}

	return objects
}
