package runtime

import (
	"io"
	"time"

	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
)

type ImageInspectResult struct {
	ID          string
	RepoDigests []string
	Size        int64
}

type ContainerInspectResult struct {
	ID     string
	Image  string
	Status string
}

type Bind struct {
	HostPath      string
	ContainerPath string
}

// Convenience generator for Binds which have the same host and container path
func BindBoth(path string) *Bind {
	return &Bind{
		HostPath:      path,
		ContainerPath: path,
	}
}

type ContainerConfig struct {
	Cmd          []string
	Hostname     string
	Labels       map[string]string
	Binds        []*Bind
	CapAdds      []string
	Devices      []*Bind
	StopTimeout  uint32
	AutoRemove   bool
	NetworkMode  string
	PortBindings meta.PortMappings
	Env          []string
}

type Interface interface {
	InspectImage(image string) (*ImageInspectResult, error)
	PullImage(image string) (io.ReadCloser, error)
	ExportImage(image string) (io.ReadCloser, string, error)

	InspectContainer(container string) (*ContainerInspectResult, error)
	AttachContainer(container string) error
	RunContainer(image string, config *ContainerConfig, name string) (string, error)
	StopContainer(container string, timeout *time.Duration) error
	KillContainer(container, signal string) error
	RemoveContainer(container string) error
	ContainerLogs(container string) (io.ReadCloser, error)
	ContainerNetNS(container string) (string, error)

	RawClient() interface{}
}
