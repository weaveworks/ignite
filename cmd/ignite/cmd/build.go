package cmd

import (
	"fmt"
	"github.com/luxas/ignite/pkg/constants"
	"github.com/luxas/ignite/pkg/image"
	"github.com/luxas/ignite/pkg/util"
	"github.com/pkg/errors"
	"io"
	"os"
	"path"

	"github.com/luxas/ignite/pkg/errutils"
	"github.com/mholt/archiver"
	"github.com/spf13/cobra"
)

// buildOptions specifies the properties of a new VM
type buildOptions struct {
	tarPath string
	vmID    string
	vmDir   string
}

// NewCmdBuild builds a Firecracker VM.
func NewCmdBuild(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build [tar]",
		Short: "Build a Firecracker VM",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			err := RunBuild(out, cmd, args)
			errutils.Check(err)
		},
	}

	//cmd.Flags().StringP("output", "o", "", "Output format; available options are 'yaml', 'json' and 'short'")
	return cmd
}

// RunBuild runs when the Build command is invoked
func RunBuild(out io.Writer, cmd *cobra.Command, args []string) error {
	// Construct the build options for a new VM
	buildOptions, err := newBuildOptions(cmd, args)
	if err != nil {
		return err
	}

	// Check if given tarfile exists
	if !util.FileExists(buildOptions.tarPath) {
		return fmt.Errorf("input %q is not a file", buildOptions.tarPath)
	}

	// Decompress given file to later be extracted into the VM disk image
	if err := archiver.DecompressFile(buildOptions.tarPath, path.Join(buildOptions.vmDir, constants.VM_FS_TAR)); err != nil {
		return err
	}

	// Create the unique directory for the VM
	if err := os.MkdirAll(buildOptions.vmDir, os.ModePerm); err != nil {
		return errors.Wrap(err, "failed to create VM directory")
	}

	return nil
}

// newBuildOptions constructs a set of options for new VMs
func newBuildOptions(cmd *cobra.Command, args []string) (*buildOptions, error) {
	newID, err := build.NewVMID()
	if err != nil {
		return nil, err
	}

	return &buildOptions{
		tarPath: args[0], // The tar path is given as the first argument
		vmID:    newID,
		vmDir:   path.Join(constants.VM_DIR, newID),
	}, nil
}
