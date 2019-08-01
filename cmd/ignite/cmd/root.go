package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/lithammer/dedent"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/ignite/cmd/ignite/cmd/imgcmd"
	"github.com/weaveworks/ignite/cmd/ignite/cmd/kerncmd"
	"github.com/weaveworks/ignite/cmd/ignite/cmd/vmcmd"
	"github.com/weaveworks/ignite/pkg/logs"
	logflag "github.com/weaveworks/ignite/pkg/logs/flag"
	versioncmd "github.com/weaveworks/ignite/pkg/version/cmd"
)

var logLevel = logrus.InfoLevel

// NewIgniteCommand returns the root command for ignite
func NewIgniteCommand(in io.Reader, out, err io.Writer) *cobra.Command {

	imageCmd := imgcmd.NewCmdImage(os.Stdout)
	kernelCmd := kerncmd.NewCmdKernel(os.Stdout)
	vmCmd := vmcmd.NewCmdVM(os.Stdout)

	root := &cobra.Command{
		Use:   "ignite",
		Short: "ignite: easily run Firecracker VMs",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Set the desired logging level, now that the flags are parsed
			logs.Logger.SetLevel(logLevel)
		},
		Long: dedent.Dedent(fmt.Sprintf(`
			Ignite is a containerized Firecracker microVM administration tool.
			It can build VM images, spin VMs up/down and manage multiple VMs efficiently.

			Administration is divided into three subcommands:
			  image       %s
			  kernel      %s
			  vm          %s

			Ignite also supports the same commands as the Docker CLI.
			Combining an Image and a Kernel gives you a runnable VM.

			Example usage:

				$ ignite run centos:7 \
					--cpus 2 \
					--memory 2GB \
					--ssh \
					--name my-vm
				$ ignite images
				$ ignite kernels
				$ ignite ps
				$ ignite logs my-vm
				$ ignite ssh my-vm
		`, imageCmd.Short, kernelCmd.Short, vmCmd.Short)),
	}

	addGlobalFlags(root.PersistentFlags())

	root.AddCommand(imageCmd)
	root.AddCommand(kernelCmd)
	root.AddCommand(vmCmd)

	root.AddCommand(NewCmdAttach(os.Stdout))
	root.AddCommand(NewCmdCompletion(os.Stdout, root))
	root.AddCommand(NewCmdCreate(os.Stdout))
	root.AddCommand(NewCmdKill(os.Stdout))
	root.AddCommand(NewCmdLogs(os.Stdout))
	root.AddCommand(NewCmdInspect(os.Stdout))
	root.AddCommand(NewCmdPs(os.Stdout))
	root.AddCommand(NewCmdRm(os.Stdout))
	root.AddCommand(NewCmdRmi(os.Stdout))
	root.AddCommand(NewCmdRmk(os.Stdout))
	root.AddCommand(NewCmdRun(os.Stdout))
	root.AddCommand(NewCmdSSH(os.Stdout))
	root.AddCommand(NewCmdExec(os.Stdout, os.Stderr, os.Stdin))
	root.AddCommand(NewCmdStart(os.Stdout))
	root.AddCommand(NewCmdStop(os.Stdout))
	root.AddCommand(versioncmd.NewCmdVersion(os.Stdout))
	return root
}

func addGlobalFlags(fs *pflag.FlagSet) {
	AddQuietFlag(fs)
	logflag.LogLevelFlagVar(fs, &logLevel)
}

// AddQuietFlag adds the quiet flag to a flagset
func AddQuietFlag(fs *pflag.FlagSet) {
	fs.BoolVarP(&logs.Quiet, "quiet", "q", logs.Quiet, "The quiet mode allows for machine-parsable output by printing only IDs")
}
