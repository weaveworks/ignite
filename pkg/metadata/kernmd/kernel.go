package kernmd

import (
	"fmt"
	"os"
	"path"

	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/util"
)

func (md *KernelMetadata) ImportKernel(p string) error {
	if err := util.CopyFile(p, path.Join(md.ObjectPath(), constants.KERNEL_FILE)); err != nil {
		return fmt.Errorf("failed to copy kernel file %q to kernel %q: %v", p, md.GetUID(), err)
	}

	return nil
}

func (md *KernelMetadata) Size() (int64, error) {
	fi, err := os.Stat(path.Join(md.ObjectPath(), constants.KERNEL_FILE))
	if err != nil {
		return 0, err
	}

	return fi.Size(), nil
}
