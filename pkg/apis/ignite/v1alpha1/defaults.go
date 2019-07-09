package v1alpha1

import (
	"reflect"

	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/constants"
	"k8s.io/apimachinery/pkg/runtime"
)

func addDefaultingFuncs(scheme *runtime.Scheme) error {
	return RegisterDefaults(scheme)
}

func SetDefaults_OCIImageClaim(obj *OCIImageClaim) {
	obj.Type = ImageSourceTypeDocker
}

func SetDefaults_PoolSpec(obj *PoolSpec) {
	if obj.AllocationSize == meta.EmptySize {
		obj.AllocationSize = meta.NewSizeFromSectors(constants.POOL_ALLOCATION_SIZE_SECTORS)
	}

	if obj.DataSize == meta.EmptySize {
		obj.AllocationSize = meta.NewSizeFromBytes(constants.POOL_DATA_SIZE_BYTES)
	}

	if obj.MetadataSize == meta.EmptySize {
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

	if obj.Memory == meta.EmptySize {
		obj.Memory = meta.NewSizeFromBytes(constants.VM_DEFAULT_MEMORY)
	}

	if obj.DiskSize == meta.EmptySize {
		obj.DiskSize = meta.NewSizeFromBytes(constants.VM_DEFAULT_SIZE)
	}

	if len(obj.NetworkMode) == 0 {
		obj.NetworkMode = NetworkModeDockerBridge
	}
}

func SetDefaults_VMKernelSpec(obj *VMKernelSpec) {
	// Default the kernel image if unset
	if len(obj.OCIClaim.Ref) == 0 {
		obj.OCIClaim.Ref, _ = meta.NewOCIImageRef(constants.DEFAULT_KERNEL_IMAGE)
	}

	if len(obj.CmdLine) == 0 {
		obj.CmdLine = constants.VM_DEFAULT_KERNEL_ARGS
	}
}

func SetDefaults_VMStatus(obj *VMStatus) {
	if obj.State == "" {
		obj.State = VMStateCreated
	}
}

func calcMetadataDevSize(obj *PoolSpec) meta.Size {
	// The minimum size is 2 MB and the maximum size is 16 GB
	minSize := meta.NewSizeFromBytes(2 * constants.MB)
	maxSize := meta.NewSizeFromBytes(16 * constants.GB)

	return meta.NewSizeFromBytes(48 * obj.DataSize.Bytes() / obj.AllocationSize.Bytes()).Min(maxSize).Max(minSize)
}

// TODO: Temporary hacks to populate TypeMeta until we get the generator working
func SetDefaults_VM(obj *VM) {
	setTypeMeta(obj)
}

func SetDefaults_Image(obj *Image) {
	setTypeMeta(obj)
}

func SetDefaults_Kernel(obj *Kernel) {
	setTypeMeta(obj)
}

func setTypeMeta(obj meta.Object) {
	obj.GetTypeMeta().APIVersion = SchemeGroupVersion.String()
	obj.GetTypeMeta().Kind = reflect.Indirect(reflect.ValueOf(obj)).Type().Name()
}
