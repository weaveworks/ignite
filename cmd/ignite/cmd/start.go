package cmd

import (
	"fmt"
	"github.com/luxas/ignite/pkg/constants"
	"github.com/luxas/ignite/pkg/util"
	"github.com/pkg/errors"
	"io"
	"path"

	"github.com/luxas/ignite/pkg/errutils"
	"github.com/spf13/cobra"
)

// NewCmdStart starts a Firecracker VM.
func NewCmdStart(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start [vm]",
		Short: "Start a Firecracker VM",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			err := RunStart(out, cmd, args)
			errutils.Check(err)
		},
	}
	//cmd.Flags().StringP("output", "o", "", "Output format; available options are 'yaml', 'json' and 'short'")
	return cmd
}

// RunStart runs when the start command is invoked
func RunStart(out io.Writer, cmd *cobra.Command, args []string) error {
	vmID := args[0]

	// Check if given vm exists TODO: Selection by name
	if !util.DirExists(path.Join(constants.VM_DIR, vmID)) {
		return fmt.Errorf("not a vm: %s", vmID)
	}

	md := &vmMetadata{
		ID: vmID,
	}

	if err := md.load(); err != nil {
		return err
	}

	if md.State == Running {
		return errors.New("given VM is already running")
	}

	// Start the vm in docker
	if _, err := util.ExecuteCommand("docker", "run", "-itd", "--rm", "ignite", "container", vmID); err != nil {
		return errors.Wrapf(err, "failed to start container for VM: %s", vmID)
	}

	return nil
}
