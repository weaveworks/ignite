package run

import (
	"fmt"

	"github.com/c2h5oh/datasize"
	"github.com/weaveworks/ignite/pkg/metadata/kernmd"
	"github.com/weaveworks/ignite/pkg/metadata/loader"
	"github.com/weaveworks/ignite/pkg/util"
)

type kernelsOptions struct {
	allKernels []*kernmd.Kernel
}

func NewKernelsOptions(l *loader.ResLoader) (*kernelsOptions, error) {
	io := &kernelsOptions{}

	if allKernels, err := l.Kernels(); err == nil {
		io.allKernels = allKernels.MatchAll()
	} else {
		return nil, err
	}

	return io, nil
}

func Kernels(ko *kernelsOptions) error {
	o := util.NewOutput()
	defer o.Flush()

	o.Write("KERNEL ID", "CREATED", "SIZE", "NAME")
	for _, md := range ko.allKernels {
		size, err := md.Size()
		if err != nil {
			return fmt.Errorf("failed to get size for %s %q: %v", md.Type(), md.GetUID(), err)
		}

		o.Write(md.GetUID(), md.Created, datasize.ByteSize(size).HR(), md.GetName())
	}

	return nil
}
