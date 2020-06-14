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
	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	"github.com/weaveworks/ignite/pkg/apis/ignite/scheme"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/logs"
	logflag "github.com/weaveworks/ignite/pkg/logs/flag"
	networkflag "github.com/weaveworks/ignite/pkg/network/flag"
	"github.com/weaveworks/ignite/pkg/providers"
	"github.com/weaveworks/ignite/pkg/providers/ignite"
	runtimeflag "github.com/weaveworks/ignite/pkg/runtime/flag"
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
			if cmd.Name() == "version" && cmd.Parent().Name() == "ignite" {
				return
			}

			// Ignite needs to run as root for now, see
			// https://github.com/weaveworks/ignite/issues/46
			// TODO: Remove this when ready
			util.GenericCheckErr(util.TestRoot())

			// Create the directories needed for running
			util.GenericCheckErr(util.CreateDirectories())

			var configFilePath string

			// If an ignite config flag is set, use it as the config file, else check
			// if the global config file exists.
			// If a config file path is passed, configure ignite using it.
			if configPath != "" {
				configFilePath = configPath
			} else {
				// Check the default config location.
				if _, err := os.Stat(constants.IGNITE_CONFIG_FILE); !os.IsNotExist(err) {
					log.Debugf("Found default ignite configuration file %s", constants.IGNITE_CONFIG_FILE)
					configFilePath = constants.IGNITE_CONFIG_FILE
				}
			}

			if configFilePath != "" {
				log.Debugf("Using ignite configuration file %s", configFilePath)
				var err error
				providers.ComponentConfig, err = getConfigFromFile(configFilePath)
				if err != nil {
					log.Fatal(err)
				}

				// Set providers runtime and network plugin if found in config.
				if providers.ComponentConfig.Spec.Runtime != "" {
					providers.RuntimeName = providers.ComponentConfig.Spec.Runtime
				}
				if providers.ComponentConfig.Spec.NetworkPlugin != "" {
					providers.NetworkPluginName = providers.ComponentConfig.Spec.NetworkPlugin
				}
			} else {
				log.Debugln("Using ignite default configurations")
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

func addGlobalFlags(fs *pflag.FlagSet) {
	AddQuietFlag(fs)
	logflag.LogLevelFlagVar(fs, &logLevel)
	runtimeflag.RuntimeVar(fs, &providers.RuntimeName)
	networkflag.NetworkPluginVar(fs, &providers.NetworkPluginName)
	fs.StringVar(&configPath, "ignite-config", "", "Ignite configuration path")
}

// AddQuietFlag adds the quiet flag to a flagset
func AddQuietFlag(fs *pflag.FlagSet) {
	fs.BoolVarP(&logs.Quiet, "quiet", "q", logs.Quiet, "The quiet mode allows for machine-parsable output by printing only IDs")
}

// getConfigFromFile reads a config file and returns ignite configuration.
func getConfigFromFile(configPath string) (*api.Configuration, error) {
	componentConfig := &api.Configuration{}

	// TODO: Fix libgitops DecodeFileInto to not allow empty files.
	if err := scheme.Serializer.DecodeFileInto(configPath, componentConfig); err != nil {
		return nil, err
	}

	// Ensure the read configuration is valid. If a file contains Kind and
	// APIVersion, it's a valid config file. Empty file is invalid.
	// NOTE: This is a workaround for libgitops allowing decode of empty file.
	if componentConfig.Kind == "" || componentConfig.APIVersion == "" {
		return nil, fmt.Errorf("invalid config file, Kind and APIVersion must be set")
	}

	return componentConfig, nil
}
