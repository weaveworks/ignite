package kerncmd

import (
	"io"

	"github.com/luxas/ignite/cmd/ignite/cmd/cmdutil"
	"github.com/luxas/ignite/cmd/ignite/run"
	"github.com/luxas/ignite/pkg/errutils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// NewCmdRm removes a kernel
// TODO: Support removing multiple kernels at once
func NewCmdRm(out io.Writer) *cobra.Command {
	ro := &run.RmkOptions{}

	cmd := &cobra.Command{
		Use:   "rm [kernel]...",
		Short: "Remove kernels",
		Long:  "TODO", // TODO: Long description
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			errutils.Check(func() error {
				var err error
				if ro.Kernels, err = cmdutil.MatchSingleKernels(args); err != nil {
					return err
				}
				if ro.VMs, err = cmdutil.MatchAllVMs(true); err != nil {
					return err
				}
				return run.Rmk(ro)
			}())
		},
	}

	addRmkFlags(cmd.Flags(), ro)
	return cmd
}

func addRmkFlags(fs *pflag.FlagSet, ro *run.RmkOptions) {
	cmdutil.AddForceFlag(fs, &ro.Force)
}
