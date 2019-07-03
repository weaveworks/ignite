package api

import (
	"encoding/json"
	"fmt"

	"github.com/c2h5oh/datasize"

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
func (o *ObjectMeta) GetUID() string {
	return o.UID.String()
}

// SetUID sets the UID of the Object
func (o *ObjectMeta) SetUID(uid string) {
	o.UID = UID(uid)
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

	// TODO: Use UID
	GetUID() string
	SetUID(string)

	GetCreated() *metav1.Time
	SetCreated(t *metav1.Time)
}

// UID represents an unique ID for a type
type UID string

// String returns the UID in string representation
func (u UID) String() string {
	return string(u)
}

// Size specifies a common unit for data sizes
type Size struct {
	datasize.ByteSize
}

var EmptySize = NewSizeFromBytes(0)

var _ json.Marshaler = &Size{}
var _ json.Unmarshaler = &Size{}

func NewSizeFromString(str string) (Size, error) {
	s := Size{}
	err := s.UnmarshalText([]byte(str))
	return s, err
}

func NewSizeFromBytes(bytes uint64) Size {
	return Size{
		datasize.ByteSize(bytes),
	}
}

func NewSizeFromSectors(sectors uint64) Size {
	return Size{
		datasize.ByteSize(sectors * sectorSize),
	}
}

func (s *Size) Sectors() uint64 {
	return s.Bytes() / sectorSize
}

// Override ByteSize's default string implementation which results in something similar to HR()
func (s *Size) String() string {
	b, _ := s.MarshalText()
	return string(b)
}

// Int64 returns the byte size as int64
func (s *Size) Int64() int64 {
	return int64(s.Bytes())
}

// Add returns a copy, does not modify the receiver
func (s Size) Add(other Size) Size {
	s.ByteSize += other.ByteSize
	return s
}

func (s Size) Min(other Size) Size {
	if other.ByteSize < s.ByteSize {
		return other
	}

	return s
}

func (s Size) Max(other Size) Size {
	if other.ByteSize > s.ByteSize {
		return other
	}

	return s
}

func (s *Size) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *Size) UnmarshalJSON(b []byte) error {
	var str string
	err := json.Unmarshal(b, &str)
	if err != nil {
		return err
	}

	*s, err = NewSizeFromString(str)
	return err
}

// DMID specifies the format for device mapper IDs
type DMID struct {
	index int32
}

var _ fmt.Stringer = DMID{}

func NewDMID(i int) DMID {
	// device mapper IDs are unsigned 24-bit integers
	if i < 0 || i >= 1<<24 {
		panic(fmt.Sprintf("device mapper ID out of range: %d", i))
	}

	return DMID{
		index: int32(i),
	}
}

func NewPoolDMID() DMID {
	// Internally we keep the pool ID out of range
	return DMID{
		index: -1,
	}
}

func (d *DMID) Pool() bool {
	return d.index < 0
}

func (d *DMID) Index() int {
	if !d.Pool() {
		return int(d.index)
	}

	panic("attempt to index nonexistent ID")
}

func (d DMID) String() string {
	if !d.Pool() {
		return fmt.Sprintf("%d", d.index)
	}

	return "pool"
}
