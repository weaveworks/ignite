package vmcmd

import (
	"io"

	"github.com/lithammer/dedent"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/ignite/cmd/ignite/run"
	"github.com/weaveworks/ignite/pkg/errutils"
)

// NewCmdPs lists running VMs
func NewCmdPs(out io.Writer) *cobra.Command {
	pf := &run.PsFlags{}

	cmd := &cobra.Command{
		Use:     "ps",
		Short:   "List running VMs",
		Aliases: []string{"ls", "list"},
		Long: dedent.Dedent(`
			List all running VMs. By specifying the all flag (-a, --all),
			also list VMs that are not currently running.
		`),
		Run: func(cmd *cobra.Command, args []string) {
			errutils.Check(func() error {
				po, err := pf.NewPsOptions()
				if err != nil {
					return err
				}

				return run.Ps(po)
			}())
		},
	}

	addPsFlags(cmd.Flags(), pf)
	return cmd
}

func addPsFlags(fs *pflag.FlagSet, pf *run.PsFlags) {
	fs.BoolVarP(&pf.All, "all", "a", false, "Show all VMs, not just running ones")
}
