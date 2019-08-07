package docker

import (
	log "github.com/sirupsen/logrus"
	"github.com/weaveworks/ignite/pkg/providers"
	"github.com/weaveworks/ignite/pkg/runtime/docker"
)

func SetDockerRuntime() (err error) {
	log.Trace("Initializing the Docker provider...")
	providers.Runtime, err = docker.GetDockerClient()
	return
}
