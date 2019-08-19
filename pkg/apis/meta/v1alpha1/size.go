package v1alpha1

import (
	"encoding/json"
	"fmt"

	"github.com/c2h5oh/datasize"
)

const sectorSize = 512

// Size specifies a common unit for data sizes
type Size struct {
	datasize.ByteSize
}

var _ fmt.Stringer = Size{}

var EmptySize = NewSizeFromBytes(0)

var _ json.Marshaler = &Size{}
var _ json.Unmarshaler = &Size{}

func NewSizeFromString(str string) (Size, error) {
	s := Size{}
	return s, s.UnmarshalText([]byte(str))
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

func (s Size) Sectors() uint64 {
	return s.Bytes() / sectorSize
}

// Override ByteSize's default string implementation which results in .HR() without spaces
func (s Size) String() string {
	return s.HR()
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
	b, _ := s.MarshalText()
	return json.Marshal(string(b))
}

func (s *Size) UnmarshalJSON(b []byte) error {
	var str string
	var err error

	if err = json.Unmarshal(b, &str); err != nil {
		return err
	}

	*s, err = NewSizeFromString(str)
	return err
}
