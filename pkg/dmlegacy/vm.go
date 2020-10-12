package dmlegacy

import (
	"fmt"
	"io/ioutil"
	"math"
	"net"
	"os"
	"path"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/operations/lookup"
	"github.com/weaveworks/ignite/pkg/providers"
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
func AllocateAndPopulateOverlay(vm *api.VM) error {
	requestedSize := vm.Spec.DiskSize.Bytes()
	// Truncate only accepts an int64
	if requestedSize > math.MaxInt64 {
		return fmt.Errorf("requested size %d too large, cannot truncate", requestedSize)
	}
	size := int64(requestedSize)

	// Get the image UID from the VM
	imageUID, err := lookup.ImageUIDForVM(vm, providers.Client)
	if err != nil {
		return err
	}

	// Get the size of the image ext4 file
	fi, err := os.Stat(path.Join(constants.IMAGE_DIR, imageUID.String(), constants.IMAGE_FS))
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

func copyToOverlay(vm *api.VM) (err error) {
	_, err = ActivateSnapshot(vm)
	if err != nil {
		return
	}
	defer util.DeferErr(&err, func() error { return DeactivateSnapshot(vm) })

	mp, err := util.Mount(vm.SnapshotDev())
	if err != nil {
		return
	}
	defer util.DeferErr(&err, mp.Umount)

	// Copy the kernel files to the VM. TODO: Use snapshot overlaying instead.
	if err = copyKernelToOverlay(vm, mp.Path); err != nil {
		return
	}

	// do not mutate vm.Spec.CopyFiles
	fileMappings := vm.Spec.CopyFiles

	if vm.Spec.SSH != nil {
		pubKeyPath := vm.Spec.SSH.PublicKey
		if vm.Spec.SSH.Generate {
			// generate a key if PublicKey is empty
			pubKeyPath, err = newSSHKeypair(vm)
			if err != nil {
				return
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
		if err = os.MkdirAll(path.Dir(vmFilePath), constants.DATA_DIR_PERM); err != nil {
			return
		}

		if err = util.CopyFile(mapping.HostPath, vmFilePath); err != nil {
			return
		}
	}

	ip := net.IP{127, 0, 0, 1}
	if len(vm.Status.Network.IPAddresses) > 0 {
		ip = vm.Status.Network.IPAddresses[0]
	}

	// Write /etc/hosts for the VM
	if err = writeEtcHosts(mp.Path, vm.GetUID().String(), ip); err != nil {
		return
	}

	// Write the UID to /etc/hostname for the VM
	if err = writeEtcHostname(mp.Path, vm.GetUID().String()); err != nil {
		return
	}

	// Populate /etc/fstab with the VM's volume mounts
	if err = populateFstab(vm, mp.Path); err != nil {
		return
	}

	// Set overlay root permissions
	err = os.Chmod(mp.Path, constants.DATA_DIR_PERM)

	return
}

func copyKernelToOverlay(vm *api.VM, mountPoint string) error {
	kernelUID, err := lookup.KernelUIDForVM(vm, providers.Client)
	if err != nil {
		return err
	}
	kernelTarPath := path.Join(constants.KERNEL_DIR, kernelUID.String(), constants.KERNEL_TAR)

	if !util.FileExists(kernelTarPath) {
		log.Warnf("Could not find kernel overlay files, not copying into the VM.")
		return nil
	}
	// It is important to not replace existing symlinks to directories when extracting because
	// /lib might be a symlink (usually to usr/lib)
	// WARNING: `-h` is only available on GNU and Busybox tar.  It is not present on BSD's bsdtar
	//          It means "Follow symlinks"
	_, err = util.ExecuteCommand("tar", "-h", "-xf", kernelTarPath, "-C", mountPoint)
	return err
}

// writeEtcHosts populates the /etc/hosts file to avoid errors like
// sudo: unable to resolve host 4462576f8bf5b689
func writeEtcHosts(tmpDir, hostname string, primaryIP net.IP) error {
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

func writeEtcHostname(tmpDir, hostname string) error {
	hostnameFilePath := filepath.Join(tmpDir, "/etc/hostname")
	return ioutil.WriteFile(hostnameFilePath, []byte(hostname), 0644)
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
