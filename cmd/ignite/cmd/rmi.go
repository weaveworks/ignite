package cmd

import (
	"fmt"
	"github.com/luxas/ignite/pkg/errutils"
	"github.com/luxas/ignite/pkg/metadata/imgmd"
	"github.com/spf13/cobra"
	"io"
	"os"
)

type rmiOptions struct {
	image *imgmd.ImageMetadata
}

// NewCmdRmi removes the given image
func NewCmdRmi(out io.Writer) *cobra.Command {
	ro := &rmiOptions{}

	cmd := &cobra.Command{
		Use:   "rmi [id]",
		Short: "Remove an image",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			errutils.Check(func() error {
				var err error
				if ro.image, err = matchSingleImage(args[0]); err != nil {
					return err
				}
				return RunRmi(ro)
			}())
		},
	}

	return cmd
}

func RunRmi(ro *rmiOptions) error {
	// TODO: Check that the given image is not used by any VMs
	if err := os.RemoveAll(ro.image.ObjectPath()); err != nil {
		return fmt.Errorf("unable to remove directory for %s %q: %v", ro.image.Type, ro.image.ID, err)
	}

	fmt.Println(ro.image.ID)
	return nil
}
