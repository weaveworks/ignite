package source

import (
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"

	log "github.com/sirupsen/logrus"
	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/providers"
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
		log.Infof("Docker image %q not found locally, pulling...", source)
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

	// By default parse the OCI content ID from the Docker image ID
	contentRef := res.ID
	if len(res.RepoDigests) > 0 {
		// If the image has Repo digests, use the first one of them to parse
		// the fully qualified OCI image name and digest. All of the digests
		// point to the same contents, so it doesn't matter which one we use
		// here. It will be translated to the right content by the runtime.
		contentRef = res.RepoDigests[0]
	}

	// Parse the OCI content ID based on the available reference
	ci, err := meta.ParseOCIContentID(contentRef)
	if err != nil {
		return nil, err
	}

	return &api.OCIImageSource{
		ID:   ci,
		Size: meta.NewSizeFromBytes(uint64(res.Size)),
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
