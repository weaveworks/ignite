package cmd

import (
	"io"
	"os"

	"github.com/lithammer/dedent"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/ignite/pkg/logs"
	logflag "github.com/weaveworks/ignite/pkg/logs/flag"
	networkflag "github.com/weaveworks/ignite/pkg/network/flag"
	"github.com/weaveworks/ignite/pkg/providers"
	"github.com/weaveworks/ignite/pkg/providers/ignite"
	runtimeflag "github.com/weaveworks/ignite/pkg/runtime/flag"
	versioncmd "github.com/weaveworks/ignite/pkg/version/cmd"
)

var logLevel = logrus.InfoLevel

// NewIgnitedCommand returns the root command for ignited
func NewIgnitedCommand(in io.Reader, out, err io.Writer) *cobra.Command {
	root := &cobra.Command{
		Use:   "ignited",
		Short: "ignited: run Firecracker VMs declaratively through a manifest directory or Git",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Set the desired logging level, now that the flags are parsed
			logs.Logger.SetLevel(logLevel)

			// Populate the providers after flags have been parsed
			if err := providers.Populate(ignite.Providers); err != nil {
				log.Fatal(err)
			}
		},
		Long: dedent.Dedent(`
			Ignite is a containerized Firecracker microVM administration tool.
			It can build VM images, spin VMs up/down and manage multiple VMs efficiently.

			TODO: ignited documentation
		`),
	}

	addGlobalFlags(root.PersistentFlags())

	root.AddCommand(NewCmdCompletion(os.Stdout, root))
	root.AddCommand(NewCmdGitOps(os.Stdout))
	root.AddCommand(NewCmdDaemon(os.Stdout))
	root.AddCommand(versioncmd.NewCmdVersion(os.Stdout))
	return root
}

func addGlobalFlags(fs *pflag.FlagSet) {
	logflag.LogLevelFlagVar(fs, &logLevel)
	runtimeflag.RuntimeVar(fs, &providers.RuntimeName)
	networkflag.NetworkPluginVar(fs, &providers.NetworkPluginName)
}
