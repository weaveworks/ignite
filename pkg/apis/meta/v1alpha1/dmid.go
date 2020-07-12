package v1alpha1

import "fmt"

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

func (d DMID) Index() int {
	if !d.Pool() {
		return int(d.index)
	}

	panic("attempt to index pool ID")
}

func (d DMID) String() string {
	if !d.Pool() {
		return fmt.Sprintf("%d", d.index)
	}

	return "pool"
}
