package run

import (
	"fmt"

	"github.com/c2h5oh/datasize"
	"github.com/weaveworks/ignite/cmd/ignite/run/runutil"
	"github.com/weaveworks/ignite/pkg/metadata/kernmd"
	"github.com/weaveworks/ignite/pkg/util"
)

type kernelsOptions struct {
	allKernels []*kernmd.KernelMetadata
}

func NewKernelsOptions(l *runutil.ResLoader) (*kernelsOptions, error) {
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
			return fmt.Errorf("failed to get size for %s %q: %v", md.Type, md.ID, err)
		}

		o.Write(md.ID, md.Created, datasize.ByteSize(size).HR(), md.Name.String())
	}

	return nil
}
