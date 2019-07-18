package storage

import (
	"testing"

	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/apis/ignite/scheme"
)

var s = NewGenericStorage(NewDefaultRawStorage("/tmp/bar"), scheme.Serializer)
var vmGVK = api.SchemeGroupVersion.WithKind("VM")

func TestStorageNew(t *testing.T) {
	obj, err := s.New(vmGVK)
	t.Fatal(*(obj.(*api.VM)), err)
}

func TestStorageGet(t *testing.T) {
	obj, err := s.Get(vmGVK, meta.UID("0123456789101112"))
	t.Error(err)
	t.Error(*(obj.(*api.VM)))
}

func TestStorageSet(t *testing.T) {
	vm := &api.VM{
		ObjectMeta: meta.ObjectMeta{
			Name: "foo",
			UID:  meta.UID("0123456789101112"),
		},
		Spec: api.VMSpec{
			CPUs:   2,
			Memory: meta.NewSizeFromBytes(4 * 1024 * 1024),
			Image: api.VMImageSpec{
				OCIClaim: api.OCIImageClaim{
					Ref: meta.OCIImageRef("centos:7"),
				},
			},
			Kernel: api.VMKernelSpec{
				OCIClaim: api.OCIImageClaim{
					Ref: meta.OCIImageRef("ubuntu:18.04"),
				},
			},
		},
	}
	err := s.Set(vmGVK, vm)
	t.Fatal(err)
}
/*
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
*/