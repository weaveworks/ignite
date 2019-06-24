package format

import (
	"encoding/json"
	"strconv"

	"github.com/c2h5oh/datasize"
)

type Data struct {
	datasize.ByteSize
}

var _ json.Marshaler = &Data{}
var _ json.Unmarshaler = &Data{}

func DataFrom(bytes uint64) Data {
	return Data{
		datasize.ByteSize(bytes),
	}
}

func (d *Data) Sectors() uint64 {
	return d.Bytes() / 512
}

// Override ByteSize's default string implementation which results in something similar to HR()
func (d *Data) String() string {
	return strconv.FormatUint(d.Bytes(), 10)
}

// Add returns a copy, does not modify the receiver
func (d Data) Add(other Data) Data {
	d.ByteSize += other.ByteSize
	return d
}

func (d Data) Min(other Data) Data {
	if other.ByteSize < d.ByteSize {
		return other
	}

	return d
}

func (d Data) Max(other Data) Data {
	if other.ByteSize > d.ByteSize {
		return other
	}

	return d
}

func (d *Data) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Bytes())
}

func (d *Data) UnmarshalJSON(b []byte) error {
	var i uint64
	if err := json.Unmarshal(b, &i); err != nil {
		return err
	}

	*d = DataFrom(i)

	return nil
}
