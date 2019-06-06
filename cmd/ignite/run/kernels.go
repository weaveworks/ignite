package run

import (
	"fmt"

	"github.com/c2h5oh/datasize"
	"github.com/weaveworks/ignite/pkg/metadata/kernmd"
	"github.com/weaveworks/ignite/pkg/util"
)

type KernelsOptions struct {
	Kernels []*kernmd.KernelMetadata
}

func Kernels(ko *KernelsOptions) error {
	o := util.NewOutput()
	defer o.Flush()

	o.Write("KERNEL ID", "CREATED", "SIZE", "NAME")
	for _, md := range ko.Kernels {
		size, err := md.Size()
		if err != nil {
			return fmt.Errorf("failed to get size for %s %q: %v", md.Type, md.ID, err)
		}

		o.Write(md.ID, md.Created, datasize.ByteSize(size).HR(), md.Name.String())
	}

	return nil
}
