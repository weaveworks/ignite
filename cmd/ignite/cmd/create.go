package cmd

import (
	"github.com/luxas/ignite/cmd/ignite/run"
	"github.com/luxas/ignite/pkg/constants"
	"github.com/luxas/ignite/pkg/errutils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"io"
)

// NewCmdCreate creates a new VM from an image
func NewCmdCreate(out io.Writer) *cobra.Command {
	co := &run.CreateOptions{}

	cmd := &cobra.Command{
		// TODO: ValidArgs and different metadata loading setup
		Use:   "create [image] [kernel]",
		Short: "Create a new containerized VM without starting it",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			errutils.Check(func() error {
				var err error
				if co.Image, err = matchSingleImage(args[0]); err != nil {
					return err
				}
				if co.Kernel, err = matchSingleKernel(args[1]); err != nil {
					return err
				}
				return run.Create(co)
			}())
		},
	}

	addCreateFlags(cmd.Flags(), co)
	return cmd
}

func addCreateFlags(fs *pflag.FlagSet, co *run.CreateOptions) {
	addNameFlag(fs, &co.Name)
	fs.Int64Var(&co.CPUs, "cpus", constants.VM_DEFAULT_CPUS, "VM vCPU count, 1 or even numbers between 1 and 32")
	fs.Int64Var(&co.Memory, "memory", constants.VM_DEFAULT_MEMORY, "VM RAM in MiB")
}
