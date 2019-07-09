package source

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os/exec"

	api "github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/runtime/docker"
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

func (ds *DockerSource) ID() string {
	return ds.imageID
}

func (ds *DockerSource) Parse(ociRef meta.OCIImageRef) (*api.OCIImageSource, error) {
	client, err := docker.GetDockerClient()
	if err != nil {
		return nil, err
	}

	source := ociRef.String()
	res, err := client.InspectImage(source)
	if err != nil {
		log.Printf("Docker image %q not found locally, pulling...", source)
		rc, err := client.PullImage(source)
		if err != nil {
			return nil, err
		}

		// Don't output the pull command
		io.Copy(ioutil.Discard, rc)
		rc.Close()
		res, err = client.InspectImage(source)
		if err != nil {
			return nil, err
		}
	}

	if res.Size == 0 || len(res.ID) == 0 {
		return nil, fmt.Errorf("parsing docker image %q data failed", source)
	}

	ds.imageID = res.ID
	return &api.OCIImageSource{
		ID:          res.ID,
		RepoDigests: res.RepoDigests,
		Size:        meta.NewSizeFromBytes(uint64(res.Size)),
	}, nil
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
