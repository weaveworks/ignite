# v1alpha1

`import "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"`

  - [Overview](#pkg-overview)
  - [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>

\+k8s:deepcopy-gen=package +k8s:openapi-gen=true

## <a name="pkg-index">Index</a>

  - [Variables](#pkg-variables)
  - [type DMID](#DMID)
      - [func NewDMID(i int) DMID](#NewDMID)
      - [func NewPoolDMID() DMID](#NewPoolDMID)
      - [func (d \*DMID) Index() int](#DMID.Index)
      - [func (d \*DMID) Pool() bool](#DMID.Pool)
      - [func (d DMID) String() string](#DMID.String)
  - [type IPAddresses](#IPAddresses)
      - [func (i IPAddresses) String() string](#IPAddresses.String)
  - [type OCIContentID](#OCIContentID)
      - [func ParseOCIContentID(str string) (\*OCIContentID,
        error)](#ParseOCIContentID)
      - [func (o \*OCIContentID) Digest()
        digest.Digest](#OCIContentID.Digest)
      - [func (o \*OCIContentID) Local() bool](#OCIContentID.Local)
      - [func (o \*OCIContentID) MarshalJSON() (\[\]byte,
        error)](#OCIContentID.MarshalJSON)
      - [func (o \*OCIContentID) RepoDigest() (n
        reference.Named)](#OCIContentID.RepoDigest)
      - [func (o \*OCIContentID) String() string](#OCIContentID.String)
      - [func (o \*OCIContentID) UnmarshalJSON(b \[\]byte) (err
        error)](#OCIContentID.UnmarshalJSON)
  - [type OCIImageRef](#OCIImageRef)
      - [func NewOCIImageRef(imageStr string) (OCIImageRef,
        error)](#NewOCIImageRef)
      - [func (i OCIImageRef) IsUnset() bool](#OCIImageRef.IsUnset)
      - [func (i OCIImageRef) String() string](#OCIImageRef.String)
      - [func (i \*OCIImageRef) UnmarshalJSON(b \[\]byte) (err
        error)](#OCIImageRef.UnmarshalJSON)
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

#### <a name="pkg-files">Package files</a>

[dmid.go](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/dmid.go)
[doc.go](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/doc.go)
[image.go](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/image.go)
[net.go](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/net.go)
[size.go](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/size.go)

## <a name="pkg-variables">Variables</a>

``` go
var EmptySize = NewSizeFromBytes(0)
```

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

## <a name="OCIContentID">type</a> [OCIContentID](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/image.go?s=2112:2353#L85)

``` go
type OCIContentID struct {
    // contains filtered or unexported fields
}
```

### <a name="ParseOCIContentID">func</a> [ParseOCIContentID](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/image.go?s=1683:1740#L62)

``` go
func ParseOCIContentID(str string) (*OCIContentID, error)
```

ParseOCIContentID takes in a string to parse into an \*OCIContentID If
given a local Docker SHA like
“sha256:3285f65b2651c68b5316e7a1fbabd30b5ae47914ac5791ac4bb9d59d029b924b”,
it will be parsed into the local format, encoded as “docker://<SHA>”.
Given a full repo digest, such as
“weaveworks/ignite-ubuntu@sha256:3285f65b2651c68b5316e7a1fbabd30b5ae47914ac5791ac4bb9d59d029b924b”,
it will be parsed into the OCI registry format, encoded as
“oci://<full path>@<SHA>”.

### <a name="OCIContentID.Digest">func</a> (\*OCIContentID) [Digest](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/image.go?s=3288:3333#L129)

``` go
func (o *OCIContentID) Digest() digest.Digest
```

Digest gets the digest of the content ID

### <a name="OCIContentID.Local">func</a> (\*OCIContentID) [Local](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/image.go?s=3174:3209#L124)

``` go
func (o *OCIContentID) Local() bool
```

Local returns true if the image has no repoName, i.e. it’s not available
from a registry

### <a name="OCIContentID.MarshalJSON">func</a> (\*OCIContentID) [MarshalJSON](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/image.go?s=3650:3702#L143)

``` go
func (o *OCIContentID) MarshalJSON() ([]byte, error)
```

### <a name="OCIContentID.RepoDigest">func</a> (\*OCIContentID) [RepoDigest](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/image.go?s=3452:3507#L134)

``` go
func (o *OCIContentID) RepoDigest() (n reference.Named)
```

RepoDigest returns a repo digest based on the OCIContentID if it is not
local

### <a name="OCIContentID.String">func</a> (\*OCIContentID) [String](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/image.go?s=2436:2474#L93)

``` go
func (o *OCIContentID) String() string
```

### <a name="OCIContentID.UnmarshalJSON">func</a> (\*OCIContentID) [UnmarshalJSON](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/image.go?s=3741:3799#L147)

``` go
func (o *OCIContentID) UnmarshalJSON(b []byte) (err error)
```

## <a name="OCIImageRef">type</a> [OCIImageRef](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/image.go?s=838:861#L35)

``` go
type OCIImageRef string
```

OCIImageRef is a string by which an OCI runtime can identify an image to
retrieve. It needs to have a tag and usually looks like
“weaveworks/ignite-ubuntu:latest”.

### <a name="NewOCIImageRef">func</a> [NewOCIImageRef](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/image.go?s=316:373#L19)

``` go
func NewOCIImageRef(imageStr string) (OCIImageRef, error)
```

NewOCIImageRef parses and normalizes a reference to an OCI (docker)
image.

### <a name="OCIImageRef.IsUnset">func</a> (OCIImageRef) [IsUnset](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/image.go?s=961:996#L43)

``` go
func (i OCIImageRef) IsUnset() bool
```

### <a name="OCIImageRef.String">func</a> (OCIImageRef) [String](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/image.go?s=901:937#L39)

``` go
func (i OCIImageRef) String() string
```

### <a name="OCIImageRef.UnmarshalJSON">func</a> (\*OCIImageRef) [UnmarshalJSON](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/image.go?s=1022:1079#L47)

``` go
func (i *OCIImageRef) UnmarshalJSON(b []byte) (err error)
```

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

## <a name="Size">type</a> [Size](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/size.go?s=156:195#L13)

``` go
type Size struct {
    datasize.ByteSize
}
```

Size specifies a common unit for data sizes

### <a name="NewSizeFromBytes">func</a> [NewSizeFromBytes](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/size.go?s=435:475#L29)

``` go
func NewSizeFromBytes(bytes uint64) Size
```

### <a name="NewSizeFromSectors">func</a> [NewSizeFromSectors](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/size.go?s=526:570#L35)

``` go
func NewSizeFromSectors(sectors uint64) Size
```

### <a name="NewSizeFromString">func</a> [NewSizeFromString](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/size.go?s=328:376#L24)

``` go
func NewSizeFromString(str string) (Size, error)
```

### <a name="Size.Add">func</a> (Size) [Add](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/size.go?s=896:930#L51)

``` go
func (s Size) Add(other Size) Size
```

Add returns a copy, does not modify the receiver

### <a name="Size.MarshalJSON">func</a> (\*Size) [MarshalJSON](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/size.go?s=1182:1226#L72)

``` go
func (s *Size) MarshalJSON() ([]byte, error)
```

### <a name="Size.Max">func</a> (Size) [Max](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/size.go?s=1079:1113#L64)

``` go
func (s Size) Max(other Size) Size
```

### <a name="Size.Min">func</a> (Size) [Min](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/size.go?s=976:1010#L56)

``` go
func (s Size) Min(other Size) Size
```

### <a name="Size.Sectors">func</a> (Size) [Sectors](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/size.go?s=636:666#L41)

``` go
func (s Size) Sectors() uint64
```

### <a name="Size.String">func</a> (Size) [String](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/size.go?s=794:823#L46)

``` go
func (s Size) String() string
```

Override ByteSize’s default string implementation which results in .HR()
without spaces

### <a name="Size.UnmarshalJSON">func</a> (\*Size) [UnmarshalJSON](https://github.com/weaveworks/ignite/tree/master/pkg/apis/meta/v1alpha1/size.go?s=1289:1333#L77)

``` go
func (s *Size) UnmarshalJSON(b []byte) error
```

-----

Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
