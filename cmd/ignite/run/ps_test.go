package run

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/weaveworks/gitops-toolkit/pkg/runtime"
	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/weaveworks/ignite/pkg/util"
)

// createTestVM creates a VM object, given a name and ID (optional) with default
// images.
func createTestVM(name, id string) (*api.VM, error) {
	vm := &api.VM{}
	vm.SetName(name)

	// Generate an ID if not provided.
	if id == "" {
		uid, err := util.NewUID()
		if err != nil {
			return nil, fmt.Errorf("failed to generate new UID: %v", err)
		}
		id = uid
	}
	vm.SetUID(runtime.UID(id))

	// Set a fixed time for deterministic results.
	createdTime := time.Date(2000, time.January, 1, 1, 0, 0, 0, time.UTC)
	vm.SetCreated(runtime.Time{Time: metav1.Time{Time: createdTime}})

	// Set VM image.
	ociRef, err := meta.NewOCIImageRef("foo/bar:latest")
	if err != nil {
		return nil, fmt.Errorf("failed to create new image reference: %v", err)
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
		return nil, fmt.Errorf("failed to create new image reference: %v", err)
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
		return nil, fmt.Errorf("failed to create new image reference: %v", err)
	}
	vm.Spec.Sandbox.OCI = ociRefSandbox

	return vm, nil
}

// Update the golden files with:
//   go test -v github.com/weaveworks/ignite/cmd/ignite/run -run TestPs -update
func TestPs(t *testing.T) {
	// Existing VMs with UID for deterministic results.
	// A sorted list of VMs. The VM list returned by the VM filter is sorted by
	// VM UID.
	existingVMs := []runtime.ObjectMeta{
		{Name: "vm1", UID: "20e1d566ce318ada"},
		{Name: "vm2", UID: "bfc80c948b1e2419"},
		{Name: "vm3", UID: "cddc37ba657766e3"},
	}

	cases := []struct {
		name    string
		psFlags *PsFlags
		golden  string
	}{
		{
			name:    "list in table format",
			psFlags: &PsFlags{},
			golden:  "output/ps-table.txt",
		},
		{
			name:    "filtered list in table format",
			psFlags: &PsFlags{Filter: "{{.ObjectMeta.Name}}=vm2"},
			golden:  "output/ps-filtered-table.txt",
		},
		{
			name:    "formatted filtered list",
			psFlags: &PsFlags{Filter: "{{.ObjectMeta.Name}}!=vm2", TemplateFormat: "Name: {{.ObjectMeta.Name}} Image: {{.Spec.Image.OCI}}"},
			golden:  "output/ps-formatted-table.txt",
		},
	}

	for _, rt := range cases {
		t.Run(rt.name, func(t *testing.T) {
			vms := []*api.VM{}

			// Create VMs.
			for _, eVM := range existingVMs {
				vm, err := createTestVM(eVM.Name, eVM.UID.String())
				if err != nil {
					t.Errorf("failed to create VM: %v", err)
				}
				vms = append(vms, vm)
			}

			psop := &psOptions{PsFlags: rt.psFlags, allVMs: vms}

			// Run vm list and capture stdout.
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			if err := Ps(psop); err != nil {
				t.Errorf("unexpected error while listing VMs: %v", err)
			}

			w.Close()
			os.Stdout = oldStdout

			var buf bytes.Buffer
			io.Copy(&buf, r)

			// Construct golden file path.
			goldenFilePath := fmt.Sprintf("testdata%c%s", filepath.Separator, rt.golden)

			// Update the golden file if needed.
			if *update {
				t.Log("update ps golden files")
				if err := ioutil.WriteFile(goldenFilePath, buf.Bytes(), 0644); err != nil {
					t.Fatalf("failed to update ps golden file: %s: %v", goldenFilePath, err)
				}
			}

			// Read golden file.
			wantOutput, err := ioutil.ReadFile(goldenFilePath)
			if err != nil {
				t.Fatalf("failed to read ps golden file: %s: %v", goldenFilePath, err)
			}

			// Check if the output contains expected result.
			if !bytes.Equal(buf.Bytes(), wantOutput) {
				t.Errorf("expected output to be:\n%v\ngot output:\n%v", string(wantOutput), string(buf.Bytes()))
			}
		})
	}

}
