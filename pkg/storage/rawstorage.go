package storage

import (
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/util"
)

type RawStorage interface {
	Read(key string) ([]byte, error)
	Exists(key string) bool
	Write(key string, content []byte) error
	Delete(key string) error
	List(directory string) ([]string, error)
}

func NewDefaultRawStorage(dir string) RawStorage {
	return &DefaultRawStorage{
		dir: dir,
	}
}

type DefaultRawStorage struct {
	dir string
}

func (r *DefaultRawStorage) realPath(key string) string {
	// The "/" prefix is enforced
	if !strings.HasPrefix(key, "/") {
		key = "/" + key
	}
	// If a top-level kind is described, and not a file, return the kind directory path
	if len(strings.Split(key, "/")) == 2 {
		return path.Join(r.dir, key)
	}
	// Return the file location, with the metadata.json suffix
	return path.Join(r.dir, key, constants.METADATA)
}

func (r *DefaultRawStorage) Read(key string) ([]byte, error) {
	file := r.realPath(key)
	return ioutil.ReadFile(file)
}

func (r *DefaultRawStorage) Exists(key string) bool {
	file := r.realPath(key)
	return util.FileExists(file)
}

func (r *DefaultRawStorage) Write(key string, content []byte) error {
	file := r.realPath(key)

	// Create the underlying directories if they do not exist already
	if !r.Exists(key) {
		if err := os.MkdirAll(path.Dir(file), 0755); err != nil {
			return err
		}
	}
	return ioutil.WriteFile(file, content, 0644)
}

func (r *DefaultRawStorage) Delete(key string) error {
	file := r.realPath(key)
	dir := path.Dir(file)
	return os.RemoveAll(dir)
}

func (r *DefaultRawStorage) List(parentKey string) ([]string, error) {
	realPath := r.realPath(parentKey)
	entries, err := ioutil.ReadDir(realPath)
	if err != nil {
		return nil, err
	}
	result := []string{}
	for _, entry := range entries {
		entryPath := path.Join(parentKey, entry.Name())
		result = append(result, entryPath)
	}
	return result, nil
}
