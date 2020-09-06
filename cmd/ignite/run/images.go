package run

import (
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
