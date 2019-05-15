package cmd

import (
	"fmt"
	"github.com/luxas/ignite/pkg/constants"
	"github.com/pkg/errors"
	"io"
	"os"
	"path"

	"github.com/luxas/ignite/pkg/errutils"
	"github.com/spf13/cobra"
)

// NewCmdRmi removes the given image
func NewCmdRmi(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rmi [id]",
		Short: "Remove an image",
		Run: func(cmd *cobra.Command, args []string) {
			err := RunRmi(out, cmd, args)
			errutils.Check(err)
		},
	}
	//cmd.Flags().StringP("output", "o", "", "Output format; available options are 'yaml', 'json' and 'short'")
	return cmd
}

// RunRmi runs when the Rmi command is invoked and removes the given image
func RunRmi(out io.Writer, cmd *cobra.Command, args []string) error {
	// TODO: Make this match without specifying the whole ID

	imageID := args[0]
	imageDir := path.Join(constants.IMAGE_DIR, imageID)

	if dir, err := os.Stat(imageDir); !os.IsNotExist(err) && dir.IsDir() {
		if err := os.RemoveAll(imageDir); err != nil {
			return errors.Wrapf(err, "unable to remove directory for image %s", imageID)
		}
	} else {
		return fmt.Errorf("not an image: %s", imageID)
	}

	fmt.Println(imageID)
	return nil
}
