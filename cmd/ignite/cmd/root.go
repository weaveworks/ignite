package cmd

import (
	"github.com/luxas/ignite/cmd/ignite/cmd/imgcmd"
	"github.com/luxas/ignite/cmd/ignite/cmd/kerncmd"
	"github.com/luxas/ignite/cmd/ignite/cmd/vmcmd"
	"github.com/luxas/ignite/pkg/util"
	"io"
	"os"

	"github.com/lithammer/dedent"
	"github.com/spf13/cobra"
)

// NewIgniteCommand returns the root command for ignite
func NewIgniteCommand(in io.Reader, out, err io.Writer) *cobra.Command {
	root := &cobra.Command{
		Use:   "ignite",
		Short: "ignite: easily run Firecracker VMs",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// TODO: Handle this error more softly?
			if err := util.CreateDirectories(); err != nil {
				panic(err)
			}
		},
		Long: dedent.Dedent(`
			Ignite is a containerized Firecracker microVM administration tool.
			It can build VM images, spin VMs up/down and manage multiple VMs efficiently.

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
		`),
	}

	root.AddCommand(imgcmd.NewCmdImage(os.Stdout))
	root.AddCommand(kerncmd.NewCmdKernel(os.Stdout))
	root.AddCommand(vmcmd.NewCmdVM(os.Stdout))

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
