package vmcmd

import (
	"io"

	"github.com/lithammer/dedent"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/ignite/cmd/ignite/cmd/cmdutil"
	"github.com/weaveworks/ignite/cmd/ignite/run"
	"github.com/weaveworks/ignite/pkg/errutils"
	"github.com/weaveworks/ignite/pkg/logs"
	"github.com/weaveworks/ignite/pkg/metadata"
)

// NewCmdRun creates, starts (and attaches to) a VM
func NewCmdRun(out io.Writer) *cobra.Command {
	ro := &run.RunOptions{}

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
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			errutils.Check(func() error {
				var err error
				ro.Image, err = cmdutil.MatchSingleImage(args[0])
				if err != nil {
					// Tolerate a nonexistent error, but return the error otherwise
					if _, ok := err.(*metadata.NonexistentError); !ok {
						return err
					}
					allImages, err := cmdutil.MatchAllImageNames()
					if err != nil {
						return err
					}
					// If the image doesn't exist, build it
					if _, err := run.Build(&run.BuildOptions{
						Source:     args[0],
						ImageNames: allImages,
					}); err != nil {
						return err
					}
					ro.Image, _ = cmdutil.MatchSingleImage(args[0])
				}
				// TODO: deduplicate this from create.go code
				if len(ro.KernelName) == 0 {
					ro.KernelName = args[0]
				}
				if ro.Kernel, err = cmdutil.MatchSingleKernel(ro.KernelName); err != nil {
					return err
				}
				if ro.VMNames, err = cmdutil.MatchAllVMNames(); err != nil {
					return err
				}
				return logs.PrintMachineReadableID(run.Run(ro))
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
