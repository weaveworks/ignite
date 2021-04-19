package vmcmd

import (
	"io"

	"github.com/lithammer/dedent"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/ignite/cmd/ignite/cmd/cmdutil"
	"github.com/weaveworks/ignite/cmd/ignite/run"
)

// NewCmdPs lists running VMs
func NewCmdPs(out io.Writer) *cobra.Command {
	pf := &run.PsFlags{}

	cmd := &cobra.Command{
		Use:     "ps",
		Short:   "List running VMs",
		Aliases: []string{"ls", "list"},
		Long: dedent.Dedent(`
			List all running VMs. By specifying the all flag (-a, --all),
			also list VMs that are not currently running.
			Using the -f (--filter) flag, you can give conditions VMs should fullfilled to be displayed.
			You can filter on all the underlying fields of the VM struct, see the documentation:
			https://ignite.readthedocs.io/en/stable/api/ignite_v1alpha4#VM.

			Different operators can be used:
			- "=" and "==" for the equal
			- "!=" for the is not equal
			- "=~" for the contains
			- "!~" for the not contains

			Non-exhaustive list of identifiers to apply filter on:
			- the VM name
			- CPUs usage
			- Labels
			- Image
			- Kernel
			- Memory

			Example usage:
				$ ignite ps -f "{{.ObjectMeta.Name}}=my-vm2,{{.Spec.CPUs}}!=3,{{.Spec.Image.OCI}}=~weaveworks/ignite-ubuntu"

				$ ignite ps -f "{{.Spec.Memory}}=~1024,{{.Status.Running}}=true"
		`),
		Run: func(cmd *cobra.Command, args []string) {
			// If `ps` is called via any of its aliases
			// (`ls`, `list`), list all VMs
			if cmd.CalledAs() != cmd.Name() {
				pf.All = true
			}

			cmdutil.CheckErr(func() error {
				po, err := pf.NewPsOptions()
				if err != nil {
					return err
				}

				return run.Ps(po)
			}())
		},
	}

	addPsFlags(cmd.Flags(), pf)
	return cmd
}

func addPsFlags(fs *pflag.FlagSet, pf *run.PsFlags) {
	fs.BoolVarP(&pf.All, "all", "a", false, "Show all VMs, not just running ones")
	fs.StringVarP(&pf.Filter, "filter", "f", "", "Filter the VMs")
	fs.StringVarP(&pf.TemplateFormat, "template", "t", "", "Format the output using the given Go template")
}
