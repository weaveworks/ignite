package filter

import (
	"io/ioutil"
)

type LoadFunc func(string) (Filterable, error)

type Filter interface {
	Filter(Filterable) ([]string, error)
}

type Filterable interface{}

type match struct {
	object  Filterable
	strings []string
}

type filterer struct {
	filter  Filter
	sources []Filterable
}

func NewFilterer(filter Filter, path string, loadFunc LoadFunc) (*filterer, error) {
	var sources []Filterable

	entries, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			filterable, err := loadFunc(entry.Name())
			if err != nil {
				return nil, err
			}

			sources = append(sources, filterable)
		}
	}

	return &filterer{
		filter:  filter,
		sources: sources,
	}, nil
}

func (f *filterer) match() ([]*match, error) {
	var matches []*match

	for _, object := range f.sources {
		strings, err := f.filter.Filter(object)
		if err != nil {
			return nil, err
		}

		if len(strings) > 0 {
			matches = append(matches, &match{
				object:  object,
				strings: strings,
			})
		}
	}

	return matches, nil
}

func (f *filterer) Single() (Filterable, error) {
	matches, err := f.match()
	if err != nil {
		return nil, err
	}

	if len(matches) == 0 {
		return nil, NewErrNonexistent()
	}

	if len(matches) > 1 {
		return nil, NewErrAmbiguous(matches)
	}

	return matches[0].object, nil
}

func (f *filterer) All() ([]Filterable, error) {
	matches, err := f.match()
	if err != nil {
		return nil, err
	}

	objects := make([]Filterable, 0, len(matches))
	for _, match := range matches {
		objects = append(objects, match.object)
	}

	return objects, nil
}
