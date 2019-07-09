package vmcmd

import (
	"io"

	"github.com/lithammer/dedent"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/ignite/cmd/ignite/cmd/cmdutil"
	"github.com/weaveworks/ignite/cmd/ignite/run"
	"github.com/weaveworks/ignite/pkg/errutils"
)

// NewCmdCreate creates a new VM given an image and a kernel
func NewCmdCreate(out io.Writer) *cobra.Command {
	cf := run.NewCreateFlags()

	cmd := &cobra.Command{
		Use:   "create <OCI image>",
		Short: "Create a new VM without starting it",
		Long: dedent.Dedent(`
			Create a new VM by combining the given image and kernel.
			Various configuration options can be set during creation by using
			the flags for this command.
			
			If the name flag (-n, --name) is not specified,
			the VM is given a random name. Using the copy files
			flag (-f, --copy-files), additional files can be added to
			the VM during creation with the syntax /host/path:/vm/path.

			Example usage:
				$ ignite create centos:7 \
					--name my-vm \
					--cpus 2 \
					--ssh \
					--memory 2GB \
					--size 6GB
		`),
		Args: cobra.RangeArgs(0, 1),
		Run: func(cmd *cobra.Command, args []string) {
			errutils.Check(func() error {
				co, err := cf.NewCreateOptions(args)
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
	// Register common flags
	cmdutil.AddNameFlag(fs, &cf.VM.ObjectMeta.Name)
	cmdutil.AddConfigFlag(fs, &cf.ConfigFile)

	// Register flags bound to temporary holder values
	fs.StringSliceVarP(&cf.PortMappings, "ports", "p", cf.PortMappings, "Map host ports to VM ports")
	fs.StringSliceVarP(&cf.CopyFiles, "copy-files", "f", cf.CopyFiles, "Copy files from the host to the created VM")

	// Register flags for simple types (int, string, etc.)
	fs.Uint64Var(&cf.VM.Spec.CPUs, "cpus", cf.VM.Spec.CPUs, "VM vCPU count, 1 or even numbers between 1 and 32")
	fs.StringVar(&cf.VM.Spec.Kernel.CmdLine, "kernel-args", cf.VM.Spec.Kernel.CmdLine, "Set the command line for the kernel")

	// Register more complex flags with their own flag types
	cmdutil.SizeVar(fs, &cf.VM.Spec.Memory, "memory", "Amount of RAM to allocate for the VM")
	cmdutil.SizeVarP(fs, &cf.VM.Spec.DiskSize, "size", "s", "VM filesystem size, for example 5GB or 2048MB")
	cmdutil.OCIImageRefVarP(fs, &cf.VM.Spec.Kernel.OCIClaim.Ref, "kernel-image", "k", "Specify an OCI image containing the kernel at /boot/vmlinux and optionally, modules")
	cmdutil.NetworkModeVar(fs, &cf.VM.Spec.Network.Mode)

	cf.SSH = &run.SSHFlag{}
	fs.Var(cf.SSH, "ssh", "Enable SSH for the VM. If <path> is given, it will be imported as the public key. If just '--ssh' is specified, a new keypair will be generated.")

	sshFlag := fs.Lookup("ssh")
	sshFlag.NoOptDefVal = "<path>"
	sshFlag.DefValue = "is unset, which disables SSH access to the VM"
}
