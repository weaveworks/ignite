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

// buildOptions specifies the properties of a new image
type buildOptions struct {
	tarPath  string
	imageID  string
	imageDir string
}

// NewCmdBuild builds a new VM base image
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
	// Construct the build options for a new image
	buildOptions, err := newBuildOptions(cmd, args)
	if err != nil {
		return err
	}

	// Check if given tarfile exists
	if !util.FileExists(buildOptions.tarPath) {
		return fmt.Errorf("input %q is not a file", buildOptions.tarPath)
	}

	// Create the unique directory for the VM
	if err := os.MkdirAll(buildOptions.imageDir, os.ModePerm); err != nil {
		return errors.Wrap(err, "failed to create VM directory")
	}

	// Decompress given file to later be extracted into the disk image
	// TODO: Archiver has somehow forgotten how to extract a .tar.gz
	// (format specified by source filename is not a recognized compression algorithm)
	// Figure something out, maybe we need a custom extractor?
	fmt.Println(path.Join(buildOptions.imageDir, constants.IMAGE_TAR))
	if err := archiver.DecompressFile(buildOptions.tarPath, path.Join(buildOptions.imageDir, constants.IMAGE_TAR)); err != nil {
		return err
	}

	// Create new file to host the filesystem and format it
	image := build.NewImage(buildOptions.imageID)

	if err := image.AllocateAndFormat(); err != nil {
		return err
	}

	return nil
}

// newBuildOptions constructs a set of options for new VMs
func newBuildOptions(cmd *cobra.Command, args []string) (*buildOptions, error) {
	newID, err := build.NewImageID()
	if err != nil {
		return nil, err
	}

	return &buildOptions{
		tarPath:  args[0], // The tar path is given as the first argument
		imageID:  newID,
		imageDir: path.Join(constants.IMAGE_DIR, newID),
	}, nil
}
