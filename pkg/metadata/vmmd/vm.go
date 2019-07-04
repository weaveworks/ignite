package vmmd

import (
	"fmt"
	"io/ioutil"
	"math"
	"net"
	"os"
	"path"
	"path/filepath"

	api "github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"

	"github.com/c2h5oh/datasize"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/util"
)

const (
	hostsFileTmpl = `127.0.0.1	localhost
%s	%s
# The following lines are desirable for IPv6 capable hosts
::1     ip6-localhost ip6-loopback
fe00::0 ip6-localnet
ff00::0 ip6-mcastprefix
ff02::1 ip6-allnodes
ff02::2 ip6-allrouters
`
)

func (md *VM) AllocateOverlay(requestedSize uint64) error {
	// Truncate only accepts an int64
	if requestedSize > math.MaxInt64 {
		return fmt.Errorf("requested size %d too large, cannot truncate", requestedSize)
	}

	size := int64(requestedSize)

	fi, err := os.Stat(path.Join(constants.IMAGE_DIR, md.Spec.Image.UID.String(), constants.IMAGE_FS))
	if err != nil {
		return err
	}

	// The overlay needs to be at least as large as the image
	if size < fi.Size() {
		size = fi.Size()
		// TODO: Logging
		fmt.Printf("warning: requested overlay size (%s) < image size (%s), using image size for overlay\n",
			datasize.ByteSize(requestedSize).HR(), datasize.ByteSize(size).HR())
	}

	overlayFile, err := os.Create(path.Join(md.ObjectPath(), constants.OVERLAY_FILE))
	if err != nil {
		return fmt.Errorf("failed to create overlay file for %q, %v", md.GetUID(), err)
	}
	defer overlayFile.Close()

	if err := overlayFile.Truncate(size); err != nil {
		return fmt.Errorf("failed to allocate overlay file for VM %q: %v", md.GetUID(), err)
	}

	return nil
}

func (md *VM) CopyToOverlay(fileMappings map[string]string) error {
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

	ip := net.IP{127, 0, 0, 1}
	if len(md.Status.IPAddresses) > 0 {
		ip = md.Status.IPAddresses[0]
	}

	return md.WriteEtcHosts(mp.Path, md.GetUID().String(), ip)
}

// WriteEtcHosts populates the /etc/hosts file to avoid errors like
// sudo: unable to resolve host 4462576f8bf5b689
func (md *VM) WriteEtcHosts(tmpDir, hostname string, primaryIP net.IP) error {
	hostFilePath := filepath.Join(tmpDir, "/etc/hosts")
	empty, err := util.FileIsEmpty(hostFilePath)
	if err != nil {
		return err
	}
	if !empty {
		return nil
	}
	content := []byte(fmt.Sprintf(hostsFileTmpl, primaryIP.String(), hostname))
	return ioutil.WriteFile(hostFilePath, content, 0644)
}

func (md *VM) SetState(s api.VMState) error {
	md.Status.State = s

	if err := md.Save(); err != nil {
		return err
	}

	return nil
}

func (md *VM) Running() bool {
	return md.Status.State == api.VMStateRunning
}

func (md *VM) Size() (int64, error) {
	fi, err := os.Stat(path.Join(md.ObjectPath(), constants.OVERLAY_FILE))
	if err != nil {
		return 0, err
	}

	return fi.Size(), nil
}

func (md *VM) AddIPAddress(address net.IP) {
	md.Status.IPAddresses = append(md.Status.IPAddresses, address)
}

func (md *VM) ClearIPAddresses() {
	md.Status.IPAddresses = nil
}

func (md *VM) ClearPortMappings() {
	md.Spec.Ports = nil
}

// Generate a new SSH keypair for the vm
func (md *VM) NewSSHKeypair() (string, error) {
	privKeyPath := path.Join(md.ObjectPath(), fmt.Sprintf(constants.VM_SSH_KEY_TEMPLATE, md.GetUID()))

	// Use ED25519 instead of RSA for performance (it's equally secure, but a lot faster to generate/authenticate)
	_, err := util.ExecuteCommand("ssh-keygen", "-q", "-t", "ed25519", "-N", "", "-f", privKeyPath)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s.pub", privKeyPath), nil
}
