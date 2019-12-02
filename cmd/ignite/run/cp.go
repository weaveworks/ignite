package run

import (
	"fmt"
	"path"

	log "github.com/sirupsen/logrus"
	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/util"
)

type CPFlags struct {
	Timeout      uint32
	IdentityFile string
	Recursive    bool
}

type cpOptions struct {
	*CPFlags
	vm     *api.VM
	source string
	dest   string
}

func (cf *CPFlags) NewCPOptions(vmMatch string, source string, dest string) (co *cpOptions, err error) {
	co = &cpOptions{CPFlags: cf}
	co.vm, err = getVMForMatch(vmMatch)
	co.source = source
	co.dest = dest
	return
}

func CP(co *cpOptions) error {
	// Check if the VM is running
	if !co.vm.Running() {
		return fmt.Errorf("VM %q is not running", co.vm.GetUID())
	}

	ipAddrs := co.vm.Status.IPAddresses
	if len(ipAddrs) == 0 {
		return fmt.Errorf("VM %q has no usable IP addresses", co.vm.GetUID())
	}

	// From run/ssh.go:
	// We're dealing with local VMs in a trusted (internal) subnet, disable some warnings for convenience
	// TODO: For security, track the known_hosts internally, do something about the IP collisions (if needed)
	scpOpts := []string{
		"LogLevel=ERROR", // Warning: Permanently added '<ip>' (ECDSA) to the list of known hosts.
		// We get this if the VM happens to get an address that another container has used:
		"UserKnownHostsFile=/dev/null", // WARNING: REMOTE HOST IDENTIFICATION HAS CHANGED!
		"StrictHostKeyChecking=no",     // The authenticity of host ***** can't be established
		fmt.Sprintf("ConnectTimeout=%d", co.Timeout),
	}

	scpArgs := make([]string, 0, len(scpOpts)*2+3)

	for _, opt := range scpOpts {
		scpArgs = append(scpArgs, "-o", opt)
	}

	scpArgs = append(scpArgs, "-i")

	// If an external identity file is specified, use it instead of the internal one
	if len(co.IdentityFile) > 0 {
		scpArgs = append(scpArgs, co.IdentityFile)
	} else {
		privKeyFile := path.Join(co.vm.ObjectPath(), fmt.Sprintf(constants.VM_SSH_KEY_TEMPLATE, co.vm.GetUID()))
		if !util.FileExists(privKeyFile) {
			return fmt.Errorf("no private key found for VM %q", co.vm.GetUID())
		}

		scpArgs = append(scpArgs, privKeyFile)
	}

	if co.Recursive {
		scpArgs = append(scpArgs, "-r")
	}

	// Add source, dest args
	scpArgs = append(scpArgs, co.source)
	scpArgs = append(scpArgs, fmt.Sprintf("root@%s:%s", ipAddrs[0], co.dest))

	// SSH into the VM
	if code, err := util.ExecForeground("scp", scpArgs...); err != nil {
		if code != 2 {
			return fmt.Errorf("SCP into VM %q failed: %v", co.vm.GetUID(), err)
		}

		// Code 255 is used for signaling a connection error, be it caused by
		// a failed connection attempt or disconnection by VM reboot.
		log.Warnf("SCP command terminated")
	}

	return nil
}
