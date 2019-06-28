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
	*dm.Pool
	images  []*Image
	resizes []*Resize
	kernels []*Kernel
	vms     []*VM
}

// NewSnapshotter creates a new snapshotter with a new pool
// This should only be called on the first run of Ignite
// TODO: No support for physical backing devices for now
func NewSnapshotter() (*Snapshotter, error) {
	metadataSize := calcMetadataDevSize(dataSize)

	s := &Snapshotter{
		Pool: dm.NewPool(
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

// Loader, which loads the pool and then all devices based on their metadata
func (s *Snapshotter) Load() error {
	if err := scheme.DecodeFileInto(path.Join(constants.SNAPSHOTTER_DIR, constants.METADATA), s); err != nil {
		return err
	}

	if err := s.ForDevices(func(id v1alpha1.DMID, device *dm.Device) error {
		var metadata runtime.Object = nil

		if len(device.MetadataPath) > 0 {
			if err := scheme.DecodeFileInto(device.MetadataPath, metadata); err != nil {
				return err
			}
		}

		switch device.Type {
		case v1alpha1.PoolDeviceTypeImage:
			s.images = append(s.images, NewImage(*metadata.(*v1alpha1.Image), device))
		case v1alpha1.PoolDeviceTypeResize:
			s.resizes = append(s.resizes, NewResize(device))
		case v1alpha1.PoolDeviceTypeKernel:
			s.kernels = append(s.kernels, NewKernel(*metadata.(*v1alpha1.Kernel), device))
		case v1alpha1.PoolDeviceTypeVM:
			s.vms = append(s.vms, NewVM(*metadata.(*v1alpha1.VM), device))
		default:
			return fmt.Errorf("unknown pool device type: %s", device.Type)
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
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
		s.Spec.MetadataPath: s.Spec.MetadataSize,
		s.Spec.DataPath:     s.Spec.DataSize,
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
