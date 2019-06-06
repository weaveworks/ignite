package imgcmd

import (
	"io"

	"github.com/lithammer/dedent"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/ignite/cmd/ignite/cmd/cmdutil"
	"github.com/weaveworks/ignite/cmd/ignite/run"
	"github.com/weaveworks/ignite/pkg/errutils"
)

// NewCmdRm removes an images
func NewCmdRm(out io.Writer) *cobra.Command {
	ro := &run.RmiOptions{}

	cmd := &cobra.Command{
		Use:   "rm [image]...",
		Short: "Remove VM base images",
		Long: dedent.Dedent(`
			Remove one or multiple VM base images. Images are matched by prefix based
			on their name and ID. To remove multiple images, chain the matches separated
			by spaces. The "--force" flag kills and removes any running VMs using the image.
		`),
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			errutils.Check(func() error {
				var err error
				if ro.Images, err = cmdutil.MatchSingleImages(args); err != nil {
					return err
				}
				if ro.VMs, err = cmdutil.MatchAllVMs(true); err != nil {
					return err
				}
				return run.Rmi(ro)
			}())
		},
	}

	addRmiFlags(cmd.Flags(), ro)
	return cmd
}

func addRmiFlags(fs *pflag.FlagSet, ro *run.RmiOptions) {
	cmdutil.AddForceFlag(fs, &ro.Force)
}
