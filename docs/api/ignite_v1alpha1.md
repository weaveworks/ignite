# v1alpha1

`import "github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"`

  - [Overview](#pkg-overview)
  - [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>

\+k8s:deepcopy-gen=package +k8s:defaulter-gen=TypeMeta
+k8s:openapi-gen=true
+k8s:conversion-gen=github.com/weaveworks/ignite/pkg/apis/ignite

## <a name="pkg-index">Index</a>

  - [Constants](#pkg-constants)
  - [Variables](#pkg-variables)
  - [func Convert\_ignite\_ImageSpec\_To\_v1alpha1\_ImageSpec(in
    *ignite.ImageSpec, out *ImageSpec, s conversion.Scope)
    error](#Convert_ignite_ImageSpec_To_v1alpha1_ImageSpec)
  - [func Convert\_ignite\_KernelSpec\_To\_v1alpha1\_KernelSpec(in
    *ignite.KernelSpec, out *KernelSpec, s conversion.Scope)
    error](#Convert_ignite_KernelSpec_To_v1alpha1_KernelSpec)
  - [func
    Convert\_ignite\_OCIImageSource\_To\_v1alpha1\_OCIImageSource(in
    *ignite.OCIImageSource, out *OCIImageSource, s conversion.Scope)
    error](#Convert_ignite_OCIImageSource_To_v1alpha1_OCIImageSource)
  - [func Convert\_ignite\_OCI\_To\_v1alpha1\_OCIClaim(in
    *meta.OCIImageRef, out *OCIImageClaim)
    error](#Convert_ignite_OCI_To_v1alpha1_OCIClaim)
  - [func Convert\_ignite\_VMImageSpec\_To\_v1alpha1\_VMImageSpec(in
    *ignite.VMImageSpec, out *VMImageSpec, s conversion.Scope)
    error](#Convert_ignite_VMImageSpec_To_v1alpha1_VMImageSpec)
  - [func Convert\_ignite\_VMKernelSpec\_To\_v1alpha1\_VMKernelSpec(in
    *ignite.VMKernelSpec, out *VMKernelSpec, s conversion.Scope)
    error](#Convert_ignite_VMKernelSpec_To_v1alpha1_VMKernelSpec)
  - [func Convert\_ignite\_VMSpec\_To\_v1alpha1\_VMSpec(in
    *ignite.VMSpec, out *VMSpec, s conversion.Scope)
    error](#Convert_ignite_VMSpec_To_v1alpha1_VMSpec)
  - [func Convert\_ignite\_VMStatus\_To\_v1alpha1\_VMStatus(in
    *ignite.VMStatus, out *VMStatus, s conversion.Scope)
    error](#Convert_ignite_VMStatus_To_v1alpha1_VMStatus)
  - [func Convert\_v1alpha1\_ImageSpec\_To\_ignite\_ImageSpec(in
    *ImageSpec, out *ignite.ImageSpec, s conversion.Scope)
    error](#Convert_v1alpha1_ImageSpec_To_ignite_ImageSpec)
  - [func Convert\_v1alpha1\_KernelSpec\_To\_ignite\_KernelSpec(in
    *KernelSpec, out *ignite.KernelSpec, s conversion.Scope)
    error](#Convert_v1alpha1_KernelSpec_To_ignite_KernelSpec)
  - [func Convert\_v1alpha1\_OCIClaim\_To\_ignite\_OCI(in
    *OCIImageClaim, out *meta.OCIImageRef)
    error](#Convert_v1alpha1_OCIClaim_To_ignite_OCI)
  - [func
    Convert\_v1alpha1\_OCIImageSource\_To\_ignite\_OCIImageSource(in
    *OCIImageSource, out *ignite.OCIImageSource, s conversion.Scope)
    (err
    error)](#Convert_v1alpha1_OCIImageSource_To_ignite_OCIImageSource)
  - [func Convert\_v1alpha1\_VMImageSpec\_To\_ignite\_VMImageSpec(in
    *VMImageSpec, out *ignite.VMImageSpec, s conversion.Scope)
    error](#Convert_v1alpha1_VMImageSpec_To_ignite_VMImageSpec)
  - [func Convert\_v1alpha1\_VMKernelSpec\_To\_ignite\_VMKernelSpec(in
    *VMKernelSpec, out *ignite.VMKernelSpec, s conversion.Scope)
    error](#Convert_v1alpha1_VMKernelSpec_To_ignite_VMKernelSpec)
  - [func Convert\_v1alpha1\_VMNetworkSpec\_To\_ignite\_VMNetworkSpec(in
    *VMNetworkSpec, out *ignite.VMNetworkSpec, s conversion.Scope)
    error](#Convert_v1alpha1_VMNetworkSpec_To_ignite_VMNetworkSpec)
  - [func Convert\_v1alpha1\_VMSpec\_To\_ignite\_VMSpec(in *VMSpec, out
    *ignite.VMSpec, s conversion.Scope)
    error](#Convert_v1alpha1_VMSpec_To_ignite_VMSpec)
  - [func Convert\_v1alpha1\_VMStatus\_To\_ignite\_VMStatus(in
    *VMStatus, out *ignite.VMStatus, s conversion.Scope)
    error](#Convert_v1alpha1_VMStatus_To_ignite_VMStatus)
  - [func SetDefaults\_OCIImageClaim(obj
    \*OCIImageClaim)](#SetDefaults_OCIImageClaim)
  - [func SetDefaults\_PoolSpec(obj \*PoolSpec)](#SetDefaults_PoolSpec)
  - [func SetDefaults\_VMKernelSpec(obj
    \*VMKernelSpec)](#SetDefaults_VMKernelSpec)
  - [func SetDefaults\_VMNetworkSpec(obj
    \*VMNetworkSpec)](#SetDefaults_VMNetworkSpec)
  - [func SetDefaults\_VMSpec(obj \*VMSpec)](#SetDefaults_VMSpec)
  - [func SetDefaults\_VMStatus(obj \*VMStatus)](#SetDefaults_VMStatus)
  - [type FileMapping](#FileMapping)
  - [type Image](#Image)
  - [type ImageSourceType](#ImageSourceType)
  - [type ImageSpec](#ImageSpec)
  - [type ImageStatus](#ImageStatus)
  - [type Kernel](#Kernel)
  - [type KernelSpec](#KernelSpec)
  - [type KernelStatus](#KernelStatus)
  - [type NetworkMode](#NetworkMode)
      - [func (nm NetworkMode) String() string](#NetworkMode.String)
  - [type OCIImageClaim](#OCIImageClaim)
  - [type OCIImageSource](#OCIImageSource)
  - [type Pool](#Pool)
  - [type PoolDevice](#PoolDevice)
  - [type PoolDeviceType](#PoolDeviceType)
  - [type PoolSpec](#PoolSpec)
  - [type PoolStatus](#PoolStatus)
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
  - [type VMState](#VMState)
  - [type VMStatus](#VMStatus)

#### <a name="pkg-files">Package files</a>

[conversion.go](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha1/conversion.go)
[defaults.go](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha1/defaults.go)
[doc.go](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha1/doc.go)
[json.go](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha1/json.go)
[register.go](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha1/register.go)
[types.go](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha1/types.go)

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
    Version: "v1alpha1",
}
```

SchemeGroupVersion is group version used to register these objects

## <a name="Convert_ignite_ImageSpec_To_v1alpha1_ImageSpec">func</a> [Convert\_ignite\_ImageSpec\_To\_v1alpha1\_ImageSpec](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha1/conversion.go?s=4265:4380#L107)

``` go
func Convert_ignite_ImageSpec_To_v1alpha1_ImageSpec(in *ignite.ImageSpec, out *ImageSpec, s conversion.Scope) error
```

Convert\_ignite\_ImageSpec\_To\_v1alpha1\_ImageSpec calls the
autogenerated conversion function along with custom conversion logic

## <a name="Convert_ignite_KernelSpec_To_v1alpha1_KernelSpec">func</a> [Convert\_ignite\_KernelSpec\_To\_v1alpha1\_KernelSpec](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha1/conversion.go?s=5121:5240#L125)

``` go
func Convert_ignite_KernelSpec_To_v1alpha1_KernelSpec(in *ignite.KernelSpec, out *KernelSpec, s conversion.Scope) error
```

Convert\_ignite\_KernelSpec\_To\_v1alpha1\_KernelSpec calls the
autogenerated conversion function along with custom conversion logic

## <a name="Convert_ignite_OCIImageSource_To_v1alpha1_OCIImageSource">func</a> [Convert\_ignite\_OCIImageSource\_To\_v1alpha1\_OCIImageSource](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha1/conversion.go?s=2032:2167#L50)

``` go
func Convert_ignite_OCIImageSource_To_v1alpha1_OCIImageSource(in *ignite.OCIImageSource, out *OCIImageSource, s conversion.Scope) error
```

Convert\_ignite\_OCIImageSource\_To\_v1alpha1\_OCIImageSource calls the
autogenerated conversion function along with custom conversion logic

## <a name="Convert_ignite_OCI_To_v1alpha1_OCIClaim">func</a> [Convert\_ignite\_OCI\_To\_v1alpha1\_OCIClaim](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha1/conversion.go?s=3825:3917#L96)

``` go
func Convert_ignite_OCI_To_v1alpha1_OCIClaim(in *meta.OCIImageRef, out *OCIImageClaim) error
```

## <a name="Convert_ignite_VMImageSpec_To_v1alpha1_VMImageSpec">func</a> [Convert\_ignite\_VMImageSpec\_To\_v1alpha1\_VMImageSpec](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha1/conversion.go?s=5993:6116#L143)

``` go
func Convert_ignite_VMImageSpec_To_v1alpha1_VMImageSpec(in *ignite.VMImageSpec, out *VMImageSpec, s conversion.Scope) error
```

Convert\_ignite\_VMImageSpec\_To\_v1alpha1\_VMImageSpec calls the
autogenerated conversion function along with custom conversion logic

## <a name="Convert_ignite_VMKernelSpec_To_v1alpha1_VMKernelSpec">func</a> [Convert\_ignite\_VMKernelSpec\_To\_v1alpha1\_VMKernelSpec](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha1/conversion.go?s=6881:7008#L161)

``` go
func Convert_ignite_VMKernelSpec_To_v1alpha1_VMKernelSpec(in *ignite.VMKernelSpec, out *VMKernelSpec, s conversion.Scope) error
```

Convert\_ignite\_VMKernelSpec\_To\_v1alpha1\_VMKernelSpec calls the
autogenerated conversion function along with custom conversion logic

## <a name="Convert_ignite_VMSpec_To_v1alpha1_VMSpec">func</a> [Convert\_ignite\_VMSpec\_To\_v1alpha1\_VMSpec](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha1/conversion.go?s=299:402#L10)

``` go
func Convert_ignite_VMSpec_To_v1alpha1_VMSpec(in *ignite.VMSpec, out *VMSpec, s conversion.Scope) error
```

Convert\_ignite\_VMSpec\_To\_v1alpha1\_VMSpec calls the autogenerated
conversion function along with custom conversion logic

## <a name="Convert_ignite_VMStatus_To_v1alpha1_VMStatus">func</a> [Convert\_ignite\_VMStatus\_To\_v1alpha1\_VMStatus](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha1/conversion.go?s=1098:1209#L22)

``` go
func Convert_ignite_VMStatus_To_v1alpha1_VMStatus(in *ignite.VMStatus, out *VMStatus, s conversion.Scope) error
```

Convert\_ignite\_VMStatus\_To\_v1alpha1\_VMStatus calls the
autogenerated conversion function along with custom conversion logic

## <a name="Convert_v1alpha1_ImageSpec_To_ignite_ImageSpec">func</a> [Convert\_v1alpha1\_ImageSpec\_To\_ignite\_ImageSpec](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha1/conversion.go?s=4692:4807#L116)

``` go
func Convert_v1alpha1_ImageSpec_To_ignite_ImageSpec(in *ImageSpec, out *ignite.ImageSpec, s conversion.Scope) error
```

Convert\_v1alpha1\_ImageSpec\_To\_ignite\_ImageSpec calls the
autogenerated conversion function along with custom conversion logic

## <a name="Convert_v1alpha1_KernelSpec_To_ignite_KernelSpec">func</a> [Convert\_v1alpha1\_KernelSpec\_To\_ignite\_KernelSpec](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha1/conversion.go?s=5556:5675#L134)

``` go
func Convert_v1alpha1_KernelSpec_To_ignite_KernelSpec(in *KernelSpec, out *ignite.KernelSpec, s conversion.Scope) error
```

Convert\_v1alpha1\_KernelSpec\_To\_ignite\_KernelSpec calls the
autogenerated conversion function along with custom conversion logic

## <a name="Convert_v1alpha1_OCIClaim_To_ignite_OCI">func</a> [Convert\_v1alpha1\_OCIClaim\_To\_ignite\_OCI](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha1/conversion.go?s=3592:3684#L88)

``` go
func Convert_v1alpha1_OCIClaim_To_ignite_OCI(in *OCIImageClaim, out *meta.OCIImageRef) error
```

## <a name="Convert_v1alpha1_OCIImageSource_To_ignite_OCIImageSource">func</a> [Convert\_v1alpha1\_OCIImageSource\_To\_ignite\_OCIImageSource](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha1/conversion.go?s=2775:2916#L68)

``` go
func Convert_v1alpha1_OCIImageSource_To_ignite_OCIImageSource(in *OCIImageSource, out *ignite.OCIImageSource, s conversion.Scope) (err error)
```

Convert\_v1alpha1\_OCIImageSource\_To\_ignite\_OCIImageSource calls the
autogenerated conversion function along with custom conversion logic

## <a name="Convert_v1alpha1_VMImageSpec_To_ignite_VMImageSpec">func</a> [Convert\_v1alpha1\_VMImageSpec\_To\_ignite\_VMImageSpec](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha1/conversion.go?s=6436:6559#L152)

``` go
func Convert_v1alpha1_VMImageSpec_To_ignite_VMImageSpec(in *VMImageSpec, out *ignite.VMImageSpec, s conversion.Scope) error
```

Convert\_v1alpha1\_VMImageSpec\_To\_ignite\_VMImageSpec calls the
autogenerated conversion function along with custom conversion logic

## <a name="Convert_v1alpha1_VMKernelSpec_To_ignite_VMKernelSpec">func</a> [Convert\_v1alpha1\_VMKernelSpec\_To\_ignite\_VMKernelSpec](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha1/conversion.go?s=7332:7459#L170)

``` go
func Convert_v1alpha1_VMKernelSpec_To_ignite_VMKernelSpec(in *VMKernelSpec, out *ignite.VMKernelSpec, s conversion.Scope) error
```

Convert\_v1alpha1\_VMKernelSpec\_To\_ignite\_VMKernelSpec calls the
autogenerated conversion function along with custom conversion logic

## <a name="Convert_v1alpha1_VMNetworkSpec_To_ignite_VMNetworkSpec">func</a> [Convert\_v1alpha1\_VMNetworkSpec\_To\_ignite\_VMNetworkSpec](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha1/conversion.go?s=7785:7916#L179)

``` go
func Convert_v1alpha1_VMNetworkSpec_To_ignite_VMNetworkSpec(in *VMNetworkSpec, out *ignite.VMNetworkSpec, s conversion.Scope) error
```

Convert\_v1alpha1\_VMNetworkSpec\_To\_ignite\_VMNetworkSpec calls the
autogenerated conversion function along with custom conversion logic

## <a name="Convert_v1alpha1_VMSpec_To_ignite_VMSpec">func</a> [Convert\_v1alpha1\_VMSpec\_To\_ignite\_VMSpec](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha1/conversion.go?s=697:800#L16)

``` go
func Convert_v1alpha1_VMSpec_To_ignite_VMSpec(in *VMSpec, out *ignite.VMSpec, s conversion.Scope) error
```

Convert\_ignite\_VMSpec\_To\_v1alpha1\_VMSpec calls the autogenerated
conversion function along with custom conversion logic

## <a name="Convert_v1alpha1_VMStatus_To_ignite_VMStatus">func</a> [Convert\_v1alpha1\_VMStatus\_To\_ignite\_VMStatus](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha1/conversion.go?s=1582:1693#L38)

``` go
func Convert_v1alpha1_VMStatus_To_ignite_VMStatus(in *VMStatus, out *ignite.VMStatus, s conversion.Scope) error
```

Convert\_v1alpha1\_VMStatus\_To\_ignite\_VMStatus calls the
autogenerated conversion function along with custom conversion logic

## <a name="SetDefaults_OCIImageClaim">func</a> [SetDefaults\_OCIImageClaim](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha1/defaults.go?s=307:357#L14)

``` go
func SetDefaults_OCIImageClaim(obj *OCIImageClaim)
```

## <a name="SetDefaults_PoolSpec">func</a> [SetDefaults\_PoolSpec](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha1/defaults.go?s=397:437#L18)

``` go
func SetDefaults_PoolSpec(obj *PoolSpec)
```

## <a name="SetDefaults_VMKernelSpec">func</a> [SetDefaults\_VMKernelSpec](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha1/defaults.go?s=1279:1327#L54)

``` go
func SetDefaults_VMKernelSpec(obj *VMKernelSpec)
```

## <a name="SetDefaults_VMNetworkSpec">func</a> [SetDefaults\_VMNetworkSpec](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha1/defaults.go?s=1575:1625#L65)

``` go
func SetDefaults_VMNetworkSpec(obj *VMNetworkSpec)
```

## <a name="SetDefaults_VMSpec">func</a> [SetDefaults\_VMSpec](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha1/defaults.go?s=963:999#L40)

``` go
func SetDefaults_VMSpec(obj *VMSpec)
```

## <a name="SetDefaults_VMStatus">func</a> [SetDefaults\_VMStatus](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha1/defaults.go?s=1696:1736#L71)

``` go
func SetDefaults_VMStatus(obj *VMStatus)
```

## <a name="FileMapping">type</a> [FileMapping](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha1/types.go?s=7766:7861#L207)

``` go
type FileMapping struct {
    HostPath string `json:"hostPath"`
    VMPath   string `json:"vmPath"`
}
```

FileMapping defines mappings between files on the host and VM

## <a name="Image">type</a> [Image](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha1/types.go?s=398:871#L18)

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

## <a name="ImageSourceType">type</a> [ImageSourceType](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha1/types.go?s=1060:1087#L35)

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

## <a name="ImageSpec">type</a> [ImageSpec](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha1/types.go?s=919:986#L30)

``` go
type ImageSpec struct {
    OCIClaim OCIImageClaim `json:"ociClaim"`
}
```

ImageSpec declares what the image contains

## <a name="ImageStatus">type</a> [ImageStatus](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha1/types.go?s=2385:2534#L69)

``` go
type ImageStatus struct {
    // OCISource contains the information about how this OCI image was imported
    OCISource OCIImageSource `json:"ociSource"`
}
```

ImageStatus defines the status of the image

## <a name="Kernel">type</a> [Kernel](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha1/types.go?s=4919:5395#L133)

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

## <a name="KernelSpec">type</a> [KernelSpec](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha1/types.go?s=5448:5621#L145)

``` go
type KernelSpec struct {
    OCIClaim OCIImageClaim `json:"ociClaim"`
}
```

KernelSpec describes the properties of a kernel

## <a name="KernelStatus">type</a> [KernelStatus](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha1/types.go?s=5672:5788#L152)

``` go
type KernelStatus struct {
    Version   string         `json:"version"`
    OCISource OCIImageSource `json:"ociSource"`
}
```

KernelStatus describes the status of a kernel

## <a name="NetworkMode">type</a> [NetworkMode](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha1/types.go?s=8215:8238#L222)

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

### <a name="NetworkMode.String">func</a> (NetworkMode) [String](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha1/types.go?s=8278:8315#L226)

``` go
func (nm NetworkMode) String() string
```

## <a name="OCIImageClaim">type</a> [OCIImageClaim](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha1/types.go?s=1283:1701#L43)

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

## <a name="OCIImageSource">type</a> [OCIImageSource](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha1/types.go?s=1808:2336#L55)

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

OCIImageSource specifies how the OCI image was imported. It is the
status variant of OCIImageClaim

## <a name="Pool">type</a> [Pool](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha1/types.go?s=2809:2995#L78)

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

## <a name="PoolDevice">type</a> [PoolDevice](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha1/types.go?s=4286:4676#L120)

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

## <a name="PoolDeviceType">type</a> [PoolDeviceType](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha1/types.go?s=4015:4041#L110)

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

## <a name="PoolSpec">type</a> [PoolSpec](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha1/types.go?s=3042:3755#L88)

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

## <a name="PoolStatus">type</a> [PoolStatus](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha1/types.go?s=3805:4013#L104)

``` go
type PoolStatus struct {
    // The Devices array needs to contain pointers to accommodate "holes" in the mapping
    // Where devices have been deleted, the pointer is nil
    Devices []*PoolDevice `json:"devices"`
}
```

PoolStatus defines the Pool’s current status

## <a name="SSH">type</a> [SSH](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha1/types.go?s=8081:8158#L216)

``` go
type SSH struct {
    Generate  bool   `json:"-"`
    PublicKey string `json:"-"`
}
```

SSH specifies different ways to connect via SSH to the VM SSH uses a
custom marshaller/unmarshaller. If generate is true, it marshals to true
(a JSON bool). If PublicKey is set, it marshals to that string.

### <a name="SSH.MarshalJSON">func</a> (\*SSH) [MarshalJSON](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha1/json.go?s=117:160#L9)

``` go
func (s *SSH) MarshalJSON() ([]byte, error)
```

### <a name="SSH.UnmarshalJSON">func</a> (\*SSH) [UnmarshalJSON](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha1/json.go?s=308:351#L21)

``` go
func (s *SSH) UnmarshalJSON(b []byte) error
```

## <a name="VM">type</a> [VM](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha1/types.go?s=5990:6454#L160)

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

## <a name="VMImageSpec">type</a> [VMImageSpec](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha1/types.go?s=7387:7456#L192)

``` go
type VMImageSpec struct {
    OCIClaim OCIImageClaim `json:"ociClaim"`
}
```

## <a name="VMKernelSpec">type</a> [VMKernelSpec](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha1/types.go?s=7458:7579#L196)

``` go
type VMKernelSpec struct {
    OCIClaim OCIImageClaim `json:"ociClaim"`
    CmdLine  string        `json:"cmdLine,omitempty"`
}
```

## <a name="VMNetworkSpec">type</a> [VMNetworkSpec](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha1/types.go?s=7581:7699#L201)

``` go
type VMNetworkSpec struct {
    Mode  NetworkMode       `json:"mode"`
    Ports meta.PortMappings `json:"ports,omitempty"`
}
```

## <a name="VMSpec">type</a> [VMSpec](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha1/types.go?s=6502:7385#L172)

``` go
type VMSpec struct {
    Image    VMImageSpec   `json:"image"`
    Kernel   VMKernelSpec  `json:"kernel"`
    CPUs     uint64        `json:"cpus"`
    Memory   meta.Size     `json:"memory"`
    DiskSize meta.Size     `json:"diskSize"`
    Network  VMNetworkSpec `json:"network"`

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

## <a name="VMState">type</a> [VMState](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha1/types.go?s=8712:8731#L239)

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

## <a name="VMStatus">type</a> [VMStatus](https://github.com/weaveworks/ignite/tree/master/pkg/apis/ignite/v1alpha1/types.go?s=8891:9112#L248)

``` go
type VMStatus struct {
    State       VMState          `json:"state"`
    IPAddresses meta.IPAddresses `json:"ipAddresses,omitempty"`
    Image       OCIImageSource   `json:"image"`
    Kernel      OCIImageSource   `json:"kernel"`
}
```

VMStatus defines the status of a VM

-----

Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
