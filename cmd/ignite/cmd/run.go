package cmd

import (
	"github.com/luxas/ignite/pkg/errutils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"io"
)

type runOptions struct {
	createOptions
	startOptions
}

// NewCmdRun creates, starts (and attaches to) a Firecracker VM
func NewCmdRun(out io.Writer) *cobra.Command {
	ro := &runOptions{}

	cmd := &cobra.Command{
		Use:   "run [image] [kernel]",
		Short: "Create and start a new Firecracker VM",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			errutils.Check(func() error {
				var err error
				if ro.image, err = matchSingleImage(args[0]); err != nil {
					return err
				}
				if ro.kernel, err = matchSingleKernel(args[1]); err != nil {
					return err
				}
				return RunRun(ro)
			}())
		},
	}

	addRunFlags(cmd.Flags(), ro)
	return cmd
}

func addRunFlags(fs *pflag.FlagSet, ro *runOptions) {
	addCreateFlags(fs, &ro.createOptions)
	addStartFlags(fs, &ro.startOptions)
}

func RunRun(ro *runOptions) error {
	if err := RunCreate(&ro.createOptions); err != nil {
		return err
	}

	ro.startOptions.vm = ro.createOptions.vm

	if err := RunStart(&ro.startOptions); err != nil {
		return err
	}

	return nil
}
