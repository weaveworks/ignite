package run

import (
	"fmt"
	"net"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"golang.org/x/crypto/ssh"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/weaveworks/ignite/pkg/apis/ignite"
	"github.com/weaveworks/ignite/pkg/config"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/operations"
	"github.com/weaveworks/ignite/pkg/preflight/checkers"
	"github.com/weaveworks/ignite/pkg/providers"
	"github.com/weaveworks/ignite/pkg/util"
)

type StartFlags struct {
	Interactive            bool
	Debug                  bool
	IgnoredPreflightErrors []string
}

type StartOptions struct {
	*StartFlags
	*AttachOptions
}

func (sf *StartFlags) NewStartOptions(vmMatch string) (*StartOptions, error) {
	ao, err := NewAttachOptions(vmMatch)
	if err != nil {
		return nil, err
	}

	// Disable running check as it takes a while for ignite-spawn to update the state
	ao.checkRunning = false

	return &StartOptions{sf, ao}, nil
}

func Start(so *StartOptions, fs *flag.FlagSet) error {
	// Check if the given VM is already running
	if so.vm.Running() {
		return fmt.Errorf("VM %q is already running", so.vm.GetUID())
	}

	// Stopped VMs don't contain the runtime and network information. Set the
	// default runtime and network from the providers if empty.
	if so.vm.Status.Runtime.Name == "" {
		so.vm.Status.Runtime.Name = providers.RuntimeName
	}
	if so.vm.Status.Network.Plugin == "" {
		so.vm.Status.Network.Plugin = providers.NetworkPluginName
	}

	// In case the runtime and network-plugin are specified explicitly at
	// start, set the runtime and network-plugin on the VM. This overrides the
	// global config and config on the VM object, if any.
	if fs.Changed("runtime") {
		so.vm.Status.Runtime.Name = providers.RuntimeName
	}
	if fs.Changed("network-plugin") {
		so.vm.Status.Network.Plugin = providers.NetworkPluginName
	}

	// Set the runtime and network-plugin providers from the VM status.
	if err := config.SetAndPopulateProviders(so.vm.Status.Runtime.Name, so.vm.Status.Network.Plugin); err != nil {
		return err
	}

	ignoredPreflightErrors := sets.NewString(util.ToLower(so.StartFlags.IgnoredPreflightErrors)...)
	if err := checkers.StartCmdChecks(so.vm, ignoredPreflightErrors); err != nil {
		return err
	}

	if err := operations.StartVM(so.vm, so.Debug); err != nil {
		return err
	}

	// When --ssh is enabled, wait until SSH service started on port 22 at most N seconds
	if ssh := so.vm.Spec.SSH; ssh != nil && ssh.Generate && len(so.vm.Status.Network.IPAddresses) > 0 {
		if err := waitForSSH(so.vm, constants.SSH_DEFAULT_TIMEOUT_SECONDS, constants.IGNITE_SPAWN_TIMEOUT); err != nil {
			return err
		}
	}

	// If starting interactively, attach after starting
	if so.Interactive {
		return Attach(so.AttachOptions)
	}
	return nil
}

func dialSuccess(vm *ignite.VM, seconds int) error {
	addr := vm.Status.Network.IPAddresses[0].String() + ":22"
	perSecond := 10
	delay := time.Second / time.Duration(perSecond)
	var err error
	for i := 0; i < seconds*perSecond; i++ {
		conn, dialErr := net.DialTimeout("tcp", addr, delay)
		if conn != nil {
			conn.Close()
			err = nil
			break
		}
		err = dialErr
		time.Sleep(delay)
		// Report every ten seconds that we're waiting for SSH
		if i%(10*perSecond) == 0 {
			log.Info("Waiting for the ssh daemon within the VM to start...")
		}
	}
	if err != nil {
		if err, ok := err.(*net.OpError); ok && err.Timeout() {
			return fmt.Errorf("Tried connecting to SSH but timed out %s", err)
		}
		return err
	}

	return nil
}

func waitForSSH(vm *ignite.VM, dialSeconds int, sshTimeout time.Duration) error {
	if err := dialSuccess(vm, dialSeconds); err != nil {
		return err
	}

	certCheck := &ssh.CertChecker{
		IsHostAuthority: func(auth ssh.PublicKey, address string) bool {
			return true
		},
		IsRevoked: func(cert *ssh.Certificate) bool {
			return false
		},
		HostKeyFallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	config := &ssh.ClientConfig{
		HostKeyCallback: certCheck.CheckHostKey,
		Timeout:         sshTimeout,
	}

	addr := vm.Status.Network.IPAddresses[0].String() + ":22"
	sshConn, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		if strings.Contains(err.Error(), "unable to authenticate") {
			// we connected to the ssh server and recieved the expected failure
			return nil
		}
		return err
	}

	defer sshConn.Close()
	return fmt.Errorf("waitForSSH: connected successfully with no authentication -- failure was expected")
}
