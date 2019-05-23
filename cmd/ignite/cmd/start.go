package cmd

import (
	"fmt"
	"github.com/luxas/ignite/pkg/constants"
	"github.com/luxas/ignite/pkg/filter"
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
		"--privileged",
		"--device=/dev/kvm",
		"ignite",
		md.ID,
	}

	// Start the VM in docker
	if _, err := util.ExecuteCommand("docker", dockerArgs...); err != nil {
		return errors.Wrapf(err, "failed to start container for VM: %s", md.ID)
	}

	// Print the ID of the started VM
	fmt.Println(md.ID)

	return nil
}
