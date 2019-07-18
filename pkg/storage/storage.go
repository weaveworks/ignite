package storage

import (
	"fmt"
	"path"

	"github.com/weaveworks/ignite/pkg/apis/ignite/scheme"
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/storage/serializer"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/yaml"
)

// Storage is an interface for persisting and retrieving API objects to/from a backend
// One Storage instance handles all different Kinds of Objects
type Storage interface {
	// Get populates the Object using the given pointer, based on the file content
	Get(obj meta.Object) error
	// Set saves the Object to disk. If the object does not exist, the
	// ObjectMeta.Created field is set automatically
	Set(obj meta.Object) error
	// GetByID returns a new Object for the resource at the specified kind/uid path, based on the file content
	GetByID(kind meta.Kind, uid meta.UID) (meta.Object, error)
	// Delete removes an object from the storage
	Delete(kind meta.Kind, uid meta.UID) error
	// List lists objects for the specific kind
	List(kind meta.Kind) ([]meta.Object, error)
	// ListMeta lists all objects' APIType representation. In other words,
	// only metadata about each object is unmarshalled (uid/name/kind/apiVersion).
	// This allows for faster runs (no need to unmarshal "the world"), and less
	// resource usage, when only metadata is unmarshalled into memory
	ListMeta(kind meta.Kind) ([]meta.Object, error)
	// Count returns the amount of available Objects of a specific kind
	// This is used by Caches to check if all objects are cached to perform a List
	Count(kind meta.Kind) (uint64, error)
}

// DefaultStorage is the default storage implementation
var DefaultStorage = NewCache(NewGenericStorage(NewDefaultRawStorage(constants.DATA_DIR), scheme.Serializer))

// NewGenericStorage constructs a new Storage
func NewGenericStorage(rawStorage RawStorage, serializer serializer.Serializer) Storage {
	return &GenericStorage{rawStorage, serializer}
}

// GenericStorage implements the Storage interface
type GenericStorage struct {
	raw        RawStorage
	serializer serializer.Serializer
}

var _ Storage = &GenericStorage{}

// Get populates the pointer to the Object given, based on the file content
func (s *GenericStorage) Get(obj meta.Object) error {
	storageKey, err := s.keyForObj(obj)
	if err != nil {
		return err
	}

	content, err := s.raw.Read(storageKey)
	if err != nil {
		return err
	}

	return s.serializer.DecodeInto(content, obj)
}

// GetByID returns a new Object for the resource at the specified kind/uid path, based on the file content
func (s *GenericStorage) GetByID(kind meta.Kind, uid meta.UID) (meta.Object, error) {
	storageKey := KeyForUID(kind, uid)
	content, err := s.raw.Read(storageKey)
	if err != nil {
		return nil, err
	}

	// Decode the bytes to the internal version of the object
	internalobj, err := s.serializer.Decode(content, true)
	if err != nil {
		return nil, err
	}

	obj, ok := internalobj.(meta.Object)
	if !ok {
		return nil, fmt.Errorf("cannot convert to ignite Object")
	}

	return obj, nil
}

// Set saves the Object to disk
func (s *GenericStorage) Set(obj meta.Object) error {
	storageKey, err := s.keyForObj(obj)
	if err != nil {
		return err
	}

	b, err := s.serializer.EncodeJSON(obj)
	if err != nil {
		return err
	}

	return s.raw.Write(storageKey, b)
}

// Delete removes an object from the storage
func (s *GenericStorage) Delete(kind meta.Kind, uid meta.UID) error {
	storageKey := KeyForUID(kind, uid)
	return s.raw.Delete(storageKey)
}

// List lists objects for the specific kind
func (s *GenericStorage) List(kind meta.Kind) ([]meta.Object, error) {
	result := []meta.Object{} // TODO: Fix these initializations
	err := s.walkKind(kind, func(content []byte) error {
		// Decode the bytes to the internal version of the object
		internalobj, err := s.serializer.Decode(content, true)
		if err != nil {
			return err
		}

		obj, ok := internalobj.(meta.Object)
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
func (s *GenericStorage) ListMeta(kind meta.Kind) ([]meta.Object, error) {
	result := []meta.Object{}
	err := s.walkKind(kind, func(content []byte) error {
		obj := meta.NewAPIType()
		// The yaml package supports both YAML and JSON
		if err := yaml.Unmarshal(content, obj); err != nil {
			return err
		}

		result = append(result, obj)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

// Count counts the objects for the specific kind
func (s *GenericStorage) Count(kind meta.Kind) (uint64, error) {
	entries, err := s.raw.List(KeyForKind(kind))
	return uint64(len(entries)), err
}

func (s *GenericStorage) walkKind(kind meta.Kind, fn func(content []byte) error) error {
	kindKey := KeyForKind(kind)
	entries, err := s.raw.List(kindKey)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		// Allow metadata.json to not exist, although the directory does exist
		if !s.raw.Exists(entry) {
			continue
		}

		content, err := s.raw.Read(entry)
		if err != nil {
			return err
		}

		if err := fn(content); err != nil {
			return err
		}
	}

	return nil
}

func (s *GenericStorage) keyForObj(obj meta.Object) (string, error) {
	gvk, err := s.gvkFromObj(obj)
	if err != nil {
		return "", err
	}

	return KeyForUID(meta.Kind(gvk.Kind), obj.GetUID()), nil
}

func (s *GenericStorage) gvkFromObj(obj runtime.Object) (*schema.GroupVersionKind, error) {
	gvks, unversioned, err := s.serializer.Scheme().ObjectKinds(obj.(runtime.Object))
	if err != nil {
		return nil, err
	}

	if unversioned {
		return nil, fmt.Errorf("unversioned")
	}

	if len(gvks) != 1 {
		return nil, fmt.Errorf("unexpected gvks %v", gvks)
	}

	return &gvks[0], nil
}

func KeyForUID(kind meta.Kind, uid meta.UID) string {
	return "/" + path.Join(kind.Lower(), uid.String())
}

func KeyForKind(kind meta.Kind) string {
	return "/" + kind.Lower()
}
