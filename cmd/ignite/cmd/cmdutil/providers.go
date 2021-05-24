package cmdutil

import (
	log "github.com/sirupsen/logrus"

	"github.com/weaveworks/ignite/pkg/providers"
)

// ResolveClientConfigDir reads various configuration to resolve the client
// configuration directory.
func ResolveClientConfigDir() {
	if providers.ComponentConfig != nil {
		// Set the providers client config dir from ignite configuration if
		// it's empty. When it's set in the providers and in the ignite
		// configuration, log about the override.
		if providers.ClientConfigDir == "" {
			providers.ClientConfigDir = providers.ComponentConfig.Spec.ClientConfigDir
		} else if providers.ComponentConfig.Spec.ClientConfigDir != "" {
			log.Debug("client-config-dir flag overriding the ignite configuration")
		}
	}
}
