package run

import (
	"fmt"
	"net"
	"time"

	"github.com/weaveworks/ignite/pkg/apis/ignite"

	"github.com/weaveworks/ignite/pkg/operations"
)

type StartFlags struct {
	Interactive bool
	Debug       bool
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

	if err := operations.StartVM(so.vm, so.Debug); err != nil {
		return err
	}

	if err := waitForSSH(so.vm, 10); err != nil {
		return err
	}

	// If starting interactively, attach after starting
	if so.Interactive {
		return Attach(so.attachOptions)
	}
	return nil
}

func waitForSSH(vm *ignite.VM, seconds int) error {
	// When --ssh is enabled, wait until SSH service started on port 22 at most N seconds
	ssh := vm.Spec.SSH
	if ssh != nil && ssh.Generate && len(vm.Status.IPAddresses) > 0 {
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
	}

	return nil
}
