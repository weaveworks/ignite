package cmd

import (
	"fmt"
	"github.com/luxas/ignite/pkg/constants"
	"github.com/luxas/ignite/pkg/errutils"
	"github.com/luxas/ignite/pkg/metadata/kernmd"
	"github.com/luxas/ignite/pkg/util"
	"github.com/spf13/cobra"
	"io"
)

// NewCmdAddKernel adds a new kernel for VM use
func NewCmdAddKernel(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "addkernel [path] [name]",
		Short: "Add an uncompressed kernel image for VM use",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			err := RunAddKernel(out, cmd, args)
			errutils.Check(err)
		},
	}

	//cmd.Flags().StringP("output", "o", "", "Output format; available options are 'yaml', 'json' and 'short'")
	return cmd
}

func RunAddKernel(out io.Writer, cmd *cobra.Command, args []string) error {
	p := args[0]

	if !util.FileExists(p) {
		return fmt.Errorf("not a kernel image: %s", p)
	}

	// Create a new ID for the VM
	kernelID, err := util.NewID(constants.KERNEL_DIR)
	if err != nil {
		return err
	}

	md := kernmd.NewKernelMetadata(kernelID, args[1])

	// Save the metadata
	if err := md.Save(); err != nil {
		return err
	}

	// Perform the image copy
	// TODO: Replace this with overlayfs
	if err := md.ImportKernel(p); err != nil {
		return err
	}

	fmt.Println(kernelID)

	return nil
}
