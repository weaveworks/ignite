package v1alpha3

import (
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	igniteNetwork "github.com/weaveworks/ignite/pkg/network"
	igniteRuntime "github.com/weaveworks/ignite/pkg/runtime"
	"github.com/weaveworks/libgitops/pkg/runtime"
)

const (
	KindImage  runtime.Kind = "Image"
	KindKernel runtime.Kind = "Kernel"
	KindVM     runtime.Kind = "VM"
)

// Image represents a cached OCI image ready to be used with Ignite
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Image struct {
	runtime.TypeMeta `json:",inline"`
	// runtime.ObjectMeta is also embedded into the struct, and defines the human-readable name, and the machine-readable ID
	// Name is available at the .metadata.name JSON path
	// ID is available at the .metadata.uid JSON path (the Go type is k8s.io/apimachinery/pkg/types.UID, which is only a typed string)
	runtime.ObjectMeta `json:"metadata"`

	Spec   ImageSpec   `json:"spec"`
	Status ImageStatus `json:"status"`
}

// ImageSpec declares what the image contains
type ImageSpec struct {
	OCI meta.OCIImageRef `json:"oci"`
}

// OCIImageSource specifies how the OCI image was imported.
// It is the status variant of OCIImageClaim
type OCIImageSource struct {
	// ID defines the source's content ID (e.g. the canonical OCI path or Docker image ID)
	ID *meta.OCIContentID `json:"id"`
	// Size defines the size of the source in bytes
	Size meta.Size `json:"size"`
}

// ImageStatus defines the status of the image
type ImageStatus struct {
	// OCISource contains the information about how this OCI image was imported
	OCISource OCIImageSource `json:"ociSource"`
}

// Pool defines device mapper pool database
// This file is managed by the snapshotter part of Ignite, and the file (existing as a singleton)
// is present at /var/lib/firecracker/snapshotter/pool.json
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Pool struct {
	runtime.TypeMeta `json:",inline"`
	// Not needed (yet)
	// runtime.ObjectMeta `json:"metadata"`

	Spec   PoolSpec   `json:"spec"`
	Status PoolStatus `json:"status"`
}

// PoolSpec defines the Pool's specification
type PoolSpec struct {
	// MetadataSize specifies the size of the pool's metadata
	MetadataSize meta.Size `json:"metadataSize"`
	// DataSize specifies the size of the pool's data
	DataSize meta.Size `json:"dataSize"`
	// AllocationSize specifies the smallest size that can be allocated at a time
	AllocationSize meta.Size `json:"allocationSize"`
	// MetadataPath points to the file where device mapper stores all metadata information
	// Defaults to constants.SNAPSHOTTER_METADATA_PATH
	MetadataPath string `json:"metadataPath"`
	// DataPath points to the backing physical device or sparse file (to be loop mounted) for the pool
	// Defaults to constants.SNAPSHOTTER_DATA_PATH
	DataPath string `json:"dataPath"`
}

// PoolStatus defines the Pool's current status
type PoolStatus struct {
	// The Devices array needs to contain pointers to accommodate "holes" in the mapping
	// Where devices have been deleted, the pointer is nil
	Devices []*PoolDevice `json:"devices"`
}

type PoolDeviceType string

const (
	PoolDeviceTypeImage  PoolDeviceType = "Image"
	PoolDeviceTypeResize PoolDeviceType = "Resize"
	PoolDeviceTypeKernel PoolDeviceType = "Kernel"
	PoolDeviceTypeVM     PoolDeviceType = "VM"
)

// PoolDevice defines one device in the pool
type PoolDevice struct {
	Size   meta.Size `json:"size"`
	Parent meta.DMID `json:"parent"`
	// Type specifies the type of the contents of the device
	Type PoolDeviceType `json:"type"`
	// MetadataPath points to the JSON/YAML file with metadata about this device
	// This is most often of the format /var/lib/firecracker/{type}/{id}/metadata.json
	MetadataPath string `json:"metadataPath"`
}

// Kernel is a serializable object that caches information about imported kernels
// This file is stored in /var/lib/firecracker/kernels/{oci-image-digest}/metadata.json
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Kernel struct {
	runtime.TypeMeta `json:",inline"`
	// runtime.ObjectMeta is also embedded into the struct, and defines the human-readable name, and the machine-readable ID
	// Name is available at the .metadata.name JSON path
	// ID is available at the .metadata.uid JSON path (the Go type is k8s.io/apimachinery/pkg/types.UID, which is only a typed string)
	runtime.ObjectMeta `json:"metadata"`

	Spec   KernelSpec   `json:"spec"`
	Status KernelStatus `json:"status"`
}

// KernelSpec describes the properties of a kernel
type KernelSpec struct {
	OCI meta.OCIImageRef `json:"oci"`
	// Optional future feature, support per-kernel specific default command lines
	// DefaultCmdLine string
}

// KernelStatus describes the status of a kernel
type KernelStatus struct {
	Version   string         `json:"version"`
	OCISource OCIImageSource `json:"ociSource"`
}

