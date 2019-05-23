package filter

import (
	"errors"
	"fmt"
)

var (
	errFilterNonexistent = errors.New("nonexistent")
	errFilterAmbiguous   = errors.New("ambiguous")
)

type Filter interface {
	Filter(Filterable) (bool, error)
}

type Filterable interface {
	//Matches(Filter) bool
}

type filterer struct {
	filter  Filter
	matches []Filterable
}

func NewFilterer(filter Filter) *filterer {
	return &filterer{
		filter:  filter,
		matches: []Filterable{},
	}
}

func (o *filterer) match(input []Filterable) error {
	// Clear previous matches
	o.matches = nil

	fmt.Printf("input: %v\n", input)

	for _, f := range input {
		ok, err := o.filter.Filter(f)
		if err != nil {
			return err
		}

		if ok {
			o.matches = append(o.matches, f)
		}
	}

	fmt.Printf("matches: %v\n", o.matches)

	return nil
}

func (o *filterer) Single(input []Filterable) (Filterable, error) {
	if err := o.match(input); err != nil {
		return nil, err
	}

	if len(o.matches) == 0 {
		return nil, errFilterNonexistent
	}

	if len(o.matches) > 1 {
		return nil, errFilterAmbiguous
	}

	return o.matches[0], nil
}

func (o *filterer) All(input []Filterable) ([]Filterable, error) {
	if err := o.match(input); err != nil {
		return nil, err
	}

	return o.matches, nil
}
