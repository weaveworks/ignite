package run

import (
	"fmt"

	"github.com/weaveworks/ignite/cmd/ignite/run/runutil"

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
	allKernels []metadata.AnyMetadata
}

func (i *ImportKernelFlags) NewImportKernelOptions(l *runutil.ResLoader) (*importKernelOptions, error) {
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

	// Verify the name
	name, err := metadata.NewNameWithLatest(ao.Name, &ao.allKernels)
	if err != nil {
		return err
	}

	// Create new kernel metadata
	md, err := kernmd.NewKernelMetadata(nil, name)
	if err != nil {
		return err
	}
	defer md.Cleanup(false) // TODO: Handle silent

	// Save the metadata
	if err := md.Save(); err != nil {
		return err
	}

	// Perform the copy
	if err := md.ImportKernel(ao.Source); err != nil {
		return err
	}

	return md.Success()
}
