package v1alpha1

import "fmt"

// GetNetworkModes gets the list of available network modes
func GetNetworkModes() []NetworkMode {
	return []NetworkMode{
		NetworkModeCNI,
		NetworkModeDockerBridge,
	}
}

// ValidateNetworkMode validates the network mode
// TODO: This should move into a dedicated validation package
func ValidateNetworkMode(mode NetworkMode) error {
	found := false
	modes := GetNetworkModes()
	for _, nm := range modes {
		if nm == mode {
			found = true
		}
	}
	if !found {
		return fmt.Errorf("invalid network mode %s, must be one of %v", mode, modes)
	}
	return nil
}

// SetImage populates relevant fields to an Image on the VM object
func (vm *VM) SetImage(image *Image) {
	vm.Spec.Image.OCIClaim = image.Spec.OCIClaim
	vm.Status.Image = image.Status.OCISource
}

// SetKernel populates relevant fields to a Kernel on the VM object
func (vm *VM) SetKernel(kernel *Kernel) {
	vm.Spec.Kernel.OCIClaim = kernel.Spec.OCIClaim
	vm.Status.Kernel = kernel.Status.OCISource
}
