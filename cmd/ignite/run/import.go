package run

import (
	"log"
	"os"
	"path"

	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/metadata"
	"github.com/weaveworks/ignite/pkg/metadata/imgmd"
	"github.com/weaveworks/ignite/pkg/metadata/loader"
	"github.com/weaveworks/ignite/pkg/source"
)

type importOptions struct {
	source    string
	resLoader *loader.ResLoader
	newImage  *imgmd.ImageMetadata
	allImages []metadata.AnyMetadata
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
	// Parse the source
	dockerSource := source.NewDockerSource()
	src, err := dockerSource.Parse(bo.source)
	if err != nil {
		return err
	}

	// Verify the name
	name, err := metadata.NewNameWithLatest(bo.source, &bo.allImages)
	if err != nil {
		return err
	}

	// Create new image metadata
	if bo.newImage, err = imgmd.NewImageMetadata(nil, name); err != nil {
		return err
	}
	defer bo.newImage.Cleanup(false) // TODO: Handle silent

	log.Println("Starting image import...")

	// Create new file to host the filesystem and format it
	if err := bo.newImage.AllocateAndFormat(src.Size.Int64()); err != nil {
		return err
	}

	// Add the files to the filesystem
	if err := bo.newImage.AddFiles(dockerSource); err != nil {
		return err
	}

	if err := bo.newImage.Save(); err != nil {
		return err
	}
	log.Printf("Created a %s filesystem of the input", src.Size.HR())

	// If the kernel already exists, don't try to import something with the same name
	if allKernels, err := bo.resLoader.Kernels(); err != nil {
		return err
	} else {
		if k, err := allKernels.MatchSingle(name.String()); k != nil && err == nil {
			return bo.newImage.Success()
		}
	}

	// Import a new kernel from the image if specified
	tmpKernelDir, err := bo.newImage.ExportKernel()
	if err == nil {
		io, err := (&ImportKernelFlags{
			Source: path.Join(tmpKernelDir, constants.KERNEL_FILE),
			Name:   name.String(),
		}).NewImportKernelOptions(bo.resLoader)
		if err != nil {
			return err
		}

		if err := ImportKernel(io); err != nil {
			return err
		}

		if err := os.RemoveAll(tmpKernelDir); err != nil {
			return err
		}

		//log.Printf("A kernel was imported from the image with name %q and ID %q", name.String(), kernelID)
	} else {
		// Tolerate the kernel to not be found
		if _, ok := err.(*imgmd.KernelNotFoundError); !ok {
			return err
		}
	}

	return bo.newImage.Success()
}
