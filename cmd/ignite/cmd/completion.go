package cmd

import (
	"io"
	"os"

	"github.com/lithammer/dedent"
	"github.com/spf13/cobra"
	"github.com/weaveworks/ignite/pkg/errutils"
)

func NewCmdCompletion(out io.Writer, rootCmd *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion",
		Short: "Output bash completion for ignite to stdout",
		Long: dedent.Dedent(`
		In order to start using the auto-completion, run:

			. <(ignite completion)
		
		To configure your bash shell to load completions for each session, run:

			echo '. <(ignite completion)' >> ~/.bashrc
		`),
		Run: func(cmd *cobra.Command, args []string) {
			err := rootCmd.GenBashCompletion(os.Stdout)
			errutils.Check(err)
		},
	}
	return cmd
}
