package run

import (
	"fmt"
	"github.com/luxas/ignite/pkg/constants"
	"github.com/luxas/ignite/pkg/metadata"
	"github.com/luxas/ignite/pkg/metadata/imgmd"
	"github.com/luxas/ignite/pkg/util"
	"github.com/mholt/archiver"
	"github.com/pkg/errors"
	"os"
	"path"
	"strings"
)

type BuildOptions struct {
	Source     string
	Name       string
	KernelName string
	image      *imgmd.ImageMetadata
	ImageNames []*metadata.Name
}

func Build(bo *BuildOptions) error {
	// Create new image metadata
	if err := bo.newImage(); err != nil {
		return err
	}

	// If the source is a file, import it as a file, otherwise check if it's a docker image
	tarFilePath := path.Join(bo.image.ObjectPath(), constants.IMAGE_TAR)
	if !util.FileExists(bo.Source) {
		// Query docker for the image
		out, err := util.ExecuteCommand("docker", "images", "-q", bo.Source)
		if err != nil {
			return err
		}

		if util.IsEmptyString(out) {
			return fmt.Errorf("docker image %s not found", bo.Source)
		}

		// Docker outputs one image per line
		dockerIDs := strings.Split(strings.TrimSpace(out), "\n")

		// Check if the image query is too broad
		if len(dockerIDs) > 1 {
			return fmt.Errorf("multiple matches, too broad docker image query: %s", bo.Source)
		}

		// Select the first (and only) match
		dockerID := dockerIDs[0]

		// Create a container from the image to be exported
		containerID, err := util.ExecuteCommand("docker", "create", dockerID, "sh")
		if err != nil {
			return errors.Wrapf(err, "failed to create docker container from image %s", dockerID)
		}

		// Export the created container to a tar archive that will be later extracted into the VM disk image
		_, err = util.ExecuteCommand("docker", "export", "-o", tarFilePath, containerID)
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

		if err := decompressor.DecompressFile(bo.Source, tarFilePath); err != nil {
			return err
		}
	}

	fo, err := os.Stat(tarFilePath)
	if err != nil {
		return err
	}

	// Create new file to host the filesystem and format it
	if err := bo.image.AllocateAndFormat(fo.Size()); err != nil {
		return err
	}

	// Add the files to the filesystem
	if err := bo.image.AddFiles(tarFilePath); err != nil {
		return err
	}

	if err := bo.image.Save(); err != nil {
		return err
	}

	// Import a new kernel from the image if specified
	if bo.KernelName != "" {
		dir, err := bo.image.ExportKernel()
		if err != nil {
			return err
		}

		if dir != "" {
			if err := ImportKernel(&ImportKernelOptions{
				Source: path.Join(dir, constants.KERNEL_FILE),
				Name:   bo.KernelName,
			}); err != nil {
				return err
			}

			if err := os.RemoveAll(dir); err != nil {
				return err
			}
		}
	}

	//if err := container.ExportToDocker(image); err != nil {
	//	return err
	//}

	// Print the ID of the newly generated image
	fmt.Println(bo.image.ID)

	return nil
}

// newImage creates the new image metadata
func (bo *BuildOptions) newImage() error {
	newID, err := util.NewID(constants.IMAGE_DIR)
	if err != nil {
		return err
	}

	// Verify the name
	name, err := metadata.NewName(bo.Name, &bo.ImageNames)
	if err != nil {
		return err
	}

	bo.image = imgmd.NewImageMetadata(newID, name)

	return nil
}
