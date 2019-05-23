package metadata

import (
	"fmt"
	"io/ioutil"
)

type LoadFunc func(string) (Filterable, error)

type objectMatcher struct {
	filter  Filter
	matches []*Filterable
}

func NewObjectMatcher(filter Filter) *objectMatcher {
	return &objectMatcher{
		filter:  filter,
		matches: []*Filterable{},
	}
}

// TODO: Filters, for example if a VM is running
// MatchObject gets the full IDs of matching objects based on the given name/ID sample
func (o *objectMatcher) match(objectType ObjectType, lf LoadFunc) error {
	entries, err := ioutil.ReadDir(objectType.Path())
	if err != nil {
		return err
	}

	// Clear previous matches
	o.matches = nil

	for _, entry := range entries {
		if entry.IsDir() {
			md, err := lf(entry.Name())
			if err != nil {
				return fmt.Errorf("failed to load metadata for %s %q: %v", objectType, entry.Name(), err)
			}

			if md.Matches(o.filter) {
				o.matches = append(o.matches, &md)
			}
		}
	}

	return nil
}

func (o *objectMatcher) Single(objectType ObjectType, lf func(string) (Filterable, error)) (*Filterable, error) {
	if err := o.match(objectType, lf); err != nil {
		return nil, err
	}

	if len(o.matches) == 0 {
		return nil, fmt.Errorf("nonexistent %s: %s", objectType, o.filter)
	}

	if len(o.matches) > 1 {
		return nil, fmt.Errorf("ambiguous %s: %s", objectType, o.filter)
	}

	return o.matches[0], nil
}

func (o *objectMatcher) All(objectType ObjectType, lf func(string) (Filterable, error)) ([]*Filterable, error) {
	if err := o.match(objectType, lf); err != nil {
		return nil, err
	}

	return o.matches, nil
}
