package cmd

import (
	"github.com/luxas/ignite/pkg/errutils"
	"github.com/spf13/cobra"
	"io"
)

// NewCmdRun starts and attaches to a Firecracker VM
func NewCmdRun(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run [vm]",
		Short: "Start and attach to a Firecracker VM",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			err := RunRun(out, cmd, args)
			errutils.Check(err)
		},
	}
	//cmd.Flags().StringP("output", "o", "", "Output format; available options are 'yaml', 'json' and 'short'")
	return cmd
}

func RunRun(out io.Writer, cmd *cobra.Command, args []string) error {
	if err := RunStart(out, cmd, args); err == nil {
		if err = RunAttach(out, cmd, args, false); err != nil {
			return err
		}
	} else {
		return err
	}

	return nil
}
