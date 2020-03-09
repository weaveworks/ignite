package v1alpha1

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	// "github.com/opencontainers/go-digest" requires us to load the algorithms
	// that we want to use into the binary. (it calls algorithm.Available)
	_ "crypto/sha256"

	"github.com/containers/image/docker/reference"
	"github.com/opencontainers/go-digest"
)

const (
	ociSchemeRegistry = "oci://"
	ociSchemeLocal    = "docker://"
)

// NewOCIImageRef parses and normalizes a reference to an OCI (docker) image.
func NewOCIImageRef(imageStr string) (o OCIImageRef, err error) {
	named, err := reference.ParseDockerRef(imageStr)
	if err != nil {
		return
	}

	if namedTagged, ok := named.(reference.NamedTagged); ok {
		o.name, o.tag = namedTagged.Name(), namedTagged.Tag()
	} else {
		err = fmt.Errorf("could not parse image %q with a tag", imageStr)
	}

	return
}

// OCIImageRef is a struct containing a names and tagged reference
// by which an OCI runtime can identify an image to retrieve.
type OCIImageRef struct {
	name string
	tag  string
}

var _ fmt.Stringer = OCIImageRef{}

// Ref parses the internal strings to a reference.NamedTagged
func (i OCIImageRef) Ref() reference.NamedTagged {
	r, _ := reference.ParseDockerRef(fmt.Sprintf("%s:%s", i.name, i.tag))
	return r.(reference.NamedTagged)
}

// String returns the familiar form of the reference, e.g. "weaveworks/ignite-ubuntu:latest"
func (i OCIImageRef) String() string {
	return reference.FamiliarString(i.Ref())
}

// Normalized returns the normalized reference, e.g. "docker.io/weaveworks/ignite-ubuntu:latest"
func (i OCIImageRef) Normalized() string {
	return i.Ref().String()
}

func (i OCIImageRef) IsUnset() bool {
	return len(i.name) == 0
}

func (i OCIImageRef) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

func (i *OCIImageRef) UnmarshalJSON(b []byte) (err error) {
	var s string
	if err = json.Unmarshal(b, &s); err != nil {
		return err
	}

	*i, err = NewOCIImageRef(s)
	return err
}

// ParseOCIContentID takes in a string to parse into an *OCIContentID
// If given a local Docker SHA like "sha256:3285f65b2651c68b5316e7a1fbabd30b5ae47914ac5791ac4bb9d59d029b924b",
// it will be parsed into the local format, encoded as "docker://<SHA>". Given a full repo digest, such as
// "weaveworks/ignite-ubuntu@sha256:3285f65b2651c68b5316e7a1fbabd30b5ae47914ac5791ac4bb9d59d029b924b", it will
// be parsed into the OCI registry format, encoded as "oci://<full path>@<SHA>".
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

func parseOCIString(s string) (*OCIContentID, error) {
	// Check if it's a local docker image.
	// Example: docker://sha256:fce289e99eb9bca977dae136fbe2a82b6b7d4c372474c9235adc1741675f587e
	if strings.HasPrefix(s, ociSchemeLocal) {
		return ParseOCIContentID(strings.TrimPrefix(s, ociSchemeLocal))
	}

	// For full repo digest with repo name, url parse the string and obtain the
	// url components.
	u, err := url.Parse(s)
	if err != nil {
		return nil, fmt.Errorf("failed to parse oci string %q, err: %v", s, err)
	}

	// Remove the "oci://" scheme by only caring about the host and path.
	return ParseOCIContentID(u.Host + u.Path)
}

// String returns the string representation for either format
func (o *OCIContentID) String() string {
	var sb strings.Builder
	if !o.Local() {
		sb.WriteString(o.repoName + "@")
	}

	sb.WriteString(o.digest)
	return sb.String()
}

// Scheme returns the string representation with the scheme prefix
func (o *OCIContentID) SchemeString() string {
	scheme := ociSchemeRegistry
	if o.Local() {
		scheme = ociSchemeLocal
	}

	return scheme + o.String()
}

// Local returns true if the image has no repoName, i.e. it's not available from a registry
func (o *OCIContentID) Local() bool {
	return len(o.repoName) == 0
}

// Digest gets the digest of the content ID
func (o *OCIContentID) Digest() digest.Digest {
	return digest.Digest(o.digest)
}

// RepoDigest returns a repo digest based on the OCIContentID if it is not local
func (o *OCIContentID) RepoDigest() (n reference.Named) {
	if !o.Local() {
		// Were parsing already validated data, ignore the error
		n, _ = reference.ParseDockerRef(o.String())
	}

	return
}

func (o *OCIContentID) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.SchemeString())
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
