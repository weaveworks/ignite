package run

import (
	"fmt"

	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/metadata"
	"github.com/weaveworks/ignite/pkg/metadata/kernmd"
	"github.com/weaveworks/ignite/pkg/util"
)

type ImportKernelOptions struct {
	Source      string
	Name        string
	KernelNames []*metadata.Name
}

func ImportKernel(ao *ImportKernelOptions) error {
	if !util.FileExists(ao.Source) {
		return fmt.Errorf("not a kernel image: %s", ao.Source)
	}

	// Create a new ID and directory for the kernel
	idHandler, err := util.NewID(constants.KERNEL_DIR)
	if err != nil {
		return err
	}
	defer idHandler.Remove()

	// Verify the name
	name, err := metadata.NewName(ao.Name, &ao.KernelNames)
	if err != nil {
		return err
	}

	md := kernmd.NewKernelMetadata(idHandler.ID, name)

	// Save the metadata
	if err := md.Save(); err != nil {
		return err
	}

	// Perform the copy
	if err := md.ImportKernel(ao.Source); err != nil {
		return err
	}

	fmt.Println(md.ID)

	idHandler.Success()
	return nil
}
