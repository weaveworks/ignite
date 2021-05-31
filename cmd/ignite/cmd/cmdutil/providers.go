package cmdutil

import (
	log "github.com/sirupsen/logrus"

	"github.com/weaveworks/ignite/pkg/providers"
)

// ResolveRegistryConfigDir reads various configuration to resolve the registry
// configuration directory.
func ResolveRegistryConfigDir() {
	if providers.ComponentConfig != nil {
		// Set the providers registry config dir from ignite configuration if
		// it's empty. When it's set in the providers and in the ignite
		// configuration, log about the override.
		if providers.RegistryConfigDir == "" {
			providers.RegistryConfigDir = providers.ComponentConfig.Spec.RegistryConfigDir
		} else if providers.ComponentConfig.Spec.RegistryConfigDir != "" {
			log.Debug("registry-config-dir flag overriding the ignite configuration")
		}
	}
}
