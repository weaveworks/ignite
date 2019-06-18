package source

import (
	"fmt"
	"github.com/weaveworks/ignite/pkg/dm"
	"io"
	"log"
	"os/exec"
	"strconv"
	"strings"

	"github.com/weaveworks/ignite/pkg/util"
)

type DockerSource struct {
	dockerImage string
	dockerID    string
	size        int64
	containerID string
	exportCmd   *exec.Cmd
}

// Compile-time assert to verify interface compatibility
var _ Source = &DockerSource{}

func NewDockerSource(src string) (*DockerSource, error) {
	// Query Docker for the image
	out, err := util.ExecuteCommand("docker", "images", "-q", src)
	if err != nil {
		return nil, err
	}

	// If the Docker image isn't found, try to pull it
	if util.IsEmptyString(out) {
		log.Printf("Docker image %q not found locally, pulling...", src)
		if _, err := util.ExecForeground("docker", "pull", src); err != nil {
			return nil, err
		}

		out, err = util.ExecuteCommand("docker", "images", "-q", src)
		if err != nil {
			return nil, err
		}

		if util.IsEmptyString(out) {
			return nil, fmt.Errorf("docker image %s could not be found", src)
		}
	}

	// Docker outputs one image per line
	dockerIDs := strings.Split(strings.TrimSpace(out), "\n")

	// Check if the image query is too broad
	if len(dockerIDs) > 1 {
		return nil, fmt.Errorf("multiple matches, Docker image query too broad: %q", src)
	}

	// Select the first (and only) match
	dockerID := dockerIDs[0]

	// Get the size of the Docker image
	out, err = util.ExecuteCommand("docker", "inspect", src, "-f", "{{.Size}}")
	if err != nil {
		return nil, err
	}

	// Parse the size from the output
	size, err := strconv.Atoi(out)
	if err != nil {
		return nil, err
	}

	return &DockerSource{
		dockerImage: src,
		dockerID:    dockerID,
		size:        int64(size),
	}, nil
}

func (ds *DockerSource) Reader() (io.ReadCloser, error) {
	// Create a container from the image to be exported
	var err error
	if ds.containerID, err = util.ExecuteCommand("docker", "create", ds.dockerID, "sh"); err != nil {
		return nil, fmt.Errorf("failed to create Docker container from image %q: %v", ds.dockerID, err)
	}

	// Open a tar file stream to be extracted straight into the VM disk image
	ds.exportCmd = exec.Command("docker", "export", ds.containerID)
	reader, err := ds.exportCmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := ds.exportCmd.Start(); err != nil {
		return nil, err
	}

	return reader, nil
}

func (ds *DockerSource) DockerImage() string {
	return ds.dockerImage
}

func (ds *DockerSource) SizeBytes() int64 {
	return ds.size
}

// Get the size as 512-byte sectors
func (ds *DockerSource) SizeSectors() dm.Sectors {
	return dm.SectorsFromBytes(ds.size)
}

func (ds *DockerSource) ID() string {
	return ds.dockerID
}

func (ds *DockerSource) Cleanup() error {
	if len(ds.containerID) > 0 {
		// Remove the temporary container
		if _, err := util.ExecuteCommand("docker", "rm", ds.containerID); err != nil {
			return fmt.Errorf("failed to remove container %q: %v", ds.containerID, err)
		}
	}

	return nil
}
