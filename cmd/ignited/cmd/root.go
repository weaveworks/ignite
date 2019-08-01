package cmd

import (
	"io"
	"os"

	"github.com/lithammer/dedent"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/ignite/pkg/logs"
	logflag "github.com/weaveworks/ignite/pkg/logs/flag"
	versioncmd "github.com/weaveworks/ignite/pkg/version/cmd"
)

var logLevel = logrus.InfoLevel

// NewIgnitedCommand returns the root command for ignited
func NewIgnitedCommand(in io.Reader, out, err io.Writer) *cobra.Command {
	// Set the default logging level
	logs.Logger.SetLevel(logLevel)

	root := &cobra.Command{
		Use:   "ignited",
		Short: "ignited: run Firecracker VMs declaratively through a manifest directory or Git",
		Long: dedent.Dedent(`
			Ignite is a containerized Firecracker microVM administration tool.
			It can build VM images, spin VMs up/down and manage multiple VMs efficiently.

			TODO: ignited documentation
		`),
	}

	addGlobalFlags(root.PersistentFlags())

	root.AddCommand(NewCmdGitOps(os.Stdout))
	root.AddCommand(NewCmdDaemon(os.Stdout))
	root.AddCommand(versioncmd.NewCmdVersion(os.Stdout))
	return root
}

func addGlobalFlags(fs *pflag.FlagSet) {
	AddQuietFlag(fs)
	logflag.LogLevelFlagVar(fs, &logLevel)
}

// AddQuietFlag adds the quiet flag to a flagset
func AddQuietFlag(fs *pflag.FlagSet) {
	fs.BoolVarP(&logs.Quiet, "quiet", "q", logs.Quiet, "The quiet mode allows for machine-parsable output, by printing only IDs")
}
