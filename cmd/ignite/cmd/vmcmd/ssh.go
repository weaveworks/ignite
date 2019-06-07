package vmcmd

import (
	"github.com/spf13/pflag"
	"io"

	"github.com/lithammer/dedent"

	"github.com/spf13/cobra"
	"github.com/weaveworks/ignite/cmd/ignite/cmd/cmdutil"
	"github.com/weaveworks/ignite/cmd/ignite/run"
	"github.com/weaveworks/ignite/pkg/errutils"
)

// NewCmdSSH ssh's into a running VM
func NewCmdSSH(out io.Writer) *cobra.Command {
	so := &run.SSHOptions{}

	cmd := &cobra.Command{
		Use:   "ssh [vm]",
		Short: "SSH into a running VM",
		Long: dedent.Dedent(`
			SSH into the running VM using the private key created for it during generation.
			If no private key was created or wanting to use a different identity file,
			use the identity file flag (-i, --identity) to override the used identity file.
			The given VM is matched by prefix based on its ID and name.
		`),
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			errutils.Check(func() error {
				var err error
				if so.VM, err = cmdutil.MatchSingleVM(args[0]); err != nil {
					return err
				}
				return run.SSH(so)
			}())
		},
	}

	addSSHFlags(cmd.Flags(), so)
	return cmd
}

func addSSHFlags(fs *pflag.FlagSet, so *run.SSHOptions) {
	fs.StringVarP(&so.IdentityFile, "identity", "i", "", "Override the VM's default identity file")
}
