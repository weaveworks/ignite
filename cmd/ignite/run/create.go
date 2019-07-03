package run

import (
	"fmt"
	"path"
	"strings"

	"github.com/spf13/pflag"
	"github.com/weaveworks/ignite/pkg/apis/ignite/scheme"
	api "github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"
	"github.com/weaveworks/ignite/pkg/metadata"
	"github.com/weaveworks/ignite/pkg/metadata/imgmd"
	"github.com/weaveworks/ignite/pkg/metadata/kernmd"
	"github.com/weaveworks/ignite/pkg/metadata/loader"
	"github.com/weaveworks/ignite/pkg/metadata/vmmd"
	"github.com/weaveworks/ignite/pkg/util"
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
		VM: &api.VM{},
	}

	scheme.Scheme.Default(cf.VM)

	return cf
}

type CreateFlags struct {
	// TODO: Also respect CopyFiles, SSH, PortMappings, Networking mode, and kernel stuff from the config file
	CopyFiles  []string
	KernelName string
	SSH        *SSHFlag
	ConfigFile string
	VM         *api.VM
}

type createOptions struct {
	*CreateFlags
	image        *imgmd.ImageMetadata
	kernel       *kernmd.KernelMetadata
	allVMs       []metadata.Metadata
	newVM        *vmmd.VMMetadata
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

func (cf *CreateFlags) NewCreateOptions(l *loader.ResLoader, args []string) (*createOptions, error) {
	err := cf.parseArgsAndConfig(args)
	if err != nil {
		return nil, err
	}

	co := &createOptions{CreateFlags: cf}

	if allImages, err := l.Images(); err == nil {
		if co.image, err = allImages.MatchSingle(cf.VM.Spec.Image.Ref); err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}

	if len(cf.KernelName) == 0 {
		cf.KernelName = cf.VM.Spec.Image.Ref
	}

	if allKernels, err := l.Kernels(); err == nil {
		if co.kernel, err = allKernels.MatchSingle(cf.KernelName); err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}

	// The VM metadata needs the image and kernel IDs to be saved for now
	// TODO: Replace with pool/snapshotter
	cf.VM.Spec.Image.ID = co.image.GetUID()
	cf.VM.Spec.Kernel.ID = co.kernel.GetUID()

	if allVMs, err := l.VMs(); err == nil {
		co.allVMs = *allVMs
	} else {
		return nil, err
	}

	// Parse the --copy-files flag
	if co.fileMappings, err = parseFileMappings(co.CopyFiles); err != nil {
		return nil, err
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
	if co.newVM, err = vmmd.NewVMMetadata("", &name, co.VM); err != nil {
		return err
	}
	defer metadata.Cleanup(co.newVM, false) // TODO: Handle silent

	// Parse SSH key importing
	if err := co.parseSSH(&co.fileMappings); err != nil {
		return err
	}

	// Save the metadata
	if err := co.newVM.Save(); err != nil {
		return err
	}

	// Allocate the overlay file
	if err := co.newVM.AllocateOverlay(co.VM.Spec.DiskSize.Bytes()); err != nil {
		return err
	}

	// Copy the additional files to the overlay
	if err := co.newVM.CopyToOverlay(co.fileMappings); err != nil {
		return err
	}

	return metadata.Success(co.newVM)
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
