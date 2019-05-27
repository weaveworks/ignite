package cmd

import (
	"fmt"
	"github.com/luxas/ignite/pkg/constants"
	"github.com/luxas/ignite/pkg/filter"
	"github.com/luxas/ignite/pkg/metadata"
	"github.com/luxas/ignite/pkg/metadata/vmmd"
	"github.com/luxas/ignite/pkg/util"
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
	var md *vmmd.VMMetadata

	// Match a single VM using the VMFilter
	if matches, err := filter.NewFilterer(vmmd.NewVMFilter(args[0]), metadata.VM.Path(), vmmd.LoadVMMetadata); err == nil {
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

	// Check if the given VM is already running
	if md.Running() {
		return fmt.Errorf("%s is already running", md.ID)
	}

	igniteBinary, _ := filepath.Abs(os.Args[0])

	dockerArgs := []string{
		"run",
		"-itd",
		"--rm",
		"--name",
		md.ID,
		fmt.Sprintf("-v=%s:/ignite/ignite", igniteBinary),
		fmt.Sprintf("-v=%s:%s", constants.DATA_DIR, constants.DATA_DIR),
		fmt.Sprintf("--stop-timeout=%d", constants.STOP_TIMEOUT+constants.IGNITE_TIMEOUT),
		"--privileged",
		"--device=/dev/kvm",
		"ignite",
		md.ID,
	}

	// Start the VM in docker
	if _, err := util.ExecuteCommand("docker", dockerArgs...); err != nil {
		return fmt.Errorf("failed to start container for VM %q: %v", md.ID, err)
	}

	// Print the ID of the started VM
	fmt.Println(md.ID)

	return nil
}
