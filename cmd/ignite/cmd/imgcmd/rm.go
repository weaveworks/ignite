package imgcmd

import (
	"github.com/spf13/pflag"
	"github.com/luxas/ignite/cmd/ignite/cmd/cmdutil"
	"github.com/luxas/ignite/cmd/ignite/run"
	"github.com/luxas/ignite/pkg/errutils"
	"github.com/spf13/cobra"
	"io"
)

// NewCmdRm removes an images
// TODO: Support removing multiple images at once
func NewCmdRm(out io.Writer) *cobra.Command {
	ro := &run.RmiOptions{}

	cmd := &cobra.Command{
		Use:   "rm [image]...",
		Short: "Remove VM base images",
		Long:  "TODO", // TODO: Long description
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			errutils.Check(func() error {
				var err error
				if ro.Images, err = cmdutil.MatchSingleImages(args); err != nil {
					return err
				}
				if ro.VMs, err = cmdutil.MatchAllVMs(true); err != nil {
					return err
				}
				return run.Rmi(ro)
			}())
		},
	}

	addRmiFlags(cmd.Flags(), ro)

	return cmd
}

func addRmiFlags(fs *pflag.FlagSet, ro *run.RmiOptions) {
	cmdutil.AddForceFlag(fs, &ro.Force)
}