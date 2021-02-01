package vmcmd

import (
	"io"

	"github.com/lithammer/dedent"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/ignite/cmd/ignite/cmd/cmdutil"
	"github.com/weaveworks/ignite/cmd/ignite/run"
	networkflag "github.com/weaveworks/ignite/pkg/network/flag"
	"github.com/weaveworks/ignite/pkg/providers"
	runtimeflag "github.com/weaveworks/ignite/pkg/runtime/flag"
)

// NewCmdStart starts a VM
func NewCmdStart(out io.Writer) *cobra.Command {
	sf := &run.StartFlags{}

	cmd := &cobra.Command{
		Use:   "start <vm>",
		Short: "Start a VM",
		Long: dedent.Dedent(`
			Start the given VM. The VM is matched by prefix based on its ID and name.
			If the interactive flag (-i, --interactive) is specified, attach to the
			VM after starting.
		`),
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(func() error {
				so, err := sf.NewStartOptions(args[0])
				if err != nil {
					return err
				}

				return run.Start(so, cmd.Flags())
			}())
		},
	}

	addStartFlags(cmd.Flags(), sf)

	// NOTE: Since the run command combines the create and start command flags,
	// to avoid redefining runtime, network, and id-prefix flags in the run command,
	// they are defined separately here, and re-used from addCreateFlags.
	runtimeflag.RuntimeVar(cmd.Flags(), &providers.RuntimeName)
	networkflag.NetworkPluginVar(cmd.Flags(), &providers.NetworkPluginName)

	return cmd
}

func addStartFlags(fs *pflag.FlagSet, sf *run.StartFlags) {
	cmdutil.AddInteractiveFlag(fs, &sf.Interactive)
	fs.BoolVarP(&sf.Debug, "debug", "d", false, "Debug mode, keep container after VM shutdown")
	fs.StringSliceVar(&sf.IgnoredPreflightErrors, "ignore-preflight-checks", []string{}, "A list of checks whose errors will be shown as warnings. Example: 'BinaryInPath,Port,ExistingFile'. Value 'all' ignores errors from all checks.")
}
