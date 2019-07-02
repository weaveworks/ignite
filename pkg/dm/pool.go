package dm

import (
	"fmt"
	"path"

	"github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"
	ignitemeta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/util"
)

var (
	poolName = util.NewPrefixer().Prefix("pool")
)

type Pool struct {
	v1alpha1.Pool

	free         ignitemeta.DMID
	freeComputed bool
}

// NewPool creates a new Pool object
// The Pool is stateless, so this can also be used to load a configuration
func NewPool(poolMeta *v1alpha1.Pool) *Pool {
	return &Pool{
		Pool: *poolMeta,
	}
}

func (p *Pool) getID(device *Device) ignitemeta.DMID {
	// If the querying for nil, return the pool's ID
	if device == nil {
		return ignitemeta.NewPoolDMID()
	}

	for i, d := range p.Status.Devices {
		if d == device.PoolDevice {
			return ignitemeta.NewDMID(i)
		}
	}

	// This should never happen, unless you try to get
	// the ID of a device residing in another pool
	panic("pool getID: device not found")
}

// GetDevice dynamically spawns a device from a v1alpha1.PoolDevice
func (p *Pool) GetDevice(id ignitemeta.DMID) *Device {
	// If querying for the pool's ID, return nil
	if id.Pool() {
		return nil
	}

	if id.Index() >= len(p.Status.Devices) {
		// This should never happen, unless you try
		// to get a device residing in another pool
		panic("pool GetDevice: index out of range")
	}

	spec := p.Status.Devices[id.Index()]
	if spec == nil {
		panic("pool GetDevice: nonexistent device")
	}

	return &Device{
		PoolDevice: spec,
		pool:       p,
	}
}

// This is a custom iterator to iterate over existing devices only (it skips nil slots)
func (p *Pool) ForDevices(iterFunc func(ignitemeta.DMID, *Device) error) error {
	for i := 0; i < len(p.Status.Devices); i++ {
		spec := p.Status.Devices[i]
		if spec != nil {
			id := ignitemeta.NewDMID(i)
			if err := iterFunc(id, p.GetDevice(id)); err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *Pool) allocate() error {
	// Allocate the backing files (if not allocated already)
	if err := allocateBackingFile(p.Spec.MetadataPath, p.Spec.MetadataSize); err != nil {
		return fmt.Errorf("failed to allocate metadata backing file: %v", err)
	}

	if err := allocateBackingFile(p.Spec.DataPath, p.Spec.DataSize); err != nil {
		return fmt.Errorf("failed to allocate data backing file: %v", err)
	}
}

func (p *Pool) activate() error {
	// Don't try to activate an already active pool
	if p.active() {
		return nil
	}

	// Trigger allocation
	if err := p.allocate(); err != nil {
		return err
	}

	// Activate the backing devices
	metadataDev, err := activateBackingDevice(p.Spec.MetadataPath, false)
	if err != nil {
		return err
	}

	dataDev, err := activateBackingDevice(p.Spec.DataPath, false)
	if err != nil {
		return err
	}

	dmTable := fmt.Sprintf("0 %d thin-pool %s %s %d 0",
		p.Spec.DataSize.Sectors(),
		metadataDev.Path(),
		dataDev.Path(),
		p.Spec.AllocationSize.Sectors(),
	)

	if err := dmsetup("create", poolName, "--table", dmTable); err != nil {
		return err
	}

	return nil
}

func (p *Pool) Path() string {
	return path.Join("/dev/mapper", poolName)
}

func (p *Pool) Size() int {
	var size int

	for i := 0; i < len(p.Status.Devices); i++ {
		if p.Status.Devices[i] != nil {
			size++
		}
	}

	return size
}

// If /dev/mapper/<name> exists the pool is active
func (p *Pool) active() bool {
	return util.FileExists(p.Path())
}

func (p *Pool) getFree() ignitemeta.DMID {
	computeFree := func() int {
		for i, device := range p.Status.Devices {
			if device == nil {
				return i
			}
		}

		return len(p.Status.Devices)
	}

	if !p.freeComputed {
		p.free = ignitemeta.NewDMID(computeFree())
		p.freeComputed = true
	}

	return p.free
}

// This returns a free ID in the pool
// TODO: Check that this works correctly
func (p *Pool) newID() ignitemeta.DMID {
	index := p.getFree().Index()
	nDevices := len(p.Status.Devices)

	if index < nDevices {
		for i := index + 1; i <= nDevices; i++ {
			if i == nDevices || p.Status.Devices[i] == nil {
				p.free = ignitemeta.NewDMID(i)
				break
			}
		}
	} else {
		p.Status.Devices = append(p.Status.Devices, nil)
		p.free = ignitemeta.NewDMID(len(p.Status.Devices))
	}

	return ignitemeta.NewDMID(index)
}

func (p *Pool) newDevice(genFunc func(ignitemeta.DMID) (*Device, error)) (*Device, error) {
	free := p.getFree()
	id := p.newID()

	device, err := genFunc(id)
	if err != nil {
		p.free = free
	} else {
		p.Status.Devices[id.Index()] = device.PoolDevice
	}

	return device, nil
}

func (p *Pool) Remove(id ignitemeta.DMID) {
	if p.GetDevice(id) != nil {
		p.Status.Devices[id.Index()] = nil

		if p.getFree().Index() > id.Index() {
			p.free = id
		}
	}
}
