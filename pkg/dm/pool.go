package dm

import (
	"fmt"
	"path"

	"github.com/weaveworks/ignite/pkg/format"
	"github.com/weaveworks/ignite/pkg/util"
)

// blockSize specifies the data block size of the pool,
// it should be between 128 (64KB) and 2097152 (1GB).
// 128 is recommended if snapshotting a lot (like we do with layers).
type Pool struct {
	name        string
	devices     []*Device
	blocks      format.DataSize
	blockSize   format.DataSize
	metadataDev blockDevice
	dataDev     blockDevice
	free        int
}

const (
	idPool = -1
)

func NewPool(name string, blocks, blockSize format.DataSize, metadataDev, dataDev blockDevice) *Pool {
	return &Pool{
		name:        name,
		blocks:      blocks,
		blockSize:   blockSize,
		metadataDev: metadataDev,
		dataDev:     dataDev,
	}
}

//func (p *Pool) Get(name string) (*Device, error) {
//	fmt.Printf("%#v\n", p.devices)
//
//	for _, device := range p.devices {
//		if device.name == name {
//			return device, nil
//		}
//	}
//
//	return nil, fmt.Errorf("device %q not found in pool", name)
//}

func (p *Pool) getID(device *Device) int {
	// If the querying for nil, return the pool's ID
	if device == nil {
		return idPool
	}

	for i, d := range p.devices {
		if d == device {
			return i
		}
	}

	// This should never happen, unless you try to get
	// the ID of a device residing in another pool
	panic(fmt.Sprintf("pool %q getID failed!", p.name))
}

func (p *Pool) getDevice(id int) *Device {
	// If querying for the pool's ID, return nil
	if id == idPool {
		return nil
	}

	if id < 0 || id >= len(p.devices) {
		// This should never happen, unless you try
		// to get a device residing in another pool
		panic(fmt.Sprintf("pool %q getDevice failed!", p.name))
	}

	return p.devices[id]
}

func (p *Pool) activate() error {
	// Don't try to activate an already active pool
	if p.active() {
		return nil
	}

	// Activate the backing devices
	if err := p.metadataDev.activate(); err != nil {
		return err
	}

	if err := p.dataDev.activate(); err != nil {
		return err
	}

	dmTable := fmt.Sprintf("0 %d thin-pool %s %s %d 0",
		p.blocks.Sectors(),
		p.metadataDev.Path(),
		p.dataDev.Path(),
		p.blockSize.Sectors(),
	)

	if err := dmsetup("create", p.name, "--table", dmTable); err != nil {
		return err
	}

	return nil
}

func (p *Pool) Path() string {
	return path.Join("/dev/mapper", p.name)
}

// If /dev/mapper/<name> exists the pool is active
func (p *Pool) active() bool {
	return util.FileExists(p.Path())
}

// This returns a free ID in the pool
// TODO: Verify that this works
func (p *Pool) newID() int {
	if p.free < len(p.devices) {
		returnID := p.free
		for i := p.free + 1; i <= len(p.devices); i++ {
			if i == len(p.devices) || p.devices[i] == nil {
				p.free = i
				break
			}
		}

		return returnID
	}

	p.devices = append(p.devices, nil)
	p.free = len(p.devices)
	return p.free - 1
}

func (p *Pool) newDevice(genFunc func(int) (*Device, error)) (*Device, error) {
	var err error

	id := p.newID()
	p.devices[id], err = genFunc(id)
	if err != nil {
		p.Remove(id)
	}

	return p.devices[id], err
}

func (p *Pool) Remove(id int) {
	if p.getDevice(id) != nil {
		p.devices[id] = nil

		if p.free > id {
			p.free = id
		}
	}
}
