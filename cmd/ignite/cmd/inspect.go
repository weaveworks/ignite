package cmd

import (
	"io"

	"github.com/lithammer/dedent"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/ignite/cmd/ignite/cmd/cmdutil"
	"github.com/weaveworks/ignite/cmd/ignite/run"
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
			be overridden with the output flag (-o, --output).

			Example usage:
				$ ignite inspect vm my-vm

				$ ignite inspect vm my-vm -t {{.Status.IPAddresses}}

				$ ignite inspect vm my-vm -t {{.ObjectMeta.Name}}

				$ ignite inspect vm my-vm -t {{.Spec.Image.OCI}}
		`),
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(func() error {
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
	fs.StringVarP(&i.OutputFormat, "output", "o", "json", "Output the object in the specified format")
	fs.StringVarP(&i.TemplateFormat, "template", "t", "", "Format the output using the given Go template")
}
