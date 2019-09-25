package vmcmd

import (
	"io"

	"github.com/lithammer/dedent"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/ignite/cmd/ignite/cmd/cmdutil"
	"github.com/weaveworks/ignite/cmd/ignite/run"
)

func NewCmdPreflight(out io.Writer) *cobra.Command {
	preflightFlags := &run.PreflightFlags{}
	cmd := &cobra.Command{
		Use:   "preflight",
		Short: "Checks dependencies are fullfilled",
		Long: dedent.Dedent(`
			Run preflight checkers to verify all the required dependencies are present
		`),
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(func() error {
				po := preflightFlags.NewPreflightOptions()
				return run.Preflight(po)
			}())
		},
	}

	addPreflightFlags(cmd.Flags(), preflightFlags)
	return cmd
}

func addPreflightFlags(ps *pflag.FlagSet, pf *run.PreflightFlags) {
	ps.StringSliceVar(&pf.IgnoredPreflightErrors, "ignore-preflight", pf.IgnoredPreflightErrors, "ignore listed preflights")
}
