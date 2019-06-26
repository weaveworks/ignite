package image

import "github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"

type source struct {
	v1alpha1.ImageSource
}

func NewImageSource() *source {
	return &source{}
}
