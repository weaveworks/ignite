package run

import (
	"fmt"
	"path"
	"strings"

	"github.com/spf13/pflag"
	"github.com/weaveworks/ignite/pkg/apis/ignite/scheme"
	"github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"
	"github.com/weaveworks/ignite/pkg/metadata"
	"github.com/weaveworks/ignite/pkg/metadata/imgmd"
	ignitemeta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/metadata/loader"
	"github.com/weaveworks/ignite/pkg/metadata/vmmd"
	"github.com/weaveworks/ignite/pkg/snapshotter"
	"github.com/weaveworks/ignite/pkg/source"
	"github.com/weaveworks/ignite/pkg/util"
	"github.com/weaveworks/ignite/pkg/filtering/filter"
)

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

var _ pflag.Value = &SSHFlag{}

const vmAuthorizedKeys = "/root/.ssh/authorized_keys"

func NewCreateFlags() *CreateFlags {
	cf := &CreateFlags{
		VM: &v1alpha1.VM{},
	}
	scheme.Scheme.Default(cf.VM)
	return cf
}

type CreateFlags struct {
	// TODO: Also respect CopyFiles, SSH, PortMappings, Networking mode, and kernel stuff from the config file
	CopyFiles  []string
	KernelName string
	KernelCmd  string
	SSH        *SSHFlag
	ConfigFile string
	VM         *v1alpha1.VM
}

type createOptions struct {
	*CreateFlags
	image        *snapshotter.Image
	kernel       *v1alpha1.ImageSource
	newVM        *snapshotter.VM
	fileMappings map[string]string
}

func (cf *CreateFlags) NewCreateOptions(ss *snapshotter.Snapshotter, args []string) (*createOptions, error) {
	co := &createOptions{CreateFlags: cf}
	
	if len(args) == 1 {
		co.VM.Spec.Image = &v1alpha1.ImageClaim{
			Type: v1alpha1.ImageSourceTypeDocker,
			Ref:  args[0],
		}
	}

	// Decode the config file if given
	if len(co.ConfigFile) != 0 {
		// Marshal into a "clean" object, discard all flag input
		co.VM = &v1alpha1.VM{}
		if err := scheme.DecodeFileInto(co.ConfigFile, co.VM); err != nil {
			return nil, err
		}
	}

	if co.VM.Spec.Image == nil || len(co.VM.Spec.Image.Ref) == 0 {
		return nil, fmt.Errorf("you must specify an image to run either via CLI args or a config file")
	}

	var err error
	if co.image, err = ss.GetImage(filter.NewIDNameFilter(co.VM.Spec.Image.Ref)); err != nil {
		return nil, err
	}

	// Parse the --copy-files flag
	if co.fileMappings, err = parseFileMappings(co.CopyFiles); err != nil {
		return nil, err
	}

	if len(co.KernelName) == 0 {
		co.KernelName = constants.DEFAULT_KERNEL
	}

	// TODO: Filter that checks if image contains wanted kernel and if it has the correct size
	if kernel, err := ss.GetKernel(filter.NewKernelFilter(co.KernelName, co.image, co.size)); err != nil {
		switch err.(type) {
		case snapshotter.ErrNonexistent:

		}
	}

	co.kernel = &v1alpha1.ImageSource{
		Name: co.KernelName,
	}

	return co, nil
}

func Create(co *createOptions) error {
	// Verify the name
	name, err := metadata.NewName(co.VM.Name, &co.allVMs)
	if err != nil {
		return err
	}

	// Create new metadata for the VM
	if co.newVM, err = vmmd.NewVMMetadata(nil, name,
		vmmd.NewVMObjectData(co.image.ID, metadata.IDFromSource(co.kernel), co.VM.Spec.DiskSize, int64(co.CPUs), co.VM.Spec.Memory, co.KernelCmd)); err != nil {
		return err
	}
	defer co.newVM.Cleanup(false) // TODO: Handle silent

	// Parse SSH key importing
	if err := co.parseSSH(&co.fileMappings); err != nil {
		return err
	}

	// Import the kernel and create the overlay
	_, err = co.image.CreateOverlay(co.kernel, co.VM.Spec.DiskSize, co.newVM.ID)
	if err != nil {
		return err
	}

	// Copy the additional files to the overlay
	// TODO: Support this for the overlay in the image
	if err := co.newVM.CopyToOverlay(co.fileMappings); err != nil {
		return err
	}

	// Save the metadata
	if err := co.newVM.Save(); err != nil {
		return err
	}

	// Save the image metadata to register the new overlays
	if err := co.image.Save(); err != nil {
		return err
	}

	return co.newVM.Success()
}

func parseFileMappings(fileMappings []string) (map[string]string, error) {
	result := map[string]string{}

	for _, fileMapping := range fileMappings {
		files := strings.Split(fileMapping, ":")
		if len(files) != 2 {
			return nil, fmt.Errorf("--copy-files requires the /host/path:/vm/path form")
		}

		src, dest := files[0], files[1]
		if !path.IsAbs(src) || !path.IsAbs(dest) {
			return nil, fmt.Errorf("--copy-files path arguments must be absolute")
		}

		result[src] = dest
	}

	return result, nil
}

// If we're requested to import/generate an SSH key, add that to fileMappings
func (co *createOptions) parseSSH(fileMappings *map[string]string) error {
	importKey := co.SSH.String()

	if co.SSH.Generate() {
		pubKeyPath, err := co.newVM.NewSSHKeypair()
		if err != nil {
			return err
		}

		(*fileMappings)[pubKeyPath] = vmAuthorizedKeys
	} else if len(importKey) > 0 {
		// Always digest the public key
		if !strings.HasSuffix(importKey, ".pub") {
			importKey = fmt.Sprintf("%s.pub", importKey)
		}

		if !util.FileExists(importKey) {
			return fmt.Errorf("invalid SSH key: %s", importKey)
		}

		(*fileMappings)[importKey] = vmAuthorizedKeys
	}

	return nil
}
