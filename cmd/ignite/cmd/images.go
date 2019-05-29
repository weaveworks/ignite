package cmd

import (
	"github.com/luxas/ignite/cmd/ignite/run"
	"github.com/luxas/ignite/pkg/errutils"
	"github.com/spf13/cobra"
	"io"
)

// NewCmdImages lists the images for your Firecracker VM.
func NewCmdImages(out io.Writer) *cobra.Command {
	io := &run.ImagesOptions{}

	cmd := &cobra.Command{
		Use:   "images",
		Short: "List available images",
		Run: func(cmd *cobra.Command, args []string) {
			errutils.Check(func() error {
				var err error
				if io.Images, err = matchAllImages(); err != nil {
					return err
				}
				return run.Images(io)
			}())
		},
	}

	return cmd
}
