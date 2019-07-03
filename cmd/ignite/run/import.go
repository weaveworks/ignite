package run

import (
	"log"
	"os"
	"path"

	api "github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"

	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/metadata"
	"github.com/weaveworks/ignite/pkg/metadata/imgmd"
	"github.com/weaveworks/ignite/pkg/metadata/loader"
	"github.com/weaveworks/ignite/pkg/source"
)

type importOptions struct {
	source    string
	resLoader *loader.ResLoader
	newImage  *imgmd.Image
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
	// Parse the source
	dockerSource := source.NewDockerSource()
	src, err := dockerSource.Parse(bo.source)
	if err != nil {
		return err
	}

	image := &api.Image{
		Spec: api.ImageSpec{
			Source: *src,
		},
	}

	// Verify the name
	name, err := metadata.NewNameWithLatest(bo.source, &bo.allImages)
	if err != nil {
		return err
	}

	// Create new image metadata
	if bo.newImage, err = imgmd.NewImage("", &name, image); err != nil {
		return err
	}
	defer metadata.Cleanup(bo.newImage, false) // TODO: Handle silent

	log.Println("Starting image import...")

	// Create new file to host the filesystem and format it
	if err := bo.newImage.AllocateAndFormat(); err != nil {
		return err
	}

	// Add the files to the filesystem
	if err := bo.newImage.AddFiles(dockerSource); err != nil {
		return err
	}

	if err := bo.newImage.Save(); err != nil {
		return err
	}
	log.Printf("Created imported a %s filesystem", image.Spec.Source.Size.HR())

	// If the kernel already exists, don't try to import something with the same name
	if allKernels, err := bo.resLoader.Kernels(); err != nil {
		return err
	} else {
		if k, err := allKernels.MatchSingle(name); k != nil && err == nil {
			return metadata.Success(bo.newImage)
		}
	}

	// Import a new kernel from the image if specified
	tmpKernelDir, err := bo.newImage.ExportKernel()
	if err == nil {
		io, err := (&ImportKernelFlags{
			Source: path.Join(tmpKernelDir, constants.KERNEL_FILE),
			Name:   name,
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

	return metadata.Success(bo.newImage)
}
