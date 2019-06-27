package snapshotter

import (
	"fmt"
	"os"

	"github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/dm"
	"github.com/weaveworks/ignite/pkg/util"
)

var (
	dataSize       = v1alpha1.NewSizeFromBytes(constants.POOL_DATA_SIZE_BYTES)           // Size of the common pool
	allocationSize = v1alpha1.NewSizeFromSectors(constants.POOL_ALLOCATION_SIZE_SECTORS) // Pool allocation block size
	extraSize      = v1alpha1.NewSizeFromBytes(constants.POOL_VOLUME_EXTRA_SIZE)         // Additional space to add to image volumes for the ext4 partition
)

// Snapshotter abstracts the device mapper pool and provides convenience methods
// It's also responsible for (de)serializing the pool
type Snapshotter struct {
	*dm.Pool
}

// NewSnapshotter creates a new snapshotter with a new pool
// This should only be called on the first run of Ignite
// TODO: No support for physical backing devices for now
func NewSnapshotter() (*Snapshotter, error) {
	metadataSize := calcMetadataDevSize(dataSize)

	s := &Snapshotter{
		dm.NewPool(
			dataSize,
			metadataSize,
			allocationSize,
			constants.SNAPSHOTTER_METADATA_PATH,
			constants.SNAPSHOTTER_DATA_PATH),
	}

	if err := s.initialize(); err != nil {
		return nil, err
	}

	return s, nil
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
