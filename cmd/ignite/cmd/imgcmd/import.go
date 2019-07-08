package imgcmd

import (
	"io"

	"github.com/lithammer/dedent"
	"github.com/spf13/cobra"
	"github.com/weaveworks/ignite/cmd/ignite/run"
	"github.com/weaveworks/ignite/pkg/errutils"
)

// NewCmdImport imports  a new VM image
func NewCmdImport(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import <source>",
		Short: "Import a new base image for VMs",
		Long: dedent.Dedent(`
			Import a new base image for VMs, takes in a Docker image as the source.
			The base image is an ext4 block device file, which contains a root filesystem.

			If a kernel is found in the image, /boot/vmlinux is extracted from it
			and imported to a kernel with the same name.

			Example usage:
			    $ ignite image import luxas/ubuntu-base:18.04
		`),
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			errutils.Check(func() error {
				io, err := run.NewImportOptions(args[0])
				if err != nil {
					return err
				}
				return run.Import(io)
			}())
		},
	}
	return cmd
}
