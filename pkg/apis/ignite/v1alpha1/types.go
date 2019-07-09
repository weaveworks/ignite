package v1alpha1

import (
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
)

const (
	KindImage  meta.Kind = "Image"
	KindKernel meta.Kind = "Kernel"
	KindVM     meta.Kind = "VM"
)

// Image represents a cached OCI image ready to be used with Ignite
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Image struct {
	meta.TypeMeta `json:",inline"`
	// meta.ObjectMeta is also embedded into the struct, and defines the human-readable name, and the machine-readable ID
	// Name is available at the .metadata.name JSON path
	// ID is available at the .metadata.uid JSON path (the Go type is k8s.io/apimachinery/pkg/types.UID, which is only a typed string)
	meta.ObjectMeta `json:"metadata"`

	Spec   ImageSpec   `json:"spec"`
	Status ImageStatus `json:"status"`
}

// ImageSpec declares what the image contains
type ImageSpec struct {
	OCIClaim OCIImageClaim `json:"ociClaim"`
}

// ImageSourceType is an enum of different supported Image Source Types
type ImageSourceType string

const (
	// ImageSourceTypeDocker defines that the image is imported from Docker
	ImageSourceTypeDocker ImageSourceType = "Docker"
)

// OCIImageClaim defines a claim for importing an OCI image
type OCIImageClaim struct {
	// Type defines how the image should be imported
	Type ImageSourceType `json:"type"`
	// Ref defines the reference to use when talking to the backend.
	// This is most commonly the image name, followed by a tag.
	// Other supported ways are $registry/$user/$image@sha256:$digest
	// This ref is also used as ObjectMeta.Name for kinds Images and Kernels
	Ref meta.OCIImageRef `json:"ref"`
}

// OCIImageSource specifies how the OCI image was imported.
// It is the status variant of OCIImageClaim
type OCIImageSource struct {
	// ID defines the source's ID (e.g. the Docker image ID)
	ID string `json:"id"`
	// Size defines the size of the source in bytes
	Size meta.Size `json:"size"`
	// RepoDigests defines the image name as it was when pulled
	// from a repository, and the digest of the image
	// The format is $registry/$user/$image@sha256:$digest
	// This field is unpopulated if the image used as the source
	// has never been pushed to or pulled from a registry
	RepoDigests []string `json:"repoDigests,omitempty"`
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
	meta.TypeMeta `json:",inline"`
	// Not needed (yet)
	// meta.ObjectMeta `json:"metadata"`

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
	meta.TypeMeta `json:",inline"`
	// meta.ObjectMeta is also embedded into the struct, and defines the human-readable name, and the machine-readable ID
	// Name is available at the .metadata.name JSON path
	// ID is available at the .metadata.uid JSON path (the Go type is k8s.io/apimachinery/pkg/types.UID, which is only a typed string)
	meta.ObjectMeta `json:"metadata"`

	Spec   KernelSpec   `json:"spec"`
	Status KernelStatus `json:"status"`
}

// KernelSpec describes the properties of a kernel
type KernelSpec struct {
	OCIClaim OCIImageClaim `json:"ociClaim"`
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
	meta.TypeMeta `json:",inline"`
	// meta.ObjectMeta is also embedded into the struct, and defines the human-readable name, and the machine-readable ID
	// Name is available at the .metadata.name JSON path
	// ID is available at the .metadata.uid JSON path (the Go type is k8s.io/apimachinery/pkg/types.UID, which is only a typed string)
	meta.ObjectMeta `json:"metadata"`

	Spec   VMSpec   `json:"spec"`
	Status VMStatus `json:"status"`
}

// VMSpec describes the configuration of a VM
type VMSpec struct {
	Image    VMImageSpec   `json:"image"`
	Kernel   VMKernelSpec  `json:"kernel"`
	CPUs     uint64        `json:"cpus"`
	Memory   meta.Size     `json:"memory"`
	DiskSize meta.Size     `json:"diskSize"`
	Network  VMNetworkSpec `json:"network"`

	// This will be done at either "ignite start" or "ignite create" time
	// TODO: We might to revisit this later
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
	OCIClaim OCIImageClaim `json:"ociClaim"`
}

type VMKernelSpec struct {
	OCIClaim OCIImageClaim `json:"ociClaim"`
	CmdLine  string        `json:"cmdLine,omitempty"`
}

type VMNetworkSpec struct {
	Mode  NetworkMode       `json:"mode"`
	Ports meta.PortMappings `json:"ports,omitempty"`
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

// NetworkMode defines different states a VM can be in
type NetworkMode string

func (nm NetworkMode) String() string {
	return string(nm)
}

const (
	// NetworkModeCNI specifies the network mode where CNI is used
	NetworkModeCNI NetworkMode = "cni"
	// NetworkModeDockerBridge specifies the default docker bridge network is used
	NetworkModeDockerBridge NetworkMode = "docker-bridge"
	// Whenever updating this list, also update GetNetworkModes in helpers.go
)

// VMState defines different states a VM can be in
type VMState string

const (
	VMStateCreated VMState = "Created"
	VMStateRunning VMState = "Running"
	VMStateStopped VMState = "Stopped"
)

// VMStatus defines the status of a VM
type VMStatus struct {
	State       VMState          `json:"state"`
	IPAddresses meta.IPAddresses `json:"ipAddresses,omitempty"`
	Image       OCIImageSource   `json:"image"`
	Kernel      OCIImageSource   `json:"kernel"`
}
