package run

import (
	"fmt"
	"path"
	"strings"

	"github.com/c2h5oh/datasize"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/metadata"
	"github.com/weaveworks/ignite/pkg/metadata/imgmd"
	"github.com/weaveworks/ignite/pkg/metadata/kernmd"
	"github.com/weaveworks/ignite/pkg/metadata/vmmd"
	"github.com/weaveworks/ignite/pkg/util"
)

type CreateOptions struct {
	Image     *imgmd.ImageMetadata
	Kernel    *kernmd.KernelMetadata
	vm        *vmmd.VMMetadata
	Name      string
	CPUs      int64
	Memory    int64
	Size      string
	CopyFiles []string
	KernelCmd string
	VMNames   []*metadata.Name
}

func Create(co *CreateOptions) error {
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
	name, err := metadata.NewName(co.Name, &co.VMNames)
	if err != nil {
		return err
	}

	// Create new metadata for the VM and add to createOptions for further processing
	// This enables the generated VM metadata to pass straight to start and attach via run
	co.vm = vmmd.NewVMMetadata(idHandler.ID, name, vmmd.NewVMObjectData(co.Image.ID, co.Kernel.ID, co.CPUs, co.Memory, co.KernelCmd))

	// Save the metadata
	if err := co.vm.Save(); err != nil {
		return err
	}

	// Parse the given overlay size
	var size datasize.ByteSize
	if err := size.UnmarshalText([]byte(co.Size)); err != nil {
		return err
	}

	// Allocate the overlay file
	if err := co.vm.AllocateOverlay(size.Bytes()); err != nil {
		return err
	}

	// Copy the additional files to the overlay
	if err := co.vm.CopyToOverlay(fileMappings); err != nil {
		return err
	}

	// Print the ID of the created VM
	fmt.Println(co.vm.ID)

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
