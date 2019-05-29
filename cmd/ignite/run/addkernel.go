package run

import (
	"fmt"
	"github.com/luxas/ignite/pkg/constants"
	"github.com/luxas/ignite/pkg/metadata/kernmd"
	"github.com/luxas/ignite/pkg/util"
)

type AddKernelOptions struct {
	Source string
	Name   string
}

func AddKernel(ao *AddKernelOptions) error {
	if !util.FileExists(ao.Source) {
		return fmt.Errorf("not a kernel image: %s", ao.Source)
	}

	// Create a new ID for the VM
	kernelID, err := util.NewID(constants.KERNEL_DIR)
	if err != nil {
		return err
	}

	md := kernmd.NewKernelMetadata(kernelID, ao.Name)

	// Save the metadata
	if err := md.Save(); err != nil {
		return err
	}

	// Perform the image copy
	if err := md.ImportKernel(ao.Source); err != nil {
		return err
	}

	fmt.Println(md.ID)

	return nil
}
