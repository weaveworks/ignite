package patch

import (
	"testing"

	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
)

var testbytes = []byte(`
{
	"kind": "VM",
	"apiVersion": "ignite.weave.works/v1alpha1",
	"metadata": {
	  "name": "foo",
	  "uid": "0123456789101112"
	},
	"spec": {
	  "image": {
		"ociClaim": {
		  "ref": "centos:7"
		}
	  },
	  "kernel": {
		"ociClaim": {
		  "ref": "ubuntu:18.04"
		}
	  },
	  "cpus": 2,
	  "memory": "4MB"
	}
}`)

var vmGVK = api.SchemeGroupVersion.WithKind("VM")

func TestCreatePatch(t *testing.T) {
	vm := &api.VM{
		Spec: api.VMSpec{
			CPUs: 2,
		},
		Status: api.VMStatus{
			State: api.VMStateCreated,
		},
	}
	vm.SetGroupVersionKind(vmGVK)
	bytes, err := Create(vm, func(obj meta.Object) error {
		vm2 := obj.(*api.VM)
		vm2.Status.State = api.VMStateRunning
		return nil
	})
	t.Error(string(bytes), err, vm.Status.State)
}

func TestApplyPatch(t *testing.T) {
	patch := []byte(`{"status":{"state":"Running"}}`)
	result, err := Apply(testbytes, patch, vmGVK)
	t.Error(string(result), err)
}
