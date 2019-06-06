package vmcmd

import (
	"io"

	"github.com/weaveworks/ignite/cmd/ignite/cmd/cmdutil"
	"github.com/weaveworks/ignite/cmd/ignite/run"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/errutils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// NewCmdCreate creates a new VM given an image and a kernel
func NewCmdCreate(out io.Writer) *cobra.Command {
	co := &run.CreateOptions{}

	cmd := &cobra.Command{
		// TODO: ValidArgs and different metadata loading setup?
		Use:   "create [image] [kernel]",
		Short: "Create a new VM without starting it",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			errutils.Check(func() error {
				var err error
				if co.Image, err = cmdutil.MatchSingleImage(args[0]); err != nil {
					return err
				}
				if co.Kernel, err = cmdutil.MatchSingleKernel(args[1]); err != nil {
					return err
				}
				if co.VMNames, err = cmdutil.MatchAllVMNames(); err != nil {
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
	cmdutil.AddNameFlag(fs, &co.Name)
	fs.Int64Var(&co.CPUs, "cpus", constants.VM_DEFAULT_CPUS, "VM vCPU count, 1 or even numbers between 1 and 32")
	fs.Int64Var(&co.Memory, "memory", constants.VM_DEFAULT_MEMORY, "VM RAM in MiB")
	fs.StringVarP(&co.Size, "size", "s", constants.VM_DEFAULT_SIZE, "VM filesystem size, for example 5GB or 2048MB")
	fs.StringSliceVarP(&co.CopyFiles, "copy-files", "f", nil, "Copy files from the host to the created VM")
	fs.StringVar(&co.KernelCmd, "kernel-args", constants.VM_KERNEL_ARGS, "Set the command line for the kernel")
}
