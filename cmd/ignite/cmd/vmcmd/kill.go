package vmcmd

import (
	"io"

	"github.com/lithammer/dedent"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/ignite/cmd/ignite/cmd/cmdutil"
	"github.com/weaveworks/ignite/cmd/ignite/run"
	"github.com/weaveworks/ignite/pkg/util"
)

// NewCmdKill kills running VMs
func NewCmdKill(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kill <vm>...",
		Short: "Kill running VMs",
		Long: dedent.Dedent(`
			Kill (force stop) one or multiple VMs. The VMs are matched by prefix based
			on their ID and name. To kill multiple VMs, chain the matches separated
			by spaces.
		`),
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(func() error {
				so, err := (&run.StopFlags{Kill: true}).NewStopOptions(args)
				if err != nil {
					return err
				}

				return run.Stop(so)
			}())
		},
	}

	addKillFlags(cmd.Flags())
	return cmd
}

func addKillFlags(fs *pflag.FlagSet) {
	cmdutil.AddNamePrefixFlag(fs, &util.NamePrefix)
}
