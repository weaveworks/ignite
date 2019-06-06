package vmcmd

import (
	"io"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/ignite/cmd/ignite/cmd/cmdutil"
	"github.com/weaveworks/ignite/cmd/ignite/run"
	"github.com/weaveworks/ignite/pkg/errutils"
)

// NewCmdStart starts a VM
func NewCmdStart(out io.Writer) *cobra.Command {
	so := &run.StartOptions{}

	cmd := &cobra.Command{
		Use:   "start [vm]",
		Short: "Start a VM",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			errutils.Check(func() error {
				var err error
				if so.VM, err = cmdutil.MatchSingleVM(args[0]); err != nil {
					return err
				}
				return run.Start(so)
			}())
		},
	}

	addStartFlags(cmd.Flags(), so)
	return cmd
}

func addStartFlags(fs *pflag.FlagSet, so *run.StartOptions) {
	cmdutil.AddInteractiveFlag(fs, &so.Interactive)
	fs.StringSliceVarP(&so.PortMappings, "ports", "p", nil, "Map host ports to VM ports")
	fs.BoolVarP(&so.Debug, "debug", "d", false, "Debug mode, keep container after VM shutdown")
}
