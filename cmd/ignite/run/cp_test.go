package run

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/weaveworks/libgitops/pkg/runtime"
	"github.com/weaveworks/libgitops/pkg/storage"
	"github.com/weaveworks/libgitops/pkg/storage/cache"
	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	"github.com/weaveworks/ignite/pkg/apis/ignite/scheme"
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/client"
	"github.com/weaveworks/ignite/pkg/providers"
	"github.com/weaveworks/ignite/pkg/util"
)

func TestNewCPOptions(t *testing.T) {
	testVMName := "my-vm"

	cases := []struct {
		name            string
		sourceArg       string
		destArg         string
		wantSource      string
		wantDest        string
		wantCPDirection CopyDirection
		err             bool
	}{
		{
			name:            "cp from host to VM",
			sourceArg:       "afile.txt",
			destArg:         "my-vm:bfile.txt",
			wantSource:      "afile.txt",
			wantDest:        "bfile.txt",
			wantCPDirection: CopyDirectionHostToVM,
		},
		{
			name:            "cp from VM to host",
			sourceArg:       "my-vm:afile.txt",
			destArg:         "bfile.txt",
			wantSource:      "afile.txt",
			wantDest:        "bfile.txt",
			wantCPDirection: CopyDirectionVMToHost,
		},
		{
			name:      "cp without VM reference",
			sourceArg: "afile.txt",
			destArg:   "bfile.txt",
			err:       true,
		},
		{
			name:      "cp with VM reference on both source and destination",
			sourceArg: "my-vm1:afile.txt",
			destArg:   "my-vm2:bfile.txt",
			err:       true,
		},
	}

	for _, rt := range cases {
		t.Run(rt.name, func(t *testing.T) {
			// Setup storage backend.
			dir, err := ioutil.TempDir("", "ignite")
			if err != nil {
				t.Fatalf("failed to create storage for ignite: %v", err)
			}
			defer os.RemoveAll(dir)

			storage := cache.NewCache(
				storage.NewGenericStorage(
					storage.NewGenericRawStorage(dir), scheme.Serializer))

			ic := client.NewClient(storage)

			// Create a test vm.
			vm := &api.VM{}
			vm.SetName(testVMName)
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

			// Set kernel image.
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

			// Set provider client used in copy to find VM matches.
			providers.Client = ic

			// Create new copy options.
			cf := &CPFlags{}
			co, err := cf.NewCPOptions(rt.sourceArg, rt.destArg)
			if (err != nil) != rt.err {
				t.Fatalf("expected error %t, actual: %v", rt.err, err)
			}

			// Continue checking if no error from copy options was expected.
			if !rt.err {
				// Check if the wanted copy options are set.
				if co.vm.Name != testVMName {
					t.Errorf("expected match vm to be %q, actual: %q", testVMName, co.vm.Name)
				}

				if co.source != rt.wantSource {
					t.Errorf("expected source to be %q, actual: %q", rt.wantSource, co.source)
				}

				if co.dest != rt.wantDest {
					t.Errorf("expected destination to be %q, actual: %q", rt.wantDest, co.dest)
				}

				if co.copyDirection != rt.wantCPDirection {
					t.Errorf("expected copy direction to be %v, actual: %v", rt.wantCPDirection, co.copyDirection)
				}
			}
		})
	}
}
