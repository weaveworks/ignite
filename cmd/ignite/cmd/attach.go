package cmd

import (
	"fmt"
	"github.com/luxas/ignite/pkg/errutils"
	"github.com/luxas/ignite/pkg/filter"
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
	// Load all VM metadata as Filterable objects
	mdf, err := metadata.LoadVMMetadataFilterable()
	if err != nil {
		return err
	}

	// Create an IDNameFilter to match a single VM
	d, err := filter.NewFilterer(metadata.NewVMFilter(args[0])).Single(mdf)
	if err != nil {
		return err
	}

	// Convert the result Filterable to a VMMetadata
	md, err := metadata.ToVMMetadata(d)
	if err != nil {
		return err
	}

	// Check if the VM is running
	if !md.Running() {
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
