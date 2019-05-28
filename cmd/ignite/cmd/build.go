package cmd

import (
	"fmt"
	"github.com/luxas/ignite/pkg/constants"
	"github.com/luxas/ignite/pkg/metadata/imgmd"
	"github.com/luxas/ignite/pkg/util"
	"github.com/pkg/errors"
	"io"
	"path"
	"strings"

	"github.com/luxas/ignite/pkg/errutils"
	"github.com/mholt/archiver"
	"github.com/spf13/cobra"
)

// buildOptions specifies the properties of a new image
type buildOptions struct {
	source    string
	imageID   string
	imageDir  string
	imageName string
}

// NewCmdBuild builds a new VM base image
func NewCmdBuild(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build [source] [name]",
		Short: "Build a Firecracker VM base image",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			err := RunBuild(out, cmd, args)
			errutils.Check(err)
		},
	}

	//cmd.Flags().StringP("output", "o", "", "Output format; available options are 'yaml', 'json' and 'short'")
	return cmd
}

func RunBuild(out io.Writer, cmd *cobra.Command, args []string) error {
	// Construct the build options for a new image
	buildOptions, err := newBuildOptions(cmd, args)
	if err != nil {
		return err
	}

	// If the source is a file, import it as a file, otherwise check if it's a docker image
	if !util.FileExists(buildOptions.source) {
		// Query docker for the image
		out, err := util.ExecuteCommand("docker", "images", "-q", buildOptions.source)
		if err != nil {
			return err
		}

		if util.IsEmptyString(out) {
			return fmt.Errorf("docker image %s not found", buildOptions.source)
		}

		// Docker outputs one image per line
		dockerIDs := strings.Split(strings.TrimSpace(out), "\n")

		// Check if the image query is too broad
		if len(dockerIDs) > 1 {
			return fmt.Errorf("multiple matches, too broad docker image query: %s", buildOptions.source)
		}

		// Select the first (and only) match
		dockerID := dockerIDs[0]

		// Create a container from the image to be exported
		containerID, err := util.ExecuteCommand("docker", "create", dockerID, "sh")
		if err != nil {
			return errors.Wrapf(err, "failed to create docker container from image %s", dockerID)
		}

		// Export the created container to a tar archive that will be later extracted into the VM disk image
		_, err = util.ExecuteCommand("docker", "export", "-o", path.Join(buildOptions.imageDir, constants.IMAGE_TAR), containerID)
		if err != nil {
			return errors.Wrapf(err, "failed to export created container %s:", containerID)
		}

		// Remove the temporary container
		_, err = util.ExecuteCommand("docker", "rm", containerID)
		if err != nil {
			return errors.Wrapf(err, "failed to remove container %s:", containerID)
		}

	} else {
		// Decompress given file to later be extracted into the disk image
		// TODO: Either extract directly into image or intermediate tarfile from different formats
		decompressor := archiver.FileCompressor{
			Decompressor: archiver.NewGz(),
		}

		if err := decompressor.DecompressFile(buildOptions.source, path.Join(buildOptions.imageDir, constants.IMAGE_TAR)); err != nil {
			return err
		}
	}

	// Create new image metadata
	md := imgmd.NewImageMetadata(buildOptions.imageID, buildOptions.imageName)

	// Create new file to host the filesystem and format it
	if err := md.AllocateAndFormat(); err != nil {
		return err
	}

	// Add the files to the filesystem
	if err := md.AddFiles3(path.Join(buildOptions.imageDir, constants.IMAGE_TAR)); err != nil {
		return err
	}

	if err := md.Save(); err != nil {
		return err
	}

	//if err := container.ExportToDocker(image); err != nil {
	//	return err
	//}

	// Print the ID of the newly generated image
	fmt.Println(buildOptions.imageID)

	return nil
}

// newBuildOptions constructs a set of options for new VMs
func newBuildOptions(cmd *cobra.Command, args []string) (*buildOptions, error) {
	newID, err := util.NewID(constants.IMAGE_DIR)
	if err != nil {
		return nil, err
	}

	return &buildOptions{
		source:    args[0], // The source (tar path/docker image) is given as the first argument
		imageID:   newID,
		imageDir:  path.Join(constants.IMAGE_DIR, newID),
		imageName: args[1], // The name is given as the second argument
	}, nil
}
