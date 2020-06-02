package vmcmd

import (
	"io"

	"github.com/lithammer/dedent"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/ignite/cmd/ignite/cmd/cmdutil"
	"github.com/weaveworks/ignite/cmd/ignite/run"
)

// NewCmdExec exec's into a running VM
func NewCmdExec(out io.Writer, err io.Writer, in io.Reader) *cobra.Command {
	ef := &run.ExecFlags{}

	cmd := &cobra.Command{
		Use:   "exec <vm> <command...>",
		Short: "execute a command in a running VM",
		Long: dedent.Dedent(`
			Execute a command in a running VM using SSH and the private key created for it during generation.
			If no private key was created or wanting to use a different identity file,
			use the identity file flag (-i, --identity) to override the used identity file.
			The given VM is matched by prefix based on its ID and name.
		`),
		Args: cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(func() error {
				eo, err := ef.NewExecOptions(args[0], args[1:]...)
				if err != nil {
					return err
				}

				return run.Exec(eo)
			}())
		},
	}

	addExecFlags(cmd.Flags(), ef)
	return cmd
}

func addExecFlags(fs *pflag.FlagSet, ef *run.ExecFlags) {
	cmdutil.AddSSHFlags(fs, &ef.IdentityFile, &ef.Timeout)
	fs.BoolVarP(&ef.Tty, "tty", "t", false, "Allocate a pseudo-TTY")
}
