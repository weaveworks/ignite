package cmd

import (
	"io"
	"os"

	"github.com/lithammer/dedent"
	"github.com/spf13/cobra"
	"github.com/weaveworks/ignite/cmd/ignite/cmd/cmdutil"
)

func NewCmdCompletion(out io.Writer, rootCmd *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion",
		Short: "Output bash completion for ignited to stdout",
		Long: dedent.Dedent(`
		In order to start using the auto-completion, run:

			. <(ignited completion)
		
		To configure your bash shell to load completions for each session, run:

			echo '. <(ignited completion)' >> ~/.bashrc
		`),
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(rootCmd.GenBashCompletion(os.Stdout))
		},
	}
	return cmd
}
