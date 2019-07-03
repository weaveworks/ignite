package storage

import (
	"testing"

	"github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"
	ignitemeta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
)

var s = NewStorage("/tmp/foo")

func TestStorageGet(t *testing.T) {
	obj := &v1alpha1.VM{
		ObjectMeta: ignitemeta.ObjectMeta{
			UID: types.UID("1234"),
		},
	}
	err := s.Get(obj)
	t.Fatal(*obj, err)
}

func TestStorageSet(t *testing.T) {
	err := s.Set(&v1alpha1.VM{
		ObjectMeta: ignitemeta.ObjectMeta{
			Name: "foo",
			UID:  types.UID("1234"),
		},
		Spec: v1alpha1.VMSpec{
			CPUs:   2,
			Memory: ignitemeta.NewSizeFromBytes(4 * 1024 * 1024),
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
		vm, ok := vmobj.(*v1alpha1.VM)
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
