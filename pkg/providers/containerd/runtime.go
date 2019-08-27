package containerd

import (
	log "github.com/sirupsen/logrus"
	"github.com/weaveworks/ignite/pkg/providers"
	containerdruntime "github.com/weaveworks/ignite/pkg/runtime/containerd"
)

func SetContainerdRuntime() (err error) {
	log.Trace("Initializing the containerd runtime provider...")
	providers.Runtime, err = containerdruntime.GetContainerdClient()
	return
}
