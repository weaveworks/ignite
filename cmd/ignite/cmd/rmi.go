package cmd

import (
	"fmt"
	"github.com/luxas/ignite/pkg/errutils"
	"github.com/luxas/ignite/pkg/filter"
	"github.com/luxas/ignite/pkg/metadata"
	"github.com/luxas/ignite/pkg/metadata/imgmd"
	"github.com/spf13/cobra"
	"io"
	"os"
)

// NewCmdRmi removes the given image
func NewCmdRmi(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rmi [id]",
		Short: "Remove an image",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			err := RunRmi(out, cmd, args)
			errutils.Check(err)
		},
	}
	//cmd.Flags().StringP("output", "o", "", "Output format; available options are 'yaml', 'json' and 'short'")
	return cmd
}

func RunRmi(out io.Writer, cmd *cobra.Command, args []string) error {
	var image *imgmd.ImageMetadata

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

	// TODO: Check that the given image is not used by any VMs
	if err := os.RemoveAll(image.ObjectPath()); err != nil {
		return fmt.Errorf("unable to remove directory for %s %q: %v", image.Type, image.ID, err)
	}

	fmt.Println(image.ID)
	return nil
}
