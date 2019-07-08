package source

import (
	"io"

	api "github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"
)

// Source represents a source for VM images
type Source interface {
	// ID returns the ID of the source
	ID() string

	// Parse verifies the ImageSource, fills in any missing fields and prepares the reader
	Parse(src string) (*api.OCIImageSource, error)

	// Reader provides a tar stream reader
	Reader() (io.ReadCloser, error)

	// Cleanup cleans up any temporary assets after reading
	Cleanup() error
}
