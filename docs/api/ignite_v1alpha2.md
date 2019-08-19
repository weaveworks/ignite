# v1alpha2

`import "github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha2"`

  - [Overview](#pkg-overview)
  - [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>

\+k8s:deepcopy-gen=package +k8s:defaulter-gen=TypeMeta
+k8s:openapi-gen=true
+k8s:conversion-gen=github.com/weaveworks/ignite/pkg/apis/ignite

## <a name="pkg-index">Index</a>

  - [Constants](#pkg-constants)
  - [Variables](#pkg-variables)
  - [func SetDefaults\_PoolSpec(obj \*PoolSpec)](#SetDefaults_PoolSpec)
  - [func SetDefaults\_VMKernelSpec(obj
    \*VMKernelSpec)](#SetDefaults_VMKernelSpec)
  - [func SetDefaults\_VMSpec(obj \*VMSpec)](#SetDefaults_VMSpec)
  - [type BlockDeviceVolume](#BlockDeviceVolume)
  - [type FileMapping](#FileMapping)
  - [type Image](#Image)
  - [type ImageSpec](#ImageSpec)
  - [type ImageStatus](#ImageStatus)
  - [type Kernel](#Kernel)
  - [type KernelSpec](#KernelSpec)
  - [type KernelStatus](#KernelStatus)
  - [type OCIImageSource](#OCIImageSource)
  - [type Pool](#Pool)
  - [type PoolDevice](#PoolDevice)
  - [type PoolDeviceType](#PoolDeviceType)
  - [type PoolSpec](#PoolSpec)
  - [type PoolStatus](#PoolStatus)
  - [type Runtime](#Runtime)
  - [type SSH](#SSH)
      - [func (s \*SSH) MarshalJSON() (\[\]byte,
        error)](#SSH.MarshalJSON)
      - [func (s \*SSH) UnmarshalJSON(b \[\]byte)
        error](#SSH.UnmarshalJSON)
  - [type VM](#VM)
  - [type VMImageSpec](#VMImageSpec)
  - [type VMKernelSpec](#VMKernelSpec)
  - [type VMNetworkSpec](#VMNetworkSpec)
  - [type VMSpec](#VMSpec)
  - [type VMStatus](#VMStatus)
  - [type VMStorageSpec](#VMStorageSpec)
  - [type Volume](#Volume)
  - [type VolumeMount](#VolumeMount)

#### <a name="pkg-files">Package files</a>

[defaults.go](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha2/defaults.go)
[doc.go](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha2/doc.go)
[json.go](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha2/json.go)
[register.go](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha2/register.go)
[types.go](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha2/types.go)

## <a name="pkg-constants">Constants</a>

``` go
const (
    KindImage  runtime.Kind = "Image"
    KindKernel runtime.Kind = "Kernel"
    KindVM     runtime.Kind = "VM"
)
```

``` go
const (
    // GroupName is the group name use in this package
    GroupName = "ignite.weave.works"
)
```

## <a name="pkg-variables">Variables</a>

``` go
var (
    // SchemeBuilder the schema builder
    SchemeBuilder = runtime.NewSchemeBuilder(
        addKnownTypes,
        addDefaultingFuncs,
    )

    AddToScheme = localSchemeBuilder.AddToScheme
)
```

``` go
var SchemeGroupVersion = schema.GroupVersion{
    Group:   GroupName,
    Version: "v1alpha2",
}
```

SchemeGroupVersion is group version used to register these objects

## <a name="SetDefaults_PoolSpec">func</a> [SetDefaults\_PoolSpec](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha2/defaults.go?s=263:303#L13)

``` go
func SetDefaults_PoolSpec(obj *PoolSpec)
```

## <a name="SetDefaults_VMKernelSpec">func</a> [SetDefaults\_VMKernelSpec](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha2/defaults.go?s=1145:1193#L49)

``` go
func SetDefaults_VMKernelSpec(obj *VMKernelSpec)
```

## <a name="SetDefaults_VMSpec">func</a> [SetDefaults\_VMSpec](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha2/defaults.go?s=829:865#L35)

``` go
func SetDefaults_VMSpec(obj *VMSpec)
```

## <a name="BlockDeviceVolume">type</a> [BlockDeviceVolume](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha2/types.go?s=7433:7493#L195)

``` go
type BlockDeviceVolume struct {
    Path string `json:"path"`
}
```

BlockDeviceVolume defines a block device on the host

## <a name="FileMapping">type</a> [FileMapping](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha2/types.go?s=7728:7823#L206)

``` go
type FileMapping struct {
    HostPath string `json:"hostPath"`
    VMPath   string `json:"vmPath"`
}
```

FileMapping defines mappings between files on the host and VM

## <a name="Image">type</a> [Image](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha2/types.go?s=395:868#L16)

``` go
type Image struct {
    runtime.TypeMeta `json:",inline"`
    // runtime.ObjectMeta is also embedded into the struct, and defines the human-readable name, and the machine-readable ID
    // Name is available at the .metadata.name JSON path
    // ID is available at the .metadata.uid JSON path (the Go type is k8s.io/apimachinery/pkg/types.UID, which is only a typed string)
    runtime.ObjectMeta `json:"metadata"`

    Spec   ImageSpec   `json:"spec"`
    Status ImageStatus `json:"status"`
}
```

Image represents a cached OCI image ready to be used with Ignite
+k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

## <a name="ImageSpec">type</a> [ImageSpec](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha2/types.go?s=916:976#L28)

``` go
type ImageSpec struct {
    OCI meta.OCIImageRef `json:"oci"`
}
```

ImageSpec declares what the image contains

## <a name="ImageStatus">type</a> [ImageStatus](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha2/types.go?s=1364:1513#L42)

``` go
type ImageStatus struct {
    // OCISource contains the information about how this OCI image was imported
    OCISource OCIImageSource `json:"ociSource"`
}
```

ImageStatus defines the status of the image

## <a name="Kernel">type</a> [Kernel](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha2/types.go?s=3898:4374#L106)

``` go
type Kernel struct {
    runtime.TypeMeta `json:",inline"`
    // runtime.ObjectMeta is also embedded into the struct, and defines the human-readable name, and the machine-readable ID
    // Name is available at the .metadata.name JSON path
    // ID is available at the .metadata.uid JSON path (the Go type is k8s.io/apimachinery/pkg/types.UID, which is only a typed string)
    runtime.ObjectMeta `json:"metadata"`

    Spec   KernelSpec   `json:"spec"`
    Status KernelStatus `json:"status"`
}
```

Kernel is a serializable object that caches information about imported
kernels This file is stored in
/var/lib/firecracker/kernels/{oci-image-digest}/metadata.json
+k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

## <a name="KernelSpec">type</a> [KernelSpec](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha2/types.go?s=4427:4593#L118)

``` go
type KernelSpec struct {
    OCI meta.OCIImageRef `json:"oci"`
}
```

KernelSpec describes the properties of a kernel

## <a name="KernelStatus">type</a> [KernelStatus](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha2/types.go?s=4644:4760#L125)

``` go
type KernelStatus struct {
    Version   string         `json:"version"`
    OCISource OCIImageSource `json:"ociSource"`
}
```

KernelStatus describes the status of a kernel

## <a name="OCIImageSource">type</a> [OCIImageSource](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha2/types.go?s=1083:1315#L34)

``` go
type OCIImageSource struct {
    // ID defines the source's content ID (e.g. the canonical OCI path or Docker image ID)
    ID *meta.OCIContentID `json:"id"`
    // Size defines the size of the source in bytes
    Size meta.Size `json:"size"`
}
```

OCIImageSource specifies how the OCI image was imported. It is the
status variant of OCIImageClaim

## <a name="Pool">type</a> [Pool](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha2/types.go?s=1788:1974#L51)

``` go
type Pool struct {
    runtime.TypeMeta `json:",inline"`

    Spec   PoolSpec   `json:"spec"`
    Status PoolStatus `json:"status"`
}
```

Pool defines device mapper pool database This file is managed by the
snapshotter part of Ignite, and the file (existing as a singleton) is
present at /var/lib/firecracker/snapshotter/pool.json
+k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

## <a name="PoolDevice">type</a> [PoolDevice](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha2/types.go?s=3265:3655#L93)

``` go
type PoolDevice struct {
    Size   meta.Size `json:"size"`
    Parent meta.DMID `json:"parent"`
    // Type specifies the type of the contents of the device
    Type PoolDeviceType `json:"type"`
    // MetadataPath points to the JSON/YAML file with metadata about this device
    // This is most often of the format /var/lib/firecracker/{type}/{id}/metadata.json
    MetadataPath string `json:"metadataPath"`
}
```

PoolDevice defines one device in the pool

## <a name="PoolDeviceType">type</a> [PoolDeviceType](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha2/types.go?s=2994:3020#L83)

``` go
type PoolDeviceType string
```

``` go
const (
    PoolDeviceTypeImage  PoolDeviceType = "Image"
    PoolDeviceTypeResize PoolDeviceType = "Resize"
    PoolDeviceTypeKernel PoolDeviceType = "Kernel"
    PoolDeviceTypeVM     PoolDeviceType = "VM"
)
```

## <a name="PoolSpec">type</a> [PoolSpec](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha2/types.go?s=2021:2734#L61)

``` go
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
```

PoolSpec defines the Pool’s specification

## <a name="PoolStatus">type</a> [PoolStatus](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha2/types.go?s=2784:2992#L77)

``` go
type PoolStatus struct {
    // The Devices array needs to contain pointers to accommodate "holes" in the mapping
    // Where devices have been deleted, the pointer is nil
    Devices []*PoolDevice `json:"devices"`
}
```

PoolStatus defines the Pool’s current status

## <a name="Runtime">type</a> [Runtime](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha2/types.go?s=8172:8218#L221)

``` go
type Runtime struct {
    ID string `json:"id"`
}
```

Runtime specifies the VM’s runtime information

## <a name="SSH">type</a> [SSH](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha2/types.go?s=8043:8120#L215)

``` go
type SSH struct {
    Generate  bool   `json:"-"`
    PublicKey string `json:"-"`
}
```

SSH specifies different ways to connect via SSH to the VM SSH uses a
custom marshaller/unmarshaller. If generate is true, it marshals to true
(a JSON bool). If PublicKey is set, it marshals to that string.

### <a name="SSH.MarshalJSON">func</a> (\*SSH) [MarshalJSON](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha2/json.go?s=117:160#L9)

``` go
func (s *SSH) MarshalJSON() ([]byte, error)
```

### <a name="SSH.UnmarshalJSON">func</a> (\*SSH) [UnmarshalJSON](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha2/json.go?s=308:351#L21)

``` go
func (s *SSH) UnmarshalJSON(b []byte) error
```

## <a name="VM">type</a> [VM](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha2/types.go?s=4962:5426#L133)

``` go
type VM struct {
    runtime.TypeMeta `json:",inline"`
    // runtime.ObjectMeta is also embedded into the struct, and defines the human-readable name, and the machine-readable ID
    // Name is available at the .metadata.name JSON path
    // ID is available at the .metadata.uid JSON path (the Go type is k8s.io/apimachinery/pkg/types.UID, which is only a typed string)
    runtime.ObjectMeta `json:"metadata"`

    Spec   VMSpec   `json:"spec"`
    Status VMStatus `json:"status"`
}
```

VM represents a virtual machine run by Firecracker These files are
stored in /var/lib/firecracker/vm/{vm-id}/metadata.json
+k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

## <a name="VMImageSpec">type</a> [VMImageSpec](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha2/types.go?s=6733:6795#L169)

``` go
type VMImageSpec struct {
    OCI meta.OCIImageRef `json:"oci"`
}
```

## <a name="VMKernelSpec">type</a> [VMKernelSpec](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha2/types.go?s=6797:6917#L173)

``` go
type VMKernelSpec struct {
    OCI     meta.OCIImageRef `json:"oci"`
    CmdLine string           `json:"cmdLine,omitempty"`
}
```

## <a name="VMNetworkSpec">type</a> [VMNetworkSpec](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha2/types.go?s=6919:6998#L178)

``` go
type VMNetworkSpec struct {
    Ports meta.PortMappings `json:"ports,omitempty"`
}
```

## <a name="VMSpec">type</a> [VMSpec](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha2/types.go?s=5474:6731#L145)

``` go
type VMSpec struct {
    Image    VMImageSpec  `json:"image"`
    Kernel   VMKernelSpec `json:"kernel"`
    CPUs     uint64       `json:"cpus"`
    Memory   meta.Size    `json:"memory"`
    DiskSize meta.Size    `json:"diskSize"`
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
```

VMSpec describes the configuration of a VM

## <a name="VMStatus">type</a> [VMStatus](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha2/types.go?s=8259:8598#L226)

``` go
type VMStatus struct {
    Running     bool             `json:"running"`
    Runtime     *Runtime         `json:"runtime,omitempty"`
    StartTime   *runtime.Time    `json:"startTime,omitempty"`
    IPAddresses meta.IPAddresses `json:"ipAddresses,omitempty"`
    Image       OCIImageSource   `json:"image"`
    Kernel      OCIImageSource   `json:"kernel"`
}
```

VMStatus defines the status of a VM

## <a name="VMStorageSpec">type</a> [VMStorageSpec](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha2/types.go?s=7059:7203#L183)

``` go
type VMStorageSpec struct {
    Volumes      []Volume      `json:"volumes,omitempty"`
    VolumeMounts []VolumeMount `json:"volumeMounts,omitempty"`
}
```

VMStorageSpec defines the VM’s Volumes and VolumeMounts

## <a name="Volume">type</a> [Volume](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha2/types.go?s=7244:7375#L189)

``` go
type Volume struct {
    Name        string             `json:"name"`
    BlockDevice *BlockDeviceVolume `json:"blockDevice,omitempty"`
}
```

Volume defines named storage volume

## <a name="VolumeMount">type</a> [VolumeMount](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha2/types.go?s=7565:7661#L200)

``` go
type VolumeMount struct {
    Name      string `json:"name"`
    MountPath string `json:"mountPath"`
}
```

VolumeMount defines the mount point for a named volume inside a VM

-----

Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
