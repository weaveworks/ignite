package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/lithammer/dedent"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/ignite/cmd/ignite/cmd/imgcmd"
	"github.com/weaveworks/ignite/cmd/ignite/cmd/kerncmd"
	"github.com/weaveworks/ignite/cmd/ignite/cmd/vmcmd"
	"github.com/weaveworks/ignite/pkg/config"
	"github.com/weaveworks/ignite/pkg/logs"
	logflag "github.com/weaveworks/ignite/pkg/logs/flag"
	"github.com/weaveworks/ignite/pkg/providers"
	"github.com/weaveworks/ignite/pkg/providers/ignite"
	"github.com/weaveworks/ignite/pkg/util"
	versioncmd "github.com/weaveworks/ignite/pkg/version/cmd"
)

var logLevel = log.InfoLevel

// Ignite config file path flag variable.
var configPath string

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

			// TODO Some commands do not need to check root
			// Currently it seems to be only ignite version that does not require root
			if isNonRootCommand(cmd.Name(), cmd.Parent().Name()) {
				return
			}

			// Ignite needs to run as root for now, see
			// https://github.com/weaveworks/ignite/issues/46
			// TODO: Remove this when ready
			util.GenericCheckErr(util.TestRoot())

			// Create the directories needed for running
			util.GenericCheckErr(util.CreateDirectories())

			if err := config.ApplyConfiguration(configPath); err != nil {
				log.Fatal(err)
			}

			// Populate the providers after flags have been parsed
			if err := providers.Populate(ignite.Providers); err != nil {
				log.Fatal(err)
			}
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

				$ ignite run weaveworks/ignite-ubuntu \
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
	root.AddCommand(NewCmdCP(os.Stdout))
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

func isNonRootCommand(cmd string, parentCmd string) bool {
	if parentCmd != "ignite" {
		return false
	}

	switch cmd {
	case "version", "help", "image", "kernel", "completion", "inspect", "ps":
		return true
	}

	return false
}

func addGlobalFlags(fs *pflag.FlagSet) {
	AddQuietFlag(fs)
	logflag.LogLevelFlagVar(fs, &logLevel)
	fs.StringVar(&configPath, "ignite-config", "", "Ignite configuration path; refer to the 'Ignite Configuration' docs for more details")
}

// AddQuietFlag adds the quiet flag to a flagset
func AddQuietFlag(fs *pflag.FlagSet) {
	fs.BoolVarP(&logs.Quiet, "quiet", "q", logs.Quiet, "The quiet mode allows for machine-parsable output by printing only IDs")
}
