package imgmd

import (
	"os"
	"path/filepath"

	"github.com/weaveworks/ignite/pkg/source"
	"github.com/weaveworks/ignite/pkg/util"
)

// AddFiles copies the contents of the tar file a new pool volume
func (md *ImageMetadata) AddFiles(src source.Source) error {
	mountPoint, err := md.newImageVolume(src)
	if err != nil {
		return err
	}
	defer mountPoint.Umount()

	return md.setupResolvConf(mountPoint.Path)
}

// setupResolvConf makes sure there is a resolv.conf file, otherwise
// name resolution won't work. The kernel uses DHCP by default, and
// puts the nameservers in /proc/net/pnp at runtime. Hence, as a default,
// if /etc/resolv.conf doesn't exist, we can use /proc/net/pnp as /etc/resolv.conf
func (md *ImageMetadata) setupResolvConf(dir string) error {
	resolvConf := filepath.Join(dir, "/etc/resolv.conf")
	empty, err := util.FileIsEmpty(resolvConf)
	if err != nil {
		return err
	}

	if !empty {
		return nil
	}

	//fmt.Println("Symlinking /etc/resolv.conf to /proc/net/pnp")
	return os.Symlink("../proc/net/pnp", resolvConf)
}

// TODO: This
func (md *ImageMetadata) Size() (int64, error) {
	return 0, nil
}
