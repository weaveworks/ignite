package vmcmd

import (
	"fmt"
	"io"

	"github.com/lithammer/dedent"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/ignite/cmd/ignite/cmd/cmdutil"
	"github.com/weaveworks/ignite/cmd/ignite/run"
	"github.com/weaveworks/ignite/pkg/constants"
)

// NewCmdStop stops VMs
func NewCmdStop(out io.Writer) *cobra.Command {
	sf := &run.StopFlags{}

	cmd := &cobra.Command{
		Use:   "stop <vm>...",
		Short: "Stop running VMs",
		Long: dedent.Dedent(fmt.Sprintf(`
			Stop one or multiple VMs. The VMs are matched by prefix based on their
			ID and name. To stop multiple VMs, chain the matches separated by spaces.
			The force flag (-f, --force) kills VMs instead of trying to stop them
			gracefully.

			The VMs are given a %d second grace period to shut down before they
			will be forcibly killed.
		`, constants.STOP_TIMEOUT)),
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(func() error {
				so, err := sf.NewStopOptions(args)
				if err != nil {
					return err
				}

				return run.Stop(so)
			}())
		},
	}

	addStopFlags(cmd.Flags(), sf)
	return cmd
}

func addStopFlags(fs *pflag.FlagSet, sf *run.StopFlags) {
	fs.BoolVarP(&sf.Kill, "force-kill", "f", false, "Force kill the VM")
}
