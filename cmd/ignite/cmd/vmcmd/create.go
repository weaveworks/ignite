package vmcmd

import (
	"io"

	"github.com/lithammer/dedent"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/ignite/cmd/ignite/cmd/cmdutil"
	"github.com/weaveworks/ignite/cmd/ignite/run"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/errutils"
)

// NewCmdCreate creates a new VM given an image and a kernel
func NewCmdCreate(out io.Writer) *cobra.Command {
	co := &run.CreateOptions{}

	cmd := &cobra.Command{
		Use:   "create [image]",
		Short: "Create a new VM without starting it",
		Long: dedent.Dedent(`
			Create a new VM by combining the given image and kernel.
			Various VM tunables can be set during creation by using
			the flags for this command. The image and kernel are
			matched by prefix based on their ID and name.
			
			If the name flag (-n, --name) is not specified,
			the VM is given a random name. Using the copy files
			flag (-f, --copy-files), additional files can be added to
			the VM during creation with the syntax /host/path:/vm/path.

			Example usage:
				$ ignite create my-image my-kernel \
					--name my-vm \
					--cpus 2 \
					--memory 2048 \
					--size 6GB
		`),
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			errutils.Check(func() error {
				var err error
				if co.Image, err = cmdutil.MatchSingleImage(args[0]); err != nil {
					return err
				}
				if len(co.KernelName) == 0 {
					co.KernelName = args[0]
				}
				if co.Kernel, err = cmdutil.MatchSingleKernel(co.KernelName); err != nil {
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
	//cmd.Flag("test").
	return cmd
}

func addCreateFlags(fs *pflag.FlagSet, co *run.CreateOptions) {
	cmdutil.AddNameFlag(fs, &co.Name)
	fs.Int64Var(&co.CPUs, "cpus", constants.VM_DEFAULT_CPUS, "VM vCPU count, 1 or even numbers between 1 and 32")
	fs.Int64Var(&co.Memory, "memory", constants.VM_DEFAULT_MEMORY, "VM RAM in MiB")
	fs.StringVarP(&co.Size, "size", "s", constants.VM_DEFAULT_SIZE, "VM filesystem size, for example 5GB or 2048MB")
	fs.StringSliceVarP(&co.CopyFiles, "copy-files", "f", nil, "Copy files from the host to the created VM")
	fs.StringVarP(&co.KernelName, "kernel", "k", "", "Specify a kernel to use. By default this equals the image name")
	fs.StringVar(&co.KernelCmd, "kernel-args", constants.VM_KERNEL_ARGS, "Set the command line for the kernel")

	co.SSH = &run.SSHFlag{}
	fs.Var(co.SSH, "ssh", "Enable SSH for the VM. If <path> is given, it will be imported as the public key. If just '--ssh' is specified, a new keypair will be generated.")

	sshFlag := fs.Lookup("ssh")
	sshFlag.NoOptDefVal = "<path>"
	sshFlag.DefValue = "is unset, which disables SSH access to the VM"
}
