package imgcmd

import (
	"io"

	"github.com/lithammer/dedent"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/ignite/cmd/ignite/cmd/cmdutil"
	"github.com/weaveworks/ignite/cmd/ignite/run"
)

// NewCmdRm removes images
func NewCmdRm(out io.Writer) *cobra.Command {
	rf := &run.RmiFlags{}

	cmd := &cobra.Command{
		Use:   "rm <image>...",
		Short: "Remove VM base images",
		Long: dedent.Dedent(`
			Remove one or multiple VM base images. Images are matched by prefix based on
			their ID and name. To remove multiple images, chain the matches separated by spaces.
			The force flag (-f, --force) kills and removes any running VMs using the image.
		`),
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(func() error {
				ro, err := rf.NewRmiOptions(args)
				if err != nil {
					return err
				}

				return run.Rmi(ro)
			}())
		},
	}

	addRmiFlags(cmd.Flags(), rf)
	return cmd
}

func addRmiFlags(fs *pflag.FlagSet, rf *run.RmiFlags) {
	cmdutil.AddForceFlag(fs, &rf.Force)
}
