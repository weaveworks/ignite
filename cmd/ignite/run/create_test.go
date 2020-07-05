package run

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	flag "github.com/spf13/pflag"
	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	"github.com/weaveworks/ignite/pkg/apis/ignite/scheme"
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/client"
	"github.com/weaveworks/ignite/pkg/providers"
	"github.com/weaveworks/libgitops/pkg/runtime"
	"github.com/weaveworks/libgitops/pkg/storage"
	"github.com/weaveworks/libgitops/pkg/storage/cache"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Update the golden files with:
//   go test -v github.com/weaveworks/ignite/cmd/ignite/run -run TestApplyVMConfigFile -update
func TestApplyVMConfigFile(t *testing.T) {
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

	// Set provider client.
	providers.Client = ic

	// Helper for getting size from string.
	sizeFromString := func(strSize string) meta.Size {
		size, err := meta.NewSizeFromString(strSize)
		if err != nil {
			t.Fatalf("failed creating new size from string")
		}
		return size
	}

	testImage := "weaveworks/ignite-testimg"
	ociRef, err := meta.NewOCIImageRef(testImage)
	if err != nil {
		t.Fatalf("failed to create image ref for %q", testImage)
	}

	// Create a base VM.
	baseVM := providers.Client.VMs().New()
	// Set default image.
	baseVM.Spec.Image.OCI = ociRef

	// NOTE: When defining a base VM spec, ensure all the images are defined,
	// or the input vm config file has those images to avoid panic due to empty
	// image at object serialization. The default VM object contains default
	// images, but at the start of subtest, the vm spec is replaced with the
	// test input spec. The test checks the actual patching behavior.
	tests := []struct {
		name       string
		baseSpec   *api.VMSpec
		configFile string
		err        bool
		golden     string
	}{
		{
			name: "yaml VM config",
			baseSpec: &api.VMSpec{
				Memory: sizeFromString("500MB"),
				CPUs:   uint64(4),
				Image: api.VMImageSpec{
					OCI: ociRef,
				},
				Sandbox: api.VMSandboxSpec{
					OCI: ociRef,
				},
				Kernel: api.VMKernelSpec{
					OCI: ociRef,
				},
			},
			configFile: "input/apply-vm-config.yaml",
			golden:     "output/apply-vm-config-yaml.json",
		},
		{
			name: "json vm config",
			baseSpec: &api.VMSpec{
				SSH: &api.SSH{Generate: true},
				Image: api.VMImageSpec{
					OCI: ociRef,
				},
				Sandbox: api.VMSandboxSpec{
					OCI: ociRef,
				},
				Kernel: api.VMKernelSpec{
					OCI: ociRef,
				},
			},
			configFile: "input/apply-vm-config.json",
			golden:     "output/apply-vm-config-json.json",
		},
		{
			name:       "empty vm config",
			configFile: "input/apply-vm-config-empty.yaml",
			golden:     "output/apply-vm-config-empty.json",
		},
		{
			name:       "invalid config",
			configFile: "input/apply-vm-config-invalid.yaml",
			err:        true,
		},
	}

	for _, rt := range tests {
		t.Run(rt.name, func(t *testing.T) {
			// Copy the common base VM and set the input VM spec.
			newVM := baseVM.DeepCopy()
			if rt.baseSpec != nil {
				newVM.Spec = *rt.baseSpec
			}

			// Apply the input vm config on the base VM.
			configFilePath := fmt.Sprintf("testdata%c%s", filepath.Separator, rt.configFile)
			err = applyVMConfigFile(newVM, configFilePath)
			if (err != nil) != rt.err {
				t.Errorf("expected error %t, actual: %v", rt.err, err)
			}

			if !rt.err {
				// Check if the resulting VM config is as expected.

				// Set a fixed created time to avoid the result differences due to
				// creation time.
				createdTime := time.Date(2000, time.January, 1, 1, 0, 0, 0, time.UTC)
				newVM.SetCreated(runtime.Time{Time: metav1.Time{Time: createdTime}})

				// Convert VM object into json.
				newVMBytes, err := scheme.Serializer.EncodeJSON(newVM)
				if err != nil {
					t.Errorf("unexpected error while encoding object to json: %v", err)
				}

				// Construct golden file path.
				goldenFilePath := fmt.Sprintf("testdata%c%s", filepath.Separator, rt.golden)

				// Update the golden file if needed.
				if *update {
					t.Logf("updating golden file %s", goldenFilePath)
					if err := ioutil.WriteFile(goldenFilePath, newVMBytes, 0644); err != nil {
						t.Fatalf("failed to update apply-vm-config golden file: %s: %v", goldenFilePath, err)
					}
				}

				// Read golden file.
				wantOutput, err := ioutil.ReadFile(goldenFilePath)
				if err != nil {
					t.Fatalf("failed to read apply-vm-config golden file: %s: %v", goldenFilePath, err)
				}

				if !bytes.Equal(newVMBytes, wantOutput) {
					t.Errorf("expected VM config to be:\n%v\ngot VM config:\n%v", string(wantOutput), string(newVMBytes))
				}
			}
		})
	}
}

