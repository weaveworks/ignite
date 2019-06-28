package dm

import (
	"fmt"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/layer"
	"log"
	"os/exec"
	"path"

	"github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"

	"github.com/weaveworks/ignite/pkg/source"
	"github.com/weaveworks/ignite/pkg/util"
)

type Device struct {
	*v1alpha1.PoolDevice
	pool *Pool

	// These flags are for filesystem and snapshot creation
	mkfs   bool
	resize bool
}

var _ blockDevice = &Device{}

// Additional space to add to volumes to compensate for the ext4 partition
var extraSize = v1alpha1.NewSizeFromBytes(constants.POOL_VOLUME_EXTRA_SIZE)

func (p *Pool) CreateVolume(layer layer.Layer) (*Device, error) {
	// The pool needs to be active for this
	if err := p.activate(); err != nil {
		return nil, err
	}

	if volume, err := p.newDevice(func(id v1alpha1.DMID) (*Device, error) {
		if err := dmsetup("message", p.Path(), "0", fmt.Sprintf("create_thin %s", id)); err != nil {
			return nil, err
		}

		// Assign the volume ID to the layer
		*layer.ID() = id

		return &Device{
			PoolDevice: &v1alpha1.PoolDevice{
				Size:         layer.Size().Add(extraSize),
				Parent:       nil,
				MetadataPath: layer.MetadataPath(),
			},
			pool: p,
			mkfs: true, // This is a new volume, create a new filesystem for it on activation
		}, nil
	}); err != nil {
		return nil, err
	} else {
		return volume, volume.activate()
	}
}

func (d *Device) CreateSnapshot(layer layer.Layer) (*Device, error) {
	// The device needs to be active for this
	if err := d.activate(); err != nil {
		return nil, err
	}

	if snapshot, err := d.pool.newDevice(func(id v1alpha1.DMID) (*Device, error) {
		if err := dmsetup("suspend", d.Path()); err != nil {
			return nil, err
		}

		if err := dmsetup("message", d.pool.Path(), "0",
			fmt.Sprintf("create_snap %s %s", id, d.pool.getID(d))); err != nil {
			return nil, err
		}

		if err := dmsetup("resume", d.Path()); err != nil {
			return nil, err
		}

		// Assign the snapshot ID to the layer
		*layer.ID() = id

		// TODO: Prevent snapshots smaller than their parents?
		return &Device{
			PoolDevice: &v1alpha1.PoolDevice{
				Size:         layer.Size(),
				Parent:       d.pool.getID(d),
				MetadataPath: layer.MetadataPath(),
			},
			pool:   d.pool,
			resize: layer.Size() != d.Size, // Set the resize flag if the size differs from the parent
		}, nil
	}); err != nil {
		return nil, err
	} else {
		return snapshot, snapshot.activate()
	}
}

func (d *Device) activate() error {
	id := d.pool.getID(d)
	parent := d.pool.GetDevice(d.Parent)

	if parent == nil {
		// Activate the pool as the base device
		if err := d.pool.activate(); err != nil {
			return err
		}
	} else {
		// Check if all parents are active
		// TODO: Reference count this for deactivation
		if err := parent.activate(); err != nil {
			return err
		}
	}

	// Don't try to activate an already active device
	if d.active() {
		return nil
	}

	dmTable := fmt.Sprintf("0 %d thin %s %s",
		d.Size.Sectors(),
		d.pool.Path(),
		id,
	)

	log.Printf("Activate device: %s\n", id)
	if err := dmsetup("create", d.name(id), "--table", dmTable); err != nil {
		return err
	}

	if d.mkfs {
		log.Printf("Creating new filesystem on device: %s\n", id)
		if err := mkfs(d); err != nil {
			return err
		}
	} else if d.resize { // If the resize flag has been set, resize the filesystem to fill the device
		log.Printf("Resizing filesystem on device: %s\n", id)
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

func (d *Device) name(id v1alpha1.DMID) string {
	return util.NewPrefixer().Prefix(id.String())
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
