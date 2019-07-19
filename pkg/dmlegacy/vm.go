package dmlegacy

import (
	"fmt"
	"io/ioutil"
	"math"
	"net"
	"os"
	"path"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/metadata/vmmd"
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

// AllocateAndPopulateOverlay creates the overlay.dm file on top of an image, and
// configures the snapshot in all ways needed. It also copies in contents from the
// host as needed, and configures networking.
func AllocateAndPopulateOverlay(vm *vmmd.VM) error {
	size := int64(vm.Spec.DiskSize.Bytes())
	// Truncate only accepts an int64
	if size > math.MaxInt64 {
		return fmt.Errorf("requested size %d too large, cannot truncate", size)
	}

	// Get the size of the image ext4 file
	fi, err := os.Stat(path.Join(constants.IMAGE_DIR, vm.GetImageUID().String(), constants.IMAGE_FS))
	if err != nil {
		return err
	}
	imageSize := fi.Size()

	// The overlay needs to be at least as large as the image
	if size < imageSize {
		log.Warnf("warning: requested overlay size (%s) < image size (%s), using image size for overlay\n",
			vm.Spec.DiskSize.String(), meta.NewSizeFromBytes(uint64(imageSize)).String())
		size = imageSize
	}

	// Make sure the all directories above the snapshot directory exists
	if err := os.MkdirAll(path.Dir(vm.OverlayFile()), constants.DATA_DIR_PERM); err != nil {
		return err
	}

	overlayFile, err := os.Create(vm.OverlayFile())
	if err != nil {
		return fmt.Errorf("failed to create overlay file for %q, %v", vm.GetUID(), err)
	}
	defer overlayFile.Close()

	if err := overlayFile.Truncate(size); err != nil {
		return fmt.Errorf("failed to allocate overlay file for VM %q: %v", vm.GetUID(), err)
	}

	// populate the filesystem
	return copyToOverlay(vm)
}

func copyToOverlay(vm *vmmd.VM) error {
	if err := ActivateSnapshot(vm); err != nil {
		return err
	}
	defer DeactivateSnapshot(vm)

	mp, err := util.Mount(vm.SnapshotDev())
	if err != nil {
		return err
	}
	defer mp.Umount()

	// Copy the kernel files to the VM. TODO: Use snapshot overlaying instead.
	if err := copyKernelToOverlay(vm, mp.Path); err != nil {
		return err
	}

	// do not mutate vm.Spec.CopyFiles
	fileMappings := vm.Spec.CopyFiles

	if vm.Spec.SSH != nil {
		pubKeyPath := vm.Spec.SSH.PublicKey
		if vm.Spec.SSH.Generate {
			// generate a key if PublicKey is empty
			pubKeyPath, err = newSSHKeypair(vm.VM)
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
	if len(vm.Status.IPAddresses) > 0 {
		ip = vm.Status.IPAddresses[0]
	}

	if err := writeEtcHosts(vm.VM, mp.Path, vm.GetUID().String(), ip); err != nil {
		return err
	}

	// TODO: This code seems to be flaky and not always copy over the files?
	time.Sleep(500 * time.Millisecond)
	return nil
}

func copyKernelToOverlay(vm *vmmd.VM, mountPoint string) error {
	kernelUID := vm.GetKernelUID()
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
func writeEtcHosts(vm *api.VM, tmpDir, hostname string, primaryIP net.IP) error {
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

// Generate a new SSH keypair for the vm
func newSSHKeypair(vm *api.VM) (string, error) {
	privKeyPath := path.Join(vm.ObjectPath(), fmt.Sprintf(constants.VM_SSH_KEY_TEMPLATE, vm.GetUID()))
	// TODO: In future versions, let the user specify what key algorithm to use through the API types
	sshKeyAlgorithm := "ed25519"
	if util.FIPSEnabled() {
		// Use rsa on FIPS machines
		sshKeyAlgorithm = "rsa"
	}
	_, err := util.ExecuteCommand("ssh-keygen", "-q", "-t", sshKeyAlgorithm, "-N", "", "-f", privKeyPath)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s.pub", privKeyPath), nil
}
