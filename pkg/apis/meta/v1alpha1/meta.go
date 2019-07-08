package v1alpha1

import (
	"bytes"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	sectorSize = 512
)

// APIType is a struct implementing Object, used for
// unmarshalling unknown objects into this intermediate type
// where .Name, .UID, .Kind and .APIVersion become easily available
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type APIType struct {
	TypeMeta   `json:",inline"`
	ObjectMeta `json:"metadata"`
}

var _ Object = &APIType{}

// APITypeList is a list of many pointers APIType objects
type APITypeList []*APIType

// TypeMeta is an alias for the k8s/apimachinery TypeMeta with some additional methods
type TypeMeta struct {
	metav1.TypeMeta
}

func (t *TypeMeta) GetKind() Kind {
	return Kind(t.Kind)
}

type Kind string

var _ fmt.Stringer = Kind("")

const (
	KindImage  Kind = "Image"
	KindKernel Kind = "Kernel"
	KindVM     Kind = "VM"
)

// Returns a lowercase string representation of the Kind
func (k Kind) String() string {
	b := []byte(k)

	// Ignore TLAs
	if len(b) > 3 {
		b[0] = bytes.ToLower(b[:1])[0]
	}

	return string(b)
}

// Returns a uppercase string representation of the Kind
func (k Kind) Upper() string {
	return string(k)
}

func (k Kind) Lower() string {
	return string(bytes.ToLower([]byte(k)))
}

// ObjectMeta have to be embedded into any serializable object.
// It provides the .GetName() and .GetUID() methods that help
// implement the Object interface
type ObjectMeta struct {
	Name    string `json:"name"`
	UID     UID    `json:"uid,omitempty"`
	Created *Time  `json:"created,omitempty"`
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
func (o *ObjectMeta) GetCreated() *Time {
	return o.Created
}

// SetCreated returns when the Object was created
func (o *ObjectMeta) SetCreated(t *Time) {
	o.Created = t
}

// Object extends k8s.io/apimachinery's runtime.Object with
// extra GetName() and GetUID() methods from ObjectMeta
type Object interface {
	runtime.Object

	GetKind() Kind

	GetName() string
	SetName(string)

	GetUID() UID
	SetUID(UID)

	GetCreated() *Time
	SetCreated(t *Time)
}
