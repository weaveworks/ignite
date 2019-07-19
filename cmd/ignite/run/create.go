package run

import (
	"fmt"
	"path"
	"strings"

	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	"github.com/weaveworks/ignite/pkg/apis/ignite/scheme"
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/metadata"
	"github.com/weaveworks/ignite/pkg/metadata/imgmd"
	"github.com/weaveworks/ignite/pkg/metadata/kernmd"
	"github.com/weaveworks/ignite/pkg/metadata/vmmd"
	"github.com/weaveworks/ignite/pkg/operations"
	"github.com/weaveworks/ignite/pkg/providers"
)

func NewCreateFlags() *CreateFlags {
	return &CreateFlags{
		VM: client.VMs().New(),
	}
}

type CreateFlags struct {
	PortMappings []string
	CopyFiles    []string
	// This is a placeholder value here for now.
	// If it was set using flags, it will be copied over to
	// the API type. TODO: When we later have internal types
	// this can go away
	SSH        api.SSH
	ConfigFile string
	VM         *api.VM
}

type createOptions struct {
	*CreateFlags
	image  *imgmd.Image
	kernel *kernmd.Kernel
	newVM  *vmmd.VM
}

func (cf *CreateFlags) constructVMFromCLI(args []string) error {
	if len(args) == 1 {
		ociRef, err := meta.NewOCIImageRef(args[0])
		if err != nil {
			return err
		}

		cf.VM.Spec.Image.OCIClaim.Ref = ociRef
	}

	// Parse the --copy-files flag
	var err error
	cf.VM.Spec.CopyFiles, err = parseFileMappings(cf.CopyFiles)
	if err != nil {
		return err
	}

	// Parse the given port mappings
	if cf.VM.Spec.Network.Ports, err = meta.ParsePortMappings(cf.PortMappings); err != nil {
		return err
	}

	// If the SSH flag was set, copy it over to the API type
	if cf.SSH.Generate || cf.SSH.PublicKey != "" {
		cf.VM.Spec.SSH = &cf.SSH
	}

	return nil
}

func (cf *CreateFlags) NewCreateOptions(args []string) (*createOptions, error) {
	// Decode the config file if given, or construct the VM based off flags and args
	if len(cf.ConfigFile) != 0 {
		// Marshal into a "clean" object, discard all flag input
		cf.VM = &api.VM{}
		if err := scheme.Serializer.DecodeFileInto(cf.ConfigFile, cf.VM); err != nil {
			return nil, err
		}
	} else {
		if err := cf.constructVMFromCLI(args); err != nil {
			return nil, err
		}
	}

	// Specifying an image either way is mandatory
	if cf.VM.Spec.Image.OCIClaim.Ref.IsUnset() {
		return nil, fmt.Errorf("you must specify an image to run either via CLI args or a config file")
	}

	co := &createOptions{CreateFlags: cf}

	// Get the image, or import it if it doesn't exist
	var err error
	co.image, err = operations.FindOrImportImage(providers.Client, cf.VM.Spec.Image.OCIClaim.Ref)
	if err != nil {
		return nil, err
	}

	// Populate relevant data from the Image on the VM object
	cf.VM.SetImage(co.image.Image)

	// Get the kernel, or import it if it doesn't exist
	co.kernel, err = operations.FindOrImportKernel(providers.Client, cf.VM.Spec.Kernel.OCIClaim.Ref)
	if err != nil {
		return nil, err
	}

	// Populate relevant data from the Kernel on the VM object
	cf.VM.SetKernel(co.kernel.Kernel)
	return co, nil
}

func Create(co *createOptions) error {
	// Create new metadata for the VM
	var err error
	if co.newVM, err = vmmd.NewVM(co.VM, providers.Client); err != nil {
		return err
	}
	defer metadata.Cleanup(co.newVM, false) // TODO: Handle silent

	// Save the metadata
	if err := co.newVM.Save(); err != nil {
		return err
	}

	// Allocate and populate the overlay file
	if err := co.newVM.AllocateAndPopulateOverlay(); err != nil {
		return err
	}

	return metadata.Success(co.newVM)
}

// TODO: Move this to meta, or an helper in api
func parseFileMappings(fileMappings []string) ([]api.FileMapping, error) {
	result := []api.FileMapping{}

	for _, fileMapping := range fileMappings {
		files := strings.Split(fileMapping, ":")
		if len(files) != 2 {
			return nil, fmt.Errorf("--copy-files requires the /host/path:/vm/path form")
		}

		src, dest := files[0], files[1]
		if !path.IsAbs(src) || !path.IsAbs(dest) {
			return nil, fmt.Errorf("--copy-files path arguments must be absolute")
		}

		result = append(result, api.FileMapping{
			HostPath: src,
			VMPath:   dest,
		})
	}

	return result, nil
}
