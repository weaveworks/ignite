package run

import (
	"fmt"

	api "github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/metadata/loader"
	"github.com/weaveworks/ignite/pkg/operations"
)

type StartFlags struct {
	PortMappings []string
	Interactive  bool
	Debug        bool
	// TODO: Make a dedicated flag for networkMode, so we can bind to the custom type directly
	NetworkMode string
}

type startOptions struct {
	*StartFlags
	*attachOptions
}

func (sf *StartFlags) NewStartOptions(l *loader.ResLoader, vmMatch string) (*startOptions, error) {
	ao, err := NewAttachOptions(l, vmMatch)
	if err != nil {
		return nil, err
	}

	// Disable running check as it takes a while for the in-container Ignite to update the state
	ao.checkRunning = false

	return &startOptions{sf, ao}, nil
}

func Start(so *startOptions) error {
	// Parse the given port mappings
	var err error
	if so.vm.Spec.Ports, err = meta.ParsePortMappings(so.PortMappings); err != nil {
		return err
	}

	// Validate and set the desired networking mode
	nm := api.NetworkMode(so.NetworkMode)
	if err := api.ValidateNetworkMode(nm); err != nil {
		return err
	}
	so.vm.Spec.NetworkMode = nm
	fmt.Println("nm", nm)

	// Save the port mappings into the VM metadata
	if err := so.vm.Save(); err != nil {
		return err
	}

	// Check if the given VM is already running
	if so.vm.Running() {
		return fmt.Errorf("VM %q is already running", so.vm.GetUID())
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
