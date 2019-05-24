package cmd

import (
	"fmt"
	"github.com/luxas/ignite/pkg/errutils"
	"github.com/luxas/ignite/pkg/filter"
	"github.com/luxas/ignite/pkg/metadata"
	"github.com/luxas/ignite/pkg/metadata/imgmd"
	"github.com/luxas/ignite/pkg/util"
	"github.com/spf13/cobra"
	"io"
	"time"
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

func RunImages(out io.Writer, cmd *cobra.Command) error {
	var mds []*imgmd.ImageMetadata

	// Match all Images using the ImageFilter
	if matches, err := filter.NewFilterer(imgmd.NewImageFilter(""), metadata.Image.Path(), imgmd.LoadImageMetadata); err == nil {
		if all, err := matches.All(); err == nil {
			if mds, err = imgmd.ToImageMetadataAll(all); err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		return err
	}

	o := util.NewOutput()
	defer o.Flush()

	o.Write("IMAGE ID\tCREATED\tSIZE\tNAME")
	for _, md := range mds {
		size, err := md.Size()
		if err != nil {
			return fmt.Errorf("failed to get size for %s %q: %v", md.Type, md.ID, err)
		}

		o.Write(fmt.Sprintf("%s\t%s\t%d\t%s", md.ID, md.Created.Format(time.UnixDate), size, md.Name))
	}

	return nil
}
