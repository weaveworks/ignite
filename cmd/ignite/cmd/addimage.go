package cmd

import (
	"fmt"
	"github.com/luxas/ignite/pkg/constants"
	"github.com/luxas/ignite/pkg/errutils"
	"github.com/luxas/ignite/pkg/metadata/imgmd"
	"github.com/luxas/ignite/pkg/util"
	"github.com/spf13/cobra"
	"io"
)

// NewCmdAddImage imports an image for VM use
func NewCmdAddImage(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "addimage [path] [name]",
		Short: "Import an existing VM base image",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			err := RunAddImage(out, cmd, args)
			errutils.Check(err)
		},
	}

	//cmd.Flags().StringP("output", "o", "", "Output format; available options are 'yaml', 'json' and 'short'")
	return cmd
}

func RunAddImage(out io.Writer, cmd *cobra.Command, args []string) error {
	p := args[0]

	if !util.FileExists(p) {
		return fmt.Errorf("not an image file: %s", p)
	}

	// Create a new ID for the VM
	imageID, err := util.NewID(constants.IMAGE_DIR)
	if err != nil {
		return err
	}

	md := imgmd.NewImageMetadata(imageID, args[1])

	// Save the metadata
	if err := md.Save(); err != nil {
		return err
	}

	// Perform the image copy
	if err := md.ImportImage(p); err != nil {
		return err
	}

	fmt.Println(md.ID)

	return nil
}
