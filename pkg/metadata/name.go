package metadata

import (
	"fmt"
	"regexp"
	"strings"

	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/client"
	"github.com/weaveworks/ignite/pkg/filter"
	"github.com/weaveworks/ignite/pkg/storage/filterer"
	"github.com/weaveworks/ignite/pkg/util"
)

const (
	nameRegex = `^[a-z-_0-9.:/]*$`
)

func NewNameWithLatest(input string, kind meta.Kind) (string, error) {
	// Enforce a latest tag for images and kernels
	if !strings.Contains(input, ":") {
		input += ":latest"
	}

	return NewName(input, kind)
}

func NewName(input string, kind meta.Kind) (string, error) {
	matched, err := regexp.MatchString(nameRegex, input)
	if err != nil {
		return "", fmt.Errorf("failed to validate name input %q: %v", input, err)
	}

	if !matched {
		return "", fmt.Errorf("invalid name %q: does not match required format %s", input, nameRegex)
	}

	_, err = client.Dynamic(kind).Find(filter.NewNameFilter(input))
	switch err.(type) {
	case *filterer.NonexistentError:
		// The name is unique, no error
	case nil, *filterer.AmbiguousError:
		// The ambiguous error can only occur if someone manually created two Objects with the same name
		return "", fmt.Errorf("invalid name %q: already exists", input)
	default:
		return "", err
	}

	return input, nil
}

func InitName(md Metadata, input *string) {
	if input == nil { // If the input is nil (for loading purposes), create a temporary unset name
		md.SetName("<unset>")
	} else if *input == "" {
		md.SetName(util.RandomName()) // Otherwise if the input is unset, create a new random name
	} else {
		md.SetName(*input)
	}
}
