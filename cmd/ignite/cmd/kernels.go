package cmd

import (
	"fmt"
	"github.com/luxas/ignite/pkg/constants"
	"io"
	"io/ioutil"

	"github.com/luxas/ignite/pkg/errutils"
	"github.com/spf13/cobra"
)

// NewCmdImages lists the images for your Firecracker VM.
func NewCmdKernels(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kernels",
		Short: "List available kernels",
		Run: func(cmd *cobra.Command, args []string) {
			err := RunKernels(out, cmd)
			errutils.Check(err)
		},
	}
	//cmd.Flags().StringP("output", "o", "", "Output format; available options are 'yaml', 'json' and 'short'")
	return cmd
}

// RunImages runs when the Images command is invoked and lists the images
func RunKernels(out io.Writer, cmd *cobra.Command) error {
	ids, err := ioutil.ReadDir(constants.KERNEL_DIR)
	if err != nil {
		return err
	}

	for _, id := range ids {
		fmt.Println(id.Name())
	}

	return nil
}
