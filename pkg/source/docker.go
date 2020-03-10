package source

import (
	"fmt"
	"io"
	"os/exec"

	log "github.com/sirupsen/logrus"
	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/providers"
)

// TODO: Make this a generic "OCISource" as it now only depends on the generic providers.Runtime
type DockerSource struct {
	imageRef    meta.OCIImageRef
	cleanupFunc func() error
	exportCmd   *exec.Cmd
}

// Compile-time assert to verify interface compatibility
var _ Source = &DockerSource{}

func NewDockerSource() *DockerSource {
	return &DockerSource{}
}

func (ds *DockerSource) Ref() meta.OCIImageRef {
	return ds.imageRef
}

func (ds *DockerSource) SetRef(imageRef meta.OCIImageRef) {
	ds.imageRef = imageRef
}

func (ds *DockerSource) Parse(ociRef meta.OCIImageRef) (*api.OCIImageSource, error) {
	res, err := providers.Runtime.InspectImage(ociRef)
	if err != nil {
		log.Infof("%s image %q not found locally, pulling...", providers.Runtime.Name(), ociRef)
		if err := providers.Runtime.PullImage(ociRef); err != nil {
			return nil, err
		}

		if res, err = providers.Runtime.InspectImage(ociRef); err != nil {
			return nil, err
		}
	}

	if res.Size == 0 || res.ID == nil {
		return nil, fmt.Errorf("parsing %s image %q data failed", providers.Runtime.Name(), ociRef)
	}

	ds.imageRef = ociRef

	return &api.OCIImageSource{
		ID:   res.ID,
		Size: meta.NewSizeFromBytes(uint64(res.Size)),
	}, nil
}

func (ds *DockerSource) Import(imagePath string) (map[meta.OCIImageRef]*api.OCIImageSource, error) {
	imgRefs, err := providers.Runtime.ImportImage(imagePath)
	if err != nil {
		return nil, err
	}

	img := map[meta.OCIImageRef]*api.OCIImageSource{}
	for _, imgRef := range imgRefs {
		// Inspect image and get image ID from inspect result. Use ID to create
		// image source. Image source can be used to create ignite image.
		res, err := providers.Runtime.InspectImage(imgRef)
		if err != nil {
			return nil, err
		}

		img[imgRef] = &api.OCIImageSource{
			ID:   res.ID,
			Size: meta.NewSizeFromBytes(uint64(res.Size)),
		}
	}
	return img, nil
}

func (ds *DockerSource) Reader() (rc io.ReadCloser, err error) {
	// Export the image
	rc, ds.cleanupFunc, err = providers.Runtime.ExportImage(ds.imageRef)
	return
}

func (ds *DockerSource) Cleanup() (err error) {
	if ds.cleanupFunc != nil {
		// Perform any post-export cleanup
		err = ds.cleanupFunc()
	}

	return
}
