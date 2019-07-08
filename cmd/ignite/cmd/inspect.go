package cmd

import (
	"github.com/lithammer/dedent"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/ignite/cmd/ignite/run"
	"github.com/weaveworks/ignite/pkg/errutils"
	"io"
)

// NewCmdInspect inspects an Ignite Object
func NewCmdInspect(out io.Writer) *cobra.Command {
	i := &run.InspectFlags{}

	cmd := &cobra.Command{
		Use:   "inspect <kind> <object>",
		Short: "Inspect an Ignite Object",
		Long: dedent.Dedent(`
			Retrieve information about the given object of the given kind.
			The kind can be "image", "kernel" or "vm". The object is matched
			by prefix based on its ID and name. Outputs JSON by default, can
			be overridden to YAML with the yaml flag (-y, --yaml).
		`),
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			errutils.Check(func() error {
				io, err := i.NewInspectOptions(args[0], args[1])
				if err != nil {
					return err
				}

				return run.Inspect(io)
			}())
		},
	}

	addInspectFlags(cmd.Flags(), i)
	return cmd
}

func addInspectFlags(fs *pflag.FlagSet, i *run.InspectFlags) {
	fs.BoolVarP(&i.YAMLOutput, "yaml", "y", false, "Output the object in YAML format instead of JSON")
}
