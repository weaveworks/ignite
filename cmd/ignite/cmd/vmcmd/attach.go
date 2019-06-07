package vmcmd

import (
	"io"

	"github.com/lithammer/dedent"

	"github.com/spf13/cobra"
	"github.com/weaveworks/ignite/cmd/ignite/cmd/cmdutil"
	"github.com/weaveworks/ignite/cmd/ignite/run"
	"github.com/weaveworks/ignite/pkg/errutils"
)

// NewCmdAttach attaches to a running VM
func NewCmdAttach(out io.Writer) *cobra.Command {
	// checkRunning can be used to skip the running check, this is used by Start and Run
	// as the in-container ignite takes some time to start up and update the state
	ao := &run.AttachOptions{CheckRunning: true}

	cmd := &cobra.Command{
		Use:   "attach <vm>",
		Short: "Attach to a running VM",
		Long: dedent.Dedent(`
			Connect the current terminal to the running VM's TTY.
			To detach from the VM's TTY, type ^P^Q (Ctrl + P + Q).
			The given VM is matched by prefix based on its ID and name.
		`),
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			errutils.Check(func() error {
				var err error
				if ao.VM, err = cmdutil.MatchSingleVM(args[0]); err != nil {
					return err
				}
				return run.Attach(ao)
			}())
		},
	}

	return cmd
}
