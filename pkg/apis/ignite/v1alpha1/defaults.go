package v1alpha1

import (
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
	// TODO: These might be nil instead of EmptySize
	if obj.AllocationSize == EmptySize {
		obj.AllocationSize = NewSizeFromSectors(constants.POOL_ALLOCATION_SIZE_SECTORS)
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

	// TODO: These might be nil instead of EmptySize
	if obj.Memory == EmptySize {
		obj.Memory = NewSizeFromBytes(constants.VM_DEFAULT_MEMORY)
	}

	if obj.Size == EmptySize {
		obj.Size = NewSizeFromBytes(constants.VM_DEFAULT_SIZE)
	}
}

func SetDefaults_VMStatus(obj *VMStatus) {
	if obj.State == "" {
		obj.State = VMStateCreated
	}
}
