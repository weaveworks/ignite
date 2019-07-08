package vmcmd

import (
	"io"

	"github.com/lithammer/dedent"
	"github.com/spf13/pflag"

	"github.com/spf13/cobra"
	"github.com/weaveworks/ignite/cmd/ignite/run"
	"github.com/weaveworks/ignite/pkg/errutils"
)

// NewCmdSSH ssh's into a running vm
func NewCmdSSH(out io.Writer) *cobra.Command {
	sf := &run.SSHFlags{}

	cmd := &cobra.Command{
		Use:   "ssh <vm>",
		Short: "SSH into a running vm",
		Long: dedent.Dedent(`
			SSH into the running VM using the private key created for it during generation.
			If no private key was created or wanting to use a different identity file,
			use the identity file flag (-i, --identity) to override the used identity file.
			The given VM is matched by prefix based on its ID and name.
		`),
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			errutils.Check(func() error {
				so, err := sf.NewSSHOptions(args[0])
				if err != nil {
					return err
				}

				return run.SSH(so)
			}())
		},
	}

	addSSHFlags(cmd.Flags(), sf)
	return cmd
}

func addSSHFlags(fs *pflag.FlagSet, sf *run.SSHFlags) {
	fs.StringVarP(&sf.IdentityFile, "identity", "i", "", "Override the vm's default identity file")
}
