package run

import (
	"log"

	"github.com/weaveworks/ignite/pkg/source"

	"github.com/weaveworks/ignite/cmd/ignite/run/runutil"

	"github.com/weaveworks/ignite/pkg/metadata"
	"github.com/weaveworks/ignite/pkg/metadata/kernmd"
)

type ImportKernelFlags struct {
	Name string
}

type importKernelOptions struct {
	*ImportKernelFlags
	source     string
	allKernels []metadata.AnyMetadata
}

func (i *ImportKernelFlags) NewImportKernelOptions(l *runutil.ResLoader, source string) (*importKernelOptions, error) {
	io := &importKernelOptions{ImportKernelFlags: i, source: source}

	if allKernels, err := l.Kernels(); err == nil {
		io.allKernels = *allKernels
	} else {
		return nil, err
	}

	return io, nil
}

func ImportKernel(io *importKernelOptions) error {
	// Parse the source
	kernelSrc, err := source.NewDockerSource(io.source)
	if err != nil {
		return err
	}

	// Verify the name
	name, err := metadata.NewNameWithLatest(io.Name, &io.allKernels)
	if err != nil {
		return err
	}

	// Create new kernel metadata
	md, err := kernmd.NewKernelMetadata(nil, name)
	if err != nil {
		return err
	}
	defer md.Cleanup(false) // TODO: Handle silent

	log.Println("Starting kernel import...")

	// Create a new image file to host the filesystem and format it
	if imageFile, err := md.CreateImageFile(kernelSrc.SizeBytes()); err == nil {
		// Add the files to the filesystem and export vmlinux
		if err := md.AddFiles(imageFile, kernelSrc); err != nil {
			return err
		}
	} else {
		return err
	}

	// Save the metadata
	if err := md.Save(); err != nil {
		return err
	}

	return md.Success()
}
