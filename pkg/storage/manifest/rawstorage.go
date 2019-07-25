package manifest

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/weaveworks/ignite/pkg/storage"
	"github.com/weaveworks/ignite/pkg/util"
	"io/ioutil"
	"os"
	"strconv"
)

// TODO: Support automatic change detection via inotify on the MappedRawStorage level?
type MappedRawStorage interface {
	storage.RawStorage

	AddMapping(key storage.Key, path string)
	GetMapping(path string) (storage.Key, error)
	RemoveMapping(key storage.Key)
}

func NewManifestRawStorage(dir string) MappedRawStorage {
	return &ManifestRawStorage{
		dir:          dir,
		fileMappings: make(map[storage.Key]string),
	}
}

type ManifestRawStorage struct {
	dir          string
	fileMappings map[storage.Key]string
}

func (r *ManifestRawStorage) realPath(key storage.Key) (path string, err error) {
	path, ok := r.fileMappings[key]
	if !ok {
		err = fmt.Errorf("ManifestRawStorage: %q not tracked", key)
	}

	return
}

func (r *ManifestRawStorage) Read(key storage.Key) ([]byte, error) {
	file, err := r.realPath(key)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadFile(file)
}

func (r *ManifestRawStorage) Exists(key storage.Key) bool {
	file, err := r.realPath(key)
	if err != nil {
		return false
	}

	return util.FileExists(file)
}

func (r *ManifestRawStorage) Write(key storage.Key, content []byte) error {
	// ManifestRawStorage isn't going to generate files itself,
	// only write if the file is already known
	file, err := r.realPath(key)
	if err != nil {
		return nil
	}

	return ioutil.WriteFile(file, content, 0644)
}

func (r *ManifestRawStorage) Delete(key storage.Key) (err error) {
	file, err := r.realPath(key)
	if err != nil {
		return
	}

	err = os.Remove(file)
	if err == nil {
		r.RemoveMapping(key)
	}

	return
}

func (r *ManifestRawStorage) List(kind storage.KindKey) ([]storage.Key, error) {
	result := make([]storage.Key, 0)

	for key := range r.fileMappings {
		if key.KindKey == kind {
			result = append(result, key)
		}
	}

	return result, nil
}

// This returns the modification time as a UnixNano string
// If the file doesn't exist, return blank
func (r *ManifestRawStorage) Checksum(key storage.Key) (s string, err error) {
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

func (r *ManifestRawStorage) Dir() string {
	return r.dir
}

func (r *ManifestRawStorage) AddMapping(key storage.Key, path string) {
	log.Debugf("ManifestRawStorage: AddMapping: %q -> %q", key, path)
	r.fileMappings[key] = path
}

func (r *ManifestRawStorage) GetMapping(path string) (storage.Key, error) {
	for key, p := range r.fileMappings {
		if p == path {
			return key, nil
		}
	}

	return storage.Key{}, fmt.Errorf("no mapping found for path %q", path)
}

func (r *ManifestRawStorage) RemoveMapping(key storage.Key) {
	log.Debugf("ManifestRawStorage: RemoveMapping: %q", key)
	delete(r.fileMappings, key)
}
