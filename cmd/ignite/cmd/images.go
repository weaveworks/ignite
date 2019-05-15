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
		Short: "Imagesute a command in a Firecracker VM",
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
	images, err := ioutil.ReadDir(constants.IMAGE_DIR)
	if err != nil {
		return err
	}

	for _, image := range images {
		fmt.Println(image.Name())
	}

	return nil
}
