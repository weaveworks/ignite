package dm

import (
	"encoding/json"

	"github.com/weaveworks/ignite/pkg/format"
)

// Pool serialization
type poolSerial struct {
	Name      string
	Devices   []*deviceSerial
	Blocks    format.Sectors
	BlockSize format.Sectors
	Free      int
}

// Device serialization
type deviceSerial struct {
	Name     string
	Blocks   format.Sectors
	ParentID int
}

func (d *Device) encode() *deviceSerial {
	return &deviceSerial{
		Name:     d.name,
		Blocks:   d.blocks,
		ParentID: d.pool.getID(d.parent),
	}
}

func (dj *deviceSerial) decode(p *Pool) *Device {
	return &Device{
		pool:   p,
		name:   dj.Name,
		blocks: dj.Blocks,
	}
}

func (p *Pool) encode() *poolSerial {
	devices := make([]*deviceSerial, len(p.devices))
	for i := range p.devices {
		if p.devices[i] != nil {
			devices[i] = p.devices[i].encode()
		}
	}

	return &poolSerial{
		Name:      p.name,
		Devices:   devices,
		Blocks:    p.blocks,
		BlockSize: p.blockSize,
		Free:      p.free,
	}
}

func (ps *poolSerial) decode(p *Pool) {
	p.name = ps.Name
	p.devices = ps.decodeDevices(p)
	p.blocks = ps.Blocks
	p.blockSize = ps.BlockSize
	p.free = ps.computeFree()
	// TODO: Handle metadataDev and dataDev
}

func (ps *poolSerial) decodeDevices(p *Pool) []*Device {
	devices := make([]*Device, len(ps.Devices))

	// Decode the pool devices
	ps.iterateDevices(func(i int) {
		devices[i] = ps.Devices[i].decode(p)
	})

	// Associate device parents
	ps.iterateDevices(func(i int) {
		parentID := ps.Devices[i].ParentID
		if parentID != idPool {
			devices[i].parent = devices[parentID]
		}
	})

	return devices
}

func (ps *poolSerial) iterateDevices(iterateFunc func(int)) {
	for i := range ps.Devices {
		if ps.Devices[i] != nil {
			iterateFunc(i)
		}
	}
}

func (ps *poolSerial) computeFree() int {
	for i, device := range ps.Devices {
		if device == nil {
			return i
		}
	}

	return len(ps.Devices)
}

var _ blockDevice = &Pool{}
var _ json.Marshaler = &Pool{}
var _ json.Unmarshaler = &Pool{}

func (p *Pool) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.encode())
}

// We use this custom unmarshaller to abstract the serializable pool and devices,
// associate the devices with the pool and to resolve the parents of each device
func (p *Pool) UnmarshalJSON(b []byte) error {
	poolSerial := poolSerial{}
	if err := json.Unmarshal(b, &poolSerial); err != nil {
		return err
	}

	poolSerial.decode(p)

	return nil
}