// VM represents a virtual machine run by Firecracker
// These files are stored in /var/lib/firecracker/vm/{vm-id}/metadata.json
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type VM struct {
	runtime.TypeMeta `json:",inline"`
	// runtime.ObjectMeta is also embedded into the struct, and defines the human-readable name, and the machine-readable ID
	// Name is available at the .metadata.name JSON path
	// ID is available at the .metadata.uid JSON path (the Go type is k8s.io/apimachinery/pkg/types.UID, which is only a typed string)
	runtime.ObjectMeta `json:"metadata"`

	Spec   VMSpec   `json:"spec"`
	Status VMStatus `json:"status"`
}

// VMSpec describes the configuration of a VM
type VMSpec struct {
	Image    VMImageSpec   `json:"image"`
	Sandbox  VMSandboxSpec `json:"sandbox"`
	Kernel   VMKernelSpec  `json:"kernel"`
	CPUs     uint64        `json:"cpus"`
	Memory   meta.Size     `json:"memory"`
	DiskSize meta.Size     `json:"diskSize"`
	// TODO: Implement working omitempty without pointers for the following entries
	// Currently both will show in the JSON output as empty arrays. Making them
	// pointers requires plenty of nil checks (as their contents are accessed directly)
	// and is very risky for stability. APIMachinery potentially has a solution.
	Network VMNetworkSpec `json:"network,omitempty"`
	Storage VMStorageSpec `json:"storage,omitempty"`
	// This will be done at either "ignite start" or "ignite create" time
	// TODO: We might revisit this later
	CopyFiles []FileMapping `json:"copyFiles,omitempty"`
	// SSH specifies how the SSH setup should be done
	// nil here means "don't do anything special"
	// If SSH.Generate is set, Ignite will generate a new SSH key and copy it in to authorized_keys in the VM
	// Specifying a path in SSH.Generate means "use this public key"
	// If SSH.PublicKey is set, this struct will marshal as a string using that path
	// If SSH.Generate is set, this struct will marshal as a bool => true
	SSH *SSH `json:"ssh,omitempty"`
}

type VMImageSpec struct {
	OCI meta.OCIImageRef `json:"oci"`
}

type VMKernelSpec struct {
	OCI     meta.OCIImageRef `json:"oci"`
	CmdLine string           `json:"cmdLine,omitempty"`
}

// VMSandboxSpec is the spec of the sandbox used for the VM.
type VMSandboxSpec struct {
	OCI meta.OCIImageRef `json:"oci"`
}

type VMNetworkSpec struct {
	Ports meta.PortMappings `json:"ports,omitempty"`
}

// VMStorageSpec defines the VM's Volumes and VolumeMounts
type VMStorageSpec struct {
	Volumes      []Volume      `json:"volumes,omitempty"`
	VolumeMounts []VolumeMount `json:"volumeMounts,omitempty"`
}

// Volume defines named storage volume
type Volume struct {
	Name        string             `json:"name"`
	BlockDevice *BlockDeviceVolume `json:"blockDevice,omitempty"`
}

// BlockDeviceVolume defines a block device on the host
type BlockDeviceVolume struct {
	Path string `json:"path"`
}

// VolumeMount defines the mount point for a named volume inside a VM
type VolumeMount struct {
	Name      string `json:"name"`
	MountPath string `json:"mountPath"`
}

// FileMapping defines mappings between files on the host and VM
type FileMapping struct {
	HostPath string `json:"hostPath"`
	VMPath   string `json:"vmPath"`
}

// SSH specifies different ways to connect via SSH to the VM
// SSH uses a custom marshaller/unmarshaller. If generate is true,
// it marshals to true (a JSON bool). If PublicKey is set, it marshals
// to that string.
type SSH struct {
	Generate  bool   `json:"-"`
	PublicKey string `json:"-"`
}

// Runtime specifies the VM's runtime information
type Runtime struct {
	ID   string             `json:"id"`
	Name igniteRuntime.Name `json:"name"`
}

// Network specifies the VM's network information.
type Network struct {
	Plugin      igniteNetwork.PluginName `json:"plugin"`
	IPAddresses meta.IPAddresses         `json:"ipAddresses"`
}

// VMStatus defines the status of a VM
type VMStatus struct {
	Running   bool           `json:"running"`
	Runtime   *Runtime       `json:"runtime,omitempty"`
	StartTime *runtime.Time  `json:"startTime,omitempty"`
	Network   *Network       `json:"network,omitempty"`
	Image     OCIImageSource `json:"image"`
	Kernel    OCIImageSource `json:"kernel"`
	IDPrefix  string         `json:"idPrefix"`
}

// Configuration represents the ignite runtime configuration.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Configuration struct {
	runtime.TypeMeta   `json:",inline"`
	runtime.ObjectMeta `json:"metadata"`

	Spec ConfigurationSpec `json:"spec"`
}

// ConfigurationSpec defines the ignite configuration.
type ConfigurationSpec struct {
	Runtime       igniteRuntime.Name       `json:"runtime,omitempty"`
	NetworkPlugin igniteNetwork.PluginName `json:"networkPlugin,omitempty"`
	VMDefaults    VMSpec                   `json:"vmDefaults,omitempty"`
	IDPrefix      string                   `json:"idPrefix,omitempty"`
}
