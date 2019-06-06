package cmd

import (
	"io"

	"github.com/weaveworks/ignite/cmd/ignite/cmd/cmdutil"
	"github.com/weaveworks/ignite/cmd/ignite/run"
	"github.com/weaveworks/ignite/pkg/errutils"
	"github.com/spf13/cobra"
)

// NewCmdContainer runs the DHCP server and sets up routing inside Docker
func NewCmdContainer(out io.Writer) *cobra.Command {
	co := &run.ContainerOptions{}

	cmd := &cobra.Command{
		Use:    "container [id]",
		Hidden: true,
		Args:   cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			errutils.Check(func() error {
				var err error
				if co.VM, err = cmdutil.MatchSingleVM(args[0]); err != nil {
					return err
				}
				return run.Container(co)
			}())
		},
	}

	return cmd
}
