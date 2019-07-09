

# v1alpha1
`import "/go/src/github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>
+k8s:deepcopy-gen=package
+k8s:defaulter-gen=TypeMeta
+k8s:openapi-gen=true




## <a name="pkg-index">Index</a>
* [Constants](#pkg-constants)
* [Variables](#pkg-variables)
* [func SetDefaults_Image(obj *Image)](#SetDefaults_Image)
* [func SetDefaults_Kernel(obj *Kernel)](#SetDefaults_Kernel)
* [func SetDefaults_OCIImageClaim(obj *OCIImageClaim)](#SetDefaults_OCIImageClaim)
* [func SetDefaults_PoolSpec(obj *PoolSpec)](#SetDefaults_PoolSpec)
* [func SetDefaults_VM(obj *VM)](#SetDefaults_VM)
* [func SetDefaults_VMKernelSpec(obj *VMKernelSpec)](#SetDefaults_VMKernelSpec)
* [func SetDefaults_VMNetworkSpec(obj *VMNetworkSpec)](#SetDefaults_VMNetworkSpec)
* [func SetDefaults_VMSpec(obj *VMSpec)](#SetDefaults_VMSpec)
* [func SetDefaults_VMStatus(obj *VMStatus)](#SetDefaults_VMStatus)
* [func ValidateNetworkMode(mode NetworkMode) error](#ValidateNetworkMode)
* [type FileMapping](#FileMapping)
* [type Image](#Image)
* [type ImageSourceType](#ImageSourceType)
* [type ImageSpec](#ImageSpec)
* [type ImageStatus](#ImageStatus)
* [type Kernel](#Kernel)
* [type KernelSpec](#KernelSpec)
* [type KernelStatus](#KernelStatus)
* [type NetworkMode](#NetworkMode)
  * [func GetNetworkModes() []NetworkMode](#GetNetworkModes)
  * [func (nm NetworkMode) String() string](#NetworkMode.String)
* [type OCIImageClaim](#OCIImageClaim)
* [type OCIImageSource](#OCIImageSource)
* [type Pool](#Pool)
* [type PoolDevice](#PoolDevice)
* [type PoolDeviceType](#PoolDeviceType)
* [type PoolSpec](#PoolSpec)
* [type PoolStatus](#PoolStatus)
* [type SSH](#SSH)
  * [func (s *SSH) MarshalJSON() ([]byte, error)](#SSH.MarshalJSON)
  * [func (s *SSH) UnmarshalJSON(b []byte) error](#SSH.UnmarshalJSON)
* [type VM](#VM)
  * [func (vm *VM) SetImage(image *Image)](#VM.SetImage)
  * [func (vm *VM) SetKernel(kernel *Kernel)](#VM.SetKernel)
* [type VMImageSpec](#VMImageSpec)
* [type VMKernelSpec](#VMKernelSpec)
* [type VMNetworkSpec](#VMNetworkSpec)
* [type VMSpec](#VMSpec)
* [type VMState](#VMState)
* [type VMStatus](#VMStatus)


#### <a name="pkg-files">Package files</a>
[defaults.go](/pkg/apis/ignite/v1alpha1/defaults.go) [doc.go](/pkg/apis/ignite/v1alpha1/doc.go) [helpers.go](/pkg/apis/ignite/v1alpha1/helpers.go) [json.go](/pkg/apis/ignite/v1alpha1/json.go) [register.go](/pkg/apis/ignite/v1alpha1/register.go) [types.go](/pkg/apis/ignite/v1alpha1/types.go) 


## <a name="pkg-constants">Constants</a>
``` go
const (
    // GroupName is the group name use in this package
    GroupName = "ignite.weave.works"

    // VMKind returns the kind for the VM API type
    VMKind = "VM"
    // KernelKind returns the kind for the Kernel API type
    KernelKind = "Kernel"
    // PoolKind returns the kind for the Pool API type
    PoolKind = "Pool"
    // ImageKind returns the kind for the Image API type
    ImageKind = "Image"
)
```
``` go
const (
    KindImage  meta.Kind = "Image"
    KindKernel meta.Kind = "Kernel"
    KindVM     meta.Kind = "VM"
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
    Version: "v1alpha1",
}
```
SchemeGroupVersion is group version used to register these objects



## <a name="SetDefaults_Image">func</a> [SetDefaults_Image](/pkg/apis/ignite/v1alpha1/defaults.go?s=2218:2252#L91)
``` go
func SetDefaults_Image(obj *Image)
```


## <a name="SetDefaults_Kernel">func</a> [SetDefaults_Kernel](/pkg/apis/ignite/v1alpha1/defaults.go?s=2276:2312#L95)
``` go
func SetDefaults_Kernel(obj *Kernel)
```


## <a name="SetDefaults_OCIImageClaim">func</a> [SetDefaults_OCIImageClaim](/pkg/apis/ignite/v1alpha1/defaults.go?s=275:325#L15)
``` go
func SetDefaults_OCIImageClaim(obj *OCIImageClaim)
```


## <a name="SetDefaults_PoolSpec">func</a> [SetDefaults_PoolSpec](/pkg/apis/ignite/v1alpha1/defaults.go?s=365:405#L19)
``` go
func SetDefaults_PoolSpec(obj *PoolSpec)
```


## <a name="SetDefaults_VM">func</a> [SetDefaults_VM](/pkg/apis/ignite/v1alpha1/defaults.go?s=2166:2194#L87)
``` go
func SetDefaults_VM(obj *VM)
```
TODO: Temporary hacks to populate TypeMeta until we get the generator working



## <a name="SetDefaults_VMKernelSpec">func</a> [SetDefaults_VMKernelSpec](/pkg/apis/ignite/v1alpha1/defaults.go?s=1247:1295#L55)
``` go
func SetDefaults_VMKernelSpec(obj *VMKernelSpec)
```


## <a name="SetDefaults_VMNetworkSpec">func</a> [SetDefaults_VMNetworkSpec](/pkg/apis/ignite/v1alpha1/defaults.go?s=1532:1582#L66)
``` go
func SetDefaults_VMNetworkSpec(obj *VMNetworkSpec)
```


## <a name="SetDefaults_VMSpec">func</a> [SetDefaults_VMSpec](/pkg/apis/ignite/v1alpha1/defaults.go?s=931:967#L41)
``` go
func SetDefaults_VMSpec(obj *VMSpec)
```


## <a name="SetDefaults_VMStatus">func</a> [SetDefaults_VMStatus](/pkg/apis/ignite/v1alpha1/defaults.go?s=1653:1693#L72)
``` go
func SetDefaults_VMStatus(obj *VMStatus)
```


## <a name="ValidateNetworkMode">func</a> [ValidateNetworkMode](/pkg/apis/ignite/v1alpha1/helpers.go?s=317:365#L15)
``` go
func ValidateNetworkMode(mode NetworkMode) error
```
ValidateNetworkMode validates the network mode
TODO: This should move into a dedicated validation package




## <a name="FileMapping">type</a> [FileMapping](/pkg/apis/ignite/v1alpha1/types.go?s=7672:7767#L204)
``` go
type FileMapping struct {
    HostPath string `json:"hostPath"`
    VMPath   string `json:"vmPath"`
}

