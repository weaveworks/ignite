package v1alpha1

import (
	"encoding/json"
	"fmt"

	"github.com/c2h5oh/datasize"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
)

const (
	SectorSize = 512
)

type ObjectMeta struct {
	Name string    `json:"name"`
	UID  types.UID `json:"uid"`
}

func (o *ObjectMeta) GetName() string {
	return o.Name
}

func (o *ObjectMeta) GetUID() types.UID {
	return o.UID
}

// All types implementing Object conform to this
// interface, it's mainly used for filtering
type Object interface {
	runtime.Object
	GetName() string
	GetUID() types.UID
}

// Size specifies a common unit for data sizes
type Size struct {
	datasize.ByteSize
}

var EmptySize = NewSizeFromBytes(0)

var _ json.Marshaler = &Size{}
var _ json.Unmarshaler = &Size{}

func NewSizeFromBytes(bytes uint64) Size {
	return Size{
		datasize.ByteSize(bytes),
	}
}

func NewSizeFromSectors(sectors uint64) Size {
	return Size{
		datasize.ByteSize(sectors * SectorSize),
	}
}

func (s *Size) Sectors() uint64 {
	return s.Bytes() / SectorSize
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
	return json.Marshal(s.Bytes())
}

func (s *Size) UnmarshalJSON(b []byte) error {
	var i uint64
	if err := json.Unmarshal(b, &i); err != nil {
		return err
	}

	*s = NewSizeFromBytes(i)
	return nil
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
