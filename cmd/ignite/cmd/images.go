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
func NewCmdImages(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "images",
		Short: "List available images",
		Run: func(cmd *cobra.Command, args []string) {
			err := RunImages(out, cmd)
			errutils.Check(err)
		},
	}
	//cmd.Flags().StringP("output", "o", "", "Output format; available options are 'yaml', 'json' and 'short'")
	return cmd
}

// RunImages runs when the Images command is invoked and lists the images
func RunImages(out io.Writer, cmd *cobra.Command) error {
	ids, err := ioutil.ReadDir(constants.IMAGE_DIR)
	if err != nil {
		return err
	}

	for _, id := range ids {
		fmt.Println(id.Name())
	}

	return nil
}
