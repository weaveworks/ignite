package run

import (
	"fmt"
	"path"
	"strings"

	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	"github.com/weaveworks/ignite/pkg/apis/ignite/scheme"
	"github.com/weaveworks/ignite/pkg/apis/ignite/validation"
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/dmlegacy"
	"github.com/weaveworks/ignite/pkg/metadata"
	"github.com/weaveworks/ignite/pkg/operations"
	"github.com/weaveworks/ignite/pkg/providers"
)

// NewCreateFlags returns an initialized CreateFlags instance.
func NewCreateFlags() *CreateFlags {
	return &CreateFlags{}
}

// CreateFlags contains flag variables for create command.
type CreateFlags struct {
	PortMappings []string
	CopyFiles    []string
	// This is a placeholder value here for now.
	// If it was set using flags, it will be copied over to
	// the API type. TODO: When we later have internal types
	// this can go away
	SSH           api.SSH
	ConfigFile    string
	Labels        []string
	VMName        string
	VMCPUs        uint64
	KernelCmdLine string
	VMMemory      meta.Size
	VMDiskSize    meta.Size
	// Use string for KernelOCI instead of OCIImageRef because OCIImageRef can't
	// be created without setting a valid value. Need zero valued variable to
	// know if the flag was set or not.
	KernelOCI string
	VMStorage api.VMStorageSpec
}

type createOptions struct {
	*CreateFlags
	VM *api.VM
}

func (cf *CreateFlags) constructVMFromCLI(vm *api.VM, args []string) error {
	if len(args) == 1 {
		ociRef, err := meta.NewOCIImageRef(args[0])
		if err != nil {
			return err
		}

		// This overwrites any default VM image or from component config file.
		vm.Spec.Image.OCI = ociRef
	}

	// Parse the --copy-files flag
	var err error
	if vm.Spec.CopyFiles, err = parseFileMappings(cf.CopyFiles); err != nil {
		return err
	}

	// Parse the given port mappings
	if vm.Spec.Network.Ports, err = meta.ParsePortMappings(cf.PortMappings); err != nil {
		return err
	}

	// If the SSH flag was set, copy it over to the API type
	if cf.SSH.Generate || cf.SSH.PublicKey != "" {
		vm.Spec.SSH = &cf.SSH
	}

	if cf.VMName != "" {
		vm.SetName(cf.VMName)
	}

	if cf.VMCPUs > 0 {
		vm.Spec.CPUs = cf.VMCPUs
	}

	if cf.KernelCmdLine != "" {
		vm.Spec.Kernel.CmdLine = cf.KernelCmdLine
	}

	if cf.VMMemory != meta.EmptySize {
		vm.Spec.Memory = cf.VMMemory
	}

	if cf.VMDiskSize != meta.EmptySize {
		vm.Spec.DiskSize = cf.VMDiskSize
	}

	if cf.KernelOCI != "" {
		if vm.Spec.Kernel.OCI, err = meta.NewOCIImageRef(cf.KernelOCI); err != nil {
			return err
		}
	}

	if len(cf.VMStorage.Volumes) > 0 || len(cf.VMStorage.VolumeMounts) > 0 {
		vm.Spec.Storage = cf.VMStorage
	}

	return nil
}

// NewCreateOptions constructs a VM create option based on the create flags and
// arguments.
func (cf *CreateFlags) NewCreateOptions(args []string) (*createOptions, error) {
	co := &createOptions{CreateFlags: cf}
	// Initialize a VM instance.
	co.VM = providers.Client.VMs().New()
	co.VM.SetName(cf.VMName)

	// If component config is defined, populate VM spec with it.
	if providers.ComponentConfig != nil {
		// When no VM config is present in component config file, VM spec
		// contains all the defaults of a VM object, which is the same as the
		// new VM object created above. It's safe to overwrite the VM spec.
		co.VM.Spec = providers.ComponentConfig.Spec.VM
	}

	// Decode the config file if given, or construct the VM based off flags and args
	if len(cf.ConfigFile) != 0 {
		// Marshal into a "clean" object, discard all flag input
		co.VM = &api.VM{}
		if err := scheme.Serializer.DecodeFileInto(cf.ConfigFile, co.VM); err != nil {
			return nil, err
		}
	} else {
		if err := cf.constructVMFromCLI(co.VM, args); err != nil {
			return nil, err
		}
	}

	// Validate the VM object
	if err := validation.ValidateVM(co.VM).ToAggregate(); err != nil {
		return nil, err
	}

	// Get the image, or import it if it doesn't exist
	var err error
	if _, err = operations.FindOrImportImage(providers.Client, co.VM.Spec.Image.OCI); err != nil {
		return nil, err
	}

	// Get the kernel, or import it if it doesn't exist
	if _, err = operations.FindOrImportKernel(providers.Client, co.VM.Spec.Kernel.OCI); err != nil {
		return nil, err
	}

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
