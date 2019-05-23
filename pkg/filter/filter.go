package filter

import (
	"errors"
	"io/ioutil"
)

var (
	errFilterNonexistent = errors.New("nonexistent")
	errFilterAmbiguous   = errors.New("ambiguous")
)

type LoadFunc func(string) (Filterable, error)

type Filter interface {
	Filter(Filterable) (bool, error)
}

type Filterable interface{}

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

func (o *filterer) match() ([]Filterable, error) {
	var matches []Filterable

	for _, f := range o.sources {
		ok, err := o.filter.Filter(f)
		if err != nil {
			return nil, err
		}

		if ok {
			matches = append(matches, f)
		}
	}

	return matches, nil
}

func (o *filterer) Single() (Filterable, error) {
	matches, err := o.match()
	if err != nil {
		return nil, err
	}

	if len(matches) == 0 {
		return nil, errFilterNonexistent
	}

	if len(matches) > 1 {
		return nil, errFilterAmbiguous
	}

	return matches[0], nil
}

func (o *filterer) All() ([]Filterable, error) {
	matches, err := o.match()
	if err != nil {
		return nil, err
	}

	return matches, nil
}
