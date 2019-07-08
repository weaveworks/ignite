package run

import (
	"github.com/weaveworks/ignite/pkg/client"
	"github.com/weaveworks/ignite/pkg/metadata"
	"github.com/weaveworks/ignite/pkg/metadata/imgmd"
	"github.com/weaveworks/ignite/pkg/metadata/kernmd"
	"github.com/weaveworks/ignite/pkg/operations"
)

func ImportImage(source string) (*imgmd.Image, error) {
	runImage, err := operations.FindOrImportImage(client.DefaultClient, source)
	if err != nil {
		return nil, err
	}
	defer metadata.Cleanup(runImage, false) // TODO: Handle silent
	return runImage, metadata.Success(runImage)
}

func ImportKernel(source string) (*kernmd.Kernel, error) {
	runKernel, err := operations.FindOrImportKernel(client.DefaultClient, source)
	if err != nil {
		return nil, err
	}
	defer metadata.Cleanup(runKernel, false) // TODO: Handle silent
	return runKernel, metadata.Success(runKernel)
}
