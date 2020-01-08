package vmcmd

import (
	"fmt"
	"io"

	"github.com/lithammer/dedent"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/ignite/cmd/ignite/cmd/cmdutil"
	"github.com/weaveworks/ignite/cmd/ignite/run"
	"github.com/weaveworks/ignite/pkg/constants"
)

// NewCmdCreate creates a new VM given an image and a kernel
func NewCmdCreate(out io.Writer) *cobra.Command {
	cf := run.NewCreateFlags()

	cmd := &cobra.Command{
		Use:   "create <OCI image>",
		Short: "Create a new VM without starting it",
		Long: dedent.Dedent(fmt.Sprintf(`
			Create a new VM by combining the given image with a kernel. If no
			kernel is given using the kernel flag (-k, --kernel-image), use the
			default kernel (%s).

			Various configuration options can be set during creation by using
			the flags for this command.
			
			If the name flag (-n, --name) is not specified,
			the VM is given a random name. Using the copy files
			flag (-f, --copy-files), additional files/directories
			can be added to the VM during creation with the syntax
			/host/path:/vm/path.

			Example usage:
				$ ignite create weaveworks/ignite-ubuntu \
					--name my-vm \
					--cpus 2 \
					--ssh \
					--dns 8.8.8.8 \
					--memory 2GB \
					--size 6GB
		`, constants.DEFAULT_KERNEL_IMAGE)),
		Args: cobra.RangeArgs(0, 1),
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(func() error {
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
	fs.StringSliceVarP(&cf.CopyFiles, "copy-files", "f", cf.CopyFiles, "Copy files/directories from the host to the created VM")
	fs.StringSliceVarP(&cf.DNS, "dns", "", cf.DNS, "Override the default name servers in VM /etc/resolv.conf")

	// Register flags for simple types (int, string, etc.)
	fs.Uint64Var(&cf.VM.Spec.CPUs, "cpus", cf.VM.Spec.CPUs, "VM vCPU count, 1 or even numbers between 1 and 32")
	fs.StringVar(&cf.VM.Spec.Kernel.CmdLine, "kernel-args", cf.VM.Spec.Kernel.CmdLine, "Set the command line for the kernel")

	// Register more complex flags with their own flag types
	cmdutil.SizeVar(fs, &cf.VM.Spec.Memory, "memory", "Amount of RAM to allocate for the VM")
	cmdutil.SizeVarP(fs, &cf.VM.Spec.DiskSize, "size", "s", "VM filesystem size, for example 5GB or 2048MB")
	cmdutil.OCIImageRefVarP(fs, &cf.VM.Spec.Kernel.OCI, "kernel-image", "k", "Specify an OCI image containing the kernel at /boot/vmlinux and optionally, modules")
	cmdutil.SSHVar(fs, &cf.SSH)
	cmdutil.VolumeVarP(fs, &cf.VM.Spec.Storage, "volumes", "v", "Expose block devices from the host inside the VM")
}
