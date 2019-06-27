package v1alpha1

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/c2h5oh/datasize"

	"k8s.io/apimachinery/pkg/types"
)

type ObjectMeta struct {
	Name string    `json:"name"`
	UID  types.UID `json:"uid"`
}

func (om *ObjectMeta) GetName() string {
	return om.Name
}

func (om *ObjectMeta) GetUID() types.UID {
	return om.UID
}

// TODO: We should maybe move this to a metav1 package once we're happy with this

// Size specifies a common unit for data sizes
type Size struct {
	datasize.ByteSize
}

const sectorSize = 512

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
		datasize.ByteSize(sectors * sectorSize),
	}
}

func (s *Size) Sectors() uint64 {
	return s.Bytes() / sectorSize
}

// Override ByteSize's default string implementation which results in something similar to HR()
func (s *Size) String() string {
	return strconv.FormatUint(s.Bytes(), 10)
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
