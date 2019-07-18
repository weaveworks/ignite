package source

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os/exec"

	"github.com/weaveworks/ignite/pkg/providers"

	api "github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
)

// TODO: Make this a generic "OCISource" as it now only depends on the generic providers.Runtime
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
	source := ociRef.String()
	res, err := providers.Runtime.InspectImage(source)
	if err != nil {
		log.Printf("Docker image %q not found locally, pulling...", source)
		rc, err := providers.Runtime.PullImage(source)
		if err != nil {
			return nil, err
		}

		// Don't output the pull command
		io.Copy(ioutil.Discard, rc)
		rc.Close()
		res, err = providers.Runtime.InspectImage(source)
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
	// Export the image
	rc, ds.containerID, err = providers.Runtime.ExportImage(ds.imageID)
	return
}

func (ds *DockerSource) Cleanup() (err error) {
	if len(ds.containerID) > 0 {
		// Remove the temporary container
		err = providers.Runtime.RemoveContainer(ds.containerID)
	}

	return
}