```
FileMapping defines mappings between files on the host and VM










## <a name="Image">type</a> [Image](/pkg/apis/ignite/v1alpha1/types.go?s=334:798#L15)
``` go
type Image struct {
    meta.TypeMeta `json:",inline"`
    // meta.ObjectMeta is also embedded into the struct, and defines the human-readable name, and the machine-readable ID
    // Name is available at the .metadata.name JSON path
    // ID is available at the .metadata.uid JSON path (the Go type is k8s.io/apimachinery/pkg/types.UID, which is only a typed string)
    meta.ObjectMeta `json:"metadata"`

    Spec   ImageSpec   `json:"spec"`
    Status ImageStatus `json:"status"`
}

```
Image represents a cached OCI image ready to be used with Ignite
+k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object










## <a name="ImageSourceType">type</a> [ImageSourceType](/pkg/apis/ignite/v1alpha1/types.go?s=987:1014#L32)
``` go
type ImageSourceType string
```
ImageSourceType is an enum of different supported Image Source Types


``` go
const (
    // ImageSourceTypeDocker defines that the image is imported from Docker
    ImageSourceTypeDocker ImageSourceType = "Docker"
)
```









## <a name="ImageSpec">type</a> [ImageSpec](/pkg/apis/ignite/v1alpha1/types.go?s=846:913#L27)
``` go
type ImageSpec struct {
    OCIClaim OCIImageClaim `json:"ociClaim"`
}

```
ImageSpec declares what the image contains










## <a name="ImageStatus">type</a> [ImageStatus](/pkg/apis/ignite/v1alpha1/types.go?s=2312:2461#L66)
``` go
type ImageStatus struct {
    // OCISource contains the information about how this OCI image was imported
    OCISource OCIImageSource `json:"ociSource"`
}

