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
		Use:   "cp <vm> <source> <dest>",
		Short: "Copy a file into a running vm",
		Long: dedent.Dedent(`
			Copy a file from host into running VM.
			Uses SCP to SSH into the running VM using the private key created for it during generation.
			If no private key was created or wanting to use a different identity file,
			use the identity file flag (-i, --identity) to override the used identity file.
			The given VM is matched by prefix based on its ID and name.
		`),
		Args: cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(func() error {
				co, err := cf.NewCPOptions(args[0], args[1], args[2])
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
	fs.StringVarP(&cf.IdentityFile, "identity", "i", "", "Override the vm's default identity file")
	fs.Uint32VarP(&cf.Timeout, "timeout", "t", 10, "Timeout waiting for connection in seconds")
}
