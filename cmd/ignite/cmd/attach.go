package cmd

import (
	"fmt"
	"github.com/luxas/ignite/pkg/errutils"
	"github.com/luxas/ignite/pkg/metadata"
	"github.com/luxas/ignite/pkg/util"
	"github.com/spf13/cobra"
	"io"
)

// NewCmdAttach attaches to a running Firecracker VM
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

	return cmd
}

func RunAttach(out io.Writer, cmd *cobra.Command, args []string) error {
	// Get the VM ID
	vmID, err := metadata.MatchObject(args[0], metadata.VM)
	if err != nil {
		return err
	}

	// Load the VM metadata
	md, err := loadVMMetadata(vmID)
	if err != nil {
		return err
	}

	// Check if the VM is running
	if !md.running() {
		return fmt.Errorf("%s is not running", vmID)
	}

	// Print the ID before attaching
	fmt.Println(vmID)

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
