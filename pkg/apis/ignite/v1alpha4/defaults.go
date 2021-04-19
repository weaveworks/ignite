package v1alpha4

import (
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/constants"
	igniteNetwork "github.com/weaveworks/ignite/pkg/network"
	igniteRuntime "github.com/weaveworks/ignite/pkg/runtime"
	"github.com/weaveworks/ignite/pkg/version"

	"k8s.io/apimachinery/pkg/runtime"
)

func addDefaultingFuncs(scheme *runtime.Scheme) error {
	return RegisterDefaults(scheme)
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
}

func SetDefaults_VMKernelSpec(obj *VMKernelSpec) {
	// Default the kernel image if unset.
	if obj.OCI.IsUnset() {
		obj.OCI, _ = meta.NewOCIImageRef(version.GetIgnite().KernelImage.String())
	}

	if len(obj.CmdLine) == 0 {
		obj.CmdLine = constants.VM_DEFAULT_KERNEL_ARGS
	}
}

func SetDefaults_VMSandboxSpec(obj *VMSandboxSpec) {
	// Default the sandbox image if unset.
	if obj.OCI.IsUnset() {
		obj.OCI, _ = meta.NewOCIImageRef(version.GetIgnite().SandboxImage.String())
	}
}

func SetDefaults_ConfigurationSpec(obj *ConfigurationSpec) {
	// Default the runtime and network plugin if not set.
	if obj.Runtime == "" {
		obj.Runtime = igniteRuntime.RuntimeContainerd
	}
	if obj.NetworkPlugin == "" {
		obj.NetworkPlugin = igniteNetwork.PluginCNI
	}
}

func calcMetadataDevSize(obj *PoolSpec) meta.Size {
	// The minimum size is 2 MB and the maximum size is 16 GB
	minSize := meta.NewSizeFromBytes(2 * constants.MB)
	maxSize := meta.NewSizeFromBytes(16 * constants.GB)

	return meta.NewSizeFromBytes(48 * obj.DataSize.Bytes() / obj.AllocationSize.Bytes()).Min(maxSize).Max(minSize)
}

func SetDefaults_VMStatus(obj *VMStatus) {
	if obj.Runtime == nil {
		obj.Runtime = &Runtime{}
	}
	if obj.Network == nil {
		obj.Network = &Network{}
	}
}
