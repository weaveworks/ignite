package runtime

import (
	"fmt"

	"github.com/weaveworks/ignite/pkg/providers"
	containerdprovider "github.com/weaveworks/ignite/pkg/providers/containerd"
	dockerprovider "github.com/weaveworks/ignite/pkg/providers/docker"
	"github.com/weaveworks/ignite/pkg/runtime"
)

func SetRuntime() error {
	switch providers.RuntimeName {
	case runtime.RuntimeDocker:
		return dockerprovider.SetDockerRuntime() // Use the Docker runtime
	case runtime.RuntimeContainerd:
		return containerdprovider.SetContainerdRuntime() // Use the containerd runtime
	}

	return fmt.Errorf("unknown runtime %q", providers.RuntimeName)
}
