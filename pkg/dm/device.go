package dm

import (
	"fmt"
	"github.com/weaveworks/ignite/pkg/util"
	"log"
	"path"
)

type dmDevice struct {
	pool   *DMPool
	parent *dmDevice
	name   string
	blocks Sectors
}

var _ blockDevice = &dmDevice{}

func (p *DMPool) CreateVolume(name string, blocks Sectors) (*dmDevice, error) {
	// The pool needs to be active for this
	if err := p.activate(); err != nil {
		return nil, err
	}

	if volume, err := p.newDevice(func(id int) (*dmDevice, error) {
		if err := dmsetup("message", p.Path(), "0", fmt.Sprintf("create_thin %d", id)); err != nil {
			return nil, err
		}

		return &dmDevice{
			pool:   p,
			parent: nil,
			name:   name,
			blocks: blocks,
		}, nil
	}); err != nil {
		return nil, err
	} else {
		return volume, volume.activate()
	}
}

func (d *dmDevice) CreateSnapshot(name string) (*dmDevice, error) {
	// The device needs to be active for this
	if err := d.activate(); err != nil {
		return nil, err
	}

	if snapshot, err := d.pool.newDevice(func(id int) (*dmDevice, error) {
		if err := dmsetup("suspend", d.Path()); err != nil {
			return nil, err
		}

		if err := dmsetup("message", d.pool.Path(), "0",
			fmt.Sprintf("create_snap %d %d", id, d.pool.getID(d))); err != nil {
			return nil, err
		}

		if err := dmsetup("resume", d.Path()); err != nil {
			return nil, err
		}

		return &dmDevice{
			pool:   d.pool,
			parent: d,
			name:   name,
			blocks: d.blocks,
		}, nil
	}); err != nil {
		return nil, err
	} else {
		return snapshot, snapshot.activate()
	}
}

func (d *dmDevice) activate() error {
	log.Printf("Activate device: %s\n", d.name)
	if d.parent == nil {
		// Activate the pool as the base device
		if err := d.pool.activate(); err != nil {
			return err
		}
	} else {
		// Check if all parents are active
		// TODO: Reference count this for deactivation
		if err := d.parent.activate(); err != nil {
			return err
		}
	}

	// Don't try to activate an already active device
	if d.active() {
		return nil
	}

	dmTable := fmt.Sprintf("0 %d thin %s %d",
		d.blocks,
		d.pool.Path(),
		d.pool.getID(d),
	)

	if err := dmsetup("create", d.name, "--table", dmTable); err != nil {
		return err
	}

	return nil
}

func (d *dmDevice) Path() string {
	return path.Join("/dev/mapper", d.name)
}

// If /dev/mapper/<name> exists the device is active
func (d *dmDevice) active() bool {
	return util.FileExists(d.Path())
}
