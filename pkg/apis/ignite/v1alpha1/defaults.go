package v1alpha1

import (
	ignitemeta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/constants"
	"k8s.io/apimachinery/pkg/runtime"
)

func addDefaultingFuncs(scheme *runtime.Scheme) error {
	return RegisterDefaults(scheme)
}

func SetDefaults_ImageSource(obj *ImageSource) {
	obj.Type = ImageSourceTypeDocker
}

func SetDefaults_PoolSpec(obj *PoolSpec) {
	if obj.AllocationSize == ignitemeta.EmptySize {
		obj.AllocationSize = ignitemeta.NewSizeFromSectors(constants.POOL_ALLOCATION_SIZE_SECTORS)
	}

	if obj.DataSize == ignitemeta.EmptySize {
		obj.AllocationSize = ignitemeta.NewSizeFromBytes(constants.POOL_DATA_SIZE_BYTES)
	}

	if obj.MetadataSize == ignitemeta.EmptySize {
		obj.AllocationSize = calcMetadataDevSize(obj)
	}

	if len(obj.MetadataPath) == 0 {
		obj.MetadataPath = constants.SNAPSHOTTER_METADATA_PATH
	}

	if len(obj.DataPath) == 0 {
		obj.DataPath = constants.SNAPSHOTTER_DATA_PATH
	}
}

func SetDefaults_VMSpec(obj *VMSpec) {
	if obj.CPUs == 0 {
		obj.CPUs = constants.VM_DEFAULT_CPUS
	}

	// TODO: These might be nil instead of ignitemeta.EmptySize
	if obj.Memory == ignitemeta.EmptySize {
		obj.Memory = ignitemeta.NewSizeFromBytes(constants.VM_DEFAULT_MEMORY)
	}

	if obj.Size == ignitemeta.EmptySize {
		obj.Size = ignitemeta.NewSizeFromBytes(constants.VM_DEFAULT_SIZE)
	}
}

func SetDefaults_VMStatus(obj *VMStatus) {
	if obj.State == "" {
		obj.State = VMStateCreated
	}
}

func calcMetadataDevSize(obj *PoolSpec) ignitemeta.Size {
	// The minimum size is 2 MB and the maximum size is 16 GB
	minSize := ignitemeta.NewSizeFromBytes(2 * constants.MB)
	maxSize := ignitemeta.NewSizeFromBytes(16 * constants.GB)

	return ignitemeta.NewSizeFromBytes(48 * obj.DataSize.Bytes() / obj.AllocationSize.Bytes()).Min(maxSize).Max(minSize)
}
