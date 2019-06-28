package source

import (
	"fmt"
	"io"
	"k8s.io/apimachinery/pkg/util/json"
	"log"
	"os/exec"

	"github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"

	"github.com/weaveworks/ignite/pkg/util"
)

type DockerSource struct {
	imageID     string
	containerID string
	exportCmd   *exec.Cmd
}

// Compile-time assert to verify interface compatibility
var _ Source = &DockerSource{}

func NewDockerSource() *DockerSource {
	return &DockerSource{}
}

type inspect struct {
	ID       string   `json:"Id"`
	Size     uint64   `json:"Size"`
	RepoTags []string `json:"RepoTags"`
}

type errNotFound struct {
	source string
}

// Compile-time assert to verify interface compatibility
var _ error = &errNotFound{}

func newErrNotFound(source string) *errNotFound {
	return &errNotFound{source}
}

func (e *errNotFound) Error() string {
	return fmt.Sprintf("docker image %q could not be found", e.source)
}

func parseInspect(source string) (*inspect, error) {
	out, err := util.ExecuteCommand("docker", "inspect", source)
	if err != nil {
		return nil, err
	}

	if util.IsEmptyString(out) {
		return nil, newErrNotFound(source)
	}

	// Docker inspect outputs an array containing the struct we need
	var result []*inspect

	if err := json.Unmarshal([]byte(out), &result); err != nil {
		return nil, err
	}

	// Extract the struct from the array
	data := result[0]

	if data.Size == 0 || len(data.ID) == 0 || len(data.RepoTags) == 0 {
		return nil, fmt.Errorf("parsing docker image %q data failed", source)
	}

	return data, nil
}

func (ds *DockerSource) Parse(input *v1alpha1.ImageSource) error {
	// Use the ID to match the image
	// If it's not given, fall back to the name
	source := input.ID
	if len(source) == 0 {
		source = input.Name
	}

	var err error
	var imageData *inspect

	// Query Docker for the image
	for imageData == nil {
		imageData, err = parseInspect(source)

		switch err.(type) {
		case nil:
			// Success
		case *errNotFound:
			source = input.Name // Fall back to the name, as docker pull doesn't accept IDs

			log.Printf("Docker image %q not found locally, pulling...", source)
			if _, err := util.ExecForeground("docker", "pull", source); err != nil {
				return err
			}
		default:
			return err
		}
	}

	// Update the fields
	// TODO: Add just missing fields
	input.Type = v1alpha1.ImageSourceTypeDocker
	input.ID = imageData.ID
	input.Name = imageData.RepoTags[0]
	input.Size = v1alpha1.NewSizeFromBytes(imageData.Size)

	ds.imageID = input.ID
	return nil
}

func (ds *DockerSource) Reader() (io.ReadCloser, error) {
	// Create a container from the image to be exported
	var err error
	if ds.containerID, err = util.ExecuteCommand("docker", "create", ds.imageID, "sh"); err != nil {
		return nil, fmt.Errorf("failed to create Docker container from image %q: %v", ds.imageID, err)
	}

	// Open a tar file stream to be extracted straight into the image volume
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

func (ds *DockerSource) Cleanup() error {
	if len(ds.containerID) > 0 {
		// Remove the temporary container
		if _, err := util.ExecuteCommand("docker", "rm", ds.containerID); err != nil {
			return fmt.Errorf("failed to remove container %q: %v", ds.containerID, err)
		}
	}

	return nil
}
