package vmcmd

import (
	"io"

	"github.com/lithammer/dedent"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/ignite/cmd/ignite/cmd/cmdutil"
	"github.com/weaveworks/ignite/cmd/ignite/run"
)

// NewCmdCP CP's a file into a running vm
func NewCmdCP(out io.Writer) *cobra.Command {
	cf := &run.CPFlags{}

	cmd := &cobra.Command{
		Use:   "cp <source> <dest>",
		Short: "Copy files/folders between a running vm and the local filesystem",
		Long: dedent.Dedent(`
			Copy a file between host and a running VM.
			Creates an SFTP connection to the running VM using the private key created for
			it during generation, and transfers files between the host and VM. If no
			private key was created or wanting to use a different identity file, use the
			identity file flag (-i, --identity) to override the used identity file.

			Example usage:
				$ ignite cp localfile.txt my-vm:remotefile.txt
				$ ignite cp my-vm:remotefile.txt localfile.txt
		`),
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(func() error {
				co, err := cf.NewCPOptions(args[0], args[1])
				if err != nil {
					return err
				}

				return run.CP(co)
			}())
		},
	}

	addCPFlags(cmd.Flags(), cf)
	return cmd
}

func addCPFlags(fs *pflag.FlagSet, cf *run.CPFlags) {
	cmdutil.AddSSHFlags(fs, &cf.IdentityFile, &cf.Timeout)
}
