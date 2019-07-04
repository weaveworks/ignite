package loader

import (
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/filter"
	"github.com/weaveworks/ignite/pkg/metadata"
	"github.com/weaveworks/ignite/pkg/metadata/imgmd"
	"github.com/weaveworks/ignite/pkg/metadata/kernmd"
	"github.com/weaveworks/ignite/pkg/metadata/vmmd"
)

type allVMs []metadata.Metadata
type allImages []metadata.Metadata
type allKernels []metadata.Metadata

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
		if l.vm, err = vmmd.LoadAllVM(); err != nil {
			return nil, err
		}
	}

	return &l.vm, nil
}

func (l *ResLoader) Images() (*allImages, error) {
	if l.image == nil {
		var err error
		if l.image, err = imgmd.LoadAllImage(); err != nil {
			return nil, err
		}
	}

	return &l.image, nil
}

func (l *ResLoader) Kernels() (*allKernels, error) {
	if l.kernel == nil {
		var err error
		if l.kernel, err = kernmd.LoadAllKernel(); err != nil {
			return nil, err
		}
	}

	return &l.kernel, nil
}

func single(f metadata.Filter, sources []metadata.Metadata) (metadata.Metadata, error) {
	var result metadata.Metadata

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

func matchIndividual(fs []metadata.Filter, sources []metadata.Metadata) ([]metadata.Metadata, error) {
	results := make([]metadata.Metadata, 0, len(sources))

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
func (l *allVMs) MatchSingle(match string) (*vmmd.VM, error) {
	md, err := single(vmmd.NewVMFilter(match), *l)
	if err != nil {
		return nil, err
	}

	return vmmd.ToVM(md), nil
}

// Match multiple individual VMs with different filter strings
func (l *allVMs) MatchMultiple(matches []string) ([]*vmmd.VM, error) {
	filters := make([]metadata.Filter, 0, len(matches))
	for _, match := range matches {
		filters = append(filters, vmmd.NewVMFilter(match))
	}

	results, err := matchIndividual(filters, *l)
	if err != nil {
		return nil, err
	}

	return vmmd.ToVMAll(results), nil
}

func (l *allVMs) MatchFilter(all bool) ([]*vmmd.VM, error) {
	matches, err := filter.NewFilterer(vmmd.NewVMFilterAll("", all), *l)
	if err != nil {
		return nil, err
	}

	return vmmd.ToVMAll(matches.All()), nil
}

func (l *allVMs) MatchAll() []*vmmd.VM {
	return vmmd.ToVMAll(*l)
}

// Match a single image using the IDNameFilter
func (l *allImages) MatchSingle(match string) (*imgmd.Image, error) {
	md, err := single(metadata.NewIDNameFilter(match, meta.KindImage), *l)
	if err != nil {
		return nil, err
	}

	return imgmd.ToImage(md), nil
}

// Match multiple individual images with different filter strings
func (l *allImages) MatchMultiple(matches []string) ([]*imgmd.Image, error) {
	filters := make([]metadata.Filter, 0, len(matches))
	for _, match := range matches {
		filters = append(filters, metadata.NewIDNameFilter(match, meta.KindImage))
	}

	results, err := matchIndividual(filters, *l)
	if err != nil {
		return nil, err
	}

	return imgmd.ToImageAll(results), nil
}

func (l *allImages) MatchAll() []*imgmd.Image {
	return imgmd.ToImageAll(*l)
}

// Match a single kernel using the IDNameFilter
func (l *allKernels) MatchSingle(match string) (*kernmd.Kernel, error) {
	md, err := single(metadata.NewIDNameFilter(match, meta.KindKernel), *l)
	if err != nil {
		return nil, err
	}

	return kernmd.ToKernel(md), nil
}

// Match multiple individual kernels with different filter strings
func (l *allKernels) MatchMultiple(matches []string) ([]*kernmd.Kernel, error) {
	filters := make([]metadata.Filter, 0, len(matches))
	for _, match := range matches {
		filters = append(filters, metadata.NewIDNameFilter(match, meta.KindKernel))
	}

	results, err := matchIndividual(filters, *l)
	if err != nil {
		return nil, err
	}

	return kernmd.ToKernelAll(results), nil
}

func (l *allKernels) MatchAll() []*kernmd.Kernel {
	return kernmd.ToKernelAll(*l)
}
