package storage

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/weaveworks/gitops-toolkit/pkg/util"
)

// MappedRawStorage is an interface for RawStorages which store their
// data in a flat/unordered directory format like manifest directories.
type MappedRawStorage interface {
	RawStorage

	// AddMapping binds a Key's virtual path to a physical file path
	AddMapping(key Key, path string)
	// RemoveMapping removes the physical file
	// path mapping matching the given Key
	RemoveMapping(key Key)
}

func NewGenericMappedRawStorage(dir string) MappedRawStorage {
	return &GenericMappedRawStorage{
		dir:          dir,
		fileMappings: make(map[Key]string),
	}
}

// GenericMappedRawStorage is the default implementation of a MappedRawStorage,
// it stores files in the given directory via a path translation map.
type GenericMappedRawStorage struct {
	dir          string
	fileMappings map[Key]string
}

func (r *GenericMappedRawStorage) realPath(key Key) (path string, err error) {
	path, ok := r.fileMappings[key]
	if !ok {
		err = fmt.Errorf("GenericMappedRawStorage: %q not tracked", key)
	}

	return
}

func (r *GenericMappedRawStorage) Read(key Key) ([]byte, error) {
	file, err := r.realPath(key)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadFile(file)
}

func (r *GenericMappedRawStorage) Exists(key Key) bool {
	file, err := r.realPath(key)
	if err != nil {
		return false
	}

	return util.FileExists(file)
}

func (r *GenericMappedRawStorage) Write(key Key, content []byte) error {
	// GenericMappedRawStorage isn't going to generate files itself,
	// only write if the file is already known
	file, err := r.realPath(key)
	if err != nil {
		return nil
	}

	return ioutil.WriteFile(file, content, 0644)
}

func (r *GenericMappedRawStorage) Delete(key Key) (err error) {
	file, err := r.realPath(key)
	if err != nil {
		return
	}

	// GenericMappedRawStorage files can be deleted
	// externally, check that the file exists first
	if util.FileExists(file) {
		err = os.Remove(file)
	}

	if err == nil {
		r.RemoveMapping(key)
	}

	return
}

func (r *GenericMappedRawStorage) List(kind KindKey) ([]Key, error) {
	result := make([]Key, 0)

	for key := range r.fileMappings {
		if key.KindKey == kind {
			result = append(result, key)
		}
	}

	return result, nil
}

// This returns the modification time as a UnixNano string
// If the file doesn't exist, return blank
func (r *GenericMappedRawStorage) Checksum(key Key) (s string, err error) {
	file, err := r.realPath(key)
	if err != nil {
		return
	}

	var fi os.FileInfo
	if r.Exists(key) {
		if fi, err = os.Stat(file); err == nil {
			s = strconv.FormatInt(fi.ModTime().UnixNano(), 10)
		}
	}

	return
}

func (r *GenericMappedRawStorage) Format(key Key) (f Format) {
	if file, err := r.realPath(key); err == nil {
		f = Formats[filepath.Ext(file)] // Retrieve the correct format based on the extension
	}

	return
}

func (r *GenericMappedRawStorage) WatchDir() string {
	return r.dir
}

func (r *GenericMappedRawStorage) GetKey(path string) (Key, error) {
	for key, p := range r.fileMappings {
		if p == path {
			return key, nil
		}
	}

	return Key{}, fmt.Errorf("no mapping found for path %q", path)
}

func (r *GenericMappedRawStorage) AddMapping(key Key, path string) {
	log.Debugf("GenericMappedRawStorage: AddMapping: %q -> %q", key, path)
	r.fileMappings[key] = path
}

func (r *GenericMappedRawStorage) RemoveMapping(key Key) {
	log.Debugf("GenericMappedRawStorage: RemoveMapping: %q", key)
	delete(r.fileMappings, key)
}
