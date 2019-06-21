package format

import "fmt"

// Sectors is a data size unit for device mapper
// 1 sector = 512 bytes, Sectors basically divides via 512
// TODO: This needs to be a struct embedding an uint64?
// TODO: Custom datatype that extends datasize.ByteSize with Sector capabilities
type Sectors uint64

var _ fmt.Stringer = Sectors(0)

func SectorsFromBytes(bytes interface{}) Sectors {
	switch bytes.(type) {
	case uint64:
		return Sectors(bytes.(uint64) / 512)
	case int64:
		return Sectors(bytes.(int64) / 512)
	case int:
		return Sectors(bytes.(int) / 512)
	}

	panic(fmt.Sprintf("invalid Sectors type: %T", bytes))
}

func (s Sectors) ToBytes() uint64 {
	return uint64(s * 512)
}

func (s Sectors) String() string {
	return fmt.Sprintf("%d", s)
}
