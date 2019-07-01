package run

import (
	"log"
	"os"
	"path"

	"github.com/weaveworks/ignite/pkg/metadata/loader"

	"github.com/c2h5oh/datasize"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/metadata"
	"github.com/weaveworks/ignite/pkg/metadata/imgmd"
)

type ImportFlags struct {
	Name       string
	KernelName string
}

type importOptions struct {
	*ImportFlags
	source    string
	resLoader *loader.ResLoader
	newImage  *imgmd.ImageMetadata
	allImages []metadata.AnyMetadata
}

func (i *ImportFlags) NewImportOptions(l *loader.ResLoader, source string) (*importOptions, error) {
	io := &importOptions{ImportFlags: i, resLoader: l, source: source}

	if allImages, err := l.Images(); err == nil {
		io.allImages = *allImages
	} else {
		return nil, err
	}

	return io, nil
}

func Import(bo *importOptions) error {
	// Parse the source
	imageSrc, err := imgmd.NewSource(bo.source)
	if err != nil {
		return err
	}

	nameStr := bo.Name
	if len(imageSrc.DockerImage()) > 0 {
		nameStr = imageSrc.DockerImage()
	}

	// Verify the name
	name, err := metadata.NewNameWithLatest(nameStr, &bo.allImages)
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
	if err := bo.newImage.AllocateAndFormat(imageSrc.Size()); err != nil {
		return err
	}

	// Add the files to the filesystem
	if err := bo.newImage.AddFiles(imageSrc); err != nil {
		return err
	}

	if err := bo.newImage.Save(); err != nil {
		return err
	}
	hrsize := datasize.ByteSize(imageSrc.Size()).HR()
	log.Printf("Created a %s filesystem of the input", hrsize)

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
