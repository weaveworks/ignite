package kerncmd

import (
	"io"

	"github.com/lithammer/dedent"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/ignite/cmd/ignite/cmd/cmdutil"
	"github.com/weaveworks/ignite/cmd/ignite/run"
)

// NewCmdRm removes kernels
func NewCmdRm(out io.Writer) *cobra.Command {
	rf := &run.RmkFlags{}

	cmd := &cobra.Command{
		Use:   "rm <kernel>...",
		Short: "Remove kernels",
		Long: dedent.Dedent(`
			Remove one or multiple VM kernels. Kernels are matched by prefix based on their
			ID and name. To remove multiple kernels, chain the matches separated by spaces.
			The force flag (-f, --force) kills and removes any running VMs using the kernel.
		`),
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(func() error {
				ro, err := rf.NewRmkOptions(args)
				if err != nil {
					return err
				}

				return run.Rmk(ro)
			}())
		},
	}

	addRmkFlags(cmd.Flags(), rf)
	return cmd
}

func addRmkFlags(fs *pflag.FlagSet, rf *run.RmkFlags) {
	cmdutil.AddForceFlag(fs, &rf.Force)
}
