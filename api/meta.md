

# v1alpha1
`import "/go/src/github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>
+k8s:deepcopy-gen=package
+k8s:openapi-gen=true




## <a name="pkg-index">Index</a>
* [Variables](#pkg-variables)
* [type APIType](#APIType)
* [type APITypeList](#APITypeList)
* [type DMID](#DMID)
  * [func NewDMID(i int) DMID](#NewDMID)
  * [func NewPoolDMID() DMID](#NewPoolDMID)
  * [func (d *DMID) Index() int](#DMID.Index)
  * [func (d *DMID) Pool() bool](#DMID.Pool)
  * [func (d DMID) String() string](#DMID.String)
* [type IPAddresses](#IPAddresses)
  * [func (i IPAddresses) String() string](#IPAddresses.String)
* [type Kind](#Kind)
  * [func (k Kind) Lower() string](#Kind.Lower)
  * [func (k Kind) String() string](#Kind.String)
  * [func (k Kind) Upper() string](#Kind.Upper)
* [type Object](#Object)
* [type ObjectMeta](#ObjectMeta)
  * [func (o *ObjectMeta) GetCreated() *Time](#ObjectMeta.GetCreated)
  * [func (o *ObjectMeta) GetName() string](#ObjectMeta.GetName)
  * [func (o *ObjectMeta) GetUID() UID](#ObjectMeta.GetUID)
  * [func (o *ObjectMeta) SetCreated(t *Time)](#ObjectMeta.SetCreated)
  * [func (o *ObjectMeta) SetName(name string)](#ObjectMeta.SetName)
  * [func (o *ObjectMeta) SetUID(uid UID)](#ObjectMeta.SetUID)
* [type PortMapping](#PortMapping)
  * [func (p PortMapping) String() string](#PortMapping.String)
* [type PortMappings](#PortMappings)
  * [func ParsePortMappings(input []string) (PortMappings, error)](#ParsePortMappings)
  * [func (p PortMappings) String() string](#PortMappings.String)
* [type Size](#Size)
  * [func NewSizeFromBytes(bytes uint64) Size](#NewSizeFromBytes)
  * [func NewSizeFromSectors(sectors uint64) Size](#NewSizeFromSectors)
  * [func NewSizeFromString(str string) (Size, error)](#NewSizeFromString)
  * [func (s Size) Add(other Size) Size](#Size.Add)
  * [func (s *Size) MarshalJSON() ([]byte, error)](#Size.MarshalJSON)
  * [func (s Size) Max(other Size) Size](#Size.Max)
  * [func (s Size) Min(other Size) Size](#Size.Min)
  * [func (s *Size) Sectors() uint64](#Size.Sectors)
  * [func (s *Size) String() string](#Size.String)
  * [func (s *Size) UnmarshalJSON(b []byte) error](#Size.UnmarshalJSON)
* [type Time](#Time)
  * [func Timestamp() Time](#Timestamp)
  * [func (t *Time) String() string](#Time.String)
* [type TypeMeta](#TypeMeta)
  * [func (t *TypeMeta) GetKind() Kind](#TypeMeta.GetKind)
* [type UID](#UID)
  * [func (u UID) String() string](#UID.String)


#### <a name="pkg-files">Package files</a>
[dmid.go](/src/target/dmid.go) [doc.go](/src/target/doc.go) [meta.go](/src/target/meta.go) [net.go](/src/target/net.go) [size.go](/src/target/size.go) [time.go](/src/target/time.go) [uid.go](/src/target/uid.go) 



## <a name="pkg-variables">Variables</a>
``` go
var EmptySize = NewSizeFromBytes(0)
```



## <a name="APIType">type</a> [APIType](/src/target/meta.go?s=411:493#L19)
``` go
type APIType struct {
    TypeMeta   `json:",inline"`
    ObjectMeta `json:"metadata"`
}

```
APIType is a struct implementing Object, used for
unmarshalling unknown objects into this intermediate type
where .Name, .UID, .Kind and .APIVersion become easily available
+k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object










## <a name="APITypeList">type</a> [APITypeList](/src/target/meta.go?s=580:607#L27)
``` go
type APITypeList []*APIType
```
APITypeList is a list of many pointers APIType objects










## <a name="DMID">type</a> [DMID](/src/target/dmid.go?s=83:116#L6)
``` go
type DMID struct {
    // contains filtered or unexported fields
}

