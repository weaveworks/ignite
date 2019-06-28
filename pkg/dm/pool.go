package dm

import (
	"fmt"
	"path"

	"github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"

	"github.com/weaveworks/ignite/pkg/util"
)

type Pool struct {
	v1alpha1.Pool
	free v1alpha1.DMID
}

var poolName = util.NewPrefixer().Prefix("pool")

func NewPool(metadataSize, dataSize, allocationSize v1alpha1.Size, metadataPath, dataPath string) *Pool {
	return &Pool{
		Pool: v1alpha1.Pool{
			Spec: v1alpha1.PoolSpec{
				MetadataSize:   metadataSize,
				DataSize:       dataSize,
				AllocationSize: allocationSize,
				MetadataPath:   metadataPath,
				DataPath:       dataPath,
			},
		},
		free: v1alpha1.NewDMID(0),
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

func (p *Pool) getID(device *Device) v1alpha1.DMID {
	// If the querying for nil, return the pool's ID
	if device == nil {
		return v1alpha1.NewPoolDMID()
	}

	for i, d := range p.Status.Devices {
		if d == device.PoolDevice {
			return v1alpha1.NewDMID(i)
		}
	}

	// This should never happen, unless you try to get
	// the ID of a device residing in another pool
	panic("pool getID: device not found")
}

// GetDevice dynamically spawns a device from a v1alpha1.PoolDevice
func (p *Pool) GetDevice(id v1alpha1.DMID) *Device {
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
func (p *Pool) ForDevices(iterFunc func(v1alpha1.DMID, *Device) error) error {
	for i := 0; i < len(p.Status.Devices); i++ {
		spec := p.Status.Devices[i]
		if spec != nil {
			id := v1alpha1.NewDMID(i)
			if err := iterFunc(id, p.GetDevice(id)); err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *Pool) activate() error {
	// Don't try to activate an already active pool
	if p.active() {
		return nil
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

// If /dev/mapper/<name> exists the pool is active
func (p *Pool) active() bool {
	return util.FileExists(p.Path())
}

// This returns a free ID in the pool
// TODO: Check that this works correctly
func (p *Pool) newID() v1alpha1.DMID {
	index := p.free.Index()
	nDevices := len(p.Status.Devices)

	if index < nDevices {
		for i := index + 1; i <= nDevices; i++ {
			if i == nDevices || p.Status.Devices[i] == nil {
				p.free = v1alpha1.NewDMID(i)
				break
			}
		}
	} else {
		p.Status.Devices = append(p.Status.Devices, nil)
		p.free = v1alpha1.NewDMID(len(p.Status.Devices))
	}

	return v1alpha1.NewDMID(index)
}

func (p *Pool) newDevice(genFunc func(v1alpha1.DMID) (*Device, error)) (*Device, error) {
	free := p.free
	id := p.newID()

	device, err := genFunc(id)
	if err != nil {
		p.free = free
	} else {
		p.Status.Devices[id.Index()] = device.PoolDevice
	}

	return device, nil
}

func (p *Pool) Remove(id v1alpha1.DMID) {
	if p.GetDevice(id) != nil {
		p.Status.Devices[id.Index()] = nil

		if p.free.Index() > id.Index() {
			p.free = id
		}
	}
}
