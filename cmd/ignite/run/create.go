package run

import (
	"fmt"
	"github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"
	"github.com/weaveworks/ignite/pkg/metadata"
	"github.com/weaveworks/ignite/pkg/metadata/vmmd"
	"github.com/weaveworks/ignite/pkg/snapshotter"
	"path"
	"strings"

	"github.com/weaveworks/ignite/pkg/source"

	"github.com/weaveworks/ignite/cmd/ignite/run/runutil"

	"github.com/spf13/pflag"

	"github.com/weaveworks/ignite/pkg/constants"
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
	Name         string
	CPUs         int64
	MemoryString string
	SizeString   string
	CopyFiles    []string
	KernelName   string
	KernelCmd    string
	SSH          *SSHFlag
}

type createOptions struct {
	*CreateFlags
	image        *snapshotter.Image
	kernel       *v1alpha1.ImageSource
	newVM        *snapshotter.VM
	size         v1alpha1.Size
	memory       v1alpha1.Size
	fileMappings map[string]string
}

func (cf *CreateFlags) NewCreateOptions(l *runutil.ResLoader, imageMatch string) (*createOptions, error) {
	var err error
	co := &createOptions{CreateFlags: cf}

	if allImages, err := l.Images(); err == nil {
		if co.image, err = allImages.MatchSingle(imageMatch); err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}

	if len(co.KernelName) == 0 {
		co.KernelName = constants.DEFAULT_KERNEL
	}

	co.kernel = &v1alpha1.ImageSource{
		Name: co.KernelName,
	}

	// Parse the given overlay size
	if err := co.size.UnmarshalText([]byte(co.SizeString)); err != nil {
		return nil, err
	}

	// Parse the given memory amount
	if err := co.memory.UnmarshalText([]byte(co.MemoryString)); err != nil {
		return nil, err
	}

	// Parse the --copy-files flag
	if co.fileMappings, err = parseFileMappings(co.CopyFiles); err != nil {
		return nil, err
	}

	// Parse SSH key importing
	if err = co.parseSSH(&co.fileMappings); err != nil {
		return nil, err
	}

	return co, nil
}

func Create(co *createOptions) error {
	// Verify the name
	name, err := metadata.NewName(co.Name, &co.allVMs)
	if err != nil {
		return err
	}

	// Create new metadata for the VM
	if co.newVM, err = vmmd.NewVMMetadata(nil, name,
		vmmd.NewVMObjectData(co.image.ID, metadata.IDFromSource(co.kernel), co.size, co.CPUs, co.memory, co.KernelCmd)); err != nil {
		return err
	}
	defer co.newVM.Cleanup(false) // TODO: Handle silent

	// Import the kernel and create the overlay
	_, err = co.image.CreateOverlay(co.kernel, co.size, co.newVM.ID)
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
