package run

import (
	"fmt"
	"path"

	"github.com/weaveworks/ignite/pkg/client"
	"github.com/weaveworks/ignite/pkg/filter"

	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/metadata/vmmd"
	"github.com/weaveworks/ignite/pkg/util"
)

type SSHFlags struct {
	IdentityFile string
}

type sshOptions struct {
	*SSHFlags
	vm *vmmd.VM
}

func (sf *SSHFlags) NewSSHOptions(vmMatch string) (*sshOptions, error) {
	so := &sshOptions{SSHFlags: sf}

	if vm, err := client.VMs().Find(filter.NewIDNameFilter(vmMatch)); err == nil {
		so.vm = &vmmd.VM{vm}
	} else {
		return nil, err
	}

	return so, nil
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

	sshArgs := append(make([]string, 0, 3), fmt.Sprintf("root@%s", ipAddrs[0]), "-i")

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

	// SSH into the vm
	if _, err := util.ExecForeground("ssh", sshArgs...); err != nil {
		return fmt.Errorf("SSH into VM %q failed: %v", so.vm.GetUID(), err)
	}
	return nil
}
