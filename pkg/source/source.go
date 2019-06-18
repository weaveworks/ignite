package source

import (
	"github.com/weaveworks/ignite/pkg/dm"
	"io"
)

// Reader returns a io.ReadCloser to tar file data
type Source interface {
	Reader() (io.ReadCloser, error)
	SizeBytes() int64
	SizeSectors() dm.Sectors
	Cleanup() error
	ID() string
}
