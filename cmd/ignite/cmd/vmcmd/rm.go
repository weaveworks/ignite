package vmcmd

import (
	"github.com/luxas/ignite/cmd/ignite/cmd/cmdutil"
	"github.com/luxas/ignite/cmd/ignite/run"
	"github.com/luxas/ignite/pkg/errutils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"io"
)

// NewCmdRm removes a VM
func NewCmdRm(out io.Writer) *cobra.Command {
	ro := &run.RmOptions{}

	cmd := &cobra.Command{
		Use:   "rm [vm]...",
		Short: "Remove VMs",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			errutils.Check(func() error {
				var err error
				if ro.VMs, err = cmdutil.MatchSingleVMs(args); err != nil {
					return err
				}
				return run.Rm(ro)
			}())
		},
	}

	addRmFlags(cmd.Flags(), ro)
	return cmd
}

func addRmFlags(fs *pflag.FlagSet, ro *run.RmOptions) {
	fs.BoolVarP(&ro.Force, "force", "f", false, "Kill VM if running before removal")
}
