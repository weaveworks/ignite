package vmcmd

import (
	"io"

	"github.com/lithammer/dedent"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/ignite/cmd/ignite/cmd/cmdutil"
	"github.com/weaveworks/ignite/cmd/ignite/run"
	"github.com/weaveworks/ignite/pkg/errutils"
)

// NewCmdRun creates, starts (and attaches to) a VM
func NewCmdRun(out io.Writer) *cobra.Command {
	ro := &run.RunOptions{}

	cmd := &cobra.Command{
		Use:   "run [image] [kernel]",
		Short: "Create a new VM and start it",
		Long: dedent.Dedent(`
			Create and start a new VM immediately. The image and kernel are matched by
			prefix based on their ID and name. This command accepts all flags used to
			create and start a VM. The interactive flag (-i, --interactive) can be
			specified to immediately attach to the started VM after creation.

			Example usage:
				$ ignite run my-image my-kernel \
					--interactive \
					--name my-vm \
					--cpus 2 \
					--memory 2048 \
					--size 10G
		`),
		Args: cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			errutils.Check(func() error {
				var err error
				if ro.Image, err = cmdutil.MatchSingleImage(args[0]); err != nil {
					return err
				}
				if ro.Kernel, err = cmdutil.MatchSingleKernel(args[1]); err != nil {
					return err
				}
				if ro.VMNames, err = cmdutil.MatchAllVMNames(); err != nil {
					return err
				}
				return run.Run(ro)
			}())
		},
	}

	addRunFlags(cmd.Flags(), ro)
	return cmd
}

func addRunFlags(fs *pflag.FlagSet, ro *run.RunOptions) {
	addCreateFlags(fs, &ro.CreateOptions)
	addStartFlags(fs, &ro.StartOptions)
}
