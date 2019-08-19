package storage

import (
	"fmt"

	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/serializer"
	patchutil "github.com/weaveworks/ignite/pkg/util/patch"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/yaml"
)

// Storage is an interface for persisting and retrieving API objects to/from a backend
// One Storage instance handles all different Kinds of Objects
type Storage interface {
	// New creates a new Object for the specified kind
	New(gvk schema.GroupVersionKind) (meta.Object, error)
	// Get returns a new Object for the resource at the specified kind/uid path, based on the file content
	Get(gvk schema.GroupVersionKind, uid meta.UID) (meta.Object, error)
	// GetMeta returns a new Object's APIType representation for the resource at the specified kind/uid path
	GetMeta(gvk schema.GroupVersionKind, uid meta.UID) (meta.Object, error)
	// Set saves the Object to disk. If the Object does not exist, the
	// ObjectMeta.Created field is set automatically
	Set(gvk schema.GroupVersionKind, obj meta.Object) error
	// Patch performs a strategic merge patch on the Object with the given UID, using the byte-encoded patch given
	Patch(gvk schema.GroupVersionKind, uid meta.UID, patch []byte) error
	// Delete removes an Object from the storage
	Delete(gvk schema.GroupVersionKind, uid meta.UID) error
	// List lists Objects for the specific kind
	List(gvk schema.GroupVersionKind) ([]meta.Object, error)
	// ListMeta lists all Objects' APIType representation. In other words,
	// only metadata about each Object is unmarshalled (uid/name/kind/apiVersion).
	// This allows for faster runs (no need to unmarshal "the world"), and less
	// resource usage, when only metadata is unmarshalled into memory
	ListMeta(gvk schema.GroupVersionKind) ([]meta.Object, error)
	// Count returns the amount of available Objects of a specific kind
	// This is used by Caches to check if all Objects are cached to perform a List
	Count(gvk schema.GroupVersionKind) (uint64, error)
	// Checksum returns a string representing the state of an Object on disk
	// The checksum should change if any modifications have been made to the
	// Object on disk, it can be e.g. the Object's modification timestamp or
	// calculated checksum
	Checksum(gvk schema.GroupVersionKind, uid meta.UID) (string, error)
	// RawStorage returns the RawStorage instance backing this Storage
	RawStorage() RawStorage
	// Serializer returns the serializer
	Serializer() serializer.Serializer
	// Close closes all underlying resources (e.g. goroutines) used; before the application exits
	Close() error
}

// NewGenericStorage constructs a new Storage
func NewGenericStorage(rawStorage RawStorage, serializer serializer.Serializer) Storage {
	return &GenericStorage{rawStorage, serializer, patchutil.NewPatcher(serializer)}
}

// GenericStorage implements the Storage interface
type GenericStorage struct {
	raw        RawStorage
	serializer serializer.Serializer
	patcher    patchutil.Patcher
}

var _ Storage = &GenericStorage{}

func (s *GenericStorage) Serializer() serializer.Serializer {
	return s.serializer
}

// New creates a new Object for the specified kind
// TODO: Create better error handling if the GVK specified is not recognized
func (s *GenericStorage) New(gvk schema.GroupVersionKind) (meta.Object, error) {
	obj, err := s.serializer.Scheme().New(gvk)
	if err != nil {
		return nil, err
	}

	// Default either through the scheme, or the high-level serializer Object
	if gvk.Version == runtime.APIVersionInternal {
		if err := s.serializer.DefaultInternal(obj); err != nil {
			return nil, err
		}
	} else {
		s.serializer.Scheme().Default(obj)
	}

	// Cast to meta.Object, and make sure it works
	metaObj, ok := obj.(meta.Object)
	if !ok {
		return nil, fmt.Errorf("can't convert to ignite object")
	}
	// Set the desired gvk from the caller of this Object
	// In practice, this means, although we created an internal type,
	// from defaulting external TypeMeta information was set. Set the
	// desired gvk here so it's correctly handled in all code that gets
	// the gvk from the Object
	metaObj.SetGroupVersionKind(gvk)
	return metaObj, nil
}

// Get returns a new Object for the resource at the specified kind/uid path, based on the file content
func (s *GenericStorage) Get(gvk schema.GroupVersionKind, uid meta.UID) (meta.Object, error) {
	storageKey := KeyForUID(gvk, uid)
	content, err := s.raw.Read(storageKey)
	if err != nil {
		return nil, err
	}

	return s.decode(content, gvk)
}

// TODO: Verify this works
// GetMeta returns a new Object's APIType representation for the resource at the specified kind/uid path
func (s *GenericStorage) GetMeta(gvk schema.GroupVersionKind, uid meta.UID) (meta.Object, error) {
	storageKey := KeyForUID(gvk, uid)
	content, err := s.raw.Read(storageKey)
	if err != nil {
		return nil, err
	}

	return s.decodeMeta(content, gvk)
}

