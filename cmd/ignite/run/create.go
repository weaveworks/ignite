package run

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"github.com/weaveworks/ignite/cmd/ignite/cmd/cmdutil"
	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	"github.com/weaveworks/ignite/pkg/apis/ignite/scheme"
	"github.com/weaveworks/ignite/pkg/apis/ignite/validation"
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/config"
	"github.com/weaveworks/ignite/pkg/dmlegacy"
	"github.com/weaveworks/ignite/pkg/metadata"
	"github.com/weaveworks/ignite/pkg/operations"
	"github.com/weaveworks/ignite/pkg/providers"
	"github.com/weaveworks/ignite/pkg/util"

	flag "github.com/spf13/pflag"
	patchutil "github.com/weaveworks/libgitops/pkg/util/patch"
	"sigs.k8s.io/yaml"
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

type CreateOptions struct {
	*CreateFlags
	image  *api.Image
	kernel *api.Kernel
}

func (cf *CreateFlags) NewCreateOptions(args []string, fs *flag.FlagSet) (*CreateOptions, error) {
	// Create a new base VM and configure it by combining the component config,
	// VM config file and flags.
	baseVM := providers.Client.VMs().New()

	// If component config is in use, set the VMDefaults on the base VM.
	if providers.ComponentConfig != nil {
		baseVM.Spec = providers.ComponentConfig.Spec.VMDefaults
	}

	// Resolve registry configuration used for pulling image if required.
	cmdutil.ResolveRegistryConfigDir()

	// Initialize the VM's Prefixer
	baseVM.Status.IDPrefix = providers.IDPrefix
	// Set the runtime and network-plugin on the VM, then override the global config.
	baseVM.Status.Runtime.Name = providers.RuntimeName
	baseVM.Status.Network.Plugin = providers.NetworkPluginName
	// Populate the runtime and network-plugin providers.
	if err := config.SetAndPopulateProviders(providers.RuntimeName, providers.NetworkPluginName); err != nil {
		return nil, err
	}

	// Set the passed image argument on the new VM spec.
	// Image is necessary while serializing the VM spec.
	if len(args) == 1 {
		ociRef, err := meta.NewOCIImageRef(args[0])
		if err != nil {
			return nil, err
		}
		baseVM.Spec.Image.OCI = ociRef
	}

	// Generate a VM name and UID if not set yet.
	if err := metadata.SetNameAndUID(baseVM, providers.Client); err != nil {
		return nil, err
	}

	// Apply the VM config on the base VM, if a VM config is given.
	if len(cf.ConfigFile) != 0 {
		if err := applyVMConfigFile(baseVM, cf.ConfigFile); err != nil {
			return nil, err
		}
	}

	// Apply flag overrides.
	if err := applyVMFlagOverrides(baseVM, cf, fs); err != nil {
		return nil, err
	}

	// If --require-name is true, VM name must be provided.
	if cf.RequireName && len(baseVM.Name) == 0 {
		return nil, fmt.Errorf("must set VM name, flag --require-name set")
	}

	// Assign the new VM to the configFlag.
	cf.VM = baseVM

	// Validate the VM object.
	if err := validation.ValidateVM(cf.VM).ToAggregate(); err != nil {
		return nil, err
	}

	co := &CreateOptions{CreateFlags: cf}

	// Get the image, or import it if it doesn't exist.
	var err error
	co.image, err = operations.FindOrImportImage(providers.Client, cf.VM.Spec.Image.OCI)
	if err != nil {
		return nil, err
	}

	// Populate relevant data from the Image on the VM object.
	cf.VM.SetImage(co.image)

	// Get the kernel, or import it if it doesn't exist.
	co.kernel, err = operations.FindOrImportKernel(providers.Client, cf.VM.Spec.Kernel.OCI)
	if err != nil {
		return nil, err
	}

	// Populate relevant data from the Kernel on the VM object.
	cf.VM.SetKernel(co.kernel)
	return co, nil
}