func TestApplyVMFlagOverrides(t *testing.T) {
	testImage := "weaveworks/ubuntu"
	testOCIRef, err := meta.NewOCIImageRef(testImage)
	if err != nil {
		t.Fatalf("error parsing image: %v", err)
	}

	tests := []struct {
		name            string
		createFlag      *CreateFlags
		wantCopyFiles   []api.FileMapping
		wantPortMapping meta.PortMappings
		wantSSH         api.SSH
		err             bool
	}{
		{
			name: "valid copy files flag",
			createFlag: &CreateFlags{
				VM:        &api.VM{},
				CopyFiles: []string{"/tmp/foo:/tmp/bar"},
			},
			wantCopyFiles: []api.FileMapping{
				{
					HostPath: "/tmp/foo",
					VMPath:   "/tmp/bar",
				},
			},
		},
		{
			name: "invalid copy files syntax",
			createFlag: &CreateFlags{
				VM:        &api.VM{},
				CopyFiles: []string{"foo:bar:baz"},
			},
			err: true,
		},
		{
			name: "invalid copy files paths - not absolute paths",
			createFlag: &CreateFlags{
				VM:        &api.VM{},
				CopyFiles: []string{"foo:bar"},
			},
			err: true,
		},
		{
			name: "valid port mapping",
			createFlag: &CreateFlags{
				VM:           &api.VM{},
				PortMappings: []string{"80:80"},
			},
			wantPortMapping: meta.PortMappings{
				meta.PortMapping{
					BindAddress: net.IPv4(0, 0, 0, 0),
					HostPort:    uint64(80),
					VMPort:      uint64(80),
					Protocol:    meta.ProtocolTCP,
				},
			},
		},
		{
			name: "invalid port mapping",
			createFlag: &CreateFlags{
				VM:           &api.VM{},
				PortMappings: []string{"1.1.1.1:foo:bar"},
			},
			err: true,
		},
		{
			name: "ssh public key set",
			createFlag: &CreateFlags{
				VM: &api.VM{},
				SSH: api.SSH{
					Generate:  true,
					PublicKey: "some-pub-key",
				},
			},
			wantSSH: api.SSH{
				Generate:  true,
				PublicKey: "some-pub-key",
			},
		},
	}

	for _, rt := range tests {
		t.Run(rt.name, func(t *testing.T) {
			vm := rt.createFlag.VM
			fs := flag.NewFlagSet("test", flag.ExitOnError)

			rt.createFlag.VM.Spec.Image.OCI = testOCIRef
			err := applyVMFlagOverrides(vm, rt.createFlag, fs)
			if (err != nil) != rt.err {
				t.Errorf("expected error %t, actual: %v", rt.err, err)
			}

			if !rt.err {
				// Check if copy files are set as expected.
				if len(rt.wantCopyFiles) > 0 {
					if !reflect.DeepEqual(vm.Spec.CopyFiles, rt.wantCopyFiles) {
						t.Errorf("expected VM.Spec.CopyFiles to be %v, actual: %v", rt.wantCopyFiles, rt.createFlag.VM.Spec.CopyFiles)
					}
				} else {
					// If the copy files map is empty, compare their sizes.
					if len(rt.wantCopyFiles) != len(rt.createFlag.VM.Spec.CopyFiles) {
						t.Errorf("expected VM.Spec.CopyFiles to be %v, actual: %v", rt.wantCopyFiles, rt.createFlag.VM.Spec.CopyFiles)
					}
				}

				// Check if port mappings are set as expected.
				if reflect.DeepEqual(rt.createFlag.VM.Spec.Network.Ports, rt.wantPortMapping) {
					t.Errorf("expected VM.Spec.Network.Ports to be %v, actual: %v", rt.wantPortMapping, rt.createFlag.VM.Spec.Network.Ports)
				}

				// Check if the SSH values are set as expected.
				if reflect.DeepEqual(rt.createFlag.VM.Spec.SSH, rt.wantSSH) {
					t.Errorf("expected VM.Spec.SSH to be %v, actual: %v", rt.wantSSH, rt.createFlag.VM.Spec.SSH)
				}
			}
		})
	}
}

func TestNewCreateOptions(t *testing.T) {
	tests := []struct {
		name       string
		createFlag *CreateFlags
		err        bool
	}{
		{
			name: "require-name with no name",
			createFlag: &CreateFlags{
				VM:          &api.VM{},
				RequireName: true,
			},
			err: true,
		},
		{
			name: "require-name with VM config",
			createFlag: &CreateFlags{
				ConfigFile:  fmt.Sprintf("testdata%c%s", filepath.Separator, "input/create-config-no-name.yaml"),
				RequireName: true,
			},
			err: true,
		},
	}

	for _, rt := range tests {
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

			// Set provider client.
			providers.Client = ic

			_, err = rt.createFlag.NewCreateOptions([]string{}, flag.NewFlagSet("test", flag.ExitOnError))
			if (err != nil) != rt.err {
				t.Errorf("expected error %t, actual: %v", rt.err, err)
			}
		})
	}
}
