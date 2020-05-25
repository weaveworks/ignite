package run

import (
	"fmt"
	"net"
	"path/filepath"
	"reflect"
	"testing"

	flag "github.com/spf13/pflag"
	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/libgitops/pkg/runtime"
)

func TestConstructVMFromCLI(t *testing.T) {
	testImage := "weaveworks/ubuntu"
	testOCIRef, err := meta.NewOCIImageRef(testImage)
	if err != nil {
		t.Fatalf("error parsing image: %v", err)
	}

	testSandboxImage := "weaveworks/ignite:test"
	sandboxOCIRef, err := meta.NewOCIImageRef(testSandboxImage)
	if err != nil {
		t.Fatalf("error parsing image: %v", err)
	}

	tests := []struct {
		name             string
		createFlag       *CreateFlags
		args             []string
		wantCopyFiles    []api.FileMapping
		wantPortMapping  meta.PortMappings
		wantSSH          api.SSH
		wantSandboxImage meta.OCIImageRef
		err              bool
	}{
		{
			name: "with VM name and image arg",
			createFlag: &CreateFlags{
				VM: &api.VM{
					ObjectMeta: runtime.ObjectMeta{
						Name: "fooVM",
					},
				},
			},
			args: []string{testImage},
		},
		{
			name: "with invalid image reference",
			createFlag: &CreateFlags{
				VM: &api.VM{},
			},
			args: []string{"foo:bar:baz"},
			err:  true,
		},
		{
			name: "valid copy files flag",
			createFlag: &CreateFlags{
				VM:        &api.VM{},
				CopyFiles: []string{"/tmp/foo:/tmp/bar"},
			},
			args: []string{testImage},
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
			args: []string{testImage},
			err:  true,
		},
		{
			name: "invalid copy files paths - not absolute paths",
			createFlag: &CreateFlags{
				VM:        &api.VM{},
				CopyFiles: []string{"foo:bar"},
			},
			args: []string{testImage},
			err:  true,
		},
		{
			name: "valid port mapping",
			createFlag: &CreateFlags{
				VM:           &api.VM{},
				PortMappings: []string{"80:80"},
			},
			args: []string{testImage},
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
			args: []string{testImage},
			err:  true,
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
			args: []string{testImage},
			wantSSH: api.SSH{
				Generate:  true,
				PublicKey: "some-pub-key",
			},
		},
		{
			name: "with no VM name and --require-name flag set",
			createFlag: &CreateFlags{
				VM:          &api.VM{},
				RequireName: true,
			},
			args: []string{testImage},
			err:  true,
		},
		{
			name: "with sandbox image",
			createFlag: &CreateFlags{
				VM: &api.VM{
					Spec: api.VMSpec{
						Sandbox: api.VMSandboxSpec{
							OCI: sandboxOCIRef,
						},
					},
				},
			},
			args:             []string{testImage},
			wantSandboxImage: sandboxOCIRef,
		},
	}

	for _, rt := range tests {
		t.Run(rt.name, func(t *testing.T) {
			err := rt.createFlag.constructVMFromCLI(rt.args)
			if (err != nil) != rt.err {
				t.Errorf("expected error %t, actual: %v", rt.err, err)
			}

			if !rt.err {
				// Check if the VM image is set as expected.
				if rt.createFlag.VM.Spec.Image.OCI != testOCIRef {
					t.Errorf("expected VM.Spec.Image.OCI to be %q, actual: %q", testOCIRef.String(), rt.createFlag.VM.Spec.Image.OCI.String())
				}

				// Check if copy files are set as expected.
				if len(rt.wantCopyFiles) > 0 {
					if !reflect.DeepEqual(rt.createFlag.VM.Spec.CopyFiles, rt.wantCopyFiles) {
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

				// Check if the sandbox image is set as expected.
				if rt.createFlag.VM.Spec.Sandbox.OCI != rt.wantSandboxImage {
					t.Errorf("expected VM.Spec.Sandbox to be %v, actual: %v", rt.wantSandboxImage, rt.createFlag.VM.Spec.Sandbox.OCI)
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
			_, err := rt.createFlag.NewCreateOptions([]string{}, flag.NewFlagSet("test", flag.ExitOnError))
			if (err != nil) != rt.err {
				t.Errorf("expected error %t, actual: %v", rt.err, err)
			}
		})
	}
}
