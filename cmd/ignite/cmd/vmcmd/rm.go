package vmcmd

import (
	"io"

	"github.com/lithammer/dedent"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/ignite/cmd/ignite/cmd/cmdutil"
	"github.com/weaveworks/ignite/cmd/ignite/run"
)

// NewCmdRm removes VMs
func NewCmdRm(out io.Writer) *cobra.Command {
	rf := &run.RmFlags{}

	cmd := &cobra.Command{
		Use:   "rm <vm>...",
		Short: "Remove VMs",
		Long: dedent.Dedent(`
			Remove one or multiple VMs. The VMs are matched by prefix based
			on their ID and name. To remove multiple VMs, chain the matches
			separated by spaces. The force flag (-f, --force) kills running
			VMs before removal instead of throwing an error.
		`),
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(func() error {
				ro, err := rf.NewRmOptions(args)
				if err != nil {
					return err
				}

				return run.Rm(ro)
			}())
		},
	}

	addRmFlags(cmd.Flags(), rf)
	return cmd
}

func addRmFlags(fs *pflag.FlagSet, rf *run.RmFlags) {
	cmdutil.AddForceFlag(fs, &rf.Force)
	cmdutil.AddConfigFlag(fs, &rf.ConfigFile)
}
