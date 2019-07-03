package storage

import (
	"testing"

	api "github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
)

var s = NewStorage("/tmp/foo")

func TestStorageGet(t *testing.T) {
	obj := &api.VM{
		ObjectMeta: meta.ObjectMeta{
			UID: types.UID("1234"),
		},
	}
	err := s.Get(obj)
	t.Fatal(*obj, err)
}

func TestStorageSet(t *testing.T) {
	err := s.Set(&api.VM{
		ObjectMeta: meta.ObjectMeta{
			Name: "foo",
			UID:  types.UID("1234"),
		},
		Spec: api.VMSpec{
			CPUs:   2,
			Memory: meta.NewSizeFromBytes(4 * 1024 * 1024),
		},
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestStorageDelete(t *testing.T) {
	err := s.Delete("VM", "1234")
	t.Fatal("foo", err)
}

func TestStorageList(t *testing.T) {
	list, err := s.List("VM")
	if err != nil {
		t.Fatal(err)
	}
	for _, vmobj := range list {
		vm, ok := vmobj.(*api.VM)
		if !ok {
			t.Fatalf("can't convert")
		}
		t.Logf("name: %s, id: %s, cpus: %d, memory: %s\n", vm.Name, string(vm.UID), vm.Spec.CPUs, vm.Spec.Memory)
	}
	t.Fatal("fo")
}

func TestStorageListMeta(t *testing.T) {
	list, err := s.ListMeta("VM")
	if err != nil {
		t.Fatal(err)
	}
	for _, item := range list {
		t.Logf("name: %s, id: %s, kind: %s, apiversion: %s\n", item.Name, string(item.UID), item.Kind, item.APIVersion)
	}
	t.Fatal("fo")
}
