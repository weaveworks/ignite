package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime"
)

func addDefaultingFuncs(scheme *runtime.Scheme) error {
	return RegisterDefaults(scheme)
}

func SetDefaults_ImageSource(obj *ImageSource) {
	obj.Type = ImageSourceTypeDocker
}

func SetDefaults_PoolSpec(obj *PoolSpec) {
	// TODO: Reference these from globally-defined constants
	if obj.AllocationSize == EmptySize {
		obj.AllocationSize = NewSizeFromSectors(128)
	}
	if len(obj.DataPath) == 0 {
		obj.DataPath = "/var/lib/firecracker/snapshotter/data.dm"
	}
	if len(obj.MetadataPath) == 0 {
		obj.MetadataPath = "/var/lib/firecracker/snapshotter/metadata.dm"
	}
}

func SetDefaults_VMSpec(obj *VMSpec) {
	// TODO: Reference these from globally-defined constants
	if obj.CPUs == 0 {
		obj.CPUs = 1
	}
	if obj.Memory == EmptySize {
		obj.Memory = NewSizeFromBytes(512 * 1024 * 1024)
	}
	if obj.Size == EmptySize {
		obj.Size = NewSizeFromBytes(4 * 1024 * 1024 * 1024)
	}
}

func SetDefaults_VMStatus(obj *VMStatus) {
	if obj.State == "" {
		obj.State = VMStateCreated
	}
}
