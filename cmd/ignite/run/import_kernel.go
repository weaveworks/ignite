package run

import (
	"fmt"

	"github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"
	"github.com/weaveworks/ignite/pkg/metadata/loader"

	"github.com/weaveworks/ignite/pkg/metadata"
	"github.com/weaveworks/ignite/pkg/metadata/kernmd"
	"github.com/weaveworks/ignite/pkg/util"
)

type ImportKernelFlags struct {
	Source string
	Name   string
}

type importKernelOptions struct {
	*ImportKernelFlags
	allKernels []metadata.Metadata
}

func (i *ImportKernelFlags) NewImportKernelOptions(l *loader.ResLoader) (*importKernelOptions, error) {
	io := &importKernelOptions{ImportKernelFlags: i}

	if allKernels, err := l.Kernels(); err == nil {
		io.allKernels = *allKernels
	} else {
		return nil, err
	}

	return io, nil
}

func ImportKernel(ao *importKernelOptions) error {
	if !util.FileExists(ao.Source) {
		return fmt.Errorf("not a kernel image: %s", ao.Source)
	}

	// TODO: Kernel importing from docker when moving to pool/snapshotter
	kernel := &v1alpha1.Kernel{
		Spec: v1alpha1.KernelSpec{
			Version: "unknown",
			Source: v1alpha1.ImageSource{
				Type: "file",
				ID:   "-",
				Name: "-",
			},
		},
	}

	// Verify the name
	name, err := metadata.NewNameWithLatest(ao.Name, &ao.allKernels)
	if err != nil {
		return err
	}

	// Create new kernel metadata
	md, err := kernmd.NewKernelMetadata("", &name, kernel)
	if err != nil {
		return err
	}
	defer metadata.Cleanup(md, false) // TODO: Handle silent

	// Save the metadata
	if err := md.Save(); err != nil {
		return err
	}

	// Perform the copy
	if err := md.ImportKernel(ao.Source); err != nil {
		return err
	}

	return metadata.Success(md)
}
