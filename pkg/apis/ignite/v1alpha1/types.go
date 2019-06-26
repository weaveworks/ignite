package v1alpha1

import (
	"net"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Image represents a cached OCI image ready to be used with Ignite
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Image struct {
	metav1.TypeMeta `json:",inline"`
	// ObjectMeta is also embedded into the struct, and defines the human-readable name, and the machine-readable ID
	// Name is available at the .metadata.name JSON path
	// ID is available at the .metadata.uid JSON path (the Go type is k8s.io/apimachinery/pkg/types.UID, which is only a typed string)
	ObjectMeta `json:"metadata"`

	Spec   ImageSpec   `json:"spec"`
	Status ImageStatus `json:"status"`
}

// ImageSpec declares what the image contains
type ImageSpec struct {
	Source ImageSource `json:"source"`
}

// ImageSourceType is an enum of different supported Image Source Types
type ImageSourceType string

const (
	// ImageSourceTypeDocker defines that the image is imported from Docker
	ImageSourceTypeDocker ImageSourceType = "Docker"
)

// ImageSource defines where the image was imported from
type ImageSource struct {
	// Type defines how the image was imported
	Type ImageSourceType `json:"type"`
	// Digest defines the source contents (e.g. the Docker image ID)
	// See https://github.com/opencontainers/image-spec/blob/master/descriptor.md for more info
	Digest string `json:"digest"`
	// Name defines the user-friendly name of the imported source
	Name string `json:"name"`
	// Size defines the size of the source in bytes
	Size Size `json:"size"`
}

// ImageStatus defines the status of the image
type ImageStatus struct {
	// LayerID points to the index of the device in the DM pool
	// TODO: Make this this a dedicated DMID type or similar
	LayerID uint32 `json:"layerID"`
}

// Pool defines devicemapper pool database
// This file is managed by the snapshotter part of Ignite, and the file (existing as a singleton)
// is present at /var/lib/firecracker/snapshotter/pool.json
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Pool struct {
	metav1.TypeMeta `json:",inline"`
	// Not needed (yet)
	// ObjectMeta `json:"metadata"`

	Spec   PoolSpec   `json:"spec"`
	Status PoolStatus `json:"status"`
}

// PoolSpec defines the Pool's specification
type PoolSpec struct {
	Blocks         Size `json:"blocks"`
	AllocationSize Size `json:"allocationSize"`
	// DataPath points to the backing physical device or sparse file (to be loop mounted) for the pool
	// Defaults to /var/lib/firecracker/snapshotter/data.dm
	DataPath string `json:"dataPath"`
	// MetadataPath points to the file where dm stores all metadata information
	// Defaults to /var/lib/firecracker/snapshotter/metadata.dm
	MetadataPath string `json:"metadataPath"`
}

// PoolStatus defines the Pool's current status
type PoolStatus struct {
	Devices []PoolDevice `json:"devices"`
}

// PoolDeviceType defines what kind of DM device this is
/*type PoolDeviceType string

const (
	PoolDeviceTypeImage  PoolDeviceType = "Image"
	PoolDeviceTypeResize PoolDeviceType = "Resize"
	PoolDeviceTypeKernel PoolDeviceType = "Kernel"
	PoolDeviceTypeVM     PoolDeviceType = "VM"
)*/

// PoolDevice defines one device in the pool
type PoolDevice struct {
	Blocks Size    `json:"blocks"`
	Parent *uint32 `json:"parent"`
	// MetadataPath points to the JSON/YAML file with metadata about this device
	// This is most often of the format /var/lib/firecracker/{type}/{id}/metadata.json
	MetadataPath string `json:"metadataPath"`
	//Type   PoolDeviceType `json:"type"`
}

// Kernel is a serializable object that caches information about imported kernels
// This file is stored in /var/lib/firecracker/kernels/{oci-image-digest}/metadata.json
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Kernel struct {
	metav1.TypeMeta `json:",inline"`
	// ObjectMeta is also embedded into the struct, and defines the human-readable name, and the machine-readable ID
	// Name is available at the .metadata.name JSON path
	// ID is available at the .metadata.uid JSON path (the Go type is k8s.io/apimachinery/pkg/types.UID, which is only a typed string)
	ObjectMeta `json:"metadata"`

	Spec KernelSpec `json:"spec"`
	//Status KernelStatus `json:"status"`
}

type KernelSpec struct {
	Version string      `json:"version"`
	Source  ImageSource `json:"source"`
	// Optional future feature, support per-kernel specific default command lines
	// DefaultCmdLine string
}

// VM represents a virtual machine run by Firecracker
// These files are stored in /var/lib/firecracker/vm/{vm-id}/metadata.json
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type VM struct {
	metav1.TypeMeta `json:",inline"`
	// ObjectMeta is also embedded into the struct, and defines the human-readable name, and the machine-readable ID
	// Name is available at the .metadata.name JSON path
	// ID is available at the .metadata.uid JSON path (the Go type is k8s.io/apimachinery/pkg/types.UID, which is only a typed string)
	ObjectMeta `json:"metadata"`

	Spec   VMSpec   `json:"spec"`
	Status VMStatus `json:"status"`
}

// VMSpec de
type VMSpec struct {
	CPUs   uint64        `json:"cpus"`
	Memory Size          `json:"memory"`
	Size   Size          `json:"size"`
	Ports  []PortMapping `json:"ports"`
	// This will be done at either "ignite start" or "ignite create" time
	// TODO: We might to revisit this later
	CopyFiles []FileMapping `json:"copyFiles"`
	// SSH specifies how the SSH setup should be done
	// SSH appends to CopyFiles when active
	// nil here means "don't do anything special"
	// An empty struct means "generate a new SSH key and copy it in"
	// Specifying a path mean "use this public key"
	SSH *SSH `json:"ssh"`
}

// PortMapping defines a port mapping between the VM and the host
type PortMapping struct {
	HostPort uint64 `json:"hostPort"`
	VMPort   uint64 `json:"vmPort"`
}

// FileMapping defines mappings between files on the host and VM
type FileMapping struct {
	HostPath string `json:"hostPath"`
	VMPath   string `json:"vmPath"`
}

// SSH specifies different ways to connect via SSH to the VM
type SSH struct {
	PublicKey string `json:"publicKey"`
}

// VMState defines different states a VM can be in
type VMState string

const (
	VMStateCreated VMState = "Created"
	VMStateRunning VMState = "Running"
	VMStateStopped VMState = "Stopped"
)

// VMStatus defines the status of a VM
type VMStatus struct {
	State       VMState  `json:"state"`
	IPAddresses []net.IP `json:"ipAddresses"`
}
