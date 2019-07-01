package cmd

import (
	"io"

	"github.com/weaveworks/ignite/pkg/metadata/loader"

	"github.com/spf13/cobra"
	"github.com/weaveworks/ignite/cmd/ignite/run"
	"github.com/weaveworks/ignite/pkg/errutils"
)

// NewCmdContainer runs the DHCP server and sets up routing inside Docker
func NewCmdContainer(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "container <vm>",
		Hidden: true,
		Args:   cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			errutils.Check(func() error {
				var err error
				co, err := run.NewContainerOptions(loader.NewResLoader(), args[0])
				if err != nil {
					return err
				}
				return run.Container(co)
			}())
		},
	}

	return cmd
}
