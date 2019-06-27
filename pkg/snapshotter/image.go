package snapshotter

import (
	"github.com/weaveworks/ignite/pkg/object/image"
	"github.com/weaveworks/ignite/pkg/util"
)

func (s *Snapshotter) ImportImage(image image.Image) (*util.MountPoint, error) {
	volume, err := s.CreateVolume(image.Spec.Source.Size.Add(extraSize), image.MetadataPath())
	if err != nil {
		return nil, err
	}

	// TODO: Source loaders which parse v1alpha1.Image.Spec.Source and output a Source object for Import()
	mountPoint, err := volume.Import(src)
	if err != nil {
		return nil, err
	}

	return mountPoint, nil
}