// Set saves the Object to disk
func (s *GenericStorage) Set(gvk schema.GroupVersionKind, obj meta.Object) error {
	storageKey := KeyForUID(gvk, obj.GetUID())

	// Set the serializer based on the format given by the RawStorage
	serializeFunc := s.serializer.EncodeJSON
	if s.raw.Format(storageKey) != FormatJSON {
		serializeFunc = s.serializer.EncodeYAML
	}

	b, err := serializeFunc(obj)
	if err != nil {
		return err
	}

	return s.raw.Write(storageKey, b)
}

// Patch performs a strategic merge patch on the object with the given UID, using the byte-encoded patch given
func (s *GenericStorage) Patch(gvk schema.GroupVersionKind, uid meta.UID, patch []byte) error {
	storageKey := KeyForUID(gvk, uid)
	oldContent, err := s.raw.Read(storageKey)
	if err != nil {
		return err
	}

	newContent, err := s.patcher.Apply(oldContent, patch, gvk)
	if err != nil {
		return err
	}

	return s.raw.Write(storageKey, newContent)
}

// Delete removes an Object from the storage
func (s *GenericStorage) Delete(gvk schema.GroupVersionKind, uid meta.UID) error {
	storageKey := KeyForUID(gvk, uid)
	return s.raw.Delete(storageKey)
}

// List lists Objects for the specific kind
func (s *GenericStorage) List(gvk schema.GroupVersionKind) (result []meta.Object, walkerr error) {
	walkerr = s.walkKind(gvk, func(content []byte) error {
		obj, err := s.decode(content, gvk)
		if err != nil {
			return err
		}

		result = append(result, obj)
		return nil
	})
	return
}

// ListMeta lists all Objects' APIType representation. In other words,
// only metadata about each Object is unmarshalled (uid/name/kind/apiVersion).
// This allows for faster runs (no need to unmarshal "the world"), and less
// resource usage, when only metadata is unmarshalled into memory
func (s *GenericStorage) ListMeta(gvk schema.GroupVersionKind) (result []meta.Object, walkerr error) {
	walkerr = s.walkKind(gvk, func(content []byte) error {
		obj := meta.NewAPIType()
		// The yaml package supports both YAML and JSON
		if err := yaml.Unmarshal(content, obj); err != nil {
			return err
		}
		// Set the desired gvk from the caller of this Object
		// In practice, this means, although we got an external type,
		// we might want internal Objects later in the client. Hence,
		// set the right expectation here
		obj.SetGroupVersionKind(gvk)

		result = append(result, obj)
		return nil
	})
	return
}

// Count counts the Objects for the specific kind
func (s *GenericStorage) Count(gvk schema.GroupVersionKind) (uint64, error) {
	entries, err := s.raw.List(KeyForKind(gvk))
	return uint64(len(entries)), err
}

// Checksum returns a string representing the state of an Object on disk
func (s *GenericStorage) Checksum(gvk schema.GroupVersionKind, uid meta.UID) (string, error) {
	return s.raw.Checksum(KeyForUID(gvk, uid))
}

// RawStorage returns the RawStorage instance backing this Storage
func (s *GenericStorage) RawStorage() RawStorage {
	return s.raw
}

// Close closes all underlying resources (e.g. goroutines) used; before the application exits
func (s *GenericStorage) Close() error {
	return nil // nothing to do here for GenericStorage
}

func (s *GenericStorage) decode(content []byte, gvk schema.GroupVersionKind) (meta.Object, error) {
	// Decode the bytes to the internal version of the Object, if desired
	isInternal := gvk.Version == runtime.APIVersionInternal

	// Decode the bytes into an Object
	obj, err := s.serializer.Decode(content, isInternal)
	if err != nil {
		return nil, err
	}

	// Cast to meta.Object, and make sure it works
	metaObj, ok := obj.(meta.Object)
	if !ok {
		return nil, fmt.Errorf("can't convert to ignite object")
	}

	// Set the desired gvk of this Object from the caller
	metaObj.SetGroupVersionKind(gvk)
	return metaObj, nil
}

func (s *GenericStorage) decodeMeta(content []byte, gvk schema.GroupVersionKind) (meta.Object, error) {
	// Create a new APType object
	obj := meta.NewAPIType()

	// Decode the bytes into the APIType object
	if err := s.serializer.DecodeInto(content, obj); err != nil {
		return nil, err
	}

	// Set the desired gvk of this APIType object from the caller
	obj.SetGroupVersionKind(gvk)
	return obj, nil
}

func (s *GenericStorage) walkKind(gvk schema.GroupVersionKind, fn func(content []byte) error) error {
	kindKey := KeyForKind(gvk)
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

func KeyForUID(gvk schema.GroupVersionKind, uid meta.UID) Key {
	return NewKey(meta.Kind(gvk.Kind), uid)
}

func KeyForKind(gvk schema.GroupVersionKind) KindKey {
	return NewKindKey(meta.Kind(gvk.Kind))
}
