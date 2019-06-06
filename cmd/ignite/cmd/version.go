package cmd

import (
	"io"

	"github.com/luxas/ignite/cmd/ignite/run"

	"github.com/luxas/ignite/pkg/errutils"
	"github.com/spf13/cobra"
)

// NewCmdVersion provides the version information of ignite
func NewCmdVersion(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version of ignite",
		Run: func(cmd *cobra.Command, args []string) {
			errutils.Check(run.Version(out, cmd))
		},
	}

	cmd.Flags().StringP("output", "o", "", "Output format; available options are 'yaml', 'json' and 'short'")
	return cmd
}
