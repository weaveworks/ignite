package run

import (
	"log"
	"os"

	"github.com/c2h5oh/datasize"
	"github.com/weaveworks/ignite/pkg/metadata"
	"github.com/weaveworks/ignite/pkg/metadata/imgmd"
	"github.com/weaveworks/ignite/pkg/metadata/loader"
	"github.com/weaveworks/ignite/pkg/source"
)

type ImportFlags struct {
	Name string
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

func Import(io *importOptions) error {
	// Parse the source
	imageSrc, err := source.NewDockerSource(io.source)
	if err != nil {
		return err
	}

	nameStr := io.Name
	if len(imageSrc.DockerImage()) > 0 {
		nameStr = imageSrc.DockerImage()
	}

	// Verify the name
	name, err := metadata.NewNameWithLatest(nameStr, &io.allImages)
	if err != nil {
		return err
	}

	// Create new image metadata
	if io.newImage, err = imgmd.NewImageMetadata(metadata.IDFromSource(imageSrc), name); err != nil {
		return err
	}
	defer io.newImage.Cleanup(false) // TODO: Handle silent

	log.Println("Starting image import...")

	// TODO: Handle directory creation/removal separately
	os.MkdirAll(io.newImage.ObjectPath(), 0755)

	// Create a new DM pool for the image
	if err := io.newImage.NewDMPool(); err != nil {
		return err
	}

	// Create a new image file to host the filesystem and format it
	//imageFile, err := io.newImage.CreateImageFile(imageSrc.Size())
	//if err != nil {
	//	return err
	//}

	// Add the files to the filesystem
	if err := io.newImage.AddFiles(imageSrc); err != nil {
		return err
	}

	if err := io.newImage.Save(); err != nil {
		return err
	}

	// TODO: This
	log.Printf("Created a %s filesystem for the image", datasize.ByteSize(0).HR())

	return io.newImage.Success()
}
