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
	"github.com/weaveworks/ignite/pkg/metadata/loader"
)

// NewCmdCreate creates a new VM given an image and a kernel
func NewCmdCreate(out io.Writer) *cobra.Command {
	cf := run.NewCreateFlags()

	cmd := &cobra.Command{
		Use:   "create <image>",
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
		Args: cobra.RangeArgs(0, 1),
		Run: func(cmd *cobra.Command, args []string) {
			errutils.Check(func() error {
				co, err := cf.NewCreateOptions(loader.NewResLoader(), args)
				if err != nil {
					return err
				}

				return run.Create(co)
			}())
		},
	}

	addCreateFlags(cmd.Flags(), cf)
	return cmd
}

func addCreateFlags(fs *pflag.FlagSet, cf *run.CreateFlags) {
	cmdutil.AddNameFlag(fs, &cf.VM.ObjectMeta.Name)
	cmdutil.AddConfigFlag(fs, &cf.ConfigFile)
	fs.Uint64Var(&cf.VM.Spec.CPUs, "cpus", cf.VM.Spec.CPUs, "VM vCPU count, 1 or even numbers between 1 and 32")
	cmdutil.SizeVar(fs, &cf.VM.Spec.Memory, "memory", cf.VM.Spec.Memory, "Amount of RAM to allocate for the VM")
	cmdutil.SizeVarP(fs, &cf.VM.Spec.DiskSize, "size", "s", cf.VM.Spec.DiskSize, "VM filesystem size, for example 5GB or 2048MB")
	fs.StringSliceVarP(&cf.CopyFiles, "copy-files", "f", nil, "Copy files from the host to the created VM")
	fs.StringVarP(&cf.KernelName, "kernel", "k", "", "Specify a kernel to use. By default this equals the image name")
	fs.StringVar(&cf.VM.Spec.Kernel.CmdLine, "kernel-args", constants.VM_DEFAULT_KERNEL_ARGS, "Set the command line for the kernel")

	cf.SSH = &run.SSHFlag{}
	fs.Var(cf.SSH, "ssh", "Enable SSH for the VM. If <path> is given, it will be imported as the public key. If just '--ssh' is specified, a new keypair will be generated.")

	sshFlag := fs.Lookup("ssh")
	sshFlag.NoOptDefVal = "<path>"
	sshFlag.DefValue = "is unset, which disables SSH access to the VM"
}
