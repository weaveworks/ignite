package run

import (
	"fmt"

	"github.com/weaveworks/ignite/pkg/client"
	"github.com/weaveworks/ignite/pkg/filter"
	"github.com/weaveworks/ignite/pkg/metadata"
	"github.com/weaveworks/ignite/pkg/operations"
	"github.com/weaveworks/ignite/pkg/storage/filterer"
)

type importOptions struct {
	source string
}

func NewImportOptions(source string) (*importOptions, error) {
	return &importOptions{source: source}, nil
}

func Import(bo *importOptions) error {
	runImage, err := operations.ImportImage(bo.source)
	if err != nil {
		return err
	}
	defer metadata.Cleanup(runImage, false) // TODO: Handle silent

	// If the kernel already exists, don't try to import something with the same name
	if _, err := client.Kernels().Find(filter.NewNameFilter(runImage.GetName())); err == nil {
		return metadata.Success(runImage)
	} else {
		fmt.Printf("err %T %v", err, err)
		switch err.(type) {
		case *filterer.AmbiguousError:
			// such a kernel seem to already exist
			return metadata.Success(runImage)
		case *filterer.NonexistentError:
			// If the kernel did not exist, let's import it
			fmt.Println("kernel nonexistent, importing")
		default:
			// other, unknown error
			return err
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
		fmt.Println("no kernel in image")
		return metadata.Success(runImage)
	}
	defer metadata.Cleanup(runKernel, false) // TODO: Handle silent

	// both the image and kernel are imported successfully
	metadata.Success(runKernel)
	return metadata.Success(runImage)
}
