

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
* [func SetDefaults_OCIImageClaim(obj *OCIImageClaim)](#SetDefaults_OCIImageClaim)
* [func SetDefaults_PoolSpec(obj *PoolSpec)](#SetDefaults_PoolSpec)
* [func SetDefaults_VMKernelSpec(obj *VMKernelSpec)](#SetDefaults_VMKernelSpec)
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
* [type OCIImageClaim](#OCIImageClaim)
* [type OCIImageSource](#OCIImageSource)
* [type Pool](#Pool)
* [type PoolDevice](#PoolDevice)
* [type PoolDeviceType](#PoolDeviceType)
* [type PoolSpec](#PoolSpec)
* [type PoolStatus](#PoolStatus)
* [type SSH](#SSH)
* [type VM](#VM)
* [type VMImageSource](#VMImageSource)
* [type VMImageSpec](#VMImageSpec)
* [type VMKernelSpec](#VMKernelSpec)
* [type VMSpec](#VMSpec)
* [type VMState](#VMState)
* [type VMStatus](#VMStatus)


#### <a name="pkg-files">Package files</a>
[defaults.go](/src/target/defaults.go) [doc.go](/src/target/doc.go) [helpers.go](/src/target/helpers.go) [register.go](/src/target/register.go) [types.go](/src/target/types.go) 


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



## <a name="SetDefaults_OCIImageClaim">func</a> [SetDefaults_OCIImageClaim](/src/target/defaults.go?s=263:313#L13)
``` go
func SetDefaults_OCIImageClaim(obj *OCIImageClaim)
```


## <a name="SetDefaults_PoolSpec">func</a> [SetDefaults_PoolSpec](/src/target/defaults.go?s=353:393#L17)
``` go
func SetDefaults_PoolSpec(obj *PoolSpec)
```


## <a name="SetDefaults_VMKernelSpec">func</a> [SetDefaults_VMKernelSpec](/src/target/defaults.go?s=1315:1363#L57)
``` go
func SetDefaults_VMKernelSpec(obj *VMKernelSpec)
```


## <a name="SetDefaults_VMSpec">func</a> [SetDefaults_VMSpec](/src/target/defaults.go?s=919:955#L39)
``` go
func SetDefaults_VMSpec(obj *VMSpec)
```


## <a name="SetDefaults_VMStatus">func</a> [SetDefaults_VMStatus](/src/target/defaults.go?s=1449:1489#L63)
``` go
func SetDefaults_VMStatus(obj *VMStatus)
```


## <a name="ValidateNetworkMode">func</a> [ValidateNetworkMode](/src/target/helpers.go?s=317:365#L15)
``` go
func ValidateNetworkMode(mode NetworkMode) error
```
ValidateNetworkMode validates the network mode
TODO: This should move into a dedicated validation package




## <a name="FileMapping">type</a> [FileMapping](/src/target/types.go?s=7370:7465#L192)
``` go
type FileMapping struct {
    HostPath string `json:"hostPath"`
    VMPath   string `json:"vmPath"`
}

```
FileMapping defines mappings between files on the host and VM










## <a name="Image">type</a> [Image](/src/target/types.go?s=229:693#L9)
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










## <a name="ImageSourceType">type</a> [ImageSourceType](/src/target/types.go?s=882:909#L26)
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









## <a name="ImageSpec">type</a> [ImageSpec](/src/target/types.go?s=741:808#L21)
``` go
type ImageSpec struct {
    OCIClaim OCIImageClaim `json:"ociClaim"`
}

```
ImageSpec declares what the image contains










## <a name="ImageStatus">type</a> [ImageStatus](/src/target/types.go?s=2197:2346#L60)
``` go
type ImageStatus struct {
    // OCISource contains the information about how this OCI image was imported
    OCISource OCIImageSource `json:"ociSource"`
}

```
ImageStatus defines the status of the image










## <a name="Kernel">type</a> [Kernel](/src/target/types.go?s=4725:5192#L124)
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










## <a name="KernelSpec">type</a> [KernelSpec](/src/target/types.go?s=5245:5459#L136)
``` go
type KernelSpec struct {
    Version  string        `json:"version"`
    OCIClaim OCIImageClaim `json:"ociClaim"`
}

```
KernelSpec describes the properties of a kernel










## <a name="KernelStatus">type</a> [KernelStatus](/src/target/types.go?s=5510:5583#L144)
``` go
type KernelStatus struct {
    OCISource OCIImageSource `json:"ociSource"`
}

```
KernelStatus describes the status of a kernel










## <a name="NetworkMode">type</a> [NetworkMode](/src/target/types.go?s=7651:7674#L203)
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






### <a name="GetNetworkModes">func</a> [GetNetworkModes](/src/target/helpers.go?s=92:128#L6)
``` go
func GetNetworkModes() []NetworkMode
```
GetNetworkModes gets the list of available network modes





## <a name="OCIImageClaim">type</a> [OCIImageClaim](/src/target/types.go?s=1105:1513#L34)
``` go
type OCIImageClaim struct {
    // Type defines how the image should be imported
    Type ImageSourceType `json:"type"`
    // Ref defines the reference to use when talking to the backend.
    // This is most commonly the image name, followed by a tag.
    // Other supported ways are $registry/$user/$image@sha256:$digest
    // This ref is also used as ObjectMeta.Name for kinds Images and Kernels
    Ref string `json:"ref"`
}

```
OCIImageClaim defines a claim for importing an OCI image










## <a name="OCIImageSource">type</a> [OCIImageSource](/src/target/types.go?s=1620:2148#L46)
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










## <a name="Pool">type</a> [Pool](/src/target/types.go?s=2621:2801#L69)
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










## <a name="PoolDevice">type</a> [PoolDevice](/src/target/types.go?s=4092:4482#L111)
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










## <a name="PoolDeviceType">type</a> [PoolDeviceType](/src/target/types.go?s=3821:3847#L101)
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









## <a name="PoolSpec">type</a> [PoolSpec](/src/target/types.go?s=2848:3561#L79)
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










## <a name="PoolStatus">type</a> [PoolStatus](/src/target/types.go?s=3611:3819#L95)
``` go
type PoolStatus struct {
    // The Devices array needs to contain pointers to accommodate "holes" in the mapping
    // Where devices have been deleted, the pointer is nil
    Devices []*PoolDevice `json:"devices"`
}

```
PoolStatus defines the Pool's current status










## <a name="SSH">type</a> [SSH](/src/target/types.go?s=7528:7594#L198)
``` go
type SSH struct {
    PublicKey string `json:"publicKey,omitempty"`
}

```
SSH specifies different ways to connect via SSH to the VM










## <a name="VM">type</a> [VM](/src/target/types.go?s=5785:6240#L151)
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










## <a name="VMImageSource">type</a> [VMImageSource](/src/target/types.go?s=8566:8674#L232)
``` go
type VMImageSource struct {
    OCIImageSource `json:",inline"`
    UID            meta.UID `json:"internalUID"`
}

```
VMImageSource is a temporary wrapper around OCIImageSource to allow
passing the old UID for internal purposes










## <a name="VMImageSpec">type</a> [VMImageSpec](/src/target/types.go?s=7108:7178#L182)
``` go
type VMImageSpec struct {
    OCIClaim *OCIImageClaim `json:"ociClaim"`
}

```









## <a name="VMKernelSpec">type</a> [VMKernelSpec](/src/target/types.go?s=7180:7303#L186)
``` go
type VMKernelSpec struct {
    OCIClaim *OCIImageClaim `json:"ociClaim"`
    CmdLine  string         `json:"cmdLine,omitempty"`
}

```









## <a name="VMSpec">type</a> [VMSpec](/src/target/types.go?s=6288:7106#L163)
``` go
type VMSpec struct {
    Image       VMImageSpec       `json:"image"`
    Kernel      VMKernelSpec      `json:"kernel"`
    CPUs        uint64            `json:"cpus"`
    Memory      meta.Size         `json:"memory"`
    DiskSize    meta.Size         `json:"diskSize"`
    NetworkMode NetworkMode       `json:"networkMode"`
    Ports       meta.PortMappings `json:"ports,omitempty"`
    // This will be done at either "ignite start" or "ignite create" time
    // TODO: We might to revisit this later
    CopyFiles []FileMapping `json:"copyFiles,omitempty"`
    // SSH specifies how the SSH setup should be done
    // SSH appends to CopyFiles when active
    // nil here means "don't do anything special"
    // An empty struct means "generate a new SSH key and copy it in"
    // Specifying a path means "use this public key"
    SSH *SSH `json:"ssh,omitempty"`
}

```
VMSpec describes the configuration of a VM










## <a name="VMState">type</a> [VMState](/src/target/types.go?s=8048:8067#L214)
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









## <a name="VMStatus">type</a> [VMStatus](/src/target/types.go?s=8227:8448#L223)
``` go
type VMStatus struct {
    State       VMState          `json:"state"`
    IPAddresses meta.IPAddresses `json:"ipAddresses,omitempty"`
    Image       VMImageSource    `json:"image"`
    Kernel      VMImageSource    `json:"kernel"`
}

```
VMStatus defines the status of a VM














- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
