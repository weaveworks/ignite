package imgcmd

import (
	"io"

	"github.com/lithammer/dedent"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/weaveworks/ignite/cmd/ignite/cmd/cmdutil"
	"github.com/weaveworks/ignite/cmd/ignite/run"
	"github.com/weaveworks/ignite/pkg/providers"
	runtimeflag "github.com/weaveworks/ignite/pkg/runtime/flag"
)

// NewCmdImport imports a new VM image
func NewCmdImport(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import <OCI image>",
		Short: "Import a new base image for VMs",
		Long: dedent.Dedent(`
			Import an OCI image as a base image for VMs, takes in a Docker image identifier.
			This importing is done automatically when the "run" or "create" commands are run.
			The import step is essentially a cache for images to be used later when running VMs.
		`),
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(func() error {
				_, err := run.ImportImage(args[0])
				return err
			}())
		},
	}

	addImportFlags(cmd.Flags())
	return cmd
}

func addImportFlags(fs *pflag.FlagSet) {
	runtimeflag.RuntimeVar(fs, &providers.RuntimeName)
	cmdutil.AddClientConfigDirFlag(fs, &providers.ClientConfigDir)
}
