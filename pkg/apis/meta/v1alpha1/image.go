package v1alpha1

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/containers/image/docker/reference"
	"github.com/opencontainers/go-digest"
)

// NewOCIImageRef parses and normalizes a reference to an OCI (docker) image.
func NewOCIImageRef(imageStr string) (OCIImageRef, error) {
	named, err := reference.ParseDockerRef(imageStr)
	if err != nil {
		return "", err
	}
	namedTagged, ok := named.(reference.NamedTagged)
	if !ok {
		return "", fmt.Errorf("could not parse image name with a tag %s", imageStr)
	}
	return OCIImageRef(reference.FamiliarString(namedTagged)), nil
}

type OCIImageRef string

var _ fmt.Stringer = OCIImageRef("")

func (i OCIImageRef) String() string {
	return string(i)
}

func (i OCIImageRef) IsUnset() bool {
	return len(i) == 0
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

func ParseOCIContentID(str string) (*OCIContentID, error) {
	named, err := reference.ParseDockerRef(str)
	if err != nil {
		return nil, err
	}

	if canonical, ok := named.(reference.Canonical); ok {
		return &OCIContentID{
			repoName: named.Name(),
			digest:   canonical.Digest().String(),
		}, nil
	}

	d, err := digest.Parse(str)
	if err != nil {
		return nil, err
	}

	return &OCIContentID{
		digest: d.String(),
	}, nil
}

type OCIContentID struct {
	repoName string // Fully qualified image name, e.g. "docker.io/library/node" or blank if the image is local
	digest   string // Repo digest of the image, or sha256sum provided by the source if the image is local
}

var _ json.Marshaler = &OCIContentID{}
var _ json.Unmarshaler = &OCIContentID{}

func (o *OCIContentID) String() string {
	if len(o.repoName) > 0 {
		return fmt.Sprintf("oci://%s", o.RepoDigest())
	}

	return fmt.Sprintf("docker://%s", o.Digest())
}

func parseOCIString(s string) (*OCIContentID, error) {
	u, err := url.Parse(s)
	if err != nil {
		return nil, err
	}

	// Remove the "docker://" or "oci://" scheme by only caring about the host and path
	return ParseOCIContentID(u.Host + u.Path)
}

// Local returns true if the image has no repoName, i.e. it's not available from a registry
func (o *OCIContentID) Local() bool {
	return len(o.repoName) == 0
}

// Digest is a getter for the digest field
func (o *OCIContentID) Digest() string {
	return o.digest
}

// RepoDigest returns a repo digest based on the OCIContentID if it is not local
func (o *OCIContentID) RepoDigest() (s string) {
	if !o.Local() {
		s = fmt.Sprintf("%s@%s", o.repoName, o.digest)
	}

	return
}

func (o *OCIContentID) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.String())
}

func (o *OCIContentID) UnmarshalJSON(b []byte) (err error) {
	var s string
	if err = json.Unmarshal(b, &s); err != nil {
		return err
	}

	var id *OCIContentID
	if id, err = parseOCIString(s); err == nil {
		*o = *id
	}

	return
}