```
DMID specifies the format for device mapper IDs







### <a name="NewDMID">func</a> [NewDMID](/src/target/dmid.go?s=147:171#L12)
``` go
func NewDMID(i int) DMID
```

### <a name="NewPoolDMID">func</a> [NewPoolDMID](/src/target/dmid.go?s=355:378#L23)
``` go
func NewPoolDMID() DMID
```




### <a name="DMID.Index">func</a> (\*DMID) [Index](/src/target/dmid.go?s=514:540#L34)
``` go
func (d *DMID) Index() int
```



### <a name="DMID.Pool">func</a> (\*DMID) [Pool](/src/target/dmid.go?s=462:488#L30)
``` go
func (d *DMID) Pool() bool
```



### <a name="DMID.String">func</a> (DMID) [String](/src/target/dmid.go?s=623:652#L42)
``` go
func (d DMID) String() string
```



## <a name="IPAddresses">type</a> [IPAddresses](/src/target/net.go?s=1541:1566#L78)
``` go
type IPAddresses []net.IP
```
IPAddresses represents a list of VM IP addresses










### <a name="IPAddresses.String">func</a> (IPAddresses) [String](/src/target/net.go?s=1604:1640#L82)
``` go
func (i IPAddresses) String() string
```



## <a name="Kind">type</a> [Kind](/src/target/meta.go?s=799:815#L38)
``` go
type Kind string
```

``` go
const (
    KindImage  Kind = "Image"
    KindKernel Kind = "Kernel"
    KindVM     Kind = "VM"
)
```









### <a name="Kind.Lower">func</a> (Kind) [Lower](/src/target/meta.go?s=1244:1272#L65)
``` go
func (k Kind) Lower() string
```



### <a name="Kind.String">func</a> (Kind) [String](/src/target/meta.go?s=995:1024#L49)
``` go
func (k Kind) String() string
```
Returns a lowercase string representation of the Kind




### <a name="Kind.Upper">func</a> (Kind) [Upper](/src/target/meta.go?s=1192:1220#L61)
``` go
func (k Kind) Upper() string
```
Returns a uppercase string representation of the Kind




## <a name="Object">type</a> [Object](/src/target/meta.go?s=2345:2509#L110)
``` go
type Object interface {
    runtime.Object

    GetKind() Kind

    GetName() string
    SetName(string)

    GetUID() UID
    SetUID(UID)

    GetCreated() *Time
    SetCreated(t *Time)
}
```
Object extends k8s.io/apimachinery's runtime.Object with
extra GetName() and GetUID() methods from ObjectMeta










## <a name="ObjectMeta">type</a> [ObjectMeta](/src/target/meta.go?s=1479:1617#L72)
``` go
type ObjectMeta struct {
    Name    string `json:"name"`
    UID     UID    `json:"uid,omitempty"`
    Created *Time  `json:"created,omitempty"`
}

```
ObjectMeta have to be embedded into any serializable object.
It provides the .GetName() and .GetUID() methods that help
implement the Object interface










### <a name="ObjectMeta.GetCreated">func</a> (\*ObjectMeta) [GetCreated](/src/target/meta.go?s=2055:2094#L99)
``` go
func (o *ObjectMeta) GetCreated() *Time
```
GetCreated returns when the Object was created




### <a name="ObjectMeta.GetName">func</a> (\*ObjectMeta) [GetName](/src/target/meta.go?s=1661:1698#L79)
``` go
func (o *ObjectMeta) GetName() string
```
GetName returns the name of the Object




### <a name="ObjectMeta.GetUID">func</a> (\*ObjectMeta) [GetUID](/src/target/meta.go?s=1860:1893#L89)
``` go
func (o *ObjectMeta) GetUID() UID
```
GetUID returns the UID of the Object




### <a name="ObjectMeta.SetCreated">func</a> (\*ObjectMeta) [SetCreated](/src/target/meta.go?s=2168:2208#L104)
``` go
func (o *ObjectMeta) SetCreated(t *Time)
```
SetCreated returns when the Object was created




### <a name="ObjectMeta.SetName">func</a> (\*ObjectMeta) [SetName](/src/target/meta.go?s=1758:1799#L84)
``` go
func (o *ObjectMeta) SetName(name string)
```
SetName sets the name of the Object




### <a name="ObjectMeta.SetUID">func</a> (\*ObjectMeta) [SetUID](/src/target/meta.go?s=1950:1986#L94)
``` go
func (o *ObjectMeta) SetUID(uid UID)
```
SetUID sets the UID of the Object




## <a name="PortMapping">type</a> [PortMapping](/src/target/net.go?s=132:227#L11)
``` go
type PortMapping struct {
    HostPort uint64 `json:"hostPort"`
    VMPort   uint64 `json:"vmPort"`
}

