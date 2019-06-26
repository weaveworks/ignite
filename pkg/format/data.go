package format

import (
	"encoding/json"
	"strconv"

	"github.com/c2h5oh/datasize"
)

type DataSize struct {
	datasize.ByteSize
}

var _ json.Marshaler = &DataSize{}
var _ json.Unmarshaler = &DataSize{}

func DataFrom(bytes uint64) DataSize {
	return DataSize{
		datasize.ByteSize(bytes),
	}
}

func (d *DataSize) Sectors() uint64 {
	return d.Bytes() / 512
}

// Override ByteSize's default string implementation which results in something similar to HR()
func (d *DataSize) String() string {
	return strconv.FormatUint(d.Bytes(), 10)
}

// Add returns a copy, does not modify the receiver
func (d DataSize) Add(other DataSize) DataSize {
	d.ByteSize += other.ByteSize
	return d
}

func (d DataSize) Min(other DataSize) DataSize {
	if other.ByteSize < d.ByteSize {
		return other
	}

	return d
}

func (d DataSize) Max(other DataSize) DataSize {
	if other.ByteSize > d.ByteSize {
		return other
	}

	return d
}

func (d *DataSize) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Bytes())
}

func (d *DataSize) UnmarshalJSON(b []byte) error {
	var i uint64
	if err := json.Unmarshal(b, &i); err != nil {
		return err
	}

	*d = DataFrom(i)

	return nil
}
