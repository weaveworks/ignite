package vmmd

import (
	"fmt"
	"io/ioutil"
	"math"
	"net"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/c2h5oh/datasize"
	log "github.com/sirupsen/logrus"
	api "github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"
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
	vmAuthorizedKeys = "/root/.ssh/authorized_keys"
)

func (md *VM) AllocateAndPopulateOverlay() error {
	requestedSize := md.Spec.DiskSize.Bytes()
	// Truncate only accepts an int64
	if requestedSize > math.MaxInt64 {
		return fmt.Errorf("requested size %d too large, cannot truncate", requestedSize)
	}

	size := int64(requestedSize)

	fi, err := os.Stat(path.Join(constants.IMAGE_DIR, md.GetImageUID().String(), constants.IMAGE_FS))
	if err != nil {
		return err
	}

	// The overlay needs to be at least as large as the image
	if size < fi.Size() {
		size = fi.Size()
		log.Warnf("warning: requested overlay size (%s) < image size (%s), using image size for overlay\n",
			datasize.ByteSize(requestedSize).HR(), datasize.ByteSize(size).HR())
	}

	// Make sure the all directories above the snapshot directory exists
	if err := os.MkdirAll(path.Dir(md.OverlayFile()), constants.DATA_DIR_PERM); err != nil {
		return err
	}

	overlayFile, err := os.Create(md.OverlayFile())
	if err != nil {
		return fmt.Errorf("failed to create overlay file for %q, %v", md.GetUID(), err)
	}
	defer overlayFile.Close()

	if err := overlayFile.Truncate(size); err != nil {
		return fmt.Errorf("failed to allocate overlay file for VM %q: %v", md.GetUID(), err)
	}

	// populate the filesystem
	return md.copyToOverlay()
}

func (md *VM) copyToOverlay() error {
	if err := md.SetupSnapshot(); err != nil {
		return err
	}
	defer md.RemoveSnapshot()

	mp, err := util.Mount(md.SnapshotDev())
	if err != nil {
		return err
	}
	defer mp.Umount()

	// Copy the kernel files to the VM. TODO: Use snapshot overlaying instead.
	if err := md.copyKernelToOverlay(mp.Path); err != nil {
		return err
	}

	// do not mutate md.Spec.CopyFiles
	fileMappings := md.Spec.CopyFiles

	if md.Spec.SSH != nil {
		pubKeyPath := md.Spec.SSH.PublicKey
		if md.Spec.SSH.Generate {
			// generate a key if PublicKey is empty
			pubKeyPath, err = md.newSSHKeypair()
			if err != nil {
				return err
			}
		}

		if len(pubKeyPath) > 0 {
			fileMappings = append(fileMappings, api.FileMapping{
				HostPath: pubKeyPath,
				VMPath:   vmAuthorizedKeys,
			})
		}
	}

	// TODO: File/directory permissions?
	for _, mapping := range fileMappings {
		vmFilePath := path.Join(mp.Path, mapping.VMPath)
		if err := os.MkdirAll(path.Dir(vmFilePath), constants.DATA_DIR_PERM); err != nil {
			return err
		}

		if err := util.CopyFile(mapping.HostPath, vmFilePath); err != nil {
			return err
		}
	}

	ip := net.IP{127, 0, 0, 1}
	if len(md.Status.IPAddresses) > 0 {
		ip = md.Status.IPAddresses[0]
	}

	if err := md.writeEtcHosts(mp.Path, md.GetUID().String(), ip); err != nil {
		return err
	}

	// TODO: This code seems to be flaky and not always copy over the files?
	time.Sleep(500 * time.Millisecond)
	return nil
}

func (md *VM) copyKernelToOverlay(mountPoint string) error {
	kernelUID := md.GetKernelUID()
	kernelTarPath := path.Join(constants.KERNEL_DIR, kernelUID.String(), constants.KERNEL_TAR)

	if !util.FileExists(kernelTarPath) {
		log.Warnf("Could not find kernel overlay files, not copying into the VM.")
		return nil
	}

	_, err := util.ExecuteCommand("tar", "-xf", kernelTarPath, "-C", mountPoint)
	return err
}

// writeEtcHosts populates the /etc/hosts file to avoid errors like
// sudo: unable to resolve host 4462576f8bf5b689
func (md *VM) writeEtcHosts(tmpDir, hostname string, primaryIP net.IP) error {
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

func (md *VM) OverlayFile() string {
	return path.Join(md.ObjectPath(), constants.OVERLAY_FILE)
}

func (md *VM) AddIPAddress(address net.IP) {
	md.Status.IPAddresses = append(md.Status.IPAddresses, address)
}

func (md *VM) ClearIPAddresses() {
	md.Status.IPAddresses = nil
}

// Generate a new SSH keypair for the vm
func (md *VM) newSSHKeypair() (string, error) {
	privKeyPath := path.Join(md.ObjectPath(), fmt.Sprintf(constants.VM_SSH_KEY_TEMPLATE, md.GetUID()))

	// Use ED25519 instead of RSA for performance (it's equally secure, but a lot faster to generate/authenticate)
	_, err := util.ExecuteCommand("ssh-keygen", "-q", "-t", "ed25519", "-N", "", "-f", privKeyPath)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s.pub", privKeyPath), nil
}
