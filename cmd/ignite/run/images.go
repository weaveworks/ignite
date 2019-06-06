package run

import (
	"fmt"

	"github.com/c2h5oh/datasize"
	"github.com/luxas/ignite/pkg/metadata/imgmd"
	"github.com/luxas/ignite/pkg/util"
)

type ImagesOptions struct {
	Images []*imgmd.ImageMetadata
}

func Images(io *ImagesOptions) error {
	o := util.NewOutput()
	defer o.Flush()

	o.Write("IMAGE ID", "CREATED", "SIZE", "NAME")
	for _, md := range io.Images {
		size, err := md.Size()
		if err != nil {
			return fmt.Errorf("failed to get size for %s %q: %v", md.Type, md.ID, err)
		}

		o.Write(md.ID, md.Created, datasize.ByteSize(size).HR(), md.Name.String())
	}

	return nil
}
