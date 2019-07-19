package run

import (
	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/metadata"
	"github.com/weaveworks/ignite/pkg/operations"
	"github.com/weaveworks/ignite/pkg/providers"
)

func ImportImage(source string) (*api.Image, error) {
	ociRef, err := meta.NewOCIImageRef(source)
	if err != nil {
		return nil, err
	}

	image, err := operations.FindOrImportImage(providers.Client, ociRef)
	if err != nil {
		return nil, err
	}
	defer metadata.Cleanup(image, false) // TODO: Handle silent

	return image, metadata.Success(image)
}

func ImportKernel(source string) (*api.Kernel, error) {
	ociRef, err := meta.NewOCIImageRef(source)
	if err != nil {
		return nil, err
	}

	kernel, err := operations.FindOrImportKernel(providers.Client, ociRef)
	if err != nil {
		return nil, err
	}
	defer metadata.Cleanup(kernel, false) // TODO: Handle silent

	return kernel, metadata.Success(kernel)
}
