package run

import (
	"os"

	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	"github.com/weaveworks/ignite/pkg/providers"
	"github.com/weaveworks/ignite/pkg/util"
	"github.com/weaveworks/libgitops/pkg/filter"
)

type ImagesOptions struct {
	allImages []*api.Image
}

func NewImagesOptions() (io *ImagesOptions, err error) {
	io = &ImagesOptions{}
	io.allImages, err = providers.Client.Images().FindAll(filter.NewAllFilter())
	// If the storage is uninitialized, avoid failure and continue with empty
	// image list.
	if err != nil && os.IsNotExist(err) {
		err = nil
	}
	return
}

func Images(io *ImagesOptions) error {
	o := util.NewOutput()
	defer o.Flush()

	o.Write("IMAGE ID", "NAME", "CREATED", "SIZE")
	for _, image := range io.allImages {
		o.Write(image.GetUID(), image.GetName(), image.GetCreated(), image.Status.OCISource.Size.String())
	}

	return nil
}
