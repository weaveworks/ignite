

# v1alpha1
`import "/go/src/github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)

## <a name="pkg-overview">Overview</a>
+k8s:deepcopy-gen=package




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
* [type Object](#Object)
* [type ObjectMeta](#ObjectMeta)
  * [func (o *ObjectMeta) GetCreated() *metav1.Time](#ObjectMeta.GetCreated)
  * [func (o *ObjectMeta) GetName() string](#ObjectMeta.GetName)
  * [func (o *ObjectMeta) GetUID() string](#ObjectMeta.GetUID)
  * [func (o *ObjectMeta) SetCreated(t *metav1.Time)](#ObjectMeta.SetCreated)
  * [func (o *ObjectMeta) SetName(name string)](#ObjectMeta.SetName)
  * [func (o *ObjectMeta) SetUID(uid string)](#ObjectMeta.SetUID)
* [type Size](#Size)
  * [func NewSizeFromBytes(bytes uint64) Size](#NewSizeFromBytes)
  * [func NewSizeFromSectors(sectors uint64) Size](#NewSizeFromSectors)
  * [func NewSizeFromString(str string) (Size, error)](#NewSizeFromString)
  * [func (s Size) Add(other Size) Size](#Size.Add)
  * [func (s *Size) Int64() int64](#Size.Int64)
  * [func (s *Size) MarshalJSON() ([]byte, error)](#Size.MarshalJSON)
  * [func (s Size) Max(other Size) Size](#Size.Max)
  * [func (s Size) Min(other Size) Size](#Size.Min)
  * [func (s *Size) Sectors() uint64](#Size.Sectors)
  * [func (s *Size) String() string](#Size.String)
  * [func (s *Size) UnmarshalJSON(b []byte) error](#Size.UnmarshalJSON)
* [type UID](#UID)
  * [func (u UID) String() string](#UID.String)


#### <a name="pkg-files">Package files</a>
[doc.go](/src/target/doc.go) [meta.go](/src/target/meta.go) 



## <a name="pkg-variables">Variables</a>
``` go
var EmptySize = NewSizeFromBytes(0)
```



## <a name="APIType">type</a> [APIType](/src/target/meta.go?s=379:471#L20)
``` go
type APIType struct {
    metav1.TypeMeta `json:",inline"`
    ObjectMeta      `json:"metadata"`
}

```
APIType is a struct implementing Object, used for
unmarshalling unknown objects into this intermediate type
where .Name, .UID, .Kind and .APIVersion become easily available










## <a name="APITypeList">type</a> [APITypeList](/src/target/meta.go?s=531:558#L26)
``` go
type APITypeList []*APIType
```
APITypeList is a list of many pointers APIType objects










## <a name="DMID">type</a> [DMID](/src/target/meta.go?s=3489:3522#L172)
``` go
type DMID struct {
    // contains filtered or unexported fields
}

```
DMID specifies the format for device mapper IDs







### <a name="NewDMID">func</a> [NewDMID](/src/target/meta.go?s=3553:3577#L178)
``` go
func NewDMID(i int) DMID
```

### <a name="NewPoolDMID">func</a> [NewPoolDMID](/src/target/meta.go?s=3761:3784#L189)
``` go
func NewPoolDMID() DMID
```




### <a name="DMID.Index">func</a> (\*DMID) [Index](/src/target/meta.go?s=3920:3946#L200)
``` go
func (d *DMID) Index() int
```



### <a name="DMID.Pool">func</a> (\*DMID) [Pool](/src/target/meta.go?s=3868:3894#L196)
``` go
func (d *DMID) Pool() bool
```



### <a name="DMID.String">func</a> (DMID) [String](/src/target/meta.go?s=4036:4065#L208)
``` go
func (d DMID) String() string
```



## <a name="Object">type</a> [Object](/src/target/meta.go?s=1638:1823#L69)
``` go
type Object interface {
    runtime.Object

    GetName() string
    SetName(string)

    // TODO: Use UID
    GetUID() string
    SetUID(string)

    GetCreated() *metav1.Time
    SetCreated(t *metav1.Time)
}
```
Object extends k8s.io/apimachinery's runtime.Object with
extra GetName() and GetUID() methods from ObjectMeta










## <a name="ObjectMeta">type</a> [ObjectMeta](/src/target/meta.go?s=720:876#L31)
``` go
type ObjectMeta struct {
    Name    string       `json:"name"`
    UID     UID          `json:"uid,omitempty"`
    Created *metav1.Time `json:"created,omitempty"`
}

