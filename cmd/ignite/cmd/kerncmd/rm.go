package kerncmd

import (
	"github.com/lithammer/dedent"
	"io"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/ignite/cmd/ignite/cmd/cmdutil"
	"github.com/weaveworks/ignite/cmd/ignite/run"
	"github.com/weaveworks/ignite/pkg/errutils"
)

// NewCmdRm removes a kernel
func NewCmdRm(out io.Writer) *cobra.Command {
	ro := &run.RmkOptions{}

	cmd := &cobra.Command{
		Use:   "rm [kernel]...",
		Short: "Remove kernels",
		Long: dedent.Dedent(`
			Remove one or multiple VM kernels. Kernels are matched by prefix based
			on their ID and name. To remove multiple kernels, chain the matches separated
			by spaces. The "--force" flag kills and removes any running VMs using the kernel.
		`),
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			errutils.Check(func() error {
				var err error
				if ro.Kernels, err = cmdutil.MatchSingleKernels(args); err != nil {
					return err
				}
				if ro.VMs, err = cmdutil.MatchAllVMs(true); err != nil {
					return err
				}
				return run.Rmk(ro)
			}())
		},
	}

	addRmkFlags(cmd.Flags(), ro)
	return cmd
}

func addRmkFlags(fs *pflag.FlagSet, ro *run.RmkOptions) {
	cmdutil.AddForceFlag(fs, &ro.Force)
}
