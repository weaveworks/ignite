package metadata

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/weaveworks/ignite/pkg/util"
)

const (
	nameRegex = `^[a-z-_0-9.:/]*$`
)

func NewNameWithLatest(input string, matches *[]Metadata) (string, error) {
	// Enforce a latest tag for images and kernels
	if !strings.Contains(input, ":") {
		input += ":latest"
	}

	return NewName(input, matches)
}

func NewName(input string, matches *[]Metadata) (string, error) {
	matched, err := regexp.MatchString(nameRegex, input)
	if err != nil {
		return "", fmt.Errorf("failed to validate name input %q: %v", input, err)
	}

	if !matched {
		return "", fmt.Errorf("invalid name %q: does not match required format %s", input, nameRegex)
	}

	// Check the given matches for uniqueness
	if matches != nil {
		for _, match := range *matches {
			if input == match.GetName() {
				return "", fmt.Errorf("invalid name %q: already exists", input)
			}
		}
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
