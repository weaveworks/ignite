package providers

import (
	"github.com/weaveworks/ignite/pkg/runtime"
	"github.com/weaveworks/ignite/pkg/runtime/docker"
)

var DefaultRuntime runtime.Interface

func DefaultDockerRuntime() (err error) {
	DefaultRuntime, err = docker.GetDockerClient()
	return
}
