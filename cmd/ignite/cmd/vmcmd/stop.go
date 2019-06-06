package vmcmd

import (
	"io"

	"github.com/weaveworks/ignite/cmd/ignite/cmd/cmdutil"
	"github.com/weaveworks/ignite/cmd/ignite/run"
	"github.com/weaveworks/ignite/pkg/errutils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// NewCmdStop stops a VM
func NewCmdStop(out io.Writer) *cobra.Command {
	so := &run.StopOptions{}

	cmd := &cobra.Command{
		Use:   "stop [vm]...",
		Short: "Stop running VMs",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			errutils.Check(func() error {
				var err error
				if so.VMs, err = cmdutil.MatchSingleVMs(args); err != nil {
					return err
				}
				return run.Stop(so)
			}())
		},
	}

	addStopFlags(cmd.Flags(), so)
	return cmd
}

func addStopFlags(fs *pflag.FlagSet, so *run.StopOptions) {
	fs.BoolVarP(&so.Kill, "force-kill", "f", false, "Force kill the VM")
}
