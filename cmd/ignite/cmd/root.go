package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/weaveworks/ignite/cmd/ignite/cmd/imgcmd"
	"github.com/weaveworks/ignite/cmd/ignite/cmd/kerncmd"
	"github.com/weaveworks/ignite/cmd/ignite/cmd/vmcmd"
	"github.com/weaveworks/ignite/pkg/util"

	"github.com/lithammer/dedent"
	"github.com/spf13/cobra"
)

// NewIgniteCommand returns the root command for ignite
func NewIgniteCommand(in io.Reader, out, err io.Writer) *cobra.Command {
	imageCmd := imgcmd.NewCmdImage(os.Stdout)
	kernelCmd := kerncmd.NewCmdKernel(os.Stdout)
	vmCmd := vmcmd.NewCmdVM(os.Stdout)

	root := &cobra.Command{
		Use:   "ignite",
		Short: "ignite: easily run Firecracker VMs",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// TODO: Handle this error more softly?
			if err := util.CreateDirectories(); err != nil {
				panic(err)
			}
		},
		Long: dedent.Dedent(fmt.Sprintf(`
			Ignite is a containerized Firecracker microVM administration tool.
			It can build VM images, spin VMs up/down and manage multiple VMs efficiently.

			Administration is divided into three subcommands:
			  image       %s
			  kernel      %s
			  vm          %s

			Image + Kernel = VM

			Example usage:
				$ ignite build luxas/ubuntu-base:18.04 \
					--name my-image \
					--import-kernel my-kernel
				$ ignite images
				$ ignite kernels
				$ ignite run my-image my-kernel --name my-vm
				$ ignite ps
				$ ignite attach my-vm
				Login with user "root" and password "root".
		`, imageCmd.Short, kernelCmd.Short, vmCmd.Short)),
	}

	root.AddCommand(imageCmd)
	root.AddCommand(kernelCmd)
	root.AddCommand(vmCmd)

	root.AddCommand(NewCmdAttach(os.Stdout))
	root.AddCommand(NewCmdBuild(os.Stdout))
	root.AddCommand(NewCmdCompletion(os.Stdout, root))
	root.AddCommand(NewCmdContainer(os.Stdout))
	root.AddCommand(NewCmdCreate(os.Stdout))
	root.AddCommand(NewCmdKill(os.Stdout))
	root.AddCommand(NewCmdLogs(os.Stdout))
	root.AddCommand(NewCmdPs(os.Stdout))
	root.AddCommand(NewCmdRm(os.Stdout))
	root.AddCommand(NewCmdRmi(os.Stdout))
	root.AddCommand(NewCmdRmk(os.Stdout))
	root.AddCommand(NewCmdRun(os.Stdout))
	root.AddCommand(NewCmdStart(os.Stdout))
	root.AddCommand(NewCmdStop(os.Stdout))
	root.AddCommand(NewCmdVersion(os.Stdout))
	return root
}
