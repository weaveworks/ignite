package cmd

import (
	"fmt"
	"github.com/luxas/ignite/pkg/constants"
	"github.com/luxas/ignite/pkg/errutils"
	"github.com/luxas/ignite/pkg/filter"
	"github.com/luxas/ignite/pkg/metadata"
	"github.com/luxas/ignite/pkg/metadata/imgmd"
	"github.com/luxas/ignite/pkg/metadata/kernmd"
	"github.com/luxas/ignite/pkg/metadata/vmmd"
	"github.com/luxas/ignite/pkg/util"
	"github.com/spf13/cobra"
	"io"
)

// NewCmdCreate creates a new VM from an image
func NewCmdCreate(out io.Writer) *cobra.Command {
	co := &createOptions{}

	cmd := &cobra.Command{
		Use:   "create [image] [kernel] [name]",
		Short: "Create a new containerized VM without starting it",
		Args:  cobra.MinimumNArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			err := RunCreate(out, cmd, args[0], args[1], args[2], co, false)
			errutils.Check(err)
		},
	}

	cmd.Flags().Int64Var(&co.cpus, "cpus", constants.VM_DEFAULT_CPUS, "VM vCPU count, 1 or even numbers between 1 and 32")
	cmd.Flags().Int64Var(&co.memory, "memory", constants.VM_DEFAULT_MEMORY, "VM RAM in MiB")
	return cmd
}

func RunCreate(out io.Writer, cmd *cobra.Command, imageMatch, kernelMatch, name string, co *createOptions, start bool) error {
	var image *imgmd.ImageMetadata
	var kernel *kernmd.KernelMetadata

	// Match a single Image using the ImageFilter
	if matches, err := filter.NewFilterer(imgmd.NewImageFilter(imageMatch), metadata.Image.Path(), imgmd.LoadImageMetadata); err == nil {
		if filterable, err := matches.Single(); err == nil {
			if image, err = imgmd.ToImageMetadata(filterable); err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		return err
	}

	// Match a single Kernel using the KernelFilter
	if matches, err := filter.NewFilterer(kernmd.NewKernelFilter(kernelMatch), metadata.Kernel.Path(), kernmd.LoadKernelMetadata); err == nil {
		if filterable, err := matches.Single(); err == nil {
			if kernel, err = kernmd.ToKernelMetadata(filterable); err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		return err
	}

	// Create a new ID for the VM
	vmID, err := util.NewID(constants.VM_DIR)
	if err != nil {
		return err
	}

	md := vmmd.NewVMMetadata(vmID, name, vmmd.NewVMObjectData(image.ID, kernel.ID, co.cpus, co.memory))

	// Save the metadata
	if err := md.Save(); err != nil {
		return err
	}

	// Perform the image copy
	// TODO: Replace this with overlayfs
	if err := md.CopyImage(); err != nil {
		return err
	}

	// If start is specified, start teh VM after creation
	if start {
		if err := RunStart(out, cmd, name); err != nil {
			return err
		}
	} else {
		// Print the ID of the created VM
		fmt.Println(md.ID)
	}

	return nil
}
