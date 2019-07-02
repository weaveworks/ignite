package loader

import (
	"github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"
	"github.com/weaveworks/ignite/pkg/filtering/filter"
	"github.com/weaveworks/ignite/pkg/filtering/filterer"
	"github.com/weaveworks/ignite/pkg/snapshotter"
)

type allVMs struct{ ss *snapshotter.Snapshotter }
type allImages struct{ ss *snapshotter.Snapshotter }
type allKernels struct{ ss *snapshotter.Snapshotter }

type ResLoader struct {
	ss *snapshotter.Snapshotter

	vm     allVMs
	image  allImages
	kernel allKernels
}

func NewResLoader(ss *snapshotter.Snapshotter) *ResLoader {
	return &ResLoader{ss: ss}
}

func (l *ResLoader) VMs() (*allVMs, error) {
	if l.vm.ss == nil {
		// TODO: Partial loading
		if err := l.ss.LoadAll(); err != nil {
			return nil, err
		}

		l.vm.ss = l.ss
	}

	return &l.vm, nil
}

func (l *ResLoader) Images() (*allImages, error) {
	if l.image.ss == nil {
		// TODO: Partial loading
		if err := l.ss.LoadAll(); err != nil {
			return nil, err
		}

		l.image.ss = l.ss
	}

	return &l.image, nil
}

//func (l *ResLoader) Kernels() (*allKernels, error) {
//	if l.kernel == nil {
//		var err error
//		if l.kernel, err = kernmd.LoadAllKernelMetadata(); err != nil {
//			return nil, err
//		}
//	}
//
//	return &l.kernel, nil
//}

func single(f filterer.Filter, sources []v1alpha1.Object) (v1alpha1.Object, error) {
	var result v1alpha1.Object

	// Match a single AnyMetadata using the given filter
	if matches, err := filterer.NewFilterer(f, sources); err == nil {
		if result, err = matches.Single(); err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}

	return result, nil
}

func matchIndividual(fs []filterer.Filter, sources []v1alpha1.Object) ([]v1alpha1.Object, error) {
	results := make([]v1alpha1.Object, 0, len(sources))

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
func (l *allVMs) MatchSingle(match string) (*snapshotter.VM, error) {
	obj, err := single(filter.NewVMFilter(match), l.ss.VMObjects())
	if err != nil {
		return nil, err
	}

	return snapshotter.ObjectToVM(obj), nil
}

// Match multiple individual VMs with different filter strings
func (l *allVMs) MatchMultiple(matches []string) ([]*snapshotter.VM, error) {
	filters := make([]filterer.Filter, 0, len(matches))
	for _, match := range matches {
		filters = append(filters, filter.NewVMFilter(match))
	}

	results, err := matchIndividual(filters, l.ss.VMObjects())
	if err != nil {
		return nil, err
	}

	return snapshotter.ObjectsToVMs(results), nil
}

func (l *allVMs) MatchFilter(all bool) ([]*snapshotter.VM, error) {
	matches, err := filterer.NewFilterer(filter.NewVMFilterAll("", all), l.ss.VMObjects())
	if err != nil {
		return nil, err
	}

	return snapshotter.ObjectsToVMs(matches.All()), nil
}

func (l *allVMs) MatchAll() []*snapshotter.VM {
	return l.ss.VMs
}

// Match a single image using the IDNameFilter
func (l *allImages) MatchSingle(match string) (*snapshotter.Image, error) {
	obj, err := single(filter.NewIDNameFilter(match, v1alpha1.PoolDeviceTypeImage), l.ss.ImageObjects())
	if err != nil {
		return nil, err
	}

	return snapshotter.ObjectToImage(obj), nil
}

// Match multiple individual images with different filter strings
func (l *allImages) MatchMultiple(matches []string) ([]*snapshotter.Image, error) {
	filters := make([]filterer.Filter, 0, len(matches))
	for _, match := range matches {
		filters = append(filters, filter.NewIDNameFilter(match, v1alpha1.PoolDeviceTypeImage))
	}

	results, err := matchIndividual(filters, l.ss.ImageObjects())
	if err != nil {
		return nil, err
	}

	return snapshotter.ObjectsToImages(results), nil
}

func (l *allImages) MatchAll() []*snapshotter.Image {
	return l.ss.Images
}

// Match a single kernel using the IDNameFilter
//func (l *allKernels) MatchSingle(match string) (*kernmd.KernelMetadata, error) {
//	md, err := single(metadata.NewIDNameFilter(match, metadata.Kernel), *l)
//	if err != nil {
//		return nil, err
//	}
//
//	return kernmd.ToKernelMetadata(md), nil
//}
//
//// Match multiple individual kernels with different filter strings
//func (l *allKernels) MatchMultiple(matches []string) ([]*kernmd.KernelMetadata, error) {
//	filters := make([]metadata.Filter, 0, len(matches))
//	for _, match := range matches {
//		filters = append(filters, metadata.NewIDNameFilter(match, metadata.Kernel))
//	}
//
//	results, err := matchIndividual(filters, *l)
//	if err != nil {
//		return nil, err
//	}
//
//	return kernmd.ToKernelMetadataAll(results), nil
//}
//
//func (l *allKernels) MatchAll() []*kernmd.KernelMetadata {
//	return kernmd.ToKernelMetadataAll(*l)
//}
