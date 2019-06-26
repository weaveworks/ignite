package dm

import (
	"fmt"
	"log"
	"os/exec"
	"path"

	"github.com/weaveworks/ignite/pkg/format"
	"github.com/weaveworks/ignite/pkg/source"
	"github.com/weaveworks/ignite/pkg/util"
)

type Device struct {
	pool     *Pool
	parent   *Device
	blocks   format.DataSize
	metadata Metadata

	// These flags are for filesystem and snapshot creation
	mkfs   bool
	resize bool
}

var _ blockDevice = &Device{}

func (p *Pool) CreateVolume(blocks format.DataSize, data Metadata) (*Device, error) {
	// The pool needs to be active for this
	if err := p.activate(); err != nil {
		return nil, err
	}

	if volume, err := p.newDevice(func(id int) (*Device, error) {
		if err := dmsetup("message", p.Path(), "0", fmt.Sprintf("create_thin %d", id)); err != nil {
			return nil, err
		}

		return &Device{
			pool:     p,
			parent:   nil,
			blocks:   blocks,
			metadata: data,
			mkfs:     true, // This is a new volume, create a new filesystem for it on activation
		}, nil
	}); err != nil {
		return nil, err
	} else {
		return volume, volume.activate()
	}
}

func (d *Device) CreateSnapshot(blocks format.DataSize, data Metadata) (*Device, error) {
	// The device needs to be active for this
	if err := d.activate(); err != nil {
		return nil, err
	}

	if snapshot, err := d.pool.newDevice(func(id int) (*Device, error) {
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

		return &Device{
			pool:     d.pool,
			parent:   d,
			blocks:   blocks,
			metadata: data,
			resize:   blocks != d.blocks, // Set the resize flag if the size differs from the parent
		}, nil
	}); err != nil {
		return nil, err
	} else {
		fmt.Printf("%#v\n", d.pool.devices)
		return snapshot, snapshot.activate()
	}
}

func (d *Device) activate() error {
	id := d.pool.getID(d)
	log.Printf("Activate device: %d\n", id)
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
		d.blocks.Sectors(),
		d.pool.Path(),
		id,
	)

	if err := dmsetup("create", d.name(id), "--table", dmTable); err != nil {
		return err
	}

	if d.mkfs {
		log.Printf("Creating new filesystem on device: %d\n", id)
		if err := mkfs(d); err != nil {
			return err
		}
	} else if d.resize { // If the resize flag has been set, resize the filesystem to fill the device
		log.Printf("Resizing filesystem on device: %d\n", id)
		if err := resize2fs(d); err != nil {
			return err
		}
	}

	return nil
}

func (d *Device) Import(src source.Source) (*util.MountPoint, error) {
	mountPoint, err := util.Mount(d.Path())
	if err != nil {
		return nil, err
	}

	tarCmd := exec.Command("tar", "-x", "-C", mountPoint.Path)
	reader, err := src.Reader()
	if err != nil {
		return nil, err
	}

	tarCmd.Stdin = reader
	if err := tarCmd.Start(); err != nil {
		return nil, err
	}

	if err := tarCmd.Wait(); err != nil {
		return nil, err
	}

	if err := src.Cleanup(); err != nil {
		return nil, err
	}

	return mountPoint, nil
}

func (d *Device) name(id int) string {
	return fmt.Sprintf("%s-%d", d.pool.name, id)
}

// TODO: Path should probably activate
func (d *Device) Path() string {
	return path.Join("/dev/mapper", d.name(d.pool.getID(d)))
}

// TODO: Temporary
func (d *Device) Start() (string, error) {
	return d.Path(), d.activate()
}

// If /dev/mapper/<name> exists the device is active
func (d *Device) active() bool {
	return util.FileExists(d.Path())
}

func (d *Device) Size() format.DataSize {
	return d.blocks
}
