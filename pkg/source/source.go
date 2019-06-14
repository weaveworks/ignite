package source

import "io"

// Reader returns a io.ReadCloser to tar file data
type Source interface {
	Reader() (io.ReadCloser, error)
	Size() int64
	Cleanup() error
}
