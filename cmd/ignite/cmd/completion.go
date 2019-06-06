package cmd

import (
	"io"
	"os"

	"github.com/luxas/ignite/pkg/errutils"
	"github.com/spf13/cobra"
)

// completionCmd represents the completion command
//var completionCmd = &cobra.Command{
//	Use:   "completion",
//	Short: "Generates bash completion scripts",
//	Long: `To load completion run
//
//. <(bitbucket completion)
//
//To configure your bash shell to load completions for each session add to your bashrc
//
//# ~/.bashrc or ~/.profile
//. <(bitbucket completion)
//`,
//	Run: func(cmd *cobra.Command, args []string) {
//		rootCmd.GenBashCompletion(os.Stdout);
//	},
//}

func NewCmdCompletion(out io.Writer, rootCmd *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion",
		Short: "Output bash completion for ignite to stdout",
		Run: func(cmd *cobra.Command, args []string) {
			err := rootCmd.GenBashCompletion(os.Stdout)
			errutils.Check(err)
		},
	}

	return cmd
}
