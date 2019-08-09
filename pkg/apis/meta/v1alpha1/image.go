package v1alpha1

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/containers/image/docker/reference"
	"github.com/opencontainers/go-digest"
)

const (
	ociSchemeRegistry = "oci://"
	ociSchemeLocal    = "docker://"
)

// NewOCIImageRef parses and normalizes a reference to an OCI (docker) image.
func NewOCIImageRef(imageStr string) (OCIImageRef, error) {
	named, err := reference.ParseDockerRef(imageStr)
	if err != nil {
		return "", err
	}

	namedTagged, ok := named.(reference.NamedTagged)
	if !ok {
		return "", fmt.Errorf("could not parse image %q with a tag", imageStr)
	}

	return OCIImageRef(reference.FamiliarString(namedTagged)), nil
}

// OCIImageRef is a string by which an OCI runtime can identify an image to retrieve.
// It needs to have a tag and usually looks like "weaveworks/ignite-ubuntu:latest".
type OCIImageRef string

var _ fmt.Stringer = OCIImageRef("")

func (i OCIImageRef) String() string {
	return string(i)
}

func (i OCIImageRef) IsUnset() bool {
	return len(i) == 0
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

func (o *OCIContentID) String() string {
	scheme := ociSchemeRegistry
	if o.Local() {
		scheme = ociSchemeLocal
	}

	return scheme + o.ociString()
}

func parseOCIString(s string) (*OCIContentID, error) {
	u, err := url.Parse(s)
	if err != nil {
		return nil, err
	}

	// Remove the "docker://" or "oci://" scheme by only caring about the host and path
	return ParseOCIContentID(u.Host + u.Path)
}

// ociString returns the internal string representation for either format
func (o *OCIContentID) ociString() string {
	var sb strings.Builder
	if !o.Local() {
		sb.WriteString(o.repoName + "@")
	}

	sb.WriteString(o.digest)
	return sb.String()
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
		n, _ = reference.ParseDockerRef(o.ociString())
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
