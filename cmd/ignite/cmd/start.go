package cmd

import (
	"fmt"
	"github.com/luxas/ignite/pkg/constants"
	"github.com/luxas/ignite/pkg/metadata"
	"github.com/luxas/ignite/pkg/util"
	"github.com/pkg/errors"
	"io"
	"os"
	"path/filepath"

	"github.com/luxas/ignite/pkg/errutils"
	"github.com/spf13/cobra"
)

// NewCmdStart starts a Firecracker VM
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

func RunStart(out io.Writer, cmd *cobra.Command, args []string) error {
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

	// Check if the given VM is already running
	if md.running() {
		return fmt.Errorf("%s is already running", vmID)
	}

	igniteBinary, _ := filepath.Abs(os.Args[0])

	dockerArgs := []string{
		"run",
		"-itd",
		"--rm",
		"--name",
		vmID,
		fmt.Sprintf("-v=%s:/ignite/ignite", igniteBinary),
		fmt.Sprintf("-v=%s:%s", constants.DATA_DIR, constants.DATA_DIR),
		"--privileged",
		"--device=/dev/kvm",
		"ignite",
		vmID,
	}

	// Start the VM in docker
	if _, err := util.ExecuteCommand("docker", dockerArgs...); err != nil {
		return errors.Wrapf(err, "failed to start container for VM: %s", vmID)
	}

	// Print the ID of the started VM
	fmt.Println(vmID)

	return nil
}
