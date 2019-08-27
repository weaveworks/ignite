package docker

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/weaveworks/ignite/pkg/network"
	dockernetwork "github.com/weaveworks/ignite/pkg/network/docker"
	"github.com/weaveworks/ignite/pkg/providers"
	"github.com/weaveworks/ignite/pkg/runtime"
)

func SetDockerNetwork() error {
	log.Trace("Initializing the Docker network provider...")
	if providers.Runtime.Name() != runtime.RuntimeDocker {
		return fmt.Errorf("the %q network plugin can only be used with the %q runtime", network.PluginDockerBridge, runtime.RuntimeDocker)
	}

	providers.NetworkPlugin = dockernetwork.GetDockerNetworkPlugin(providers.Runtime)
	return nil
}
