package imgcmd

import (
	"io"

	"github.com/lithammer/dedent"
	"github.com/spf13/cobra"
	"github.com/weaveworks/ignite/cmd/ignite/cmd/cmdutil"
	"github.com/weaveworks/ignite/cmd/ignite/run"
)

// NewCmdImport imports new VM images from a tar source.
func NewCmdTarImport(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tarimport <OCI image>",
		Short: "Import new base images for VMs from a tar file",
		Long: dedent.Dedent(`
			Import OCI images as a base images for VMs from a tar file, takes in a file path to
			a tar file.
		`),
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(func() error {
				return run.ImportTarFile(args[0])
			}())
		},
	}
	return cmd
}
