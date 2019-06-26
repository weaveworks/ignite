package source

import (
	"io"

	"github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"

	"github.com/weaveworks/ignite/pkg/format"
)

// Reader returns a io.ReadCloser to tar file data
type Source interface {
	Parse(*v1alpha1.ImageSource) error
	Reader() (io.ReadCloser, error)
	Size() format.DataSize
	Cleanup() error
	ID() string
}
