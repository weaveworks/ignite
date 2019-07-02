package imgmd

import "github.com/weaveworks/ignite/pkg/dm"

type KernelMD struct {
	Version string
}

var _ dm.Metadata = KernelMD{}

func NewKernelMD() *KernelMD {
	return &KernelMD{}
}
