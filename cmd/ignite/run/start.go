package run

import (
	"fmt"

	"github.com/weaveworks/ignite/pkg/operations"
	"github.com/weaveworks/ignite/pkg/preflight"
	"github.com/weaveworks/ignite/pkg/util"
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
	if err := preflight.StartCmdChecks(so.vm, ignoredPreflightErrors); err != nil {
		return err
	}

	if err := operations.StartVM(so.vm, so.Debug); err != nil {
		return err
	}

	// If starting interactively, attach after starting
	if so.Interactive {
		return Attach(so.attachOptions)
	}
	return nil
}
