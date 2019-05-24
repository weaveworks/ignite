package vmmd

import (
	"fmt"
	"github.com/luxas/ignite/pkg/constants"
	"github.com/luxas/ignite/pkg/util"
	"os"
	"path"
)

func (md *VMMetadata) CopyImage() error {
	od := md.VMOD()

	if err := util.CopyFile(path.Join(constants.IMAGE_DIR, od.ImageID, constants.IMAGE_FS),
		path.Join(md.ObjectPath(), constants.IMAGE_FS)); err != nil {
		return fmt.Errorf("failed to copy image %q to VM %q: %v", od.ImageID, md.ID, err)
	}

	return nil
}

func (md *VMMetadata) SetState(s state) error {
	md.VMOD().State = s

	if err := md.Save(); err != nil {
		return err
	}

	return nil
}

func (md *VMMetadata) Running() bool {
	return md.VMOD().State == Running
}

func (md *VMMetadata) KernelID() string {
	return md.VMOD().KernelID
}

func (md *VMMetadata) Size() (int64, error) {
	fi, err := os.Stat(path.Join(md.ObjectPath(), constants.IMAGE_FS))
	if err != nil {
		return 0, err
	}

	return fi.Size(), nil
}
