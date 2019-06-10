package vmcmd

import (
	"github.com/weaveworks/ignite/cmd/ignite/run/runutil"
	"io"

	"github.com/lithammer/dedent"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/ignite/cmd/ignite/run"
	"github.com/weaveworks/ignite/pkg/errutils"
	"github.com/weaveworks/ignite/pkg/logs"
	"github.com/weaveworks/ignite/pkg/metadata"
)

// NewCmdRun creates, starts (and attaches to) a VM
func NewCmdRun(out io.Writer) *cobra.Command {
	rf := &run.RunFlags{}

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
				// TODO: Clean this mess up
				var err error
				ro.Image, err = runutil.MatchSingleImage(args[0])
				if err != nil {
					// Tolerate a nonexistent error, but return the error otherwise
					if _, ok := err.(*metadata.NonexistentError); !ok {
						return err
					}
					allImages, err := runutil.MatchAllImageNames()
					if err != nil {
						return err
					}
					// If the image doesn't exist, build it
					if _, err := run.Import(&run.ImportOptions{
						Source:     args[0],
						ImageNames: allImages,
					}); err != nil {
						return err
					}
					ro.Image, _ = runutil.MatchSingleImage(args[0])
				}
				// TODO: deduplicate this from create.go code
				if len(ro.KernelName) == 0 {
					ro.KernelName = args[0]
				}
				if ro.Kernel, err = runutil.MatchSingleKernel(ro.KernelName); err != nil {
					return err
				}
				if ro.VMNames, err = runutil.MatchAllVMNames(); err != nil {
					return err
				}

				ro, err := rf.NewRunOptions(runutil.NewResLoader(), args[0])
				if err != nil {
					return err
				}

				return logs.PrintMachineReadableID(run.Run(ro))
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
