package kernmd

import (
	"os"
	"path"

	"github.com/weaveworks/ignite/pkg/constants"
)

func (md *Kernel) Size() (int64, error) {
	fi, err := os.Stat(path.Join(md.ObjectPath(), constants.KERNEL_FILE))
	if err != nil {
		return 0, err
	}

	return fi.Size(), nil
}
