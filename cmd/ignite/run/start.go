package run

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/weaveworks/ignite/pkg/apis/ignite"
	"github.com/weaveworks/ignite/pkg/operations"
	"github.com/weaveworks/ignite/pkg/preflight/checkers"
	"github.com/weaveworks/ignite/pkg/util"
	"golang.org/x/crypto/ssh"
	"k8s.io/apimachinery/pkg/util/sets"
)

type StartFlags struct {
	Interactive            bool
	Debug                  bool
	IgnoredPreflightErrors []string
}

type startOptions struct {
	*StartFlags
	*attachOptions
}

func (sf *StartFlags) NewStartOptions(vmMatch string) (*startOptions, error) {
	ao, err := NewAttachOptions(vmMatch)
	if err != nil {
		return nil, err
	}

	// Disable running check as it takes a while for ignite-spawn to update the state
	ao.checkRunning = false

	return &startOptions{sf, ao}, nil
}

func Start(so *startOptions) error {
	// Check if the given VM is already running
	if so.vm.Running() {
		return fmt.Errorf("VM %q is already running", so.vm.GetUID())
	}

	ignoredPreflightErrors := sets.NewString(util.ToLower(so.StartFlags.IgnoredPreflightErrors)...)
	if err := checkers.StartCmdChecks(so.vm, ignoredPreflightErrors); err != nil {
		return err
	}

	if err := operations.StartVM(so.vm, so.Debug); err != nil {
		return err
	}

	// When --ssh is enabled, wait until SSH service started on port 22 at most N seconds
	if ssh := so.vm.Spec.SSH; ssh != nil && ssh.Generate && len(so.vm.Status.IPAddresses) > 0 {
		if err := waitForSSH(so.vm, 10, 5); err != nil {
			return err
		}
	}

	// If starting interactively, attach after starting
	if so.Interactive {
		return Attach(so.attachOptions)
	}
	return nil
}

func dialSuccess(vm *ignite.VM, seconds int) error {
	addr := vm.Status.IPAddresses[0].String() + ":22"
	perSecond := 10
	delay := time.Second / time.Duration(perSecond)
	var err error
	for i := 0; i < seconds*perSecond; i++ {
		conn, dialErr := net.DialTimeout("tcp", addr, delay)
		if conn != nil {
			defer conn.Close()
			err = nil
			break
		}
		err = dialErr
		time.Sleep(delay)
	}
	if err != nil {
		if err, ok := err.(*net.OpError); ok && err.Timeout() {
			return fmt.Errorf("Tried connecting to SSH but timed out %s", err)
		}
		return err
	}

	return nil
}

func waitForSSH(vm *ignite.VM, dialSeconds, sshTimeout int) error {
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
		Timeout:         time.Duration(sshTimeout) * time.Second,
	}

	addr := vm.Status.IPAddresses[0].String() + ":22"
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
