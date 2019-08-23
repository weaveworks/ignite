package containerd

import (
	log "github.com/sirupsen/logrus"
	"github.com/weaveworks/ignite/pkg/providers"
	runtime "github.com/weaveworks/ignite/pkg/runtime/containerd"
)

func SetContainerdRuntime() (err error) {
	log.Trace("Initializing the containerd runtime provider...")
	providers.Runtime, err = runtime.GetContainerdClient()
	return
}
