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

type flagData struct {
	cpus   int64
	memory int64
}

// NewCmdCreate creates a new VM from an image
func NewCmdCreate(out io.Writer) *cobra.Command {
	fd := &flagData{}

	cmd := &cobra.Command{
		Use:   "create [image] [kernel] [name]",
		Short: "Create a new containerized VM without starting it",
		Args:  cobra.MinimumNArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			err := RunCreate(out, cmd, args, fd)
			errutils.Check(err)
		},
	}

	cmd.Flags().Int64Var(&fd.cpus, "cpus", constants.VM_DEFAULT_CPUS, "VM vCPU count, 1 or even numbers between 1 and 32")
	cmd.Flags().Int64Var(&fd.memory, "memory", constants.VM_DEFAULT_MEMORY, "VM RAM in MiB")
	return cmd
}

func RunCreate(out io.Writer, cmd *cobra.Command, args []string, fd *flagData) error {
	var image *imgmd.ImageMetadata
	var kernel *kernmd.KernelMetadata

	// Match a single Image using the ImageFilter
	if matches, err := filter.NewFilterer(imgmd.NewImageFilter(args[0]), metadata.Image.Path(), imgmd.LoadImageMetadata); err == nil {
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
	if matches, err := filter.NewFilterer(kernmd.NewKernelFilter(args[1]), metadata.Kernel.Path(), kernmd.LoadKernelMetadata); err == nil {
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

	md := vmmd.NewVMMetadata(vmID, args[2], vmmd.NewVMObjectData(image.ID, kernel.ID, fd.cpus, fd.memory))

	// Save the metadata
	if err := md.Save(); err != nil {
		return err
	}

	// Perform the image copy
	// TODO: Replace this with overlayfs
	if err := md.CopyImage(); err != nil {
		return err
	}

	fmt.Println(vmID)

	return nil
}
