package vmcmd

import (
	"io"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/ignite/cmd/ignite/cmd/cmdutil"
	"github.com/weaveworks/ignite/cmd/ignite/run"
	"github.com/weaveworks/ignite/pkg/errutils"
)

// NewCmdPs lists running VMs
func NewCmdPs(out io.Writer) *cobra.Command {
	po := &run.PsOptions{}

	cmd := &cobra.Command{
		Use:     "ps",
		Short:   "List running VMs",
		Aliases: []string{"ls", "list"},
		Run: func(cmd *cobra.Command, args []string) {
			errutils.Check(func() error {
				var err error
				if po.VMs, err = cmdutil.MatchAllVMs(po.All); err != nil {
					return err
				}
				return run.Ps(po)
			}())
		},
	}

	addPsFlags(cmd.Flags(), po)
	return cmd
}

func addPsFlags(fs *pflag.FlagSet, po *run.PsOptions) {
	fs.BoolVarP(&po.All, "all", "a", false, "Show all VMs, not just running ones")
}
