package run

import (
	"github.com/weaveworks/ignite/pkg/client"
	"github.com/weaveworks/ignite/pkg/filter"
	"github.com/weaveworks/ignite/pkg/metadata/kernmd"
	"github.com/weaveworks/ignite/pkg/util"
)

type kernelsOptions struct {
	allKernels []*kernmd.Kernel
}

func NewKernelsOptions() (*kernelsOptions, error) {
	io := &kernelsOptions{}

	if allKernels, err := client.Kernels().FindAll(filter.NewAllFilter()); err == nil {
		io.allKernels = make([]*kernmd.Kernel, 0, len(allKernels))
		for _, kernel := range allKernels {
			io.allKernels = append(io.allKernels, &kernmd.Kernel{kernel})
		}
	} else {
		return nil, err
	}

	return io, nil
}

func Kernels(ko *kernelsOptions) error {
	o := util.NewOutput()
	defer o.Flush()

	o.Write("KERNEL ID", "NAME", "CREATED", "SIZE", "VERSION")
	for _, kernel := range ko.allKernels {
		o.Write(kernel.GetUID(), kernel.GetName(), kernel.GetCreated(), kernel.Status.OCISource.Size.String(), kernel.Status.Version)
	}

	return nil
}
