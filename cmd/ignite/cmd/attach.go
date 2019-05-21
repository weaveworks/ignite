package cmd

import (
	"errors"
	"fmt"
	"github.com/luxas/ignite/pkg/constants"
	"github.com/luxas/ignite/pkg/util"
	"io"
	"path"

	"github.com/luxas/ignite/pkg/errutils"
	"github.com/spf13/cobra"
)

// NewCmdExec execs a command in a Firecracker VM.
func NewCmdAttach(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "attach [vm]",
		Short: "Attach to a running Firecracker VM",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			err := RunAttach(out, cmd, args)
			errutils.Check(err)
		},
	}
	//cmd.Flags().StringP("output", "o", "", "Output format; available options are 'yaml', 'json' and 'short'")
	return cmd
}

// RunExec runs when the Exec command is invoked
func RunAttach(out io.Writer, cmd *cobra.Command, args []string) error {
	vmID := args[0]

	// Check if given vm exists TODO: Selection by name
	if !util.DirExists(path.Join(constants.VM_DIR, vmID)) {
		return fmt.Errorf("not a vm: %s", vmID)
	}

	md, err := loadVMMetadata(vmID)
	if err != nil {
		return err
	}

	if !md.running() {
		return errors.New("given VM is not running")
	}

	dockerArgs := []string{
		"attach",
		vmID,
	}

	// Attach to the VM in Docker
	if ec, err := util.ExecForeground("docker", dockerArgs...); err != nil {
		if ec != 1 { // Docker's detach sequence (^P^Q) has an exit code of -1
			return fmt.Errorf("failed to attach to container for VM %s: %v", vmID, err)
		}
	}

	return nil
}
