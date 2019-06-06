package cmd

import (
	"io"

	"github.com/weaveworks/ignite/cmd/ignite/run"

	"github.com/spf13/cobra"
	"github.com/weaveworks/ignite/pkg/errutils"
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