```
ImageStatus defines the status of the image










## <a name="Kernel">type</a> [Kernel](/pkg/apis/ignite/v1alpha1/types.go?s=4840:5307#L130)
``` go
type Kernel struct {
    meta.TypeMeta `json:",inline"`
    // meta.ObjectMeta is also embedded into the struct, and defines the human-readable name, and the machine-readable ID
    // Name is available at the .metadata.name JSON path
    // ID is available at the .metadata.uid JSON path (the Go type is k8s.io/apimachinery/pkg/types.UID, which is only a typed string)
    meta.ObjectMeta `json:"metadata"`

    Spec   KernelSpec   `json:"spec"`
    Status KernelStatus `json:"status"`
}

```
Kernel is a serializable object that caches information about imported kernels
This file is stored in /var/lib/firecracker/kernels/{oci-image-digest}/metadata.json
+k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object










## <a name="KernelSpec">type</a> [KernelSpec](/pkg/apis/ignite/v1alpha1/types.go?s=5360:5533#L142)
``` go
type KernelSpec struct {
    OCIClaim OCIImageClaim `json:"ociClaim"`
}

```
KernelSpec describes the properties of a kernel










## <a name="KernelStatus">type</a> [KernelStatus](/pkg/apis/ignite/v1alpha1/types.go?s=5584:5700#L149)
``` go
type KernelStatus struct {
    Version   string         `json:"version"`
    OCISource OCIImageSource `json:"ociSource"`
}

```
KernelStatus describes the status of a kernel










## <a name="NetworkMode">type</a> [NetworkMode](/pkg/apis/ignite/v1alpha1/types.go?s=8121:8144#L219)
``` go
type NetworkMode string
```
NetworkMode defines different states a VM can be in


``` go
const (
    // NetworkModeCNI specifies the network mode where CNI is used
    NetworkModeCNI NetworkMode = "cni"
    // NetworkModeDockerBridge specifies the default docker bridge network is used
    NetworkModeDockerBridge NetworkMode = "docker-bridge"
)
```






### <a name="GetNetworkModes">func</a> [GetNetworkModes](/pkg/apis/ignite/v1alpha1/helpers.go?s=92:128#L6)
``` go
func GetNetworkModes() []NetworkMode
```
GetNetworkModes gets the list of available network modes





### <a name="NetworkMode.String">func</a> (NetworkMode) [String](/pkg/apis/ignite/v1alpha1/types.go?s=8146:8183#L221)
``` go
func (nm NetworkMode) String() string
```



## <a name="OCIImageClaim">type</a> [OCIImageClaim](/pkg/apis/ignite/v1alpha1/types.go?s=1210:1628#L40)
``` go
type OCIImageClaim struct {
    // Type defines how the image should be imported
    Type ImageSourceType `json:"type"`
    // Ref defines the reference to use when talking to the backend.
    // This is most commonly the image name, followed by a tag.
    // Other supported ways are $registry/$user/$image@sha256:$digest
    // This ref is also used as ObjectMeta.Name for kinds Images and Kernels
    Ref meta.OCIImageRef `json:"ref"`
}

```
OCIImageClaim defines a claim for importing an OCI image










## <a name="OCIImageSource">type</a> [OCIImageSource](/pkg/apis/ignite/v1alpha1/types.go?s=1735:2263#L52)
``` go
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

```
OCIImageSource specifies how the OCI image was imported.
It is the status variant of OCIImageClaim










## <a name="Pool">type</a> [Pool](/pkg/apis/ignite/v1alpha1/types.go?s=2736:2916#L75)
``` go
type Pool struct {
    meta.TypeMeta `json:",inline"`

    Spec   PoolSpec   `json:"spec"`
    Status PoolStatus `json:"status"`
}

```
Pool defines device mapper pool database
This file is managed by the snapshotter part of Ignite, and the file (existing as a singleton)
is present at /var/lib/firecracker/snapshotter/pool.json
+k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object










## <a name="PoolDevice">type</a> [PoolDevice](/pkg/apis/ignite/v1alpha1/types.go?s=4207:4597#L117)
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










## <a name="PoolDeviceType">type</a> [PoolDeviceType](/pkg/apis/ignite/v1alpha1/types.go?s=3936:3962#L107)
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









## <a name="PoolSpec">type</a> [PoolSpec](/pkg/apis/ignite/v1alpha1/types.go?s=2963:3676#L85)
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
PoolSpec defines the Pool's specification










## <a name="PoolStatus">type</a> [PoolStatus](/pkg/apis/ignite/v1alpha1/types.go?s=3726:3934#L101)
``` go
type PoolStatus struct {
    // The Devices array needs to contain pointers to accommodate "holes" in the mapping
    // Where devices have been deleted, the pointer is nil
    Devices []*PoolDevice `json:"devices"`
}

