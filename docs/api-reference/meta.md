

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
  * [func (o *ObjectMeta) GetUID() UID](#ObjectMeta.GetUID)
  * [func (o *ObjectMeta) SetCreated(t *metav1.Time)](#ObjectMeta.SetCreated)
  * [func (o *ObjectMeta) SetName(name string)](#ObjectMeta.SetName)
  * [func (o *ObjectMeta) SetUID(uid UID)](#ObjectMeta.SetUID)
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
* [type UID](#UID)
  * [func (u UID) String() string](#UID.String)


#### <a name="pkg-files">Package files</a>
[dmid.go](/src/target/dmid.go) [doc.go](/src/target/doc.go) [meta.go](/src/target/meta.go) [size.go](/src/target/size.go) [uid.go](/src/target/uid.go) 



## <a name="pkg-variables">Variables</a>
``` go
var EmptySize = NewSizeFromBytes(0)
```



## <a name="APIType">type</a> [APIType](/src/target/meta.go?s=323:415#L15)
``` go
type APIType struct {
    metav1.TypeMeta `json:",inline"`
    ObjectMeta      `json:"metadata"`
}

```
APIType is a struct implementing Object, used for
unmarshalling unknown objects into this intermediate type
where .Name, .UID, .Kind and .APIVersion become easily available










## <a name="APITypeList">type</a> [APITypeList](/src/target/meta.go?s=475:502#L21)
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



## <a name="Object">type</a> [Object](/src/target/meta.go?s=1562:1723#L64)
``` go
type Object interface {
    runtime.Object

    GetName() string
    SetName(string)

    GetUID() UID
    SetUID(UID)

    GetCreated() *metav1.Time
    SetCreated(t *metav1.Time)
}
```
Object extends k8s.io/apimachinery's runtime.Object with
extra GetName() and GetUID() methods from ObjectMeta










## <a name="ObjectMeta">type</a> [ObjectMeta](/src/target/meta.go?s=664:820#L26)
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










### <a name="ObjectMeta.GetCreated">func</a> (\*ObjectMeta) [GetCreated](/src/target/meta.go?s=1258:1304#L53)
``` go
func (o *ObjectMeta) GetCreated() *metav1.Time
```
GetCreated returns when the Object was created




### <a name="ObjectMeta.GetName">func</a> (\*ObjectMeta) [GetName](/src/target/meta.go?s=864:901#L33)
``` go
func (o *ObjectMeta) GetName() string
```
GetName returns the name of the Object




### <a name="ObjectMeta.GetUID">func</a> (\*ObjectMeta) [GetUID](/src/target/meta.go?s=1063:1096#L43)
``` go
func (o *ObjectMeta) GetUID() UID
```
GetUID returns the UID of the Object




### <a name="ObjectMeta.SetCreated">func</a> (\*ObjectMeta) [SetCreated](/src/target/meta.go?s=1378:1425#L58)
``` go
func (o *ObjectMeta) SetCreated(t *metav1.Time)
```
SetCreated returns when the Object was created




### <a name="ObjectMeta.SetName">func</a> (\*ObjectMeta) [SetName](/src/target/meta.go?s=961:1002#L38)
``` go
func (o *ObjectMeta) SetName(name string)
```
SetName sets the name of the Object




### <a name="ObjectMeta.SetUID">func</a> (\*ObjectMeta) [SetUID](/src/target/meta.go?s=1153:1189#L48)
``` go
func (o *ObjectMeta) SetUID(uid UID)
```
SetUID sets the UID of the Object




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




### <a name="Size.UnmarshalJSON">func</a> (\*Size) [UnmarshalJSON](/src/target/size.go?s=1223:1267#L72)
``` go
func (s *Size) UnmarshalJSON(b []byte) error
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
