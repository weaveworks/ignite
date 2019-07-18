package providers

import (
	"github.com/weaveworks/ignite/pkg/runtime"
	"github.com/weaveworks/ignite/pkg/runtime/docker"
)

var Runtime runtime.Interface

func SetDockerRuntime() (err error) {
	Runtime, err = docker.GetDockerClient()
	return
}
