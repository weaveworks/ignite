package dm

import (
	"fmt"
	"github.com/weaveworks/ignite/pkg/util"
	"log"
	"path"
)

type dmDevice struct {
	Name     string
	Blocks   uint64
	ParentID int
}

const (
	idPool = -1
	idFail = -2
)

var _ blockDevice = &dmDevice{}

func (p *DMPool) CreateVolume(name string, blocks uint64) (int, error) {
	// The pool needs to be active for this
	// TODO: Make this a trigger
	if _, err := p.activatePool(); err != nil {
		return idFail, err
	}

	if id, err := p.newDevice(func(id int) (*dmDevice, error) {
		if err := dmsetup("message", p.Path(), "0", fmt.Sprintf("create_thin %d", id)); err != nil {
			return nil, err
		}

		return &dmDevice{
			Name:     name,
			Blocks:   blocks,
			ParentID: idPool,
		}, nil
	}); err != nil {
		return idFail, err
	} else {
		return p.activateDevice(id)
	}
}

func (p *DMPool) CreateSnapshot(parentID int, name string) (int, error) {
	parent, err := p.Get(parentID)
	if err != nil {
		return idFail, err
	}

	// The pool needs to be active for this
	if _, err := p.activatePool(); err != nil {
		return idFail, err
	}

	if id, err := p.newDevice(func(id int) (*dmDevice, error) {
		if err := dmsetup("suspend", parent.Path()); err != nil {
			return nil, err
		}

		if err := dmsetup("message", p.Path(), "0",
			fmt.Sprintf("create_snap %d %d", id, parentID)); err != nil {
			return nil, err
		}

		if err := dmsetup("resume", parent.Path()); err != nil {
			return nil, err
		}

		return &dmDevice{
			Name:     name,
			Blocks:   parent.Blocks,
			ParentID: parentID,
		}, nil
	}); err != nil {
		return idFail, err
	} else {
		return p.activateDevice(id)
	}
}

func (p *DMPool) activateDevice(id int) (int, error) {
	log.Printf("Activate device: %d\n", id)
	// Activate the pool as the base device
	if id < 0 {
		return p.activatePool()
	}

	device, err := p.Get(id)
	if err != nil {
		return idFail, err
	}

	// Check if all parents are active
	// TODO: Reference count this for deactivation
	if code, err := p.activateDevice(device.ParentID); err != nil {
		return code, err
	}

	// Don't try to activate an already active device
	if device.active() {
		return id, nil
	}

	dmTable := fmt.Sprintf("0 %d thin %s %d",
		device.Blocks,
		p.Path(),
		id,
	)

	//if d.externalDev != nil {
	//	dmTable = fmt.Sprintf("%s %s", dmTable, d.externalDev.Path())
	//}

	if err := dmsetup("create", device.Name, "--table", dmTable); err != nil {
		return idFail, err
	}

	return id, nil
}

func (d *dmDevice) Path() string {
	return path.Join("/dev/mapper", d.Name)
}

// If /dev/mapper/<name> exists the device is active
func (d *dmDevice) active() bool {
	return util.FileExists(d.Path())
}
