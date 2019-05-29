package cmd

import (
	"github.com/luxas/ignite/cmd/ignite/run"
	"github.com/luxas/ignite/pkg/errutils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"io"
)

// NewCmdVMRm removes a VM
func NewCmdVMRm(out io.Writer) *cobra.Command {
	ro := &run.RmOptions{}

	cmd := &cobra.Command{
		Use:   "rm [vm]",
		Short: "Remove a VM",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			errutils.Check(func() error {
				var err error
				if ro.VM, err = matchSingleVM(args[0]); err != nil {
					return err
				}
				return run.Rm(ro)
			}())
		},
	}

	addVMRmFlags(cmd.Flags(), ro)
	return cmd
}

func addVMRmFlags(fs *pflag.FlagSet, ro *run.RmOptions) {
	fs.BoolVarP(&ro.Force, "force", "f", false, "Kill VM if running before removal")
}
