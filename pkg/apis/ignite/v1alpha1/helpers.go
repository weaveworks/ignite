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