```
PortMapping defines a port mapping between the VM and the host










### <a name="PortMapping.String">func</a> (PortMapping) [String](/src/target/net.go?s=265:301#L18)
``` go
func (p PortMapping) String() string
```



## <a name="PortMappings">type</a> [PortMappings](/src/target/net.go?s=418:449#L23)
``` go
type PortMappings []PortMapping
```
PortMappings represents a list of port mappings







### <a name="ParsePortMappings">func</a> [ParsePortMappings](/src/target/net.go?s=488:548#L27)
``` go
func ParsePortMappings(input []string) (PortMappings, error)
```




### <a name="PortMappings.String">func</a> (PortMappings) [String](/src/target/net.go?s=1249:1286#L61)
``` go
func (p PortMappings) String() string
```



## <a name="Size">type</a> [Size](/src/target/size.go?s=125:164#L10)
``` go
type Size struct {
    datasize.ByteSize
}

```
Size specifies a common unit for data sizes







### <a name="NewSizeFromBytes">func</a> [NewSizeFromBytes](/src/target/size.go?s=375:415#L24)
``` go
func NewSizeFromBytes(bytes uint64) Size
```

### <a name="NewSizeFromSectors">func</a> [NewSizeFromSectors](/src/target/size.go?s=466:510#L30)
``` go
func NewSizeFromSectors(sectors uint64) Size
```

### <a name="NewSizeFromString">func</a> [NewSizeFromString](/src/target/size.go?s=268:316#L19)
``` go
func NewSizeFromString(str string) (Size, error)
```




### <a name="Size.Add">func</a> (Size) [Add](/src/target/size.go?s=838:872#L46)
``` go
func (s Size) Add(other Size) Size
```
Add returns a copy, does not modify the receiver




### <a name="Size.MarshalJSON">func</a> (\*Size) [MarshalJSON](/src/target/size.go?s=1124:1168#L67)
``` go
func (s *Size) MarshalJSON() ([]byte, error)
```



### <a name="Size.Max">func</a> (Size) [Max](/src/target/size.go?s=1021:1055#L59)
``` go
func (s Size) Max(other Size) Size
```



### <a name="Size.Min">func</a> (Size) [Min](/src/target/size.go?s=918:952#L51)
``` go
func (s Size) Min(other Size) Size
```



### <a name="Size.Sectors">func</a> (\*Size) [Sectors](/src/target/size.go?s=576:607#L36)
``` go
func (s *Size) Sectors() uint64
```



### <a name="Size.String">func</a> (\*Size) [String](/src/target/size.go?s=735:765#L41)
``` go
func (s *Size) String() string
```
Override ByteSize's default string implementation which results in .HR() without spaces




### <a name="Size.UnmarshalJSON">func</a> (\*Size) [UnmarshalJSON](/src/target/size.go?s=1231:1275#L72)
``` go
func (s *Size) UnmarshalJSON(b []byte) error
```



## <a name="Time">type</a> [Time](/src/target/time.go?s=134:167#L11)
``` go
type Time struct {
    metav1.Time
}

```






### <a name="Timestamp">func</a> [Timestamp](/src/target/time.go?s=460:481#L23)
``` go
func Timestamp() Time
```
Timestamp returns the current UTC time





### <a name="Time.String">func</a> (\*Time) [String](/src/target/time.go?s=299:329#L18)
``` go
func (t *Time) String() string
```
The default string for Time is a human readable difference between the Time and the current time




## <a name="TypeMeta">type</a> [TypeMeta](/src/target/meta.go?s=696:737#L30)
``` go
type TypeMeta struct {
    metav1.TypeMeta
}

```
TypeMeta is an alias for the k8s/apimachinery TypeMeta with some additional methods










### <a name="TypeMeta.GetKind">func</a> (\*TypeMeta) [GetKind](/src/target/meta.go?s=739:772#L34)
``` go
func (t *TypeMeta) GetKind() Kind
```



## <a name="UID">type</a> [UID](/src/target/uid.go?s=74:89#L6)
``` go
type UID string
```
UID represents an unique ID for a type










### <a name="UID.String">func</a> (UID) [String](/src/target/uid.go?s=172:200#L11)
``` go
func (u UID) String() string
```
String returns the UID in string representation








- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
