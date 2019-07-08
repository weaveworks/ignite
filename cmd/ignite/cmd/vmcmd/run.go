package vmcmd

import (
	"io"

	"github.com/lithammer/dedent"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/ignite/cmd/ignite/run"
	"github.com/weaveworks/ignite/pkg/errutils"
)

// NewCmdRun creates, starts (and attaches to) a VM
func NewCmdRun(out io.Writer) *cobra.Command {
	rf := &run.RunFlags{
		CreateFlags: run.NewCreateFlags(),
		StartFlags:  &run.StartFlags{},
	}

	cmd := &cobra.Command{
		Use:   "run <image>",
		Short: "Create a new VM and start it",
		Long: dedent.Dedent(`
			Create and start a new VM immediately. The image and kernel are matched by
			prefix based on their ID and name. This command accepts all flags used to
			create and start a VM. The interactive flag (-i, --interactive) can be
			specified to immediately attach to the started VM after creation.

			Example usage:
				$ ignite run weaveworks/ignite-ubuntu \
					--kernel weaveworks/ignite-ubuntu \
					--interactive \
					--name my-vm \
					--cpus 2 \
					--memory 2048 \
					--size 10G
		`),
		Args: cobra.RangeArgs(0, 1),
		Run: func(cmd *cobra.Command, args []string) {
			errutils.Check(func() error {
				ro, err := rf.NewRunOptions(args)
				if err != nil {
					return err
				}

				return run.Run(ro)
			}())
		},
	}

	addRunFlags(cmd.Flags(), rf)
	return cmd
}

func addRunFlags(fs *pflag.FlagSet, rf *run.RunFlags) {
	addCreateFlags(fs, rf.CreateFlags)
	addStartFlags(fs, rf.StartFlags)
}
