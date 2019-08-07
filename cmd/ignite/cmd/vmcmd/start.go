package vmcmd

import (
	"io"

	"github.com/lithammer/dedent"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/ignite/cmd/ignite/cmd/cmdutil"
	"github.com/weaveworks/ignite/cmd/ignite/run"
)

// NewCmdStart starts a VM
func NewCmdStart(out io.Writer) *cobra.Command {
	sf := &run.StartFlags{}

	cmd := &cobra.Command{
		Use:   "start <vm>",
		Short: "Start a VM",
		Long: dedent.Dedent(`
			Start the given VM. The VM is matched by prefix based on its ID and name.
			If the interactive flag (-i, --interactive) is specified, attach to the
			VM after starting.
		`),
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(func() error {
				so, err := sf.NewStartOptions(args[0])
				if err != nil {
					return err
				}

				return run.Start(so)
			}())
		},
	}

	addStartFlags(cmd.Flags(), sf)
	return cmd
}

func addStartFlags(fs *pflag.FlagSet, sf *run.StartFlags) {
	cmdutil.AddInteractiveFlag(fs, &sf.Interactive)
	fs.BoolVarP(&sf.Debug, "debug", "d", false, "Debug mode, keep container after VM shutdown")
}
