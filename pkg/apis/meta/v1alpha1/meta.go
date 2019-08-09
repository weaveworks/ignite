package v1alpha1

import (
	"bytes"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	sectorSize = 512
)

// APIType is a struct implementing Object, used for
// unmarshalling unknown objects into this intermediate type
// where .Name, .UID, .Kind and .APIVersion become easily available
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type APIType struct {
	*TypeMeta   `json:",inline"`
	*ObjectMeta `json:"metadata"`
}

// This constructor ensures the APIType fields are not nil
func NewAPIType() *APIType {
	return &APIType{
		&TypeMeta{},
		&ObjectMeta{},
	}
}

// APITypeFrom is used to create a bound APIType from an Object
func APITypeFrom(obj Object) *APIType {
	return &APIType{
		obj.GetTypeMeta(),
		obj.GetObjectMeta(),
	}
}

var _ Object = &APIType{}

// APITypeList is a list of many pointers APIType objects
type APITypeList []*APIType

// TypeMeta is an alias for the k8s/apimachinery TypeMeta with some additional methods
type TypeMeta struct {
	metav1.TypeMeta
}

// This is a helper for APIType generation
func (t *TypeMeta) GetTypeMeta() *TypeMeta {
	return t
}

func (t *TypeMeta) GetKind() Kind {
	return Kind(t.Kind)
}

func (t *TypeMeta) GroupVersionKind() schema.GroupVersionKind {
	return t.TypeMeta.GetObjectKind().GroupVersionKind()
}

func (t *TypeMeta) SetGroupVersionKind(gvk schema.GroupVersionKind) {
	t.TypeMeta.GetObjectKind().SetGroupVersionKind(gvk)
}

type Kind string

var _ fmt.Stringer = Kind("")

// Returns a string representation of the Kind suitable for sentences
func (k Kind) String() string {
	b := []byte(k)

	// Ignore TLAs
	if len(b) > 3 {
		b[0] = bytes.ToLower(b[:1])[0]
	}

	return string(b)
}

// Returns a title case string representation of the Kind
func (k Kind) Title() string {
	return string(k)
}

// Returns a lowercase string representation of the Kind
func (k Kind) Lower() string {
	return string(bytes.ToLower([]byte(k)))
}

// Returns a Kind parsed from the given string
func ParseKind(input string) Kind {
	b := bytes.ToUpper([]byte(input))

	// Leave TLAs as uppercase
	if len(b) > 3 {
		b = append(b[:1], bytes.ToLower(b[1:])...)
	}

	return Kind(b)
}

// ObjectMeta have to be embedded into any serializable object.
// It provides the .GetName() and .GetUID() methods that help
// implement the Object interface
type ObjectMeta struct {
	Name        string            `json:"name"`
	UID         UID               `json:"uid,omitempty"`
	Created     Time              `json:"created"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

// This is a helper for APIType generation
func (o *ObjectMeta) GetObjectMeta() *ObjectMeta {
	return o
}

// GetName returns the name of the Object
func (o *ObjectMeta) GetName() string {
	return o.Name
}

// SetName sets the name of the Object
func (o *ObjectMeta) SetName(name string) {
	o.Name = name
}

// GetUID returns the UID of the Object
func (o *ObjectMeta) GetUID() UID {
	return o.UID
}

// SetUID sets the UID of the Object
func (o *ObjectMeta) SetUID(uid UID) {
	o.UID = uid
}

// GetCreated returns when the Object was created
func (o *ObjectMeta) GetCreated() Time {
	return o.Created
}

// SetCreated sets the creation time of the Object
func (o *ObjectMeta) SetCreated(t Time) {
	o.Created = t
}

// GetLabel returns the label value for the key
func (o *ObjectMeta) GetLabel(key string) string {
	if o.Labels == nil {
		return ""
	}
	return o.Labels[key]
}

// SetLabel sets a label value for a key
func (o *ObjectMeta) SetLabel(key, value string) {
	if o.Labels == nil {
		o.Labels = map[string]string{}
	}
	o.Labels[key] = value
}

// GetAnnotation returns the label value for the key
func (o *ObjectMeta) GetAnnotation(key string) string {
	if o.Annotations == nil {
		return ""
	}
	return o.Annotations[key]
}

// SetAnnotation sets a label value for a key
func (o *ObjectMeta) SetAnnotation(key, value string) {
	if o.Annotations == nil {
		o.Annotations = map[string]string{}
	}
	o.Annotations[key] = value
}

// Object extends k8s.io/apimachinery's runtime.Object with
// extra GetName() and GetUID() methods from ObjectMeta
type Object interface {
	runtime.Object

	GetTypeMeta() *TypeMeta
	GetObjectMeta() *ObjectMeta

	GetKind() Kind
	GroupVersionKind() schema.GroupVersionKind
	SetGroupVersionKind(schema.GroupVersionKind)

	GetName() string
	SetName(string)

	GetUID() UID
	SetUID(UID)

	GetCreated() Time
	SetCreated(t Time)

	GetLabel(key string) string
	SetLabel(key, value string)

	GetAnnotation(key string) string
	SetAnnotation(key, value string)
}
