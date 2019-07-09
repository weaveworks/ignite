package run

import (
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/client"
	"github.com/weaveworks/ignite/pkg/metadata"
	"github.com/weaveworks/ignite/pkg/metadata/imgmd"
	"github.com/weaveworks/ignite/pkg/metadata/kernmd"
	"github.com/weaveworks/ignite/pkg/operations"
)

func ImportImage(source string) (*imgmd.Image, error) {
	ociRef, err := meta.NewOCIImageRef(source)
	if err != nil {
		return nil, err
	}

	runImage, err := operations.FindOrImportImage(client.DefaultClient, ociRef)
	if err != nil {
		return nil, err
	}
	defer metadata.Cleanup(runImage, false) // TODO: Handle silent

	return runImage, metadata.Success(runImage)
}

func ImportKernel(source string) (*kernmd.Kernel, error) {
	ociRef, err := meta.NewOCIImageRef(source)
	if err != nil {
		return nil, err
	}

	runKernel, err := operations.FindOrImportKernel(client.DefaultClient, ociRef)
	if err != nil {
		return nil, err
	}
	defer metadata.Cleanup(runKernel, false) // TODO: Handle silent

	return runKernel, metadata.Success(runKernel)
}
