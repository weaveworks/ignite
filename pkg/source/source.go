package source

import (
	"io"

	"github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"
)

// Reader returns a io.ReadCloser to tar file data
type Source interface {
	// Parse verifies the ImageSource, fills in any missing fields and prepares the reader
	Parse(*v1alpha1.ImageSource) error

	// Reader provides a tar stream reader
	Reader() (io.ReadCloser, error)

	// Cleanup cleans up any temporary assets after reading
	Cleanup() error
}
