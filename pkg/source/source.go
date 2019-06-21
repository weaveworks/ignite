package source

import (
	"io"

	"github.com/weaveworks/ignite/pkg/format"
)

// Reader returns a io.ReadCloser to tar file data
type Source interface {
	Reader() (io.ReadCloser, error)
	SizeBytes() int64
	SizeSectors() format.Sectors
	Cleanup() error
	ID() string
}
