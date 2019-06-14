package vmmd

import (
	"fmt"
	"io/ioutil"
	"math"
	"net"
	"os"
	"path"
	"path/filepath"

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

func (md *VMMetadata) AllocateOverlay(requestedSize uint64) error {
	// Truncate only accepts an int64
	if requestedSize > math.MaxInt64 {
		return fmt.Errorf("requested size %d too large, cannot truncate", requestedSize)
	}

	size := int64(requestedSize)

	fi, err := os.Stat(path.Join(constants.IMAGE_DIR, md.VMOD().ImageID.String(), constants.IMAGE_FS))
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

	metadataFile, err := os.Create(path.Join(md.ObjectPath(), constants.VM_METADATA_FILE))
	if err != nil {
		return fmt.Errorf("failed to create metadata file for %q, %v", md.ID, err)
	}
	defer metadataFile.Close()

	dataFile, err := os.Create(path.Join(md.ObjectPath(), constants.VM_DATA_FILE))
	if err != nil {
		return fmt.Errorf("failed to create data file for %q, %v", md.ID, err)
	}
	defer dataFile.Close()

	// TODO: Calculate the correct size
	if err := metadataFile.Truncate(1024 * 1024 * 2); err != nil {
		return fmt.Errorf("failed to allocate data file for VM %q: %v", md.ID, err)
	}

	if err := dataFile.Truncate(size); err != nil {
		return fmt.Errorf("failed to allocate data file for VM %q: %v", md.ID, err)
	}

	return nil
}

func (md *VMMetadata) CopyToOverlay(fileMappings map[string]string) error {
	// Skip the mounting process if there are no files
	if len(fileMappings) == 0 {
		return nil
	}

	if err := md.NewVMOverlay(); err != nil {
		return err
	}
	defer md.RemoveOverlay()

	mp, err := util.Mount(md.OverlayDev())
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
	if len(md.VMOD().IPAddrs) > 0 {
		ip = md.VMOD().IPAddrs[0]
	}
	return md.WriteEtcHosts(mp.Path, md.ID.String(), ip)
}

// WriteEtcHosts populates the /etc/hosts file to avoid errors like
// sudo: unable to resolve host 4462576f8bf5b689
func (md *VMMetadata) WriteEtcHosts(tmpDir, hostname string, primaryIP net.IP) error {
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
	return md.VMOD().KernelID.String()
}

func (md *VMMetadata) Size() (int64, error) {
	fi, err := os.Stat(path.Join(md.ObjectPath(), constants.VM_DATA_FILE))
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

func (md *VMMetadata) ClearPortMappings() {
	md.VMOD().PortMappings = nil
}

// Generate a new SSH keypair for the vm
func (md *VMMetadata) NewSSHKeypair() (string, error) {
	privKeyPath := path.Join(md.ObjectPath(), fmt.Sprintf(constants.VM_SSH_KEY_TEMPLATE, md.ID))

	// Use ED25519 instead of RSA for performance (it's equally secure, but a lot faster to generate/authenticate)
	_, err := util.ExecuteCommand("ssh-keygen", "-q", "-t", "ed25519", "-N", "", "-f", privKeyPath)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s.pub", privKeyPath), nil
}
