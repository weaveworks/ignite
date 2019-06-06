package vmcmd

import (
	"io"

	"github.com/lithammer/dedent"

	"github.com/weaveworks/ignite/cmd/ignite/cmd/cmdutil"
	"github.com/weaveworks/ignite/cmd/ignite/run"

	"github.com/spf13/cobra"
	"github.com/weaveworks/ignite/pkg/errutils"
)

// NewCmdLogs gets the logs for a Firecracker VM
func NewCmdLogs(out io.Writer) *cobra.Command {
	lo := &run.LogsOptions{}

	cmd := &cobra.Command{
		Use:   "logs [vm]",
		Short: "Get the logs for a running VM",
		Long: dedent.Dedent(`
			Show the logs for the given VM. The VM needs to be running (its backing
			container needs to exist). The VM is matched by prefix based on its ID and name.
		`),
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			errutils.Check(func() error {
				var err error
				if lo.VM, err = cmdutil.MatchSingleVM(args[0]); err != nil {
					return err
				}
				return run.Logs(lo)
			}())
		},
	}

	return cmd
}
