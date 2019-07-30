package main

import (
	"io"

	"github.com/lithammer/dedent"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/ignite/cmd/ignite/cmd/cmdutil"
	"github.com/weaveworks/ignite/pkg/logs"
)

var logLevel = logrus.InfoLevel

// NewIgniteSpawnCommand returns the root command for ignite-spawn
func NewIgniteSpawnCommand(in io.Reader, out, err io.Writer) *cobra.Command {
	root := &cobra.Command{
		Use: "ignite-spawn <vm>",
		Short: dedent.Dedent(`
			Start the given VM in Firecracker and manage it.
			Used internally by Ignite, don't call ignite-spawn
			by hand. Refer to Ignite for CLI usage.
		`),
		Args: cobra.ExactArgs(1),
		PersistentPreRun: func(_ *cobra.Command, _ []string) {
			logs.InitLogs(logLevel)
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(func() error {
				opts, err := NewOptions(args[0])
				if err != nil {
					return err
				}

				return StartVM(opts)
			}())
		},
	}

	addGlobalFlags(root.PersistentFlags())

	return root
}

func addGlobalFlags(fs *pflag.FlagSet) {
	cmdutil.LogLevelFlagVar(fs, &logLevel)
}
