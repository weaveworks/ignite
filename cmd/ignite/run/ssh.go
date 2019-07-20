package run

import (
	"fmt"
	"path"

	log "github.com/sirupsen/logrus"
	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/util"

	_ "golang.org/x/crypto/ssh"
)

type SSHFlags struct {
	Timeout      uint32
	IdentityFile string
}

type sshOptions struct {
	*SSHFlags
	vm *api.VM
}

func (sf *SSHFlags) NewSSHOptions(vmMatch string) (so *sshOptions, err error) {
	so = &sshOptions{SSHFlags: sf}
	so.vm, err = getVMForMatch(vmMatch)
	return
}

func SSH(so *sshOptions) error {
	// Check if the VM is running
	if !so.vm.Running() {
		return fmt.Errorf("VM %q is not running", so.vm.GetUID())
	}

	ipAddrs := so.vm.Status.IPAddresses
	if len(ipAddrs) == 0 {
		return fmt.Errorf("VM %q has no usable IP addresses", so.vm.GetUID())
	}

	// We're dealing with local VMs in a trusted (internal) subnet, disable some warnings for convenience
	// TODO: For security, track the known_hosts internally, do something about the IP collisions (if needed)
	sshOpts := []string{
		"LogLevel=ERROR", // Warning: Permanently added '<ip>' (ECDSA) to the list of known hosts.
		// We get this if the VM happens to get an address that another container has used:
		"UserKnownHostsFile=/dev/null", // WARNING: REMOTE HOST IDENTIFICATION HAS CHANGED!
		"StrictHostKeyChecking=no",     // The authenticity of host ***** can't be established
		fmt.Sprintf("ConnectTimeout=%d", so.Timeout),
	}

	sshArgs := append(make([]string, 0, len(sshOpts)*2+3),
		fmt.Sprintf("root@%s", ipAddrs[0]))

	for _, opt := range sshOpts {
		sshArgs = append(sshArgs, "-o", opt)
	}

	sshArgs = append(sshArgs, "-i")

	// If an external identity file is specified, use it instead of the internal one
	if len(so.IdentityFile) > 0 {
		sshArgs = append(sshArgs, so.IdentityFile)
	} else {
		privKeyFile := path.Join(so.vm.ObjectPath(), fmt.Sprintf(constants.VM_SSH_KEY_TEMPLATE, so.vm.GetUID()))
		if !util.FileExists(privKeyFile) {
			return fmt.Errorf("no private key found for VM %q", so.vm.GetUID())
		}

		sshArgs = append(sshArgs, privKeyFile)
	}

	// SSH into the VM
	if code, err := util.ExecForeground("ssh", sshArgs...); err != nil {
		if code != 255 {
			return fmt.Errorf("SSH into VM %q failed: %v", so.vm.GetUID(), err)
		}

		// Code 255 is used for signaling a connection error, be it caused by
		// a failed connection attempt or disconnection by VM reboot.
		log.Warnf("SSH command terminated")
	}

	return nil
}
