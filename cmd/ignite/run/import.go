package run

import (
	"github.com/weaveworks/ignite/pkg/metadata"
	"github.com/weaveworks/ignite/pkg/metadata/loader"
	"github.com/weaveworks/ignite/pkg/operations"
)

type importOptions struct {
	source    string
	resLoader *loader.ResLoader
	allImages []metadata.Metadata
}

func NewImportOptions(l *loader.ResLoader, source string) (*importOptions, error) {
	io := &importOptions{resLoader: l, source: source}

	if allImages, err := l.Images(); err == nil {
		io.allImages = *allImages
	} else {
		return nil, err
	}

	return io, nil
}

func Import(bo *importOptions) error {
	runImage, err := operations.ImportImage(bo.source, &bo.allImages)
	if err != nil {
		return err
	}
	defer metadata.Cleanup(runImage, false) // TODO: Handle silent

	// If the kernel already exists, don't try to import something with the same name
	if allKernels, err := bo.resLoader.Kernels(); err != nil {
		return err
	} else {
		// check if a kernel with that name already exists and return if there is
		if k, err := allKernels.MatchSingle(runImage.GetName()); k != nil && err == nil {
			return metadata.Success(runImage)
		}
	}
	// at this point we know that there is no kernel with the same name as the image
	// import a kernel from the image
	runKernel, err := operations.ImportKernelFromImage(runImage)
	if err != nil {
		return err
	}
	if runKernel == nil {
		// there was no kernel in the image, that's fine too
		return metadata.Success(runImage)
	}
	defer metadata.Cleanup(runKernel, false) // TODO: Handle silent

	// both the image and kernel are imported successfully
	metadata.Success(runKernel)
	return metadata.Success(runImage)
}
