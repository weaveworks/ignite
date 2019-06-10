package runutil

import (
	"github.com/weaveworks/ignite/pkg/filter"
	"github.com/weaveworks/ignite/pkg/metadata"
	"github.com/weaveworks/ignite/pkg/metadata/imgmd"
	"github.com/weaveworks/ignite/pkg/metadata/kernmd"
	"github.com/weaveworks/ignite/pkg/metadata/vmmd"
)

type ResourceLoader struct {
	vms     []metadata.AnyMetadata
	images  []metadata.AnyMetadata
	kernels []metadata.AnyMetadata
}

func NewResLoader() *ResourceLoader {
	return &ResourceLoader{}
}

func (r *ResourceLoader) loadVMs() error {
	// Don't load twice
	if r.vms != nil {
		return nil
	}

	var err error
	r.vms, err = vmmd.LoadAllVMMetadata()
	if err != nil {
		return err
	}

	return nil
}

func (r *ResourceLoader) loadImages() error {
	// Don't load twice
	if r.images != nil {
		return nil
	}

	var err error
	r.images, err = imgmd.LoadAllImageMetadata()
	if err != nil {
		return err
	}

	return nil
}

func (r *ResourceLoader) loadKernels() error {
	// Don't load twice
	if r.kernels != nil {
		return nil
	}

	var err error
	r.kernels, err = kernmd.LoadAllKernelMetadata()
	if err != nil {
		return err
	}

	return nil
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
func (r *ResourceLoader) MatchSingleVM(match string) (*vmmd.VMMetadata, error) {
	if err := r.loadVMs(); err != nil {
		return nil, err
	}

	md, err := single(vmmd.NewVMFilter(match), r.vms)
	if err != nil {
		return nil, err
	}

	return vmmd.ToVMMetadata(md), nil
}

// Match multiple individual VMs with different filter strings
func (r *ResourceLoader) MatchSingleVMs(matches []string) ([]*vmmd.VMMetadata, error) {
	if err := r.loadVMs(); err != nil {
		return nil, err
	}

	filters := make([]metadata.Filter, 0, len(matches))
	for _, match := range matches {
		filters = append(filters, vmmd.NewVMFilter(match))
	}

	results, err := matchIndividual(filters, r.vms)
	if err != nil {
		return nil, err
	}

	return vmmd.ToVMMetadataAll(results), nil
}

// Match a single image using the IDNameFilter
func (r *ResourceLoader) MatchSingleImage(match string) (*imgmd.ImageMetadata, error) {
	if err := r.loadImages(); err != nil {
		return nil, err
	}

	md, err := single(metadata.NewIDNameFilter(match, metadata.Image), r.images)
	if err != nil {
		return nil, err
	}

	return imgmd.ToImageMetadata(md), nil
}

// Match multiple individual images with different filter strings
func (r *ResourceLoader) MatchSingleImages(matches []string) ([]*imgmd.ImageMetadata, error) {
	if err := r.loadImages(); err != nil {
		return nil, err
	}

	filters := make([]metadata.Filter, 0, len(matches))
	for _, match := range matches {
		filters = append(filters, metadata.NewIDNameFilter(match, metadata.Image))
	}

	results, err := matchIndividual(filters, r.images)
	if err != nil {
		return nil, err
	}

	return imgmd.ToImageMetadataAll(results), nil
}

// Match a single kernel using the IDNameFilter
func (r *ResourceLoader) MatchSingleKernel(match string) (*kernmd.KernelMetadata, error) {
	if err := r.loadKernels(); err != nil {
		return nil, err
	}

	md, err := single(metadata.NewIDNameFilter(match, metadata.Kernel), r.kernels)
	if err != nil {
		return nil, err
	}

	return kernmd.ToKernelMetadata(md), nil
}

// Match multiple individual kernels with different filter strings
func (r *ResourceLoader) MatchSingleKernels(matches []string) ([]*kernmd.KernelMetadata, error) {
	if err := r.loadKernels(); err != nil {
		return nil, err
	}

	filters := make([]metadata.Filter, 0, len(matches))
	for _, match := range matches {
		filters = append(filters, metadata.NewIDNameFilter(match, metadata.Kernel))
	}

	results, err := matchIndividual(filters, r.kernels)
	if err != nil {
		return nil, err
	}

	return kernmd.ToKernelMetadataAll(results), nil
}

func (r *ResourceLoader) MatchAllVMs(all bool) ([]metadata.AnyMetadata, error) {
	if err := r.loadVMs(); err != nil {
		return nil, err
	}

	matches, err := filter.NewFilterer(vmmd.NewVMFilterAll("", all), r.vms)
	if err != nil {
		return nil, err
	}

	return matches.All(), nil
}

func (r *ResourceLoader) MatchAllImages() ([]metadata.AnyMetadata, error) {
	if err := r.loadImages(); err != nil {
		return nil, err
	}

	return r.images, nil
}

func (r *ResourceLoader) MatchAllKernels() ([]metadata.AnyMetadata, error) {
	if err := r.loadKernels(); err != nil {
		return nil, err
	}

	return r.kernels, nil
}