package snapshotter

import (
	"github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"
)

type ErrAmbiguous error
type ErrNonexistent error

// The Filter interface that external filters implement
type Filter interface {
	Filter(*Object) (*Object, error) // Filter returns the *Object if it matches, otherwise nil
	SetType(v1alpha1.PoolDeviceType) // Let the filter now the type of the Objects it's receiving
	ErrAmbiguous() ErrAmbiguous      // Handling of ambiguous queries
	ErrNonexistent() ErrNonexistent  // Handling of nonexistent queries
}

func (s *Snapshotter) getSingle(f Filter, t v1alpha1.PoolDeviceType) (*Object, error) {
	var result *Object

	f.SetType(t)
	for _, object := range s.objects {
		if object.device.Type == t {
			if match, err := f.Filter(object); err != nil { // Filter returns *Object if it matches, otherwise nil
				return nil, err
			} else if match != nil {
				if result != nil {
					return nil, f.ErrAmbiguous()
				} else {
					result = match
				}
			}
		}
	}

	if result == nil {
		return nil, f.ErrNonexistent()
	}

	return result, nil
}

func (s *Snapshotter) getMultiple(f Filter, t v1alpha1.PoolDeviceType) ([]*Object, error) {
	var results []*Object

	f.SetType(t)
	for _, object := range s.objects {
		if object.device().Type == t {
			if match, err := f.Filter(object); err != nil { // Filter returns *Object if it matches, otherwise nil
				return nil, err
			} else if match != nil {
				results = append(results, match)
			}
		}
	}

	return results, nil
}

func (s *Snapshotter) GetImage(f Filter) (*Image, error) {
	result, err := s.getSingle(f, v1alpha1.PoolDeviceTypeImage)
	if err != nil {
		return nil, err
	}

	return newImage(result)
}

func (s *Snapshotter) GetImages(f Filter) ([]*Image, error) {
	results, err := s.getMultiple(f, v1alpha1.PoolDeviceTypeImage)
	if err != nil {
		return nil, err
	}

	images := make([]*Image, 0, len(results))
	for _, result := range results {
		if image, err := newImage(result); err == nil {
			images = append(images, image)
		} else {
			return nil, err
		}
	}

	return images, nil
}

func (s *Snapshotter) GetKernel(f Filter) (*Kernel, error) {
	result, err := s.getSingle(f, v1alpha1.PoolDeviceTypeKernel)
	if err != nil {
		return nil, err
	}

	return newKernel(result)
}

func (s *Snapshotter) GetKernels(f Filter) ([]*Kernel, error) {
	results, err := s.getMultiple(f, v1alpha1.PoolDeviceTypeKernel)
	if err != nil {
		return nil, err
	}

	kernels := make([]*Kernel, 0, len(results))
	for _, result := range results {
		if kernel, err := newKernel(result); err == nil {
			kernels = append(kernels, kernel)
		} else {
			return nil, err
		}
	}

	return kernels, nil
}

func (s *Snapshotter) GetVM(f Filter) (*VM, error) {
	result, err := s.getSingle(f, v1alpha1.PoolDeviceTypeVM)
	if err != nil {
		return nil, err
	}

	return newVM(result)
}

func (s *Snapshotter) GetVMs(f Filter) ([]*VM, error) {
	results, err := s.getMultiple(f, v1alpha1.PoolDeviceTypeVM)
	if err != nil {
		return nil, err
	}

	vms := make([]*VM, 0, len(results))

	for _, result := range results {
		if vm, err := newVM(result); err == nil {
			vms = append(vms, vm)
		} else {
			return nil, err
		}
	}

	return vms, nil
}

func (s *Snapshotter) GetResize(f Filter) (*Resize, error) {
	result, err := s.getSingle(f, v1alpha1.PoolDeviceTypeResize)
	if err != nil {
		return nil, err
	}

	return newResize(result)
}

func (s *Snapshotter) GetResizes(f Filter) ([]*Resize, error) {
	results, err := s.getMultiple(f, v1alpha1.PoolDeviceTypeResize)
	if err != nil {
		return nil, err
	}

	resizes := make([]*Resize, 0, len(results))

	for _, result := range results {
		if resize, err := newResize(result); err == nil {
			resizes = append(resizes, resize)
		} else {
			return nil, err
		}
	}

	return resizes, nil
}
