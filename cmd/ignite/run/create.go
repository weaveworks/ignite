package run

import (
	"fmt"
	"path"
	"strings"

	flag "github.com/spf13/pflag"

	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	"github.com/weaveworks/ignite/pkg/apis/ignite/scheme"
	"github.com/weaveworks/ignite/pkg/apis/ignite/validation"
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/dmlegacy"
	"github.com/weaveworks/ignite/pkg/metadata"
	"github.com/weaveworks/ignite/pkg/operations"
	"github.com/weaveworks/ignite/pkg/providers"
)

func NewCreateFlags() *CreateFlags {
	return &CreateFlags{
		VM: providers.Client.VMs().New(),
	}
}

type CreateFlags struct {
	PortMappings []string
	CopyFiles    []string
	// This is a placeholder value here for now.
	// If it was set using flags, it will be copied over to
	// the API type. TODO: When we later have internal types
	// this can go away
	SSH         api.SSH
	ConfigFile  string
	VM          *api.VM
	Labels      []string
	RequireName bool
}

type createOptions struct {
	*CreateFlags
	image  *api.Image
	kernel *api.Kernel
}

func (cf *CreateFlags) constructVMFromCLI(args []string) error {
	if len(args) == 1 {
		ociRef, err := meta.NewOCIImageRef(args[0])
		if err != nil {
			return err
		}

		cf.VM.Spec.Image.OCI = ociRef
	}

	// Parse the --copy-files flag
	var err error
	if cf.VM.Spec.CopyFiles, err = parseFileMappings(cf.CopyFiles); err != nil {
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

	// If --require-name is true, VM name must be provided.
	if cf.RequireName && len(cf.VM.Name) == 0 {
		return fmt.Errorf("must pass a VM name, flag --require-name set")
	}

	return nil
}

func (cf *CreateFlags) NewCreateOptions(args []string, fs *flag.FlagSet) (*createOptions, error) {
	// If component config file is in use, create a new base VM, use the VM
	// spec defined in the config and based on the flags, override the configs.
	if providers.ComponentConfig != nil {
		newVM := providers.Client.VMs().New()
		newVM.Spec = providers.ComponentConfig.Spec.VM

		if fs.Changed("name") {
			newVM.Name = cf.VM.Name
		}

		if fs.Changed("cpus") {
			newVM.Spec.CPUs = cf.VM.Spec.CPUs
		}

		if fs.Changed("kernel-args") {
			newVM.Spec.Kernel.CmdLine = cf.VM.Spec.Kernel.CmdLine
		}

		if fs.Changed("memory") {
			newVM.Spec.Memory = cf.VM.Spec.Memory
		}

		if fs.Changed("size") {
			newVM.Spec.DiskSize = cf.VM.Spec.DiskSize
		}

		if fs.Changed("kernel-image") {
			newVM.Spec.Kernel.OCI = cf.VM.Spec.Kernel.OCI
		}

		if fs.Changed("sandbox-image") {
			newVM.Spec.Sandbox.OCI = cf.VM.Spec.Sandbox.OCI
		}

		if fs.Changed("volumes") {
			newVM.Spec.Storage = cf.VM.Spec.Storage
		}

		cf.VM = newVM
	}

	// Decode the config file if given, or construct the VM based off flags and args
	if len(cf.ConfigFile) != 0 {
		// Marshal into a "clean" object, discard all flag input
		cf.VM = &api.VM{}
		if err := scheme.Serializer.DecodeFileInto(cf.ConfigFile, cf.VM); err != nil {
			return nil, err
		}
		// If --require-name is true, VM name must be provided.
		if cf.RequireName && len(cf.VM.Name) == 0 {
			return nil, fmt.Errorf("must set VM name, flag --require-name set")
		}
	} else {
		if err := cf.constructVMFromCLI(args); err != nil {
			return nil, err
		}
	}

	// Validate the VM object
	if err := validation.ValidateVM(cf.VM).ToAggregate(); err != nil {
		return nil, err
	}

	co := &createOptions{CreateFlags: cf}

	// Get the image, or import it if it doesn't exist
	var err error
	co.image, err = operations.FindOrImportImage(providers.Client, cf.VM.Spec.Image.OCI)
	if err != nil {
		return nil, err
	}

	// Populate relevant data from the Image on the VM object
	cf.VM.SetImage(co.image)

	// Get the kernel, or import it if it doesn't exist
	co.kernel, err = operations.FindOrImportKernel(providers.Client, cf.VM.Spec.Kernel.OCI)
	if err != nil {
		return nil, err
	}

	// Populate relevant data from the Kernel on the VM object
	cf.VM.SetKernel(co.kernel)
	return co, nil
}

func Create(co *createOptions) error {
	// Generate a random UID and Name
	if err := metadata.SetNameAndUID(co.VM, providers.Client); err != nil {
		return err
	}
	// Set VM labels.
	if err := metadata.SetLabels(co.VM, co.Labels); err != nil {
		return err
	}
	defer metadata.Cleanup(co.VM, false) // TODO: Handle silent

	if err := providers.Client.VMs().Set(co.VM); err != nil {
		return err
	}

	// Allocate and populate the overlay file
	if err := dmlegacy.AllocateAndPopulateOverlay(co.VM); err != nil {
		return err
	}

	return metadata.Success(co.VM)
}

// TODO: Move this to meta, or a helper in API
func parseFileMappings(fileMappings []string) ([]api.FileMapping, error) {
	result := make([]api.FileMapping, 0, len(fileMappings))

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
