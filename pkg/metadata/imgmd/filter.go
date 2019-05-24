package imgmd

import (
	"fmt"
	"github.com/luxas/ignite/pkg/filter"
	"github.com/luxas/ignite/pkg/metadata"
	"strings"
)

// Compile-time assert to verify interface compatibility
var _ filter.Filter = &ImageFilter{}

type ImageFilter struct {
	prefix string
}

func NewImageFilter(p string) *ImageFilter {
	return &ImageFilter{
		prefix: p,
	}
}

func (n *ImageFilter) Filter(f filter.Filterable) (bool, error) {
	md, ok := f.(*ImageMetadata)
	if !ok {
		return false, fmt.Errorf("failed to assert Filterable %v to ImageMetadata", f)
	}

	return strings.HasPrefix(md.ID, n.prefix) || strings.HasPrefix(md.Name, n.prefix), nil
}

func LoadImageMetadata(id string) (filter.Filterable, error) {
	md := &ImageMetadata{
		Metadata: &metadata.Metadata{
			ID:         id,
			Type:       metadata.Image,
			ObjectData: &ImageObjectData{},
		},
	}

	err := md.Load()
	return md, err
}