// applyVMConfigFile patches a given base VM with the VM config in a given
// config file.
func applyVMConfigFile(baseVM *api.VM, configFile string) error {
	vmConfigBytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}

	// Marshal into a new object to extract VM image if any.
	fileVM := &api.VM{}
	if err := scheme.Serializer.DecodeInto(vmConfigBytes, fileVM); err != nil {
		return err
	}

	// Image is necessary while serializing the VM spec. Override VM image
	// provided in the VM config, in case it's not set in the base VM, to
	// avoid serialization error.
	if !fileVM.Spec.Image.OCI.IsUnset() {
		baseVM.Spec.Image.OCI = fileVM.Spec.Image.OCI
	}

	// Create a patcher to patch the base VM with the VM config.
	p := patchutil.NewPatcher(scheme.Serializer)

	// Ensure the VM config is in json. The patcher only accepts json
	// encoded bytes.
	vmConfigJSONBytes, err := yaml.YAMLToJSON(vmConfigBytes)
	if err != nil {
		return err
	}

	// Serialize the base VM into json encoded bytes.
	baseVMBytes, err := scheme.Serializer.EncodeJSON(baseVM)
	if err != nil {
		return err
	}

	// Apply the VM config on the base VM.
	resultVMBytes, err := p.Apply(baseVMBytes, vmConfigJSONBytes, baseVM.GroupVersionKind())
	if err != nil {
		return err
	}

	if err := scheme.Serializer.DecodeInto(resultVMBytes, baseVM); err != nil {
		return err
	}

	return nil
}

// applyVMFlagOverrides overrides a given VM configs with the flag options.
func applyVMFlagOverrides(baseVM *api.VM, cf *CreateFlags, fs *flag.FlagSet) error {
	var err error

	// Override configs passed through flags.
	if fs.Changed("name") {
		baseVM.Name = cf.VM.Name
	}
	if fs.Changed("cpus") {
		baseVM.Spec.CPUs = cf.VM.Spec.CPUs
	}
	if fs.Changed("kernel-args") {
		baseVM.Spec.Kernel.CmdLine = cf.VM.Spec.Kernel.CmdLine
	}
	if fs.Changed("memory") {
		baseVM.Spec.Memory = cf.VM.Spec.Memory
	}
	if fs.Changed("size") {
		baseVM.Spec.DiskSize = cf.VM.Spec.DiskSize
	}
	if fs.Changed("kernel-image") {
		baseVM.Spec.Kernel.OCI = cf.VM.Spec.Kernel.OCI
	}
	if fs.Changed("sandbox-image") {
		baseVM.Spec.Sandbox.OCI = cf.VM.Spec.Sandbox.OCI
	}
	if fs.Changed("volumes") {
		baseVM.Spec.Storage = cf.VM.Spec.Storage
	}

	if len(cf.CopyFiles) > 0 {
		// Parse the --copy-files flag.
		baseVM.Spec.CopyFiles, err = parseFileMappings(cf.CopyFiles)
		if err != nil {
			return err
		}
	}

	if len(cf.PortMappings) > 0 {
		// Parse the given port mappings.
		baseVM.Spec.Network.Ports, err = meta.ParsePortMappings(cf.PortMappings)
		if err != nil {
			return err
		}
	}

	// If the SSH flag was set, copy it over to the API type
	if cf.SSH.Generate || cf.SSH.PublicKey != "" {
		baseVM.Spec.SSH = &cf.SSH
	}

	return err
}

func Create(co *CreateOptions) (err error) {
	// Generate a random UID and Name
	if err = metadata.SetNameAndUID(co.VM, providers.Client); err != nil {
		return
	}
	// Set VM labels.
	if err = metadata.SetLabels(co.VM, co.Labels); err != nil {
		return
	}
	defer util.DeferErr(&err, func() error { return metadata.Cleanup(co.VM, false) })

	if err = providers.Client.VMs().Set(co.VM); err != nil {
		return
	}

	// Allocate and populate the overlay file
	if err = dmlegacy.AllocateAndPopulateOverlay(co.VM); err != nil {
		return
	}

	err = metadata.Success(co.VM)

	return
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
