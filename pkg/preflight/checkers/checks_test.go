package checkers

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/weaveworks/ignite/pkg/preflight"
)

type fakeChecker struct {
	info string
}

func (fc fakeChecker) Check() error {
	if fc.info == "error" {
		return fmt.Errorf("error")
	}
	return nil
}

func (fc fakeChecker) Name() string {
	return "FakeChecker"
}

func (fc fakeChecker) Type() string {
	return "FakeChecker"
}

func TestRunChecks(t *testing.T) {
	utests := []struct {
		checkers      []preflight.Checker
		ignoredErrors []string
		expectedError bool
	}{
		{
			checkers:      []preflight.Checker{fakeChecker{info: ""}},
			ignoredErrors: []string{},
			expectedError: false,
		},
		{
			checkers:      []preflight.Checker{fakeChecker{info: ""}, fakeChecker{info: "error"}},
			ignoredErrors: []string{},
			expectedError: true,
		},
		{
			checkers:      []preflight.Checker{fakeChecker{info: ""}, fakeChecker{info: "error"}},
			ignoredErrors: []string{"fakechecker"},
			expectedError: false,
		},
	}
	for _, utest := range utests {
		ignoredErrors := sets.NewString(utest.ignoredErrors...)
		err := runChecks(utest.checkers, ignoredErrors)
		assert.Equal(t, utest.expectedError, (err != nil))
	}
}

func TestIsIgnoredPreflightError(t *testing.T) {
	const (
		all             = "all"
		ignoredError    = "my-ignored-error"
		notIgnoredError = "not-ignored-error"
	)
	utests := []struct {
		name             string
		ignoredErrors    []string
		searchedError    string
		expectedResponse bool
	}{
		{
			name:             "IgnoreAll",
			ignoredErrors:    []string{all},
			searchedError:    notIgnoredError,
			expectedResponse: true,
		},
		{
			name:             "IgnoreSpecificError1",
			ignoredErrors:    []string{ignoredError},
			searchedError:    ignoredError,
			expectedResponse: true,
		},
		{
			name:             "IgnoreSpecificError2",
			ignoredErrors:    []string{ignoredError},
			searchedError:    notIgnoredError,
			expectedResponse: false,
		},
		{
			name:             "IgnoreAllAndSpecificError1",
			ignoredErrors:    []string{all, ignoredError},
			searchedError:    ignoredError,
			expectedResponse: true,
		},
		{
			name:             "IgnoreAllAndSpecificError2",
			ignoredErrors:    []string{all, ignoredError},
			searchedError:    notIgnoredError,
			expectedResponse: true,
		},
		{
			name:             "NoIgnore",
			ignoredErrors:    []string{},
			searchedError:    notIgnoredError,
			expectedResponse: false,
		},
	}

	for _, utest := range utests {
		t.Run(utest.name, func(t *testing.T) {
			ignoredErrors := sets.NewString(utest.ignoredErrors...)
			isIgnored := isIgnoredPreflightError(ignoredErrors, utest.searchedError)
			assert.Equal(t, utest.expectedResponse, isIgnored)
		})
	}
}
