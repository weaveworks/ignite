package vmmd

import (
	"fmt"
	"github.com/luxas/ignite/pkg/constants"
	"github.com/luxas/ignite/pkg/util"
	"net"
	"os"
	"path"
)

//func (md *VMMetadata) CopyImage() error {
//	od := md.VMOD()
//
//	if err := util.CopyFile(path.Join(constants.IMAGE_DIR, od.ImageID, constants.IMAGE_FS),
//		path.Join(md.ObjectPath(), constants.IMAGE_FS)); err != nil {
//		return fmt.Errorf("failed to copy image %q to VM %q: %v", od.ImageID, md.ID, err)
//	}
//
//	return nil
//}

func (md *VMMetadata) AllocateOverlay() error {
	overlayFile, err := os.Create(path.Join(md.ObjectPath(), constants.OVERLAY_FILE))
	if err != nil {
		return fmt.Errorf("failed to create overlay file for %q, %v", md.ID, err)
	}
	defer overlayFile.Close()

	// TODO: Dynamic size, for now hardcoded 4 GiB
	if err := overlayFile.Truncate(4294967296); err != nil {
		return fmt.Errorf("failed to allocate overlay file for VM %q: %v", md.ID, err)
	}

	return nil
}

func (md *VMMetadata) CopyToOverlay(fileMappings map[string]string) error {
	// Skip the mounting process if there are no files
	if len(fileMappings) == 0 {
		return nil
	}

	if err := md.SetupSnapshot(); err != nil {
		return err
	}
	defer md.RemoveSnapshot()

	mp, err := util.Mount(md.SnapshotDev())
	if err != nil {
		return err
	}
	defer mp.Umount()

	// TODO: File/directory permissions?
	for hostFile, vmFile := range fileMappings {
		vmFilePath := path.Join(mp.Path, vmFile)
		if err := os.MkdirAll(path.Dir(vmFilePath), 0755); err != nil {
			return err
		}

		if err := util.CopyFile(hostFile, vmFilePath); err != nil {
			return err
		}
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
	fi, err := os.Stat(path.Join(md.ObjectPath(), constants.OVERLAY_FILE))
	if err != nil {
		return 0, err
	}

	return fi.Size(), nil
}

func (md *VMMetadata) AddIPAddress(address net.IP) {
	od := md.VMOD()
	od.IPAddrs = append(od.IPAddrs, address)
}

func (md *VMMetadata) ClearIPAddresses() {
	md.VMOD().IPAddrs = nil
}
