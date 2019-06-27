package image

import "github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"

// This package represents the image objects, which reside in /var/lib/firecracker/image/{id}/metadata.json

type Image struct {
	v1alpha1.Image
}

// Get the metadata filename for the image
func (i *Image) MetadataPath() string {
	// TODO: This
	return ""
}
