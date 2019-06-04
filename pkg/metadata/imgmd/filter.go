package imgmd

import (
	"github.com/luxas/ignite/pkg/filter"
	"github.com/luxas/ignite/pkg/util"
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

func (n *ImageFilter) Filter(f filter.Filterable) ([]string, error) {
	md, err := ToImageMetadata(f)
	if err != nil {
		return nil, err
	}

	return util.MatchPrefix(n.prefix, md.ID, md.Name.String()), nil
}

func LoadImageMetadata(id string) (*ImageMetadata, error) {
	md := NewImageMetadata(id, nil)
	err := md.Load()
	return md, err
}

func LoadImageMetadataFilterable(id string) (filter.Filterable, error) {
	return LoadImageMetadata(id)
}
