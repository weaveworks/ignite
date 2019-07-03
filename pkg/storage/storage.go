package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	"github.com/weaveworks/ignite/pkg/apis/ignite/scheme"
	ignitemeta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/storage/serializer"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// Storage is an interface for persisting and retrieving API objects to/from a backend
type Storage interface {
	// Get populates the pointer to the Object given, based on the file content
	Get(obj ignitemeta.Object) error
	// GetByID returns a new Object for the resource at the specified kind/uid path, based on the file content
	GetByID(kind, uid string) (ignitemeta.Object, error)
	// Set saves the Object to disk. If the object does not exist, the
	// ObjectMeta.Created field is set automatically
	Set(obj ignitemeta.Object) error
	// Delete removes an object from the storage
	Delete(kind, uid string) error
	// List lists objects for the specific kind
	List(kind string) ([]ignitemeta.Object, error)
	// ListMeta lists all objects' APIType representation. In other words,
	// only metadata about each object is unmarshalled (uid/name/kind/apiVersion).
	// This allows for faster runs (no need to unmarshal "the world"), and less
	// resource usage, when only metadata is unmarshalled into memory
	ListMeta(kind string) (ignitemeta.APITypeList, error)
	// GetCache gets a new Cache implementation for the specified kind
	GetCache(kind string) (Cache, error)
}

// DefaultStorage is the default storage impl
var DefaultStorage = NewStorage(constants.DATA_DIR)

// NewStorage constructs a new Storage using the default implementation, for the specified dataDir
func NewStorage(dataDir string) Storage {
	return &storage{
		dataDir:    dataDir,
		serializer: scheme.Serializer,
	}
}

// storage implements the Storage interface
type storage struct {
	dataDir    string
	serializer serializer.Serializer
}

// Get populates the pointer to the Object given, based on the file content
func (s *storage) Get(obj ignitemeta.Object) error {
	storagePath, err := s.storagePathForObj(obj)
	if err != nil {
		return err
	}
	return s.serializer.DecodeFileInto(storagePath, obj)
}

// GetByID returns a new Object for the resource at the specified kind/uid path, based on the file content
func (s *storage) GetByID(kind, uid string) (ignitemeta.Object, error) {
	storagePath := s.storagePathForID(kind, uid)
	obj, err := s.serializer.DecodeFile(storagePath)
	if err != nil {
		return nil, err
	}
	igniteObj, ok := obj.(ignitemeta.Object)
	if !ok {
		return nil, fmt.Errorf("cannot convert ignite Object")
	}
	return igniteObj, nil
}

// Set saves the Object to disk. If the object does not exist, the
// ObjectMeta.Created field is set automatically
func (s *storage) Set(obj ignitemeta.Object) error {
	storagePath, err := s.storagePathForObj(obj)
	if err != nil {
		return err
	}
	// Make sure the parent directories exist
	if err := os.MkdirAll(path.Dir(storagePath), 0755); err != nil {
		return err
	}
	if _, err := os.Stat(storagePath); os.IsNotExist(err) {
		// Register that the object was created now
		obj.SetCreated(&metav1.Time{time.Now().UTC()})
	}

	b, err := s.serializer.EncodeJSON(obj)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(storagePath, b, 0644)
}

// Delete removes an object from the storage
func (s *storage) Delete(kind, uid string) error {
	storagePath := s.storagePathForID(kind, uid)
	// remove the whole directory, not only metadata.json
	storageDir := path.Dir(storagePath)
	return os.RemoveAll(storageDir)
}

// List lists objects for the specific kind
func (s *storage) List(kind string) ([]ignitemeta.Object, error) {
	result := []ignitemeta.Object{}
	err := s.walkDir(kind, func(content []byte) error {
		runtimeobj, err := s.serializer.Decode(content)
		if err != nil {
			return err
		}
		obj, ok := runtimeobj.(ignitemeta.Object)
		if !ok {
			return fmt.Errorf("can't convert to ignite object")
		}
		result = append(result, obj)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ListMeta lists all objects' APIType representation. In other words,
// only metadata about each object is unmarshalled (uid/name/kind/apiVersion).
// This allows for faster runs (no need to unmarshal "the world"), and less
// resource usage, when only metadata is unmarshalled into memory
func (s *storage) ListMeta(kind string) (ignitemeta.APITypeList, error) {
	result := ignitemeta.APITypeList{}
	err := s.walkDir(kind, func(content []byte) error {
		obj := &ignitemeta.APIType{}
		if err := json.Unmarshal(content, obj); err != nil {
			return err
		}
		result = append(result, obj)
		return nil
	})
	if err != nil {
		fmt.Println("listmeta error", err)
		return nil, err
	}
	return result, nil
}

// GetCache gets a new Cache implementation for the specified kind
func (s *storage) GetCache(kind string) (Cache, error) {
	list, err := s.ListMeta(kind)
	if err != nil {
		return nil, err
	}
	return NewCache(list), nil
}

func (s *storage) walkDir(kind string, fn func(content []byte) error) error {
	storageDir := path.Join(s.dataDir, strings.ToLower(kind))
	entries, err := ioutil.ReadDir(storageDir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		entryPath := path.Join(storageDir, entry.Name(), constants.METADATA)
		// Allow metadata.json to not exist, although the directory does exist
		// TODO: The ID handler should work well with this in the future
		if _, err := os.Stat(entryPath); os.IsNotExist(err) {
			continue
		} else if err != nil {
			return err
		}
		content, err := ioutil.ReadFile(entryPath)
		if err != nil {
			return err
		}
		if err := fn(content); err != nil {
			return err
		}
	}
	return nil
}

func (s *storage) storagePathForObj(obj ignitemeta.Object) (string, error) {
	gvk, err := s.gvkFromObj(obj)
	if err != nil {
		return "", err
	}
	return s.storagePathForID(gvk.Kind, string(obj.GetUID())), nil
}

func (s *storage) storagePathForID(kind, uid string) string {
	return path.Join(s.dataDir, strings.ToLower(kind), uid, constants.METADATA)
}

func (s *storage) gvkFromObj(obj runtime.Object) (*schema.GroupVersionKind, error) {
	gvks, unversioned, err := s.serializer.Scheme().ObjectKinds(obj.(runtime.Object))
	if err != nil {
		return nil, err
	}
	if unversioned {
		return nil, fmt.Errorf("unversioned")
	}
	if len(gvks) == 0 {
		return nil, fmt.Errorf("unexpected gvks")
	}
	return &gvks[0], nil
}
