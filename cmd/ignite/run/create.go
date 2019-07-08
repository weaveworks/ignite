package run

import (
	"fmt"
	"path"
	"strings"

	"github.com/weaveworks/ignite/pkg/client"
	"github.com/weaveworks/ignite/pkg/filter"

	"github.com/spf13/pflag"
	"github.com/weaveworks/ignite/pkg/apis/ignite/scheme"
	api "github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/metadata"
	"github.com/weaveworks/ignite/pkg/metadata/imgmd"
	"github.com/weaveworks/ignite/pkg/metadata/kernmd"
	"github.com/weaveworks/ignite/pkg/metadata/vmmd"
	"github.com/weaveworks/ignite/pkg/util"
)

// SSHFlag is the pflag.Value custom flag for ignite create --ssh
// TODO: Move SSHFlag somewhere else, e.g. cmdutils
type SSHFlag struct {
	value    *string
	generate bool
}

func (sf *SSHFlag) Set(x string) error {
	if x != "<path>" {
		sf.value = &x
	} else {
		sf.generate = true
	}

	return nil
}

func (sf *SSHFlag) String() string {
	if sf.value == nil {
		return ""
	}

	return *sf.value
}

func (sf *SSHFlag) Generate() bool {
	return sf.generate
}

func (sf *SSHFlag) Type() string {
	return ""
}

func (sf *SSHFlag) IsBoolFlag() bool {
	return true
}

// Parse sets the right values on the VM API object if requested to import or generate an SSH key
func (sf *SSHFlag) Parse(vm *api.VM) error {
	importKey := sf.String()
	// Check if an SSH key should be generated
	if sf.Generate() {
		// An empty struct means "generate"
		vm.Spec.SSH = &api.SSH{}
	} else if len(importKey) > 0 {
		// Always digest the public key
		if !strings.HasSuffix(importKey, ".pub") {
			importKey = fmt.Sprintf("%s.pub", importKey)
		}
		// verify the file exists
		if !util.FileExists(importKey) {
			return fmt.Errorf("invalid SSH key: %s", importKey)
		}

		// Set the SSH field on the API object
		vm.Spec.SSH = &api.SSH{
			PublicKey: importKey,
		}
	}

	return nil
}

var _ pflag.Value = &SSHFlag{}

func NewCreateFlags() *CreateFlags {
	cf := &CreateFlags{
		VM: &api.VM{},
	}

	scheme.Scheme.Default(cf.VM)

	return cf
}

type CreateFlags struct {
	// TODO: Also respect PortMappings, Networking mode, and kernel stuff from the config file
	CopyFiles  []string
	KernelName string
	SSH        *SSHFlag
	ConfigFile string
	VM         *api.VM
}

type createOptions struct {
	*CreateFlags
	image        *imgmd.Image
	kernel       *kernmd.Kernel
	newVM        *vmmd.VM
	fileMappings map[string]string
}

// parseArgsAndConfig resolves the image to use (the argument to the command)
// and the config file, if it needs to be loaded
func (cf *CreateFlags) parseArgsAndConfig(args []string) error {
	if len(args) == 1 {
		cf.VM.Spec.Image = &api.ImageClaim{
			Type: api.ImageSourceTypeDocker,
			Ref:  args[0],
		}
	}

	// Decode the config file if given
	if len(cf.ConfigFile) != 0 {
		// Marshal into a "clean" object, discard all flag input
		cf.VM = &api.VM{}
		if err := scheme.Serializer.DecodeFileInto(cf.ConfigFile, cf.VM); err != nil {
			return err
		}
	}

	// Specifying an image either way is mandatory
	if cf.VM.Spec.Image == nil || len(cf.VM.Spec.Image.Ref) == 0 {
		return fmt.Errorf("you must specify an image to run either via CLI args or a config file")
	}
	return nil
}

func (cf *CreateFlags) NewCreateOptions(args []string) (*createOptions, error) {
	err := cf.parseArgsAndConfig(args)
	if err != nil {
		return nil, err
	}

	co := &createOptions{CreateFlags: cf}

	if image, err := client.Images().Find(filter.NewIDNameFilter(cf.VM.Spec.Image.Ref)); err == nil {
		co.image = &imgmd.Image{image}
	} else {
		return nil, err
	}

	if len(cf.KernelName) == 0 {
		cf.KernelName = cf.VM.Spec.Image.Ref
	}

	if kernel, err := client.Kernels().Find(filter.NewIDNameFilter(cf.KernelName)); err == nil {
		co.kernel = &kernmd.Kernel{kernel}
	} else {
		return nil, err
	}

	// The VM metadata needs the image and kernel IDs to be saved for now
	// TODO: Replace with pool/snapshotter
	cf.VM.Spec.Image.UID = co.image.GetUID()
	cf.VM.Spec.Kernel.UID = co.kernel.GetUID()

	// Parse the --copy-files flag
	cf.VM.Spec.CopyFiles, err = parseFileMappings(co.CopyFiles)
	if err != nil {
		return nil, err
	}
	return co, nil
}

func Create(co *createOptions) error {
	// Verify the name
	name, err := metadata.NewName(co.VM.Name, meta.KindVM)
	if err != nil {
		return err
	}

	// Parse SSH key importing
	if err := co.SSH.Parse(co.VM); err != nil {
		return err
	}

	// Create new metadata for the VM
	if co.newVM, err = vmmd.NewVM("", &name, co.VM); err != nil {
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
