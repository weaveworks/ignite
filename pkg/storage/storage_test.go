package storage

import (
	"testing"

	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	"github.com/weaveworks/ignite/pkg/apis/ignite/scheme"
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
)

var s = NewGenericStorage(NewDefaultRawStorage("/tmp/foo"), scheme.Serializer)

func TestStorageGet(t *testing.T) {
	obj := &api.VM{
		ObjectMeta: meta.ObjectMeta{
			UID: meta.UID("1234"),
		},
	}

	err := s.Get(obj)
	t.Fatal(*obj, err)
}

func TestStorageSet(t *testing.T) {
	err := s.Set(&api.VM{
		ObjectMeta: meta.ObjectMeta{
			Name: "foo",
			UID:  meta.UID("1234"),
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

		t.Logf("name: %s, id: %s, cpus: %d, memory: %s\n", vm.GetName(), vm.GetUID(), vm.Spec.CPUs, vm.Spec.Memory)
	}

	t.Fatal("fo")
}

func TestStorageListMeta(t *testing.T) {
	list, err := s.ListMeta("VM")
	if err != nil {
		t.Fatal(err)
	}

	for _, item := range list {
		t.Logf("name: %s, id: %s, kind: %s, apiversion: %s\n", item.GetName(), item.GetUID(), item.GetKind(), item.GetTypeMeta().APIVersion)
	}

	t.Fatal("fo")
}
