package kernmd

import (
	"github.com/weaveworks/ignite/pkg/metadata"
)

type KernelMetadata struct {
	*metadata.Metadata
}

type KernelObjectData struct {
	// TODO: Placeholder
}

func NewKernelMetadata(id *metadata.ID, name *metadata.Name) (*KernelMetadata, error) {
	md, err := metadata.NewMetadata(id, name, metadata.Kernel, &KernelObjectData{})
	if err != nil {
		return nil, err
	}

	return &KernelMetadata{Metadata: md}, nil
}

// The md.ObjectData.(*KernelObjectData) assert won't panic as this method can only receive *KernelMetadata objects
func (md *KernelMetadata) KernelOD() *KernelObjectData {
	return md.ObjectData.(*KernelObjectData)
}
