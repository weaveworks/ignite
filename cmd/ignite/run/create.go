package run

import (
	"fmt"
	"github.com/weaveworks/ignite/cmd/ignite/run/runutil"
	"path"
	"strings"

	"github.com/spf13/pflag"

	"github.com/c2h5oh/datasize"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/metadata"
	"github.com/weaveworks/ignite/pkg/metadata/imgmd"
	"github.com/weaveworks/ignite/pkg/metadata/kernmd"
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

type CreateFlags struct {
	Name       string
	CPUs       int64
	Memory     int64
	Size       string
	CopyFiles  []string
	KernelName string
	KernelCmd  string
	SSH        *SSHFlag
}

type createOptions struct {
	*CreateFlags
	image  *imgmd.ImageMetadata
	kernel *kernmd.KernelMetadata
	vms    []metadata.AnyMetadata
	newVM  *vmmd.VMMetadata
}

func (cf *CreateFlags) NewCreateOptions(l *runutil.ResourceLoader, imageMatch string) (*createOptions, error) {
	var err error
	co := &createOptions{CreateFlags: cf}

	if co.image, err = l.MatchSingleImage(imageMatch); err != nil {
		return nil, err
	}

	if len(cf.KernelName) == 0 {
		cf.KernelName = imageMatch
	}

	if co.kernel, err = l.MatchSingleKernel(cf.KernelName); err != nil {
		return nil, err
	}

	if co.vms, err = l.MatchAllVMs(true); err != nil {
		return nil, err
	}

	return co, nil
}

func Create(co *createOptions) error {
	// Parse the --copy-files flag
	fileMappings, err := parseFileMappings(co.CopyFiles)
	if err != nil {
		return err
	}

	// Create a new ID and directory for the VM
	idHandler, err := util.NewID(constants.VM_DIR)
	if err != nil {
		return err
	}
	defer idHandler.Remove()

	// Verify the name
	name, err := metadata.NewName(co.Name, &co.vms)
	if err != nil {
		return err
	}

	// Create new metadata for the VM and add to createOptions for further processing
	// This enables the generated VM metadata to pass straight to start and attach via run
	co.newVM = vmmd.NewVMMetadata(idHandler.ID, name, vmmd.NewVMObjectData(co.image.ID, co.kernel.ID, co.CPUs, co.Memory, co.KernelCmd))

	// Save the metadata
	if err := co.newVM.Save(); err != nil {
		return err
	}

	// Parse the given overlay size
	var size datasize.ByteSize
	if err := size.UnmarshalText([]byte(co.Size)); err != nil {
		return err
	}

	// Allocate the overlay file
	if err := co.newVM.AllocateOverlay(size.Bytes()); err != nil {
		return err
	}

	// Parse SSH key importing
	if err := co.parseSSH(&fileMappings); err != nil {
		return err
	}

	// Copy the additional files to the overlay
	if err := co.newVM.CopyToOverlay(fileMappings); err != nil {
		return err
	}

	// Print the ID of the created VM
	fmt.Println(co.newVM.ID)

	idHandler.Success()
	return nil
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
