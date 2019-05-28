package cmd

import (
	"io"
	"os"

	"github.com/lithammer/dedent"
	"github.com/spf13/cobra"
)

// NewIgniteCommand returns cobra.Command to run kubeadm command
func NewIgniteCommand(in io.Reader, out, err io.Writer) *cobra.Command {
	cmds := &cobra.Command{
		Use:   "ignite",
		Short: "ignite: easily run Firecracker VMs",
		Long: dedent.Dedent(`
			Ignite helps you with foo, bar

			Example usage:

			    $ ignite foo
		`),
	}

	cmds.AddCommand(NewCmdAddImage(os.Stdout))
	cmds.AddCommand(NewCmdAddKernel(os.Stdout))
	cmds.AddCommand(NewCmdAttach(os.Stdout))
	cmds.AddCommand(NewCmdBuild(os.Stdout))
	cmds.AddCommand(NewCmdCompletion(os.Stdout, cmds))
	cmds.AddCommand(NewCmdContainer(os.Stdout))
	cmds.AddCommand(NewCmdCreate(os.Stdout))
	cmds.AddCommand(NewCmdImages(os.Stdout))
	cmds.AddCommand(NewCmdKernels(os.Stdout))
	cmds.AddCommand(NewCmdLogs(os.Stdout))
	cmds.AddCommand(NewCmdPs(os.Stdout))
	cmds.AddCommand(NewCmdRm(os.Stdout))
	cmds.AddCommand(NewCmdRmi(os.Stdout))
	cmds.AddCommand(NewCmdRmk(os.Stdout))
	cmds.AddCommand(NewCmdRun(os.Stdout))
	cmds.AddCommand(NewCmdStart(os.Stdout))
	cmds.AddCommand(NewCmdStop(os.Stdout))
	cmds.AddCommand(NewCmdVersion(os.Stdout))
	return cmds
}
