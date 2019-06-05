package vmmd

import (
	"io/ioutil"
	"path/filepath"
	"fmt"
	"github.com/luxas/ignite/pkg/constants"
	"github.com/luxas/ignite/pkg/util"
	"math"
	"net"
	"os"
	"path"
)

const (
	hostsFileTmpl = `127.0.0.1	localhost
127.0.1.1	%s
# The following lines are desirable for IPv6 capable hosts
::1     ip6-localhost ip6-loopback
fe00::0 ip6-localnet
ff00::0 ip6-mcastprefix
ff02::1 ip6-allnodes
ff02::2 ip6-allrouters
`
)

func (md *VMMetadata) AllocateOverlay(size uint64) error {
	// Truncate only accepts an int64
	if size > math.MaxInt64 {
		return fmt.Errorf("requested size %d too large, cannot truncate", size)
	}

	overlayFile, err := os.Create(path.Join(md.ObjectPath(), constants.OVERLAY_FILE))
	if err != nil {
		return fmt.Errorf("failed to create overlay file for %q, %v", md.ID, err)
	}
	defer overlayFile.Close()

	if err := overlayFile.Truncate(int64(size)); err != nil {
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

	return md.WriteEtcHosts(mp.Path, md.ID)
}

// WriteEtcHosts populates the /etc/hosts file to avoid errors like
// sudo: unable to resolve host 4462576f8bf5b689
func (md *VMMetadata) WriteEtcHosts(tmpDir, hostname string) error {
	hostFilePath := filepath.Join(tmpDir, "/etc/hosts")
	empty, err := util.FileIsEmpty(hostFilePath)
	if err != nil {
		return err
	}
	if !empty {
		return nil
	}
	content := []byte(fmt.Sprintf(hostsFileTmpl, hostname))
	return ioutil.WriteFile(hostFilePath, content, 0644)
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
