package snapshotter

import (
	"fmt"
	"github.com/weaveworks/ignite/pkg/apis/ignite/scheme"
	"k8s.io/apimachinery/pkg/runtime"
	"os"
	"path"

	"github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/dm"
	"github.com/weaveworks/ignite/pkg/util"
)

var (
	dataSize       = v1alpha1.NewSizeFromBytes(constants.POOL_DATA_SIZE_BYTES)           // Size of the common pool
	allocationSize = v1alpha1.NewSizeFromSectors(constants.POOL_ALLOCATION_SIZE_SECTORS) // Pool allocation block size
)

// Snapshotter abstracts the device mapper pool and provides convenience methods
// It's also responsible for (de)serializing the pool

type Snapshotter struct {
	pool    *dm.Pool
	objects []*Object
	Images  []*Image
	resizes []*Resize
	Kernels []*Kernel
	VMs     []*VM

	loaded bool
}

// TODO: This should store objects as "generators" like the dm.Devices
// One list of all objects, which are generic structs:

type Object struct {
	device *dm.Device
	object runtime.Object
}

// This getter loads the metadata on demand
func (o *Object) GetRuntimeObject() (runtime.Object, error) {
	if o.object == nil && len(o.device.MetadataPath) > 0 {
		if err := scheme.DecodeFileInto(o.device.MetadataPath, o.object); err != nil {
			return nil, err
		}
	}

	return o.object, nil
}

type Filter interface {
	Filter(*Object) *Object
}

// The runtime.Object is behind a private field, which enables a getter to load it on-demand

// The runtime.Object stores the v1alpha1.{Image,Kernel,VM...}, basically the loaded metadata
// Then we have a filtering framework to extract the wanted Objects:
// (there should be an AllFilter which matches everything)

func (s *Snapshotter) getSingle(f Filter, t v1alpha1.PoolDeviceType) (*Object, error) {
	var result *Object

	for _, object := range s.objects {
		if object.device.Type == t {
			if match := f.Filter(object); match != nil { // Filter returns *Object if it matches, otherwise nil
				if result != nil {
					return nil, ErrAmbiguous
				} else {
					result = match
				}
			}
		}
	}

	if result == nil {
		return nil, ErrNonexistent
	}

	return result, nil
}

func (s *Snapshotter) GetImage(f Filter) (*Image, error) {
	result, err := s.getSingle(f, v1alpha1.PoolDeviceTypeImage)
	if err != nil {
		return nil, err
	}

	return &Image{
		*result.GetRuntimeObject().(*v1alpha1.Image),
		result.device,
	}, nil
}

func (s *Snapshotter) GetKernel(f Filter) (*Kernel, error) {
	result, err := s.getSingle(f, v1alpha1.PoolDeviceTypeKernel)
	if err != nil {
		return nil, err
	}

	return &Kernel{
		*result.object.(*v1alpha1.Kernel),
		result.device,
	}, nil
}

func (s *Snapshotter) GetVM() *Image {

}

// NewSnapshotter creates a new snapshotter with a new pool
// This should only be called on the first run of Ignite
// TODO: No support for physical backing devices for now
func NewSnapshotter() (*Snapshotter, error) {
	metadataSize := calcMetadataDevSize(dataSize)

	s := &Snapshotter{
		pool: dm.NewPool(
			metadataSize,
			dataSize,
			allocationSize,
			constants.SNAPSHOTTER_METADATA_PATH,
			constants.SNAPSHOTTER_DATA_PATH),
	}

	if err := s.initialize(); err != nil {
		return nil, err
	}

	return s, nil
}

// TODO: Dependency-based loader, which loads the given device and it's parents recursively
// TODO: Even better: conditional loading on access

