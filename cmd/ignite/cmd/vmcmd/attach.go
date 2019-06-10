package vmcmd

import (
	"github.com/weaveworks/ignite/cmd/ignite/run/runutil"
	"io"

	"github.com/lithammer/dedent"

	"github.com/spf13/cobra"
	"github.com/weaveworks/ignite/cmd/ignite/run"
	"github.com/weaveworks/ignite/pkg/errutils"
)

// NewCmdAttach attaches to a running VM
func NewCmdAttach(out io.Writer) *cobra.Command {
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
				ao, err := run.NewAttachOptions(runutil.NewResLoader(), args[0])
				if err != nil {
					return err
				}

				return run.Attach(ao)
			}())
		},
	}

	return cmd
}
