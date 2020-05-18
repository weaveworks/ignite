package metadata

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/weaveworks/gitops-toolkit/pkg/runtime"
	"github.com/weaveworks/gitops-toolkit/pkg/storage"
	"github.com/weaveworks/gitops-toolkit/pkg/storage/cache"

	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	"github.com/weaveworks/ignite/pkg/apis/ignite/scheme"
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/client"
	"github.com/weaveworks/ignite/pkg/util"
)

func TestSetLabels(t *testing.T) {
	cases := []struct {
		name       string
		obj        runtime.Object
		labels     []string
		wantLabels map[string]string
		err        bool
	}{
		{
			name: "nil runtime object",
			obj:  nil,
			err:  true,
		},
		{
			name: "valid labels",
			obj:  &api.VM{},
			labels: []string{
				"label1=value1",
				"label2=value2",
				"label3=",
				"label4=value4,label5=value5",
			},
			wantLabels: map[string]string{
				"label1": "value1",
				"label2": "value2",
				"label3": "",
				"label4": "value4,label5=value5",
			},
		},
		{
			name:   "invalid label - key without value",
			obj:    &api.VM{},
			labels: []string{"key1"},
			err:    true,
		},
		{
			name:   "invalid label - empty name",
			obj:    &api.VM{},
			labels: []string{"="},
			err:    true,
		},
	}

	for _, rt := range cases {
		t.Run(rt.name, func(t *testing.T) {
			err := SetLabels(rt.obj, rt.labels)
			if (err != nil) != rt.err {
				t.Errorf("expected error %t, actual: %v", rt.err, err)
			}

			// Check the values of all the labels.
			for k, v := range rt.wantLabels {
				if rt.obj.GetLabel(k) != v {
					t.Errorf("expected label key %q to have value %q, actual: %q", k, v, rt.obj.GetLabel(k))
				}
			}
		})
	}
}

func TestVerifyUIDOrName(t *testing.T) {
	cases := []struct {
		name            string
		existingObjects []string
		newObject       string
		err             bool
	}{
		{
			name:            "create object with similar names",
			existingObjects: []string{"myvm1", "myvm11", "myvm111"},
			newObject:       "myvm",
		},
		{
			name:            "create object with existing names",
			existingObjects: []string{"myvm1", "myvm2"},
			newObject:       "myvm1",
			err:             true,
		},
	}

	for _, rt := range cases {
		t.Run(rt.name, func(t *testing.T) {
			// Create storage.
			dir, err := ioutil.TempDir("", "ignite")
			if err != nil {
				t.Fatalf("failed to create storage for ignite: %v", err)
			}
			defer os.RemoveAll(dir)

			storage := cache.NewCache(
				storage.NewGenericStorage(
					storage.NewGenericRawStorage(dir), scheme.Serializer))

			// Create ignite client with the created storage.
			ic := client.NewClient(storage)

			// Create existing VM object.
			objectKind := "VM"
			for _, objectName := range rt.existingObjects {
				vm := &api.VM{}
				vm.SetName(objectName)

				// Set UID.
				uid, err := util.NewUID()
				if err != nil {
					t.Errorf("failed to generate new UID: %v", err)
				}
				vm.SetUID(runtime.UID(uid))

				// Set VM image.
				ociRef, err := meta.NewOCIImageRef("foo/bar:latest")
				if err != nil {
					t.Errorf("failed to create new image reference: %v", err)
				}
				img := &api.Image{
					Spec: api.ImageSpec{
						OCI: ociRef,
					},
				}
				vm.SetImage(img)

				// Set Kernel image.
				ociRefKernel, err := meta.NewOCIImageRef("foo/bar:latest")
				if err != nil {
					t.Errorf("failed to create new image reference: %v", err)
				}
				kernel := &api.Kernel{
					Spec: api.KernelSpec{
						OCI: ociRefKernel,
					},
				}
				vm.SetKernel(kernel)

				// Set sandbox image without a helper.
				ociRefSandbox, err := meta.NewOCIImageRef("foo/bar:latest")
				if err != nil {
					t.Errorf("failed to create new image reference: %v", err)
				}
				vm.Spec.Sandbox.OCI = ociRefSandbox

				// Save object.
				if err := ic.VMs().Set(vm); err != nil {
					t.Errorf("failed to store VM object: %v", err)
				}
			}

			// Check if new object name exists.
			err = verifyUIDOrName(ic, rt.newObject, runtime.Kind(objectKind))
			if (err != nil) != rt.err {
				t.Errorf("expected error %t, actual: %v", rt.err, err)
			}
		})
	}
}
