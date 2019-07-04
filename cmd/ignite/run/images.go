package run

import (
	"fmt"

	"github.com/c2h5oh/datasize"
	"github.com/weaveworks/ignite/pkg/metadata/imgmd"
	"github.com/weaveworks/ignite/pkg/metadata/loader"
	"github.com/weaveworks/ignite/pkg/util"
)

type imagesOptions struct {
	allImages []*imgmd.Image
}

func NewImagesOptions(l *loader.ResLoader) (*imagesOptions, error) {
	io := &imagesOptions{}

	if allImages, err := l.Images(); err == nil {
		io.allImages = allImages.MatchAll()
	} else {
		return nil, err
	}

	return io, nil
}

func Images(io *imagesOptions) error {
	o := util.NewOutput()
	defer o.Flush()

	o.Write("IMAGE ID", "CREATED", "SIZE", "NAME")
	for _, md := range io.allImages {
		size, err := md.Size()
		if err != nil {
			return fmt.Errorf("failed to get size for %s %q: %v", md.GetKind(), md.GetUID(), err)
		}

		o.Write(md.GetUID(), md.Created, datasize.ByteSize(size).HR(), md.GetName())
	}

	return nil
}
