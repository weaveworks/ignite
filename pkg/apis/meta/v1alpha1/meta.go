package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	sectorSize = 512
)

// APIType is a struct implementing Object, used for
// unmarshalling unknown objects into this intermediate type
// where .Name, .UID, .Kind and .APIVersion become easily available
type APIType struct {
	metav1.TypeMeta `json:",inline"`
	ObjectMeta      `json:"metadata"`
}

// APITypeList is a list of many pointers APIType objects
type APITypeList []*APIType

// ObjectMeta have to be embedded into any serializable object.
// It provides the .GetName() and .GetUID() methods that help
// implement the Object interface
type ObjectMeta struct {
	Name    string       `json:"name"`
	UID     UID          `json:"uid,omitempty"`
	Created *metav1.Time `json:"created,omitempty"`
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
func (o *ObjectMeta) GetCreated() *metav1.Time {
	return o.Created
}

// SetCreated returns when the Object was created
func (o *ObjectMeta) SetCreated(t *metav1.Time) {
	o.Created = t
}

// Object extends k8s.io/apimachinery's runtime.Object with
// extra GetName() and GetUID() methods from ObjectMeta
type Object interface {
	runtime.Object

	GetName() string
	SetName(string)

	GetUID() UID
	SetUID(UID)

	GetCreated() *metav1.Time
	SetCreated(t *metav1.Time)
}
