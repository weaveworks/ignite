package ignite

// GetNetworkModes gets the list of available network modes
func GetNetworkModes() []NetworkMode {
	return []NetworkMode{
		NetworkModeCNI,
		NetworkModeDockerBridge,
	}
}

// GetImageSourceTypes gets the list of available network modes
func GetImageSourceTypes() []ImageSourceType {
	return []ImageSourceType{
		ImageSourceTypeDocker,
	}
}

// GetVMStates gets the list of available VM states
func GetVMStates() []VMState {
	return []VMState{
		VMStateCreated,
		VMStateRunning,
		VMStateStopped,
	}
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
