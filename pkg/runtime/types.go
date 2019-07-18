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

type Bind struct {
	HostPath      string
	ContainerPath string
}

type ContainerConfig struct {
	Cmd          []string
	Hostname     string
	Labels       map[string]string
	Binds        []*Bind
	CapAdds      []string
	Devices      []string
	StopTimeout  uint32
	AutoRemove   bool
	NetworkMode  string
	PortBindings meta.PortMappings
}

type Interface interface {
	InspectImage(image string) (*ImageInspectResult, error)
	PullImage(image string) (io.ReadCloser, error)
	ExportImage(image string) (io.ReadCloser, string, error)

	// TODO: AttachContainer
	RunContainer(image string, config *ContainerConfig, name string) (string, error)
	StopContainer(container string, timeout *time.Duration) error
	KillContainer(container, signal string) error
	RemoveContainer(container string) error
	ContainerLogs(container string) (io.ReadCloser, error)
	ContainerNetNS(container string) (string, error)

	RawClient() interface{}
}