```
ObjectMeta have to be embedded into any serializable object.
It provides the .GetName() and .GetUID() methods that help
implement the Object interface










### <a name="ObjectMeta.GetCreated">func</a> (\*ObjectMeta) [GetCreated](/src/target/meta.go?s=1334:1380#L58)
``` go
func (o *ObjectMeta) GetCreated() *metav1.Time
```
GetCreated returns when the Object was created




### <a name="ObjectMeta.GetName">func</a> (\*ObjectMeta) [GetName](/src/target/meta.go?s=920:957#L38)
``` go
func (o *ObjectMeta) GetName() string
```
GetName returns the name of the Object




### <a name="ObjectMeta.GetUID">func</a> (\*ObjectMeta) [GetUID](/src/target/meta.go?s=1119:1155#L48)
``` go
func (o *ObjectMeta) GetUID() string
```
GetUID returns the UID of the Object




### <a name="ObjectMeta.SetCreated">func</a> (\*ObjectMeta) [SetCreated](/src/target/meta.go?s=1454:1501#L63)
``` go
func (o *ObjectMeta) SetCreated(t *metav1.Time)
```
SetCreated returns when the Object was created




### <a name="ObjectMeta.SetName">func</a> (\*ObjectMeta) [SetName](/src/target/meta.go?s=1017:1058#L43)
``` go
func (o *ObjectMeta) SetName(name string)
```
SetName sets the name of the Object




### <a name="ObjectMeta.SetUID">func</a> (\*ObjectMeta) [SetUID](/src/target/meta.go?s=1221:1260#L53)
``` go
func (o *ObjectMeta) SetUID(uid string)
```
SetUID sets the UID of the Object




## <a name="Size">type</a> [Size](/src/target/meta.go?s=2034:2073#L92)
``` go
type Size struct {
    datasize.ByteSize
}

```
Size specifies a common unit for data sizes







### <a name="NewSizeFromBytes">func</a> [NewSizeFromBytes](/src/target/meta.go?s=2296:2336#L107)
``` go
func NewSizeFromBytes(bytes uint64) Size
```

### <a name="NewSizeFromSectors">func</a> [NewSizeFromSectors](/src/target/meta.go?s=2387:2431#L113)
``` go
func NewSizeFromSectors(sectors uint64) Size
```

### <a name="NewSizeFromString">func</a> [NewSizeFromString](/src/target/meta.go?s=2177:2225#L101)
``` go
func NewSizeFromString(str string) (Size, error)
```




### <a name="Size.Add">func</a> (Size) [Add](/src/target/meta.go?s=2891:2925#L135)
``` go
func (s Size) Add(other Size) Size
```
Add returns a copy, does not modify the receiver




### <a name="Size.Int64">func</a> (\*Size) [Int64](/src/target/meta.go?s=2780:2808#L130)
``` go
func (s *Size) Int64() int64
```
Int64 returns the byte size as int64




### <a name="Size.MarshalJSON">func</a> (\*Size) [MarshalJSON](/src/target/meta.go?s=3177:3221#L156)
``` go
func (s *Size) MarshalJSON() ([]byte, error)
```



### <a name="Size.Max">func</a> (Size) [Max](/src/target/meta.go?s=3074:3108#L148)
``` go
func (s Size) Max(other Size) Size
```



### <a name="Size.Min">func</a> (Size) [Min](/src/target/meta.go?s=2971:3005#L140)
``` go
func (s Size) Min(other Size) Size
```



### <a name="Size.Sectors">func</a> (\*Size) [Sectors](/src/target/meta.go?s=2497:2528#L119)
``` go
func (s *Size) Sectors() uint64
```



### <a name="Size.String">func</a> (\*Size) [String](/src/target/meta.go?s=2661:2691#L124)
``` go
func (s *Size) String() string
```
Override ByteSize's default string implementation which results in something similar to HR()




### <a name="Size.UnmarshalJSON">func</a> (\*Size) [UnmarshalJSON](/src/target/meta.go?s=3260:3304#L160)
``` go
func (s *Size) UnmarshalJSON(b []byte) error
```



## <a name="UID">type</a> [UID](/src/target/meta.go?s=1867:1882#L84)
``` go
type UID string
```
UID represents an unique ID for a type










### <a name="UID.String">func</a> (UID) [String](/src/target/meta.go?s=1935:1963#L87)
``` go
func (u UID) String() string
```
String returns the UID in string representation








- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
