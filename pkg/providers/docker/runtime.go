package docker

import (
	log "github.com/sirupsen/logrus"
	"github.com/weaveworks/ignite/pkg/providers"
	dockerruntime "github.com/weaveworks/ignite/pkg/runtime/docker"
)

func SetDockerRuntime() (err error) {
	log.Trace("Initializing the Docker runtime provider...")
	providers.Runtime, err = dockerruntime.GetDockerClient()
	return
}
