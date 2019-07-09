

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
  * [func APITypeFrom(obj Object) *APIType](#APITypeFrom)
  * [func NewAPIType() *APIType](#NewAPIType)
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
  * [func (k Kind) Title() string](#Kind.Title)
* [type OCIImageRef](#OCIImageRef)
  * [func NewOCIImageRef(imageStr string) (OCIImageRef, error)](#NewOCIImageRef)
  * [func (i OCIImageRef) IsUnset() bool](#OCIImageRef.IsUnset)
  * [func (i OCIImageRef) MarshalJSON() ([]byte, error)](#OCIImageRef.MarshalJSON)
  * [func (i OCIImageRef) String() string](#OCIImageRef.String)
  * [func (i *OCIImageRef) UnmarshalJSON(b []byte) error](#OCIImageRef.UnmarshalJSON)
* [type Object](#Object)
* [type ObjectMeta](#ObjectMeta)
  * [func (o *ObjectMeta) GetAnnotation(key string) string](#ObjectMeta.GetAnnotation)
  * [func (o *ObjectMeta) GetCreated() *Time](#ObjectMeta.GetCreated)
  * [func (o *ObjectMeta) GetLabel(key string) string](#ObjectMeta.GetLabel)
  * [func (o *ObjectMeta) GetName() string](#ObjectMeta.GetName)
  * [func (o *ObjectMeta) GetObjectMeta() *ObjectMeta](#ObjectMeta.GetObjectMeta)
  * [func (o *ObjectMeta) GetUID() UID](#ObjectMeta.GetUID)
  * [func (o *ObjectMeta) SetAnnotation(key, value string)](#ObjectMeta.SetAnnotation)
  * [func (o *ObjectMeta) SetCreated(t *Time)](#ObjectMeta.SetCreated)
  * [func (o *ObjectMeta) SetLabel(key, value string)](#ObjectMeta.SetLabel)
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
  * [func (t *TypeMeta) GetTypeMeta() *TypeMeta](#TypeMeta.GetTypeMeta)
* [type UID](#UID)
  * [func (u UID) String() string](#UID.String)


#### <a name="pkg-files">Package files</a>
[dmid.go](/pkg/apis/meta/v1alpha1/dmid.go) [doc.go](/pkg/apis/meta/v1alpha1/doc.go) [image.go](/pkg/apis/meta/v1alpha1/image.go) [meta.go](/pkg/apis/meta/v1alpha1/meta.go) [net.go](/pkg/apis/meta/v1alpha1/net.go) [size.go](/pkg/apis/meta/v1alpha1/size.go) [time.go](/pkg/apis/meta/v1alpha1/time.go) [uid.go](/pkg/apis/meta/v1alpha1/uid.go) 



## <a name="pkg-variables">Variables</a>
``` go
var EmptySize = NewSizeFromBytes(0)
```



## <a name="APIType">type</a> [APIType](/pkg/apis/meta/v1alpha1/meta.go?s=411:495#L19)
``` go
type APIType struct {
    *TypeMeta   `json:",inline"`
    *ObjectMeta `json:"metadata"`
}

```
APIType is a struct implementing Object, used for
unmarshalling unknown objects into this intermediate type
where .Name, .UID, .Kind and .APIVersion become easily available
+k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object







### <a name="APITypeFrom">func</a> [APITypeFrom](/pkg/apis/meta/v1alpha1/meta.go?s=705:742#L33)
``` go
func APITypeFrom(obj Object) *APIType
```
APITypeFrom is used to create a bound APIType from an Object


### <a name="NewAPIType">func</a> [NewAPIType](/pkg/apis/meta/v1alpha1/meta.go?s=556:582#L25)
``` go
func NewAPIType() *APIType
```
This constructor ensures the APIType fields are not nil





## <a name="APITypeList">type</a> [APITypeList](/pkg/apis/meta/v1alpha1/meta.go?s=898:925#L43)
``` go
type APITypeList []*APIType
```
APITypeList is a list of many pointers APIType objects










## <a name="DMID">type</a> [DMID](/pkg/apis/meta/v1alpha1/dmid.go?s=83:116#L6)
``` go
type DMID struct {
    // contains filtered or unexported fields
}

```
DMID specifies the format for device mapper IDs







### <a name="NewDMID">func</a> [NewDMID](/pkg/apis/meta/v1alpha1/dmid.go?s=147:171#L12)
``` go
func NewDMID(i int) DMID
```

### <a name="NewPoolDMID">func</a> [NewPoolDMID](/pkg/apis/meta/v1alpha1/dmid.go?s=355:378#L23)
``` go
func NewPoolDMID() DMID
```




### <a name="DMID.Index">func</a> (\*DMID) [Index](/pkg/apis/meta/v1alpha1/dmid.go?s=514:540#L34)
``` go
func (d *DMID) Index() int
```



### <a name="DMID.Pool">func</a> (\*DMID) [Pool](/pkg/apis/meta/v1alpha1/dmid.go?s=462:488#L30)
``` go
func (d *DMID) Pool() bool
```



### <a name="DMID.String">func</a> (DMID) [String](/pkg/apis/meta/v1alpha1/dmid.go?s=623:652#L42)
``` go
func (d DMID) String() string
```



## <a name="IPAddresses">type</a> [IPAddresses](/pkg/apis/meta/v1alpha1/net.go?s=1541:1566#L78)
``` go
type IPAddresses []net.IP
```
IPAddresses represents a list of VM IP addresses










### <a name="IPAddresses.String">func</a> (IPAddresses) [String](/pkg/apis/meta/v1alpha1/net.go?s=1604:1640#L82)
``` go
func (i IPAddresses) String() string
```



## <a name="Kind">type</a> [Kind](/pkg/apis/meta/v1alpha1/meta.go?s=1218:1234#L59)
``` go
type Kind string
```









### <a name="Kind.Lower">func</a> (Kind) [Lower](/pkg/apis/meta/v1alpha1/meta.go?s=1644:1672#L81)
``` go
func (k Kind) Lower() string
```
Returns a lowercase string representation of the Kind




### <a name="Kind.String">func</a> (Kind) [String](/pkg/apis/meta/v1alpha1/meta.go?s=1337:1366#L64)
``` go
func (k Kind) String() string
```
Returns a string representation of the Kind suitable for sentences




### <a name="Kind.Title">func</a> (Kind) [Title](/pkg/apis/meta/v1alpha1/meta.go?s=1535:1563#L76)
``` go
func (k Kind) Title() string
```
Returns a title case string representation of the Kind




## <a name="OCIImageRef">type</a> [OCIImageRef](/pkg/apis/meta/v1alpha1/image.go?s=562:585#L23)
``` go
type OCIImageRef string
```






### <a name="NewOCIImageRef">func</a> [NewOCIImageRef](/pkg/apis/meta/v1alpha1/image.go?s=181:238#L11)
``` go
func NewOCIImageRef(imageStr string) (OCIImageRef, error)
```
NewOCIImageRef parses and normalizes a reference to an OCI (docker) image.





### <a name="OCIImageRef.IsUnset">func</a> (OCIImageRef) [IsUnset](/pkg/apis/meta/v1alpha1/image.go?s=647:682#L29)
``` go
func (i OCIImageRef) IsUnset() bool
```



### <a name="OCIImageRef.MarshalJSON">func</a> (OCIImageRef) [MarshalJSON](/pkg/apis/meta/v1alpha1/image.go?s=708:758#L33)
``` go
func (i OCIImageRef) MarshalJSON() ([]byte, error)
```



### <a name="OCIImageRef.String">func</a> (OCIImageRef) [String](/pkg/apis/meta/v1alpha1/image.go?s=587:623#L25)
``` go
func (i OCIImageRef) String() string
```



### <a name="OCIImageRef.UnmarshalJSON">func</a> (\*OCIImageRef) [UnmarshalJSON](/pkg/apis/meta/v1alpha1/image.go?s=796:847#L37)
``` go
func (i *OCIImageRef) UnmarshalJSON(b []byte) error
```



## <a name="Object">type</a> [Object](/pkg/apis/meta/v1alpha1/meta.go?s=3735:4082#L165)
``` go
type Object interface {
    runtime.Object

    GetTypeMeta() *TypeMeta
    GetObjectMeta() *ObjectMeta

    GetKind() Kind

    GetName() string
    SetName(string)

    GetUID() UID
    SetUID(UID)

    GetCreated() *Time
    SetCreated(t *Time)

    GetLabel(key string) string
    SetLabel(key, value string)

    GetAnnotation(key string) string
    SetAnnotation(key, value string)
}
```
Object extends k8s.io/apimachinery's runtime.Object with
extra GetName() and GetUID() methods from ObjectMeta










## <a name="ObjectMeta">type</a> [ObjectMeta](/pkg/apis/meta/v1alpha1/meta.go?s=1879:2181#L88)
``` go
type ObjectMeta struct {
    Name        string            `json:"name"`
    UID         UID               `json:"uid,omitempty"`
    Created     *Time             `json:"created,omitempty"`
    Labels      map[string]string `json:"labels,omitempty"`
    Annotations map[string]string `json:"annotations,omitempty"`
}

```
ObjectMeta have to be embedded into any serializable object.
It provides the .GetName() and .GetUID() methods that help
implement the Object interface










### <a name="ObjectMeta.GetAnnotation">func</a> (\*ObjectMeta) [GetAnnotation](/pkg/apis/meta/v1alpha1/meta.go?s=3290:3343#L148)
``` go
func (o *ObjectMeta) GetAnnotation(key string) string
```
GetAnnotation returns the label value for the key




### <a name="ObjectMeta.GetCreated">func</a> (\*ObjectMeta) [GetCreated](/pkg/apis/meta/v1alpha1/meta.go?s=2726:2765#L122)
``` go
func (o *ObjectMeta) GetCreated() *Time
```
GetCreated returns when the Object was created




### <a name="ObjectMeta.GetLabel">func</a> (\*ObjectMeta) [GetLabel](/pkg/apis/meta/v1alpha1/meta.go?s=2948:2996#L132)
``` go
func (o *ObjectMeta) GetLabel(key string) string
```
GetLabel returns the label value for the key




### <a name="ObjectMeta.GetName">func</a> (\*ObjectMeta) [GetName](/pkg/apis/meta/v1alpha1/meta.go?s=2332:2369#L102)
``` go
func (o *ObjectMeta) GetName() string
```
GetName returns the name of the Object




### <a name="ObjectMeta.GetObjectMeta">func</a> (\*ObjectMeta) [GetObjectMeta](/pkg/apis/meta/v1alpha1/meta.go?s=2226:2274#L97)
``` go
func (o *ObjectMeta) GetObjectMeta() *ObjectMeta
```
This is a helper for APIType generation




### <a name="ObjectMeta.GetUID">func</a> (\*ObjectMeta) [GetUID](/pkg/apis/meta/v1alpha1/meta.go?s=2531:2564#L112)
``` go
func (o *ObjectMeta) GetUID() UID
```
GetUID returns the UID of the Object




### <a name="ObjectMeta.SetAnnotation">func</a> (\*ObjectMeta) [SetAnnotation](/pkg/apis/meta/v1alpha1/meta.go?s=3464:3517#L156)
``` go
func (o *ObjectMeta) SetAnnotation(key, value string)
```
SetAnnotation sets a label value for a key




### <a name="ObjectMeta.SetCreated">func</a> (\*ObjectMeta) [SetCreated](/pkg/apis/meta/v1alpha1/meta.go?s=2839:2879#L127)
``` go
func (o *ObjectMeta) SetCreated(t *Time)
```
SetCreated returns when the Object was created




### <a name="ObjectMeta.SetLabel">func</a> (\*ObjectMeta) [SetLabel](/pkg/apis/meta/v1alpha1/meta.go?s=3102:3150#L140)
``` go
func (o *ObjectMeta) SetLabel(key, value string)
```
SetLabel sets a label value for a key




### <a name="ObjectMeta.SetName">func</a> (\*ObjectMeta) [SetName](/pkg/apis/meta/v1alpha1/meta.go?s=2429:2470#L107)
``` go
func (o *ObjectMeta) SetName(name string)
```
SetName sets the name of the Object




### <a name="ObjectMeta.SetUID">func</a> (\*ObjectMeta) [SetUID](/pkg/apis/meta/v1alpha1/meta.go?s=2621:2657#L117)
``` go
func (o *ObjectMeta) SetUID(uid UID)
```
SetUID sets the UID of the Object




## <a name="PortMapping">type</a> [PortMapping](/pkg/apis/meta/v1alpha1/net.go?s=132:227#L11)
``` go
type PortMapping struct {
    HostPort uint64 `json:"hostPort"`
    VMPort   uint64 `json:"vmPort"`
}

```
PortMapping defines a port mapping between the VM and the host










### <a name="PortMapping.String">func</a> (PortMapping) [String](/pkg/apis/meta/v1alpha1/net.go?s=265:301#L18)
``` go
func (p PortMapping) String() string
```



## <a name="PortMappings">type</a> [PortMappings](/pkg/apis/meta/v1alpha1/net.go?s=418:449#L23)
``` go
type PortMappings []PortMapping
```
PortMappings represents a list of port mappings







### <a name="ParsePortMappings">func</a> [ParsePortMappings](/pkg/apis/meta/v1alpha1/net.go?s=488:548#L27)
``` go
func ParsePortMappings(input []string) (PortMappings, error)
```




### <a name="PortMappings.String">func</a> (PortMappings) [String](/pkg/apis/meta/v1alpha1/net.go?s=1249:1286#L61)
``` go
func (p PortMappings) String() string
```



## <a name="Size">type</a> [Size](/pkg/apis/meta/v1alpha1/size.go?s=125:164#L10)
``` go
type Size struct {
    datasize.ByteSize
}

```
Size specifies a common unit for data sizes







### <a name="NewSizeFromBytes">func</a> [NewSizeFromBytes](/pkg/apis/meta/v1alpha1/size.go?s=375:415#L24)
``` go
func NewSizeFromBytes(bytes uint64) Size
```

### <a name="NewSizeFromSectors">func</a> [NewSizeFromSectors](/pkg/apis/meta/v1alpha1/size.go?s=466:510#L30)
``` go
func NewSizeFromSectors(sectors uint64) Size
```

### <a name="NewSizeFromString">func</a> [NewSizeFromString](/pkg/apis/meta/v1alpha1/size.go?s=268:316#L19)
``` go
func NewSizeFromString(str string) (Size, error)
```




### <a name="Size.Add">func</a> (Size) [Add](/pkg/apis/meta/v1alpha1/size.go?s=838:872#L46)
``` go
func (s Size) Add(other Size) Size
```
Add returns a copy, does not modify the receiver




### <a name="Size.MarshalJSON">func</a> (\*Size) [MarshalJSON](/pkg/apis/meta/v1alpha1/size.go?s=1124:1168#L67)
``` go
func (s *Size) MarshalJSON() ([]byte, error)
```



### <a name="Size.Max">func</a> (Size) [Max](/pkg/apis/meta/v1alpha1/size.go?s=1021:1055#L59)
``` go
func (s Size) Max(other Size) Size
```



### <a name="Size.Min">func</a> (Size) [Min](/pkg/apis/meta/v1alpha1/size.go?s=918:952#L51)
``` go
func (s Size) Min(other Size) Size
```



### <a name="Size.Sectors">func</a> (\*Size) [Sectors](/pkg/apis/meta/v1alpha1/size.go?s=576:607#L36)
``` go
func (s *Size) Sectors() uint64
```



### <a name="Size.String">func</a> (\*Size) [String](/pkg/apis/meta/v1alpha1/size.go?s=735:765#L41)
``` go
func (s *Size) String() string
```
Override ByteSize's default string implementation which results in .HR() without spaces




### <a name="Size.UnmarshalJSON">func</a> (\*Size) [UnmarshalJSON](/pkg/apis/meta/v1alpha1/size.go?s=1231:1275#L72)
``` go
func (s *Size) UnmarshalJSON(b []byte) error
```



## <a name="Time">type</a> [Time](/pkg/apis/meta/v1alpha1/time.go?s=134:167#L11)
``` go
type Time struct {
    metav1.Time
}

```






### <a name="Timestamp">func</a> [Timestamp](/pkg/apis/meta/v1alpha1/time.go?s=460:481#L23)
``` go
func Timestamp() Time
```
Timestamp returns the current UTC time





### <a name="Time.String">func</a> (\*Time) [String](/pkg/apis/meta/v1alpha1/time.go?s=299:329#L18)
``` go
func (t *Time) String() string
```
The default string for Time is a human readable difference between the Time and the current time




## <a name="TypeMeta">type</a> [TypeMeta](/pkg/apis/meta/v1alpha1/meta.go?s=1014:1055#L46)
``` go
type TypeMeta struct {
    metav1.TypeMeta
}

```
TypeMeta is an alias for the k8s/apimachinery TypeMeta with some additional methods










### <a name="TypeMeta.GetKind">func</a> (\*TypeMeta) [GetKind](/pkg/apis/meta/v1alpha1/meta.go?s=1158:1191#L55)
``` go
func (t *TypeMeta) GetKind() Kind
```



### <a name="TypeMeta.GetTypeMeta">func</a> (\*TypeMeta) [GetTypeMeta](/pkg/apis/meta/v1alpha1/meta.go?s=1100:1142#L51)
``` go
func (t *TypeMeta) GetTypeMeta() *TypeMeta
```
This is a helper for APIType generation




## <a name="UID">type</a> [UID](/pkg/apis/meta/v1alpha1/uid.go?s=74:89#L6)
``` go
type UID string
```
UID represents an unique ID for a type










### <a name="UID.String">func</a> (UID) [String](/pkg/apis/meta/v1alpha1/uid.go?s=172:200#L11)
``` go
func (u UID) String() string
```
String returns the UID in string representation








- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
