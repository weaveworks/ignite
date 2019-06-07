package run

import (
	"fmt"
	"path"

	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/metadata/vmmd"
	"github.com/weaveworks/ignite/pkg/util"
)

type SSHOptions struct {
	VM           *vmmd.VMMetadata
	IdentityFile string
}

func SSH(so *SSHOptions) error {
	// Check if the VM is running
	if !so.VM.Running() {
		return fmt.Errorf("%s is not running", so.VM.ID)
	}

	ipAddrs := so.VM.VMOD().IPAddrs
	if len(ipAddrs) == 0 {
		return fmt.Errorf("VM %q has no usable IP addresses", so.VM.ID)
	}

	sshArgs := append(make([]string, 0, 3), fmt.Sprintf("root@%s", ipAddrs[0]), "-i")

	// If an external identity file is specified, use it instead of the internal one
	if len(so.IdentityFile) > 0 {
		sshArgs = append(sshArgs, so.IdentityFile)
	} else {
		privKeyFile := path.Join(so.VM.ObjectPath(), fmt.Sprintf(constants.VM_SSH_KEY_TEMPLATE, so.VM.ID))
		if !util.FileExists(privKeyFile) {
			return fmt.Errorf("no private key found for VM %q", so.VM.ID)
		}

		sshArgs = append(sshArgs, privKeyFile)
	}

	// Print the ID before calling SSH
	fmt.Println(so.VM.ID)

	// SSH into the VM
	if _, err := util.ExecForeground("ssh", sshArgs...); err != nil {
		return fmt.Errorf("SSH into VM %q failed: %v", so.VM.ID, err)
	}

	return nil
}
