# v1alpha1

`import "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"`

  - [Overview](#pkg-overview)
  - [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>

\+k8s:deepcopy-gen=package +k8s:openapi-gen=true

## <a name="pkg-index">Index</a>

  - [Variables](#pkg-variables)
  - [type APIType](#APIType)
      - [func APITypeFrom(obj Object) \*APIType](#APITypeFrom)
      - [func NewAPIType() \*APIType](#NewAPIType)
  - [type APITypeList](#APITypeList)
  - [type DMID](#DMID)
      - [func NewDMID(i int) DMID](#NewDMID)
      - [func NewPoolDMID() DMID](#NewPoolDMID)
      - [func (d \*DMID) Index() int](#DMID.Index)
      - [func (d \*DMID) Pool() bool](#DMID.Pool)
      - [func (d DMID) String() string](#DMID.String)
  - [type IPAddresses](#IPAddresses)
      - [func (i IPAddresses) String() string](#IPAddresses.String)
  - [type Kind](#Kind)
      - [func (k Kind) Lower() string](#Kind.Lower)
      - [func (k Kind) String() string](#Kind.String)
      - [func (k Kind) Title() string](#Kind.Title)
  - [type OCIContentID](#OCIContentID)
      - [func ParseOCIContentID(str string) (\*OCIContentID,
        error)](#ParseOCIContentID)
      - [func (o \*OCIContentID) Digest() string](#OCIContentID.Digest)
      - [func (o \*OCIContentID) Local() bool](#OCIContentID.Local)
      - [func (o \*OCIContentID) MarshalJSON() (\[\]byte,
        error)](#OCIContentID.MarshalJSON)
      - [func (o \*OCIContentID) RepoDigest() (s
        string)](#OCIContentID.RepoDigest)
      - [func (o \*OCIContentID) String() string](#OCIContentID.String)
      - [func (o \*OCIContentID) UnmarshalJSON(b \[\]byte) (err
        error)](#OCIContentID.UnmarshalJSON)
  - [type OCIImageRef](#OCIImageRef)
      - [func NewOCIImageRef(imageStr string) (OCIImageRef,
        error)](#NewOCIImageRef)
      - [func (i OCIImageRef) IsUnset() bool](#OCIImageRef.IsUnset)
      - [func (i OCIImageRef) String() string](#OCIImageRef.String)
      - [func (i \*OCIImageRef) UnmarshalJSON(b \[\]byte)
        error](#OCIImageRef.UnmarshalJSON)
  - [type Object](#Object)
  - [type ObjectMeta](#ObjectMeta)
      - [func (o \*ObjectMeta) GetAnnotation(key string)
        string](#ObjectMeta.GetAnnotation)
      - [func (o \*ObjectMeta) GetCreated()
        Time](#ObjectMeta.GetCreated)
      - [func (o \*ObjectMeta) GetLabel(key string)
        string](#ObjectMeta.GetLabel)
      - [func (o \*ObjectMeta) GetName() string](#ObjectMeta.GetName)
      - [func (o *ObjectMeta) GetObjectMeta()
        *ObjectMeta](#ObjectMeta.GetObjectMeta)
      - [func (o \*ObjectMeta) GetUID() UID](#ObjectMeta.GetUID)
      - [func (o \*ObjectMeta) SetAnnotation(key, value
        string)](#ObjectMeta.SetAnnotation)
      - [func (o \*ObjectMeta) SetCreated(t
        Time)](#ObjectMeta.SetCreated)
      - [func (o \*ObjectMeta) SetLabel(key, value
        string)](#ObjectMeta.SetLabel)
      - [func (o \*ObjectMeta) SetName(name
        string)](#ObjectMeta.SetName)
      - [func (o \*ObjectMeta) SetUID(uid UID)](#ObjectMeta.SetUID)
  - [type PortMapping](#PortMapping)
      - [func (p PortMapping) String() string](#PortMapping.String)
  - [type PortMappings](#PortMappings)
      - [func ParsePortMappings(input \[\]string) (PortMappings,
        error)](#ParsePortMappings)
      - [func (p PortMappings) String() string](#PortMappings.String)
  - [type Protocol](#Protocol)
      - [func (p Protocol) String() string](#Protocol.String)
      - [func (p \*Protocol) UnmarshalJSON(b \[\]byte) (err
        error)](#Protocol.UnmarshalJSON)
  - [type Size](#Size)
      - [func NewSizeFromBytes(bytes uint64) Size](#NewSizeFromBytes)
      - [func NewSizeFromSectors(sectors uint64)
        Size](#NewSizeFromSectors)
      - [func NewSizeFromString(str string) (Size,
        error)](#NewSizeFromString)
      - [func (s Size) Add(other Size) Size](#Size.Add)
      - [func (s \*Size) MarshalJSON() (\[\]byte,
        error)](#Size.MarshalJSON)
      - [func (s Size) Max(other Size) Size](#Size.Max)
      - [func (s Size) Min(other Size) Size](#Size.Min)
      - [func (s Size) Sectors() uint64](#Size.Sectors)
      - [func (s Size) String() string](#Size.String)
      - [func (s \*Size) UnmarshalJSON(b \[\]byte)
        error](#Size.UnmarshalJSON)
  - [type Time](#Time)
      - [func Timestamp() Time](#Timestamp)
      - [func (t Time) MarshalJSON() (b \[\]byte, err
        error)](#Time.MarshalJSON)
      - [func (t Time) String() string](#Time.String)
  - [type TypeMeta](#TypeMeta)
      - [func (t \*TypeMeta) GetKind() Kind](#TypeMeta.GetKind)
      - [func (t *TypeMeta) GetTypeMeta()
        *TypeMeta](#TypeMeta.GetTypeMeta)
      - [func (t \*TypeMeta) GroupVersionKind()
        schema.GroupVersionKind](#TypeMeta.GroupVersionKind)
      - [func (t \*TypeMeta) SetGroupVersionKind(gvk
        schema.GroupVersionKind)](#TypeMeta.SetGroupVersionKind)
  - [type UID](#UID)
      - [func (u UID) String() string](#UID.String)
      - [func (u \*UID) UnmarshalJSON(b \[\]byte)
        error](#UID.UnmarshalJSON)

#### <a name="pkg-files">Package files</a>

[dmid.go](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/dmid.go)
[doc.go](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/doc.go)
[image.go](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/image.go)
[meta.go](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/meta.go)
[net.go](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/net.go)
[size.go](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/size.go)
[time.go](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/time.go)
[uid.go](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/uid.go)

## <a name="pkg-variables">Variables</a>

``` go
var EmptySize = NewSizeFromBytes(0)
```

## <a name="APIType">type</a> [APIType](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/meta.go?s=453:537#L20)

``` go
type APIType struct {
    *TypeMeta   `json:",inline"`
    *ObjectMeta `json:"metadata"`
}
```

APIType is a struct implementing Object, used for unmarshalling unknown
objects into this intermediate type where .Name, .UID, .Kind and
.APIVersion become easily available
+k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

### <a name="APITypeFrom">func</a> [APITypeFrom](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/meta.go?s=747:784#L34)

``` go
func APITypeFrom(obj Object) *APIType
```

APITypeFrom is used to create a bound APIType from an Object

### <a name="NewAPIType">func</a> [NewAPIType](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/meta.go?s=598:624#L26)

``` go
func NewAPIType() *APIType
```

This constructor ensures the APIType fields are not nil

## <a name="APITypeList">type</a> [APITypeList](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/meta.go?s=940:967#L44)

``` go
type APITypeList []*APIType
```

APITypeList is a list of many pointers APIType objects

## <a name="DMID">type</a> [DMID](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/dmid.go?s=83:116#L6)

``` go
type DMID struct {
    // contains filtered or unexported fields
}
```

DMID specifies the format for device mapper IDs

### <a name="NewDMID">func</a> [NewDMID](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/dmid.go?s=147:171#L12)

``` go
func NewDMID(i int) DMID
```

### <a name="NewPoolDMID">func</a> [NewPoolDMID](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/dmid.go?s=355:378#L23)

``` go
func NewPoolDMID() DMID
```

### <a name="DMID.Index">func</a> (\*DMID) [Index](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/dmid.go?s=514:540#L34)

``` go
func (d *DMID) Index() int
```

### <a name="DMID.Pool">func</a> (\*DMID) [Pool](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/dmid.go?s=462:488#L30)

``` go
func (d *DMID) Pool() bool
```

### <a name="DMID.String">func</a> (DMID) [String](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/dmid.go?s=623:652#L42)

``` go
func (d DMID) String() string
```

## <a name="IPAddresses">type</a> [IPAddresses](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/net.go?s=3422:3447#L155)

``` go
type IPAddresses []net.IP
```

IPAddresses represents a list of VM IP addresses

### <a name="IPAddresses.String">func</a> (IPAddresses) [String](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/net.go?s=3485:3521#L159)

``` go
func (i IPAddresses) String() string
```

## <a name="Kind">type</a> [Kind](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/meta.go?s=1507:1523#L68)

``` go
type Kind string
```

### <a name="Kind.Lower">func</a> (Kind) [Lower](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/meta.go?s=1933:1961#L90)

``` go
func (k Kind) Lower() string
```

Returns a lowercase string representation of the Kind

### <a name="Kind.String">func</a> (Kind) [String](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/meta.go?s=1626:1655#L73)

``` go
func (k Kind) String() string
```

Returns a string representation of the Kind suitable for sentences

### <a name="Kind.Title">func</a> (Kind) [Title](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/meta.go?s=1824:1852#L85)

``` go
func (k Kind) Title() string
```

Returns a title case string representation of the Kind

## <a name="OCIContentID">type</a> [OCIContentID](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/image.go?s=1396:1637#L72)

``` go
type OCIContentID struct {
    // contains filtered or unexported fields
}
```

### <a name="ParseOCIContentID">func</a> [ParseOCIContentID](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/image.go?s=967:1024#L49)

``` go
func ParseOCIContentID(str string) (*OCIContentID, error)
```

### <a name="OCIContentID.Digest">func</a> (\*OCIContentID) [Digest](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/image.go?s=2344:2382#L104)

``` go
func (o *OCIContentID) Digest() string
```

Digest is a getter for the digest field

### <a name="OCIContentID.Local">func</a> (\*OCIContentID) [Local](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/image.go?s=2231:2266#L99)

``` go
func (o *OCIContentID) Local() bool
```

Local returns true if the image has no repoName, i.e. it’s not available
from a registry

### <a name="OCIContentID.MarshalJSON">func</a> (\*OCIContentID) [MarshalJSON](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/image.go?s=2616:2668#L117)

``` go
func (o *OCIContentID) MarshalJSON() ([]byte, error)
```

### <a name="OCIContentID.RepoDigest">func</a> (\*OCIContentID) [RepoDigest](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/image.go?s=2486:2532#L109)

``` go
func (o *OCIContentID) RepoDigest() (s string)
```

RepoDigest returns a repo digest based on the OCIContentID if it is not
local

### <a name="OCIContentID.String">func</a> (\*OCIContentID) [String](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/image.go?s=1720:1758#L80)

``` go
func (o *OCIContentID) String() string
```

### <a name="OCIContentID.UnmarshalJSON">func</a> (\*OCIContentID) [UnmarshalJSON](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/image.go?s=2707:2765#L121)

``` go
func (o *OCIContentID) UnmarshalJSON(b []byte) (err error)
```

## <a name="OCIImageRef">type</a> [OCIImageRef](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/image.go?s=586:609#L25)

``` go
type OCIImageRef string
```

### <a name="NewOCIImageRef">func</a> [NewOCIImageRef](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/image.go?s=231:288#L13)

``` go
func NewOCIImageRef(imageStr string) (OCIImageRef, error)
```

NewOCIImageRef parses and normalizes a reference to an OCI (docker)
image.

### <a name="OCIImageRef.IsUnset">func</a> (OCIImageRef) [IsUnset](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/image.go?s=709:744#L33)

``` go
func (i OCIImageRef) IsUnset() bool
```

### <a name="OCIImageRef.String">func</a> (OCIImageRef) [String](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/image.go?s=649:685#L29)

``` go
func (i OCIImageRef) String() string
```

### <a name="OCIImageRef.UnmarshalJSON">func</a> (\*OCIImageRef) [UnmarshalJSON](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/image.go?s=770:821#L37)

``` go
func (i *OCIImageRef) UnmarshalJSON(b []byte) error
```

## <a name="Object">type</a> [Object](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/meta.go?s=4013:4448#L174)

``` go
type Object interface {
    runtime.Object

    GetTypeMeta() *TypeMeta
    GetObjectMeta() *ObjectMeta

    GetKind() Kind
    GroupVersionKind() schema.GroupVersionKind
    SetGroupVersionKind(schema.GroupVersionKind)

    GetName() string
    SetName(string)

    GetUID() UID
    SetUID(UID)

    GetCreated() Time
    SetCreated(t Time)

    GetLabel(key string) string
    SetLabel(key, value string)

    GetAnnotation(key string) string
    SetAnnotation(key, value string)
}
```

Object extends k8s.io/apimachinery’s runtime.Object with extra GetName()
and GetUID() methods from ObjectMeta

## <a name="ObjectMeta">type</a> [ObjectMeta](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/meta.go?s=2168:2460#L97)

``` go
type ObjectMeta struct {
    Name        string            `json:"name"`
    UID         UID               `json:"uid,omitempty"`
    Created     Time              `json:"created"`
    Labels      map[string]string `json:"labels,omitempty"`
    Annotations map[string]string `json:"annotations,omitempty"`
}
```

ObjectMeta have to be embedded into any serializable object. It provides
the .GetName() and .GetUID() methods that help implement the Object
interface

### <a name="ObjectMeta.GetAnnotation">func</a> (\*ObjectMeta) [GetAnnotation](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/meta.go?s=3568:3621#L157)

``` go
func (o *ObjectMeta) GetAnnotation(key string) string
```

GetAnnotation returns the label value for the key

### <a name="ObjectMeta.GetCreated">func</a> (\*ObjectMeta) [GetCreated](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/meta.go?s=3005:3043#L131)

``` go
func (o *ObjectMeta) GetCreated() Time
```

GetCreated returns when the Object was created

### <a name="ObjectMeta.GetLabel">func</a> (\*ObjectMeta) [GetLabel](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/meta.go?s=3226:3274#L141)

``` go
func (o *ObjectMeta) GetLabel(key string) string
```

GetLabel returns the label value for the key

### <a name="ObjectMeta.GetName">func</a> (\*ObjectMeta) [GetName](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/meta.go?s=2611:2648#L111)

``` go
func (o *ObjectMeta) GetName() string
```

GetName returns the name of the Object

### <a name="ObjectMeta.GetObjectMeta">func</a> (\*ObjectMeta) [GetObjectMeta](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/meta.go?s=2505:2553#L106)

``` go
func (o *ObjectMeta) GetObjectMeta() *ObjectMeta
```

This is a helper for APIType generation

### <a name="ObjectMeta.GetUID">func</a> (\*ObjectMeta) [GetUID](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/meta.go?s=2810:2843#L121)

``` go
func (o *ObjectMeta) GetUID() UID
```

GetUID returns the UID of the Object

### <a name="ObjectMeta.SetAnnotation">func</a> (\*ObjectMeta) [SetAnnotation](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/meta.go?s=3742:3795#L165)

``` go
func (o *ObjectMeta) SetAnnotation(key, value string)
```

SetAnnotation sets a label value for a key

### <a name="ObjectMeta.SetCreated">func</a> (\*ObjectMeta) [SetCreated](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/meta.go?s=3118:3157#L136)

``` go
func (o *ObjectMeta) SetCreated(t Time)
```

SetCreated sets the creation time of the Object

### <a name="ObjectMeta.SetLabel">func</a> (\*ObjectMeta) [SetLabel](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/meta.go?s=3380:3428#L149)

``` go
func (o *ObjectMeta) SetLabel(key, value string)
```

SetLabel sets a label value for a key

### <a name="ObjectMeta.SetName">func</a> (\*ObjectMeta) [SetName](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/meta.go?s=2708:2749#L116)

``` go
func (o *ObjectMeta) SetName(name string)
```

SetName sets the name of the Object

### <a name="ObjectMeta.SetUID">func</a> (\*ObjectMeta) [SetUID](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/meta.go?s=2900:2936#L126)

``` go
func (o *ObjectMeta) SetUID(uid UID)
```

SetUID sets the UID of the Object

## <a name="PortMapping">type</a> [PortMapping](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/net.go?s=190:398#L14)

``` go
type PortMapping struct {
    BindAddress net.IP   `json:"bindAddress,omitempty"`
    HostPort    uint64   `json:"hostPort"`
    VMPort      uint64   `json:"vmPort"`
    Protocol    Protocol `json:"protocol,omitempty"`
}
```

PortMapping defines a port mapping between the VM and the host

### <a name="PortMapping.String">func</a> (PortMapping) [String](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/net.go?s=436:472#L23)

``` go
func (p PortMapping) String() string
```

## <a name="PortMappings">type</a> [PortMappings](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/net.go?s=826:857#L42)

``` go
type PortMappings []PortMapping
```

PortMappings represents a list of port mappings

### <a name="ParsePortMappings">func</a> [ParsePortMappings](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/net.go?s=896:956#L46)

``` go
func ParsePortMappings(input []string) (PortMappings, error)
```

### <a name="PortMappings.String">func</a> (PortMappings) [String](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/net.go?s=2475:2512#L104)

``` go
func (p PortMappings) String() string
```

## <a name="Protocol">type</a> [Protocol](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/net.go?s=2761:2781#L121)

``` go
type Protocol string
```

Protocol specifies a network port protocol

``` go
const (
    ProtocolTCP Protocol = "tcp"
    ProtocolUDP Protocol = "udp"
)
```

### <a name="Protocol.String">func</a> (Protocol) [String](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/net.go?s=3135:3168#L140)

``` go
func (p Protocol) String() string
```

### <a name="Protocol.UnmarshalJSON">func</a> (\*Protocol) [UnmarshalJSON](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/net.go?s=3192:3246#L144)

``` go
func (p *Protocol) UnmarshalJSON(b []byte) (err error)
```

## <a name="Size">type</a> [Size](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/size.go?s=132:171#L11)

``` go
type Size struct {
    datasize.ByteSize
}
```

Size specifies a common unit for data sizes

### <a name="NewSizeFromBytes">func</a> [NewSizeFromBytes](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/size.go?s=411:451#L27)

``` go
func NewSizeFromBytes(bytes uint64) Size
```

### <a name="NewSizeFromSectors">func</a> [NewSizeFromSectors](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/size.go?s=502:546#L33)

``` go
func NewSizeFromSectors(sectors uint64) Size
```

### <a name="NewSizeFromString">func</a> [NewSizeFromString](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/size.go?s=304:352#L22)

``` go
func NewSizeFromString(str string) (Size, error)
```

### <a name="Size.Add">func</a> (Size) [Add](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/size.go?s=872:906#L49)

``` go
func (s Size) Add(other Size) Size
```

Add returns a copy, does not modify the receiver

### <a name="Size.MarshalJSON">func</a> (\*Size) [MarshalJSON](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/size.go?s=1158:1202#L70)

``` go
func (s *Size) MarshalJSON() ([]byte, error)
```

### <a name="Size.Max">func</a> (Size) [Max](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/size.go?s=1055:1089#L62)

``` go
func (s Size) Max(other Size) Size
```

### <a name="Size.Min">func</a> (Size) [Min](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/size.go?s=952:986#L54)

``` go
func (s Size) Min(other Size) Size
```

### <a name="Size.Sectors">func</a> (Size) [Sectors](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/size.go?s=612:642#L39)

``` go
func (s Size) Sectors() uint64
```

### <a name="Size.String">func</a> (Size) [String](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/size.go?s=770:799#L44)

``` go
func (s Size) String() string
```

Override ByteSize’s default string implementation which results in .HR()
without spaces

### <a name="Size.UnmarshalJSON">func</a> (\*Size) [UnmarshalJSON](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/size.go?s=1265:1309#L75)

``` go
func (s *Size) UnmarshalJSON(b []byte) error
```

## <a name="Time">type</a> [Time](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/time.go?s=151:184#L12)

``` go
type Time struct {
    metav1.Time
}
```

### <a name="Timestamp">func</a> [Timestamp](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/time.go?s=549:570#L30)

``` go
func Timestamp() Time
```

Timestamp returns the current UTC time

### <a name="Time.MarshalJSON">func</a> (Time) [MarshalJSON](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/time.go?s=640:689#L38)

``` go
func (t Time) MarshalJSON() (b []byte, err error)
```

### <a name="Time.String">func</a> (Time) [String](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/time.go?s=347:376#L21)

``` go
func (t Time) String() string
```

The default string for Time is a human readable difference between the
Time and the current time

## <a name="TypeMeta">type</a> [TypeMeta](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/meta.go?s=1056:1097#L47)

``` go
type TypeMeta struct {
    metav1.TypeMeta
}
```

TypeMeta is an alias for the k8s/apimachinery TypeMeta with some
additional methods

### <a name="TypeMeta.GetKind">func</a> (\*TypeMeta) [GetKind](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/meta.go?s=1200:1233#L56)

``` go
func (t *TypeMeta) GetKind() Kind
```

### <a name="TypeMeta.GetTypeMeta">func</a> (\*TypeMeta) [GetTypeMeta](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/meta.go?s=1142:1184#L52)

``` go
func (t *TypeMeta) GetTypeMeta() *TypeMeta
```

This is a helper for APIType generation

### <a name="TypeMeta.GroupVersionKind">func</a> (\*TypeMeta) [GroupVersionKind](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/meta.go?s=1260:1321#L60)

``` go
func (t *TypeMeta) GroupVersionKind() schema.GroupVersionKind
```

### <a name="TypeMeta.SetGroupVersionKind">func</a> (\*TypeMeta) [SetGroupVersionKind](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/meta.go?s=1381:1448#L64)

``` go
func (t *TypeMeta) SetGroupVersionKind(gvk schema.GroupVersionKind)
```

## <a name="UID">type</a> [UID](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/uid.go?s=153:168#L12)

``` go
type UID string
```

UID represents an unique ID for a type

### <a name="UID.String">func</a> (UID) [String](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/uid.go?s=251:279#L17)

``` go
func (u UID) String() string
```

String returns the UID in string representation

### <a name="UID.UnmarshalJSON">func</a> (\*UID) [UnmarshalJSON](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/uid.go?s=445:488#L24)

``` go
func (u *UID) UnmarshalJSON(b []byte) error
```

This unmarshaler enables the UID to be passed in as an unquoted string
in JSON. Upon marshaling, quotes will be automatically added.

-----

Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
