package v1alpha1

import (
	"encoding/json"
	"fmt"

	"github.com/containers/image/docker/reference"
)

// NewOCIImageRef parses and normalizes a reference to an OCI (docker) image.
func NewOCIImageRef(imageStr string) (OCIImageRef, error) {
	named, err := reference.ParseDockerRef(imageStr)
	if err != nil {
		return OCIImageRef(""), err
	}
	namedTagged, ok := named.(reference.NamedTagged)
	if !ok {
		return OCIImageRef(""), fmt.Errorf("could not parse image name with a tag %s", imageStr)
	}
	return OCIImageRef(reference.FamiliarString(namedTagged)), nil
}

type OCIImageRef string

func (i OCIImageRef) String() string {
	return string(i)
}

func (i OCIImageRef) IsUnset() bool {
	return len(i) == 0
}

func (i OCIImageRef) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(i))
}

func (i *OCIImageRef) UnmarshalJSON(b []byte) error {
	var str string
	var err error

	if err = json.Unmarshal(b, &str); err != nil {
		return err
	}

	*i, err = NewOCIImageRef(str)
	return err
}
