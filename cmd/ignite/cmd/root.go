package cmd

import (
	"github.com/luxas/ignite/pkg/util"
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
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// TODO: Handle this error more softly?
			if err := util.CreateDirectories(); err != nil {
				panic(err)
			}
		},
		Long: dedent.Dedent(`
			Ignite helps you with foo, bar

			Example usage:

			    $ ignite foo
		`),
	}

	cmds.AddCommand(NewCmdImage(os.Stdout))
	cmds.AddCommand(NewCmdKernel(os.Stdout))
	cmds.AddCommand(NewCmdVM(os.Stdout))

	cmds.AddCommand(NewCmdAddImage(os.Stdout))
	cmds.AddCommand(NewCmdAddKernel(os.Stdout))
	cmds.AddCommand(NewCmdAttach(os.Stdout))
	cmds.AddCommand(NewCmdBuild(os.Stdout))
	cmds.AddCommand(NewCmdCompletion(os.Stdout, cmds))
	cmds.AddCommand(NewCmdContainer(os.Stdout))
	cmds.AddCommand(NewCmdCreate(os.Stdout))
	cmds.AddCommand(NewCmdKill(os.Stdout))
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
