package cmd

import (
	"github.com/luxas/ignite/cmd/ignite/run"
	"github.com/luxas/ignite/pkg/errutils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"io"
)

// NewCmdRun creates, starts (and attaches to) a Firecracker VM
func NewCmdRun(out io.Writer) *cobra.Command {
	ro := &run.RunOptions{}

	cmd := &cobra.Command{
		Use:   "run [image] [kernel]",
		Short: "Create and start a new Firecracker VM",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			errutils.Check(func() error {
				var err error
				if ro.Image, err = matchSingleImage(args[0]); err != nil {
					return err
				}
				if ro.Kernel, err = matchSingleKernel(args[1]); err != nil {
					return err
				}
				return run.Run(ro)
			}())
		},
	}

	addRunFlags(cmd.Flags(), ro)
	return cmd
}

func addRunFlags(fs *pflag.FlagSet, ro *run.RunOptions) {
	addCreateFlags(fs, &ro.CreateOptions)
	addStartFlags(fs, &ro.StartOptions)
}
