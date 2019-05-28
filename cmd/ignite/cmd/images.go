package cmd

import (
	"fmt"
	"github.com/luxas/ignite/pkg/errutils"
	"github.com/luxas/ignite/pkg/metadata/imgmd"
	"github.com/luxas/ignite/pkg/util"
	"github.com/spf13/cobra"
	"io"
)

type imagesOptions struct {
	images []*imgmd.ImageMetadata
}

// NewCmdImages lists the images for your Firecracker VM.
func NewCmdImages(out io.Writer) *cobra.Command {
	io := &imagesOptions{}

	cmd := &cobra.Command{
		Use:   "images",
		Short: "List available images",
		Run: func(cmd *cobra.Command, args []string) {
			errutils.Check(func() error {
				var err error
				if io.images, err = matchAllImages(); err != nil {
					return err
				}
				return RunImages(io)
			}())
		},
	}

	return cmd
}

func RunImages(io *imagesOptions) error {
	o := util.NewOutput()
	defer o.Flush()

	o.Write("IMAGE ID", "CREATED", "SIZE", "NAME")
	for _, md := range io.images {
		size, err := md.Size()
		if err != nil {
			return fmt.Errorf("failed to get size for %s %q: %v", md.Type, md.ID, err)
		}

		o.Write(md.ID, md.Created, util.ByteCountDecimal(size), md.Name)
	}

	return nil
}
