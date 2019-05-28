package cmd

import (
	"fmt"
	"github.com/luxas/ignite/pkg/errutils"
	"github.com/luxas/ignite/pkg/filter"
	"github.com/luxas/ignite/pkg/metadata"
	"github.com/luxas/ignite/pkg/metadata/vmmd"
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
			err := RunAttach(out, cmd, args[0], true)
			errutils.Check(err)
		},
	}

	return cmd
}

// checkRunning can be used to skip the running check, this is used by CmdRun
// As the in-container ignite takes some time to start up and update the state
func RunAttach(out io.Writer, cmd *cobra.Command, vmMatch string, checkRunning bool) error {
	var md *vmmd.VMMetadata

	// Match a single VM using the VMFilter
	if matches, err := filter.NewFilterer(vmmd.NewVMFilter(vmMatch), metadata.VM.Path(), vmmd.LoadVMMetadata); err == nil {
		if filterable, err := matches.Single(); err == nil {
			if md, err = vmmd.ToVMMetadata(filterable); err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		return err
	}

	// Check if the VM is running
	if checkRunning && !md.Running() {
		return fmt.Errorf("%s is not running", md.ID)
	}

	// Print the ID before attaching
	fmt.Println(md.ID)

	dockerArgs := []string{
		"attach",
		md.ID,
	}

	// Attach to the VM in Docker
	if ec, err := util.ExecForeground("docker", dockerArgs...); err != nil {
		if ec != 1 { // Docker's detach sequence (^P^Q) has an exit code of -1
			return fmt.Errorf("failed to attach to container for VM %s: %v", md.ID, err)
		}
	}

	return nil
}
