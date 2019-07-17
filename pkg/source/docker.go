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

func (ds *DockerSource) Reader() (rc io.ReadCloser, err error) {
	// Get the Docker client
	dc, err := docker.GetDockerClient()
	if err != nil {
		return nil, err
	}

	// Export the image
	rc, ds.containerID, err = dc.ExportImage(ds.imageID)
	return
}

func (ds *DockerSource) Cleanup() (err error) {
	if len(ds.containerID) > 0 {
		// Get the Docker client
		dc, err := docker.GetDockerClient()
		if err != nil {
			return err
		}

		// Remove the temporary container
		err = dc.RemoveContainer(ds.containerID)
	}

	return
}
