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

func NewKernelMetadata(id string, name *metadata.Name) *KernelMetadata {
	return &KernelMetadata{
		Metadata: metadata.NewMetadata(id,
			name,
			metadata.Kernel,
			&KernelObjectData{}),
	}
}

// The md.ObjectData.(*KernelObjectData) assert won't panic as this method can only receive *KernelMetadata objects
func (md *KernelMetadata) KernelOD() *KernelObjectData {
	return md.ObjectData.(*KernelObjectData)
}
