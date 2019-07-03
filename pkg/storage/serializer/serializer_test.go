package serializer_test

import (
	"testing"

	"github.com/weaveworks/ignite/pkg/apis/ignite/scheme"
	"github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"
	ignitemeta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/storage/serializer"

	"k8s.io/apimachinery/pkg/types"
)

var s = serializer.NewSerializer(scheme.Scheme, nil)
var sampleobj = &v1alpha1.VM{
	ObjectMeta: ignitemeta.ObjectMeta{
		Name: "foo",
		UID:  types.UID("1234"),
	},
	Spec: v1alpha1.VMSpec{
		CPUs: 1,
	},
}
var samplejson = []byte(`{"kind":"VM","apiVersion":"ignite.weave.works/v1alpha1","metadata":{"name":"foo","uid":"1234"},"spec":{"cpus":1}}`)
var nonstrictjson = []byte(`{"kind":"VM","apiVersion":"ignite.weave.works/v1alpha1","metadata":{"name":"foo","uid":"1234"},"spec":{"cpus":1, "foo": "bar"}}`)

func TestEncodeJSON(t *testing.T) {
	b, err := s.EncodeJSON(sampleobj)
	t.Fatal(string(b), err)
}

func TestEncodeYAML(t *testing.T) {
	b, err := s.EncodeYAML(sampleobj)
	t.Fatal(string(b), err)
}

func TestDecode(t *testing.T) {
	obj, err := s.Decode(samplejson)
	t.Fatal(obj, err)
}

func TestDecodeInto(t *testing.T) {
	vm := &v1alpha1.VM{}
	err := s.DecodeInto(samplejson, vm)
	t.Fatal(*vm, err)
}

func TestDecodeStrict(t *testing.T) {
	vm := &v1alpha1.VM{}
	err := s.DecodeInto(nonstrictjson, vm)
	t.Fatal(vm, err)
}
