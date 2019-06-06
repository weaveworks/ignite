package cmdutil

import (
	"github.com/weaveworks/ignite/pkg/filter"
	"github.com/weaveworks/ignite/pkg/metadata"
	"github.com/weaveworks/ignite/pkg/metadata/imgmd"
	"github.com/weaveworks/ignite/pkg/metadata/kernmd"
	"github.com/weaveworks/ignite/pkg/metadata/vmmd"
)

// Utility functions for metadata loading/filtering
// TODO: Make these take in pre-loaded metadata, so that it doesn't get loaded twice if calling multiple of these

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

func matchAllNames(sources []metadata.AnyMetadata) []*metadata.Name {
	names := make([]*metadata.Name, 0, len(sources))

	for _, md := range sources {
		names = append(names, &md.GetMD().Name)
	}

	return names
}

// Match a single VM using the VMFilter
func MatchSingleVM(match string) (*vmmd.VMMetadata, error) {
	mds, err := vmmd.LoadAllVMMetadata()
	if err != nil {
		return nil, err
	}

	md, err := single(vmmd.NewVMFilter(match), mds)
	if err != nil {
		return nil, err
	}

	return vmmd.ToVMMetadata(md), nil
}

// Match multiple individual VMs with different filter strings
func MatchSingleVMs(matches []string) ([]*vmmd.VMMetadata, error) {
	mds, err := vmmd.LoadAllVMMetadata()
	if err != nil {
		return nil, err
	}

	filters := make([]metadata.Filter, 0, len(matches))
	for _, match := range matches {
		filters = append(filters, vmmd.NewVMFilter(match))
	}

	results, err := matchIndividual(filters, mds)
	if err != nil {
		return nil, err
	}

	return vmmd.ToVMMetadataAll(results), nil
}

// Match a single image using the IDNameFilter
func MatchSingleImage(match string) (*imgmd.ImageMetadata, error) {
	mds, err := imgmd.LoadAllImageMetadata()
	if err != nil {
		return nil, err
	}

	md, err := single(metadata.NewIDNameFilter(match, metadata.Image), mds)
	if err != nil {
		return nil, err
	}

	return imgmd.ToImageMetadata(md), nil
}

// Match multiple individual images with different filter strings
func MatchSingleImages(matches []string) ([]*imgmd.ImageMetadata, error) {
	mds, err := imgmd.LoadAllImageMetadata()
	if err != nil {
		return nil, err
	}

	filters := make([]metadata.Filter, 0, len(matches))
	for _, match := range matches {
		filters = append(filters, metadata.NewIDNameFilter(match, metadata.Image))
	}

	results, err := matchIndividual(filters, mds)
	if err != nil {
		return nil, err
	}

	return imgmd.ToImageMetadataAll(results), nil
}

// Match a single kernel using the IDNameFilter
func MatchSingleKernel(match string) (*kernmd.KernelMetadata, error) {
	mds, err := kernmd.LoadAllKernelMetadata()
	if err != nil {
		return nil, err
	}

	md, err := single(metadata.NewIDNameFilter(match, metadata.Kernel), mds)
	if err != nil {
		return nil, err
	}

	return kernmd.ToKernelMetadata(md), nil
}

// Match multiple individual kernels with different filter strings
func MatchSingleKernels(matches []string) ([]*kernmd.KernelMetadata, error) {
	mds, err := kernmd.LoadAllKernelMetadata()
	if err != nil {
		return nil, err
	}

	filters := make([]metadata.Filter, 0, len(matches))
	for _, match := range matches {
		filters = append(filters, metadata.NewIDNameFilter(match, metadata.Kernel))
	}

	results, err := matchIndividual(filters, mds)
	if err != nil {
		return nil, err
	}

	return kernmd.ToKernelMetadataAll(results), nil
}

func MatchAllVMs(all bool) ([]*vmmd.VMMetadata, error) {
	mds, err := vmmd.LoadAllVMMetadata()
	if err != nil {
		return nil, err
	}

	matches, err := filter.NewFilterer(vmmd.NewVMFilterAll("", all), mds)
	if err != nil {
		return nil, err
	}

	return vmmd.ToVMMetadataAll(matches.All()), nil
}

func MatchAllImages() ([]*imgmd.ImageMetadata, error) {
	mds, err := imgmd.LoadAllImageMetadata()
	if err != nil {
		return nil, err
	}

	return imgmd.ToImageMetadataAll(mds), nil
}

func MatchAllKernels() ([]*kernmd.KernelMetadata, error) {
	mds, err := kernmd.LoadAllKernelMetadata()
	if err != nil {
		return nil, err
	}

	return kernmd.ToKernelMetadataAll(mds), nil
}

func MatchAllVMNames() ([]*metadata.Name, error) {
	mds, err := vmmd.LoadAllVMMetadata()
	if err != nil {
		return nil, err
	}

	return matchAllNames(mds), nil
}

func MatchAllImageNames() ([]*metadata.Name, error) {
	mds, err := imgmd.LoadAllImageMetadata()
	if err != nil {
		return nil, err
	}

	return matchAllNames(mds), nil
}

func MatchAllKernelNames() ([]*metadata.Name, error) {
	mds, err := kernmd.LoadAllKernelMetadata()
	if err != nil {
		return nil, err
	}

	return matchAllNames(mds), nil
}
