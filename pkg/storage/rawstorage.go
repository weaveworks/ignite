package storage

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"

	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/util"
)

// RawStorage is a Key-indexed low-level interface to
// store byte-encoded Objects (resources) in non-volatile
// memory.
type RawStorage interface {
	// Read returns a resource's content based on key
	Read(key Key) ([]byte, error)
	// Exists checks if the resource indicated by key exists
	Exists(key Key) bool
	// Write writes the given content to the resource indicated by key
	Write(key Key, content []byte) error
	// Delete deletes the resource indicated by key
	Delete(key Key) error
	// List returns all matching resource Keys based on the given KindKey
	List(key KindKey) ([]Key, error)
	// Checksum returns a string checksum for the resource indicated by key
	Checksum(key Key) (string, error)
	// Format returns the format of the contents of the resource indicated by key
	Format(key Key) Format
	// WatchDir returns the path for Watchers to watch changes in
	WatchDir() string
}

func NewDefaultRawStorage(dir string) RawStorage {
	return &DefaultRawStorage{
		dir: dir,
	}
}

type DefaultRawStorage struct {
	dir string
}

func (r *DefaultRawStorage) realPath(key AnyKey) string {
	var file string

	switch key.(type) {
	case KindKey:
	// KindKeys get no special treatment
	case Key:
		// Keys get the metadata filename added to the returned path
		file = constants.METADATA
	default:
		panic(fmt.Sprintf("invalid key type received: %T", key))
	}

	return path.Join(r.dir, key.String(), file)
}

func (r *DefaultRawStorage) Read(key Key) ([]byte, error) {
	return ioutil.ReadFile(r.realPath(key))
}

func (r *DefaultRawStorage) Exists(key Key) bool {
	return util.FileExists(r.realPath(key))
}

func (r *DefaultRawStorage) Write(key Key, content []byte) error {
	file := r.realPath(key)

	// Create the underlying directories if they do not exist already
	if !r.Exists(key) {
		if err := os.MkdirAll(path.Dir(file), constants.DATA_DIR_PERM); err != nil {
			return err
		}
	}

	return ioutil.WriteFile(file, content, 0644)
}

func (r *DefaultRawStorage) Delete(key Key) error {
	return os.RemoveAll(path.Dir(r.realPath(key)))
}

func (r *DefaultRawStorage) List(key KindKey) ([]Key, error) {
	entries, err := ioutil.ReadDir(r.realPath(key))
	if err != nil {
		return nil, err
	}

	result := make([]Key, 0, len(entries))
	for _, entry := range entries {
		result = append(result, NewKey(key.Kind, meta.UID(entry.Name())))
	}

	return result, nil
}

// This returns the modification time as a UnixNano string
// If the file doesn't exist, return blank
func (r *DefaultRawStorage) Checksum(key Key) (s string, err error) {
	var fi os.FileInfo

	if r.Exists(key) {
		if fi, err = os.Stat(r.realPath(key)); err == nil {
			s = strconv.FormatInt(fi.ModTime().UnixNano(), 10)
		}
	}

	return
}

func (r *DefaultRawStorage) Format(key Key) Format {
	return FormatJSON // The DefaultRawStorage always uses JSON
}

func (r *DefaultRawStorage) WatchDir() string {
	return r.dir
}
