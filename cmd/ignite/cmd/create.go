package cmd

import (
	"fmt"
	"github.com/luxas/ignite/pkg/constants"
	"github.com/luxas/ignite/pkg/errutils"
	"github.com/luxas/ignite/pkg/filter"
	"github.com/luxas/ignite/pkg/metadata"
	"github.com/luxas/ignite/pkg/metadata/vmmd"
	"github.com/luxas/ignite/pkg/util"
	"github.com/spf13/cobra"
	"io"
)

// NewCmdCreate creates a new VM from an image
func NewCmdCreate(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [image] [kernel] [name]",
		Short: "Create a new containerized VM without starting it",
		Args:  cobra.MinimumNArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			err := RunCreate(out, cmd, args)
			errutils.Check(err)
		},
	}

	//cmd.Flags().StringP("output", "o", "", "Output format; available options are 'yaml', 'json' and 'short'")
	return cmd
}

func RunCreate(out io.Writer, cmd *cobra.Command, args []string) error {
	var image *metadata.Metadata
	var kernel *metadata.Metadata

	// Match a single Image using the IDNameFilter
	if matches, err := filter.NewFilterer(metadata.NewIDNameFilter(args[0]), metadata.Image.Path(), metadata.LoadImageMetadata); err == nil {
		if filterable, err := matches.Single(); err == nil {
			if image, err = metadata.ToMetadata(filterable); err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		return err
	}

	// Match a single Image using the IDNameFilter
	if matches, err := filter.NewFilterer(metadata.NewIDNameFilter(args[1]), metadata.Kernel.Path(), metadata.LoadKernelMetadata); err == nil {
		if filterable, err := matches.Single(); err == nil {
			if kernel, err = metadata.ToMetadata(filterable); err != nil {
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

	md := vmmd.NewVMMetadata(vmID, args[2], image.ID, kernel.ID)

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
