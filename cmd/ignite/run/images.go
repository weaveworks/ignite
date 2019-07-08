package run

import (
	"github.com/weaveworks/ignite/pkg/client"
	"github.com/weaveworks/ignite/pkg/filter"
	"github.com/weaveworks/ignite/pkg/metadata/imgmd"
	"github.com/weaveworks/ignite/pkg/util"
)

type imagesOptions struct {
	allImages []*imgmd.Image
}

func NewImagesOptions() (*imagesOptions, error) {
	io := &imagesOptions{}

	if allImages, err := client.Images().FindAll(filter.NewAllFilter()); err == nil {
		io.allImages = make([]*imgmd.Image, 0, len(allImages))
		for _, image := range allImages {
			io.allImages = append(io.allImages, &imgmd.Image{image})
		}
	} else {
		return nil, err
	}

	return io, nil
}

func Images(io *imagesOptions) error {
	o := util.NewOutput()
	defer o.Flush()

	o.Write("IMAGE ID", "NAME", "CREATED", "SIZE")
	for _, image := range io.allImages {
		o.Write(image.GetUID(), image.GetName(), image.GetCreated(), image.Status.OCISource.Size.String())
	}

	return nil
}