```
PoolStatus defines the Pool's current status










## <a name="SSH">type</a> [SSH](/pkg/apis/ignite/v1alpha1/types.go?s=7987:8064#L213)
``` go
type SSH struct {
    Generate  bool   `json:"-"`
    PublicKey string `json:"-"`
}

```
SSH specifies different ways to connect via SSH to the VM
SSH uses a custom marshaller/unmarshaller. If generate is true,
it marshals to true (a JSON bool). If PublicKey is set, it marshals
to that string.










### <a name="SSH.MarshalJSON">func</a> (\*SSH) [MarshalJSON](/pkg/apis/ignite/v1alpha1/json.go?s=117:160#L9)
``` go
func (s *SSH) MarshalJSON() ([]byte, error)
```



### <a name="SSH.UnmarshalJSON">func</a> (\*SSH) [UnmarshalJSON](/pkg/apis/ignite/v1alpha1/json.go?s=306:349#L19)
``` go
func (s *SSH) UnmarshalJSON(b []byte) error
```



## <a name="VM">type</a> [VM](/pkg/apis/ignite/v1alpha1/types.go?s=5902:6357#L157)
``` go
type VM struct {
    meta.TypeMeta `json:",inline"`
    // meta.ObjectMeta is also embedded into the struct, and defines the human-readable name, and the machine-readable ID
    // Name is available at the .metadata.name JSON path
    // ID is available at the .metadata.uid JSON path (the Go type is k8s.io/apimachinery/pkg/types.UID, which is only a typed string)
    meta.ObjectMeta `json:"metadata"`

    Spec   VMSpec   `json:"spec"`
    Status VMStatus `json:"status"`
}

```
VM represents a virtual machine run by Firecracker
These files are stored in /var/lib/firecracker/vm/{vm-id}/metadata.json
+k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object










### <a name="VM.SetImage">func</a> (\*VM) [SetImage](/pkg/apis/ignite/v1alpha1/helpers.go?s=658:694#L30)
``` go
func (vm *VM) SetImage(image *Image)
```
SetImage populates relevant fields to an Image on the VM object




### <a name="VM.SetKernel">func</a> (\*VM) [SetKernel](/pkg/apis/ignite/v1alpha1/helpers.go?s=856:895#L36)
``` go
func (vm *VM) SetKernel(kernel *Kernel)
```
SetKernel populates relevant fields to a Kernel on the VM object




## <a name="VMImageSpec">type</a> [VMImageSpec](/pkg/apis/ignite/v1alpha1/types.go?s=7293:7362#L189)
``` go
type VMImageSpec struct {
    OCIClaim OCIImageClaim `json:"ociClaim"`
}

```









## <a name="VMKernelSpec">type</a> [VMKernelSpec](/pkg/apis/ignite/v1alpha1/types.go?s=7364:7485#L193)
``` go
type VMKernelSpec struct {
    OCIClaim OCIImageClaim `json:"ociClaim"`
    CmdLine  string        `json:"cmdLine,omitempty"`
}

```









## <a name="VMNetworkSpec">type</a> [VMNetworkSpec](/pkg/apis/ignite/v1alpha1/types.go?s=7487:7605#L198)
``` go
type VMNetworkSpec struct {
    Mode  NetworkMode       `json:"mode"`
    Ports meta.PortMappings `json:"ports,omitempty"`
}

```









## <a name="VMSpec">type</a> [VMSpec](/pkg/apis/ignite/v1alpha1/types.go?s=6405:7291#L169)
``` go
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

```
VMSpec describes the configuration of a VM










## <a name="VMState">type</a> [VMState](/pkg/apis/ignite/v1alpha1/types.go?s=8580:8599#L234)
``` go
type VMState string
```
VMState defines different states a VM can be in


``` go
const (
    VMStateCreated VMState = "Created"
    VMStateRunning VMState = "Running"
    VMStateStopped VMState = "Stopped"
)
```









## <a name="VMStatus">type</a> [VMStatus](/pkg/apis/ignite/v1alpha1/types.go?s=8759:8980#L243)
``` go
type VMStatus struct {
    State       VMState          `json:"state"`
    IPAddresses meta.IPAddresses `json:"ipAddresses,omitempty"`
    Image       OCIImageSource   `json:"image"`
    Kernel      OCIImageSource   `json:"kernel"`
}

```
VMStatus defines the status of a VM














- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
