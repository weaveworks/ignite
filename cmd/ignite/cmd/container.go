package cmd

import (
	"github.com/luxas/ignite/pkg/errutils"
	"github.com/spf13/cobra"
	"io"
)

// NewContainerCmd runs the dhcp server and sets up routing inside Docker
func NewCmdContainer(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Hidden: true,
		Args:   cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			err := RunContainer(out, cmd, args)
			errutils.Check(err)
		},
	}

	//cmd.Flags().StringP("output", "o", "", "Output format; available options are 'yaml', 'json' and 'short'")
	return cmd
}

// RunBuild runs when the Container command is invoked
func RunContainer(out io.Writer, cmd *cobra.Command, args []string) error {
	// The VM to run in container mode
	id := args[0]

	md := &vmMetadata{
		ID: id,
	}

	if err := md.load(); err != nil {
		return err
	}

	//util.ExecuteCommand("/firecracker")

	return nil
}
