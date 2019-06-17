package dm

import (
	"fmt"
	"github.com/weaveworks/ignite/pkg/util"
	"path"
)

// blockSize specifies the data block size of the pool,
// it should be between 128 (64KB) and 2097152 (1GB).
// 128 is recommended if snapshotting a lot (like we do with layers).
type DMPool struct {
	Name        string
	Devices     []*dmDevice
	Blocks      uint64
	BlockSize   uint64
	metadataDev blockDevice
	dataDev     blockDevice
	Free        int
}

var _ blockDevice = &DMPool{}

func NewDMPool(name string, blocks, blockSize uint64, metadataDev, dataDev blockDevice) *DMPool {
	return &DMPool{
		Name:        name,
		Blocks:      blocks,
		BlockSize:   blockSize,
		metadataDev: metadataDev,
		dataDev:     dataDev,
	}
}

func (p *DMPool) Get(id int) (*dmDevice, error) {
	if id < 0 || id >= len(p.Devices) || p.Devices[id] == nil {
		return nil, fmt.Errorf("nonexistent device: %d", id)
	}

	return p.Devices[id], nil
}

func (p *DMPool) GetByName(name string) (*dmDevice, error) {
	for _, device := range p.Devices {
		if device.Name == name {
			return device, nil
		}
	}

	return nil, fmt.Errorf("nonexistent device: %s", name)
}

func (p *DMPool) activatePool() (int, error) {
	// Don't try to activate an already active pool
	if p.active() {
		return idPool, nil
	}

	// TODO: Eliminate this assert
	if err := p.metadataDev.(*loopDevice).Attach(false); err != nil {
		return idFail, err
	}

	if err := p.dataDev.(*loopDevice).Attach(false); err != nil {
		return idFail, err
	}

	dmTable := fmt.Sprintf("0 %d thin-pool %s %s %d 0",
		p.Blocks,
		p.metadataDev.Path(),
		p.dataDev.Path(),
		p.BlockSize,
	)

	if err := dmsetup("create", p.Name, "--table", dmTable); err != nil {
		return idFail, err
	}

	return idPool, nil
}

func (p *DMPool) Path() string {
	return path.Join("/dev/mapper", p.Name)
}

// If /dev/mapper/<name> exists the pool is active
func (p *DMPool) active() bool {
	return util.FileExists(p.Path())
}

// This returns a free ID in the pool
// TODO: Verify that this works
func (p *DMPool) newID() int {
	if p.Free < len(p.Devices) {
		returnID := p.Free
		for i := p.Free + 1; i <= len(p.Devices); i++ {
			if i == len(p.Devices) || p.Devices[i] == nil {
				p.Free = i
				break
			}
		}

		return returnID
	}

	p.Devices = append(p.Devices, nil)
	return len(p.Devices) - 1
}

func (p *DMPool) newDevice(genFunc func(int) (*dmDevice, error)) (int, error) {
	var err error
	id := p.newID()

	p.Devices[id], err = genFunc(id)
	if err != nil {
		err = fmt.Errorf("%v, removal: %v", err, p.remove(id))
	}

	return id, err
}

func (p *DMPool) remove(id int) error {
	if _, err := p.Get(id); err != nil {
		return err
	}

	p.Devices[id] = nil

	if p.Free > id {
		p.Free = id
	}

	return nil
}
