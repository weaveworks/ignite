package v1alpha1

import (
	"encoding/json"
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

// Size specifies
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

func (d *Size) Sectors() uint64 {
	return d.Bytes() / sectorSize
}

// Override ByteSize's default string implementation which results in something similar to HR()
func (d *Size) String() string {
	return strconv.FormatUint(d.Bytes(), 10)
}

// Add returns a copy, does not modify the receiver
func (d Size) Add(other Size) Size {
	d.ByteSize += other.ByteSize
	return d
}

func (d Size) Min(other Size) Size {
	if other.ByteSize < d.ByteSize {
		return other
	}

	return d
}

func (d Size) Max(other Size) Size {
	if other.ByteSize > d.ByteSize {
		return other
	}

	return d
}

func (d *Size) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Bytes())
}

func (d *Size) UnmarshalJSON(b []byte) error {
	var i uint64
	if err := json.Unmarshal(b, &i); err != nil {
		return err
	}

	*d = NewSizeFromBytes(i)
	return nil
}
