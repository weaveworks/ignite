package runutil

import (
	"github.com/weaveworks/ignite/pkg/filter"
	"github.com/weaveworks/ignite/pkg/metadata"
	"github.com/weaveworks/ignite/pkg/metadata/imgmd"
	"github.com/weaveworks/ignite/pkg/metadata/kernmd"
	"github.com/weaveworks/ignite/pkg/metadata/vmmd"
)

type allVMs []metadata.AnyMetadata
type allImages []metadata.AnyMetadata
type allKernels []metadata.AnyMetadata

type ResLoader struct {
	vm     allVMs
	image  allImages
	kernel allKernels
}

func NewResLoader() *ResLoader {
	return &ResLoader{}
}

func (l *ResLoader) VMs() (*allVMs, error) {
	if l.vm == nil {
		var err error
		if l.vm, err = vmmd.LoadAllVMMetadata(); err != nil {
			return nil, err
		}
	}

	return &l.vm, nil
}

func (l *ResLoader) Images() (*allImages, error) {
	if l.image == nil {
		var err error
		if l.image, err = imgmd.LoadAllImageMetadata(); err != nil {
			return nil, err
		}
	}

	return &l.image, nil
}

func (l *ResLoader) Kernels() (*allKernels, error) {
	if l.kernel == nil {
		var err error
		if l.kernel, err = kernmd.LoadAllKernelMetadata(); err != nil {
			return nil, err
		}
	}

	return &l.kernel, nil
}

func single(f metadata.Filter, sources []metadata.AnyMetadata) (metadata.AnyMetadata, error) {
	var result metadata.AnyMetadata

	// Match a single AnyMetadata using the given filter
	if matches, err := filter.NewFilterer(f, sources); err == nil {
		if result, err = matches.Single(); err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}

	return result, nil
}

func matchIndividual(fs []metadata.Filter, sources []metadata.AnyMetadata) ([]metadata.AnyMetadata, error) {
	results := make([]metadata.AnyMetadata, 0, len(sources))

	for _, f := range fs {
		md, err := single(f, sources)
		if err != nil {
			return nil, err
		}

		results = append(results, md)
	}

	return results, nil
}

// Match a single VM using the VMFilter
func (l *allVMs) MatchSingle(match string) (*vmmd.VMMetadata, error) {
	md, err := single(vmmd.NewVMFilter(match), *l)
	if err != nil {
		return nil, err
	}

	return vmmd.ToVMMetadata(md), nil
}

// Match multiple individual VMs with different filter strings
func (l *allVMs) MatchMultiple(matches []string) ([]*vmmd.VMMetadata, error) {
	filters := make([]metadata.Filter, 0, len(matches))
	for _, match := range matches {
		filters = append(filters, vmmd.NewVMFilter(match))
	}

	results, err := matchIndividual(filters, *l)
	if err != nil {
		return nil, err
	}

	return vmmd.ToVMMetadataAll(results), nil
}

func (l *allVMs) MatchFilter(all bool) ([]*vmmd.VMMetadata, error) {
	matches, err := filter.NewFilterer(vmmd.NewVMFilterAll("", all), *l)
	if err != nil {
		return nil, err
	}

	return vmmd.ToVMMetadataAll(matches.All()), nil
}

func (l *allVMs) MatchAll() []*vmmd.VMMetadata {
	return vmmd.ToVMMetadataAll(*l)
}

// Match a single image using the IDNameFilter
func (l *allImages) MatchSingle(match string) (*imgmd.ImageMetadata, error) {
	md, err := single(metadata.NewIDNameFilter(match, metadata.Image), *l)
	if err != nil {
		return nil, err
	}

	return imgmd.ToImageMetadata(md), nil
}

// Match multiple individual images with different filter strings
func (l *allImages) MatchMultiple(matches []string) ([]*imgmd.ImageMetadata, error) {
	filters := make([]metadata.Filter, 0, len(matches))
	for _, match := range matches {
		filters = append(filters, metadata.NewIDNameFilter(match, metadata.Image))
	}

	results, err := matchIndividual(filters, *l)
	if err != nil {
		return nil, err
	}

	return imgmd.ToImageMetadataAll(results), nil
}

func (l *allImages) MatchAll() []*imgmd.ImageMetadata {
	return imgmd.ToImageMetadataAll(*l)
}

// Match a single kernel using the IDNameFilter
func (l *allKernels) MatchSingle(match string) (*kernmd.KernelMetadata, error) {
	md, err := single(metadata.NewIDNameFilter(match, metadata.Kernel), *l)
	if err != nil {
		return nil, err
	}

	return kernmd.ToKernelMetadata(md), nil
}

// Match multiple individual kernels with different filter strings
func (l *allKernels) MatchMultiple(matches []string) ([]*kernmd.KernelMetadata, error) {
	filters := make([]metadata.Filter, 0, len(matches))
	for _, match := range matches {
		filters = append(filters, metadata.NewIDNameFilter(match, metadata.Kernel))
	}

	results, err := matchIndividual(filters, *l)
	if err != nil {
		return nil, err
	}

	return kernmd.ToKernelMetadataAll(results), nil
}

func (l *allKernels) MatchAll() []*kernmd.KernelMetadata {
	return kernmd.ToKernelMetadataAll(*l)
}
