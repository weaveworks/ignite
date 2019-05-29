package cmd

import (
	"github.com/luxas/ignite/cmd/ignite/run"
	"github.com/luxas/ignite/pkg/errutils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"io"
)

// NewCmdVMStart starts a VM
func NewCmdVMStart(out io.Writer) *cobra.Command {
	so := &run.StartOptions{}

	cmd := &cobra.Command{
		Use:   "start [vm]",
		Short: "Start a VM",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			errutils.Check(func() error {
				var err error
				if so.VM, err = matchSingleVM(args[0]); err != nil {
					return err
				}
				return run.Start(so)
			}())
		},
	}

	addVMStartFlags(cmd.Flags(), so)
	return cmd
}

func addVMStartFlags(fs *pflag.FlagSet, so *run.StartOptions) {
	addInteractiveFlag(fs, &so.Interactive)
}
