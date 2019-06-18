package dm

import "encoding/json"

// DMPool serialization
type dmPoolSerial struct {
	Name      string
	Devices   []*dmDeviceSerial
	Blocks    Sectors
	BlockSize Sectors
	Free      int
}

// dmDevice serialization
type dmDeviceSerial struct {
	Name     string
	Blocks   Sectors
	ParentID int
}

func (d *dmDevice) encode() *dmDeviceSerial {
	return &dmDeviceSerial{
		Name:     d.name,
		Blocks:   d.blocks,
		ParentID: d.pool.getID(d.parent),
	}
}

func (dj *dmDeviceSerial) decode(d *dmDevice) {
	d.parent = d.pool.getDevice(dj.ParentID)
	d.name = dj.Name
	d.blocks = dj.Blocks
}

func (p *DMPool) encode() *dmPoolSerial {
	devices := make([]*dmDeviceSerial, len(p.devices))
	for i := range p.devices {
		if p.devices[i] != nil {
			devices[i] = p.devices[i].encode()
		}
	}

	return &dmPoolSerial{
		Name:      p.name,
		Devices:   devices,
		Blocks:    p.blocks,
		BlockSize: p.blockSize,
		Free:      p.free,
	}
}

func (ps *dmPoolSerial) decode(p *DMPool) {
	p.name = ps.Name
	p.devices = ps.decodeDevices(p)
	p.blocks = ps.Blocks
	p.blockSize = ps.BlockSize
	p.free = ps.computeFree()
	// TODO: Handle metadataDev and dataDev
}

func (ps *dmPoolSerial) decodeDevices(p *DMPool) []*dmDevice {
	devices := make([]*dmDevice, len(ps.Devices))

	// Generate the devices
	ps.iterateDevices(func(i int) {
		devices[i] = &dmDevice{pool: p}
	})

	// Decode the pool devices
	ps.iterateDevices(func(i int) {
		ps.Devices[i].decode(devices[i])
	})

	return devices
}

func (ps *dmPoolSerial) iterateDevices(iterateFunc func(int)) {
	for i := range ps.Devices {
		if ps.Devices[i] != nil {
			iterateFunc(i)
		}
	}
}

func (ps *dmPoolSerial) computeFree() int {
	for i, device := range ps.Devices {
		if device == nil {
			return i
		}
	}

	return len(ps.Devices)
}

var _ blockDevice = &DMPool{}
var _ json.Marshaler = &DMPool{}
var _ json.Unmarshaler = &DMPool{}

func (p *DMPool) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.encode())
}

// We use this custom unmarshaller to abstract the serializable pool and devices,
// associate the devices with the pool and to resolve the parents of each device
func (p *DMPool) UnmarshalJSON(b []byte) error {
	poolSerial := dmPoolSerial{}
	if err := json.Unmarshal(b, &poolSerial); err != nil {
		return err
	}

	poolSerial.decode(p)

	return nil
}
