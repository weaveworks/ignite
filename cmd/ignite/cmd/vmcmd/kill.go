package vmcmd

import (
	"github.com/weaveworks/ignite/cmd/ignite/run/runutil"
	"io"

	"github.com/lithammer/dedent"

	"github.com/spf13/cobra"
	"github.com/weaveworks/ignite/cmd/ignite/run"
	"github.com/weaveworks/ignite/pkg/errutils"
)

// NewCmdKill kills a VM
func NewCmdKill(out io.Writer) *cobra.Command {
	so := &run.StopOptions{Kill: true}

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
			errutils.Check(func() error {
				var err error
				if so.VMs, err = runutil.MatchSingleVMs(args); err != nil {
					return err
				}
				return run.Stop(so)
			}())
		},
	}

	return cmd
}
