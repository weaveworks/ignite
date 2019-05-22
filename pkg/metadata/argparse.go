package metadata

import (
	"fmt"
	"io/ioutil"
	"strings"
)

// TODO: Filters, for example if a VM is running
// MatchObject gets the full ID of an object based on the given name/ID sample
func MatchObject(input string, objectType ObjectType) (string, error) {
	var object string

	entries, err := ioutil.ReadDir(objectType.Path())
	if err != nil {
		return "", err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			md := &Metadata{
				ID:   entry.Name(),
				Type: objectType,
			}

			if err := md.Load(); err != nil {
				return "", fmt.Errorf("failed to load metadata for %s object %q: %v", objectType, entry.Name(), err)
			}

			if strings.HasPrefix(md.ID, input) || strings.HasPrefix(md.Name, input) {
				if object != "" {
					return "", fmt.Errorf("ambiguous %s: %s", objectType, input)
				}

				object = md.ID
			}
		}
	}

	if object == "" {
		return "", fmt.Errorf("nonexistent %s: %s", objectType, input)
	}

	return object, nil
}