// Loader, which loads the pool and then all devices based on their metadata
func (s *Snapshotter) LoadAll() error {
	// Don't load twice
	if s.loaded {
		return nil
	}

	// Load the pool configuration
	if err := scheme.DecodeFileInto(path.Join(constants.SNAPSHOTTER_DIR, constants.METADATA), s); err != nil {
		return err
	}

	if err := s.pool.ForDevices(func(id v1alpha1.DMID, device *dm.Device) error {
		var metadata runtime.Object = nil

		if len(device.MetadataPath) > 0 {
			if err := scheme.DecodeFileInto(device.MetadataPath, metadata); err != nil {
				return err
			}
		}

		switch device.Type {
		case v1alpha1.PoolDeviceTypeImage:
			s.Images = append(s.Images, NewImage(*metadata.(*v1alpha1.Image), device))
		case v1alpha1.PoolDeviceTypeResize:
			s.resizes = append(s.resizes, NewResize(device))
		case v1alpha1.PoolDeviceTypeKernel:
			s.Kernels = append(s.Kernels, NewKernel(*metadata.(*v1alpha1.Kernel), device))
		case v1alpha1.PoolDeviceTypeVM:
			s.VMs = append(s.VMs, NewVM(*metadata.(*v1alpha1.VM), device))
		default:
			return fmt.Errorf("unknown pool device type: %s", device.Type)
		}

		return nil
	}); err != nil {
		return err
	}

	s.loaded = true
	return nil
}

// TODO: These are clunky, need to figure out a better solution
func (s *Snapshotter) ImageObjects() []v1alpha1.Object {
	objects := make([]v1alpha1.Object, 0, len(s.Images))

	for _, image := range s.Images {
		objects = append(objects, image)
	}

	return objects
}

func (s *Snapshotter) VMObjects() []v1alpha1.Object {
	objects := make([]v1alpha1.Object, 0, len(s.VMs))

	for _, vm := range s.VMs {
		objects = append(objects, vm)
	}

	return objects
}

func ObjectToImage(object v1alpha1.Object) *Image {
	return object.(*Image)
}

func ObjectsToImages(objects []v1alpha1.Object) []*Image {
	images := make([]*Image, 0, len(objects))

	for _, object := range objects {
		images = append(images, ObjectToImage(object))
	}

	return images
}

func ObjectToVM(object v1alpha1.Object) *VM {
	return object.(*VM)
}

func ObjectsToVMs(objects []v1alpha1.Object) []*VM {
	vms := make([]*VM, 0, len(objects))

	for _, object := range objects {
		vms = append(vms, ObjectToVM(object))
	}

	return vms
}

// initialize creates the snapshotter directory and
// allocates the backing data/metadata devices
// TODO: Move this into dm?
func (s *Snapshotter) initialize() error {
	if err := os.MkdirAll(constants.SNAPSHOTTER_DIR, constants.DATA_DIR_PERM); err != nil {
		return err
	}

	// Allocate the thin provisioning data and metadata files
	thinFiles := map[string]v1alpha1.Size{
		s.pool.Spec.MetadataPath: s.pool.Spec.MetadataSize,
		s.pool.Spec.DataPath:     s.pool.Spec.DataSize,
	}

	for p, size := range thinFiles {
		if !util.FileExists(p) {
			file, err := os.Create(p)
			if err != nil {
				return fmt.Errorf("failed to create thin provisioning file %q: %v", p, err)
			}

			// Allocate the image file
			if err := file.Truncate(int64(size.Bytes())); err != nil {
				return fmt.Errorf("failed to allocate space for thin provisioning file %q: %v", p, err)
			}

			if err := file.Close(); err != nil {
				return err
			}
		}
	}

	return nil
}

func calcMetadataDevSize(dataDeviceSize v1alpha1.Size) v1alpha1.Size {
	// The minimum size is 2 MB and the maximum size is 16 GB
	minSize := v1alpha1.NewSizeFromBytes(2 * constants.MB)
	maxSize := v1alpha1.NewSizeFromBytes(16 * constants.GB)

	return v1alpha1.NewSizeFromBytes(48 * dataDeviceSize.Bytes() / allocationSize.Bytes()).Min(maxSize).Max(minSize)
}
