package docker

import (
	log "github.com/sirupsen/logrus"
	network "github.com/weaveworks/ignite/pkg/network/docker"
	"github.com/weaveworks/ignite/pkg/providers"
	runtime "github.com/weaveworks/ignite/pkg/runtime/docker"
)

func SetDockerRuntime() (err error) {
	log.Trace("Initializing the Docker runtime provider...")
	providers.Runtime, err = runtime.GetDockerClient()
	return
}

func SetDockerNetwork() error {
	log.Trace("Initializing the Docker network provider...")
	plugin := network.GetDockerNetworkPlugin(providers.Runtime)
	providers.NetworkPlugins[plugin.Name()] = plugin
	return nil
}
