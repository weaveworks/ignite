package source

import (
	"io"

	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
)

// Source represents a source for VM images
type Source interface {
	// Ref returns the reference of the source
	Ref() meta.OCIImageRef

	// SetRef allows setting reference of the source.
	SetRef(meta.OCIImageRef)

	// Parse verifies the ImageSource, fills in any missing fields and prepares the reader
	Parse(src meta.OCIImageRef) (*api.OCIImageSource, error)

	// Import  an image from a file source and prepares the reader.
	Import(path string) (map[meta.OCIImageRef]*api.OCIImageSource, error)

	// Reader provides a tar stream reader
	Reader() (io.ReadCloser, error)

	// Cleanup cleans up any temporary assets after reading
	Cleanup() error
}
