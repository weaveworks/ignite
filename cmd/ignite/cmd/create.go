package cmd

import (
	"fmt"
	"github.com/luxas/ignite/pkg/constants"
	"github.com/luxas/ignite/pkg/errutils"
	"github.com/luxas/ignite/pkg/filter"
	"github.com/luxas/ignite/pkg/metadata"
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

// RunCreate runs when the Create command is invoked
func RunCreate(out io.Writer, cmd *cobra.Command, args []string) error {
	// Load all image metadata as Filterable objects
	mdf, err := metadata.LoadMetadataFilterable(metadata.Image)
	if err != nil {
		return err
	}

	// Create an IDNameFilter to match a single image
	image, err := filter.NewFilterer(metadata.NewIDNameFilter(args[0])).Single(mdf)
	if err != nil {
		return err
	}

	// Load all kernel metadata as Filterable objects
	mdf, err = metadata.LoadMetadataFilterable(metadata.Kernel)
	if err != nil {
		return err
	}

	// Create an IDNameFilter to match a single kernel
	kernel, err := filter.NewFilterer(metadata.NewIDNameFilter(args[1])).Single(mdf)
	if err != nil {
		return err
	}

	// Create a new ID for the VM
	vmID, err := util.NewID(constants.VM_DIR)
	if err != nil {
		return err
	}

	md := metadata.NewVMMetadata(vmID, args[2], image.(*metadata.Metadata).ID, kernel.(*metadata.Metadata).ID)

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

//func loadVMMetadata(vmID string) (metadata.Filterable, error) {
//	md := &vmMetadata{
//		Metadata: &metadata.Metadata{
//			ID:         vmID,
//			Type:       metadata.VM,
//			ObjectData: &vmObjectData{},
//		},
//	}
//
//	if err := md.Load(); err != nil {
//		return nil, fmt.Errorf("failed to load VM metadata: %v", err)
//	}
//
//	return md, nil
//}
