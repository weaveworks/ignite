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

type Volume struct {
	HostPath      string
	ContainerPath string
}

type ContainerConfig struct {
	Cmd          string
	Hostname     string
	Labels       []string
	Volumes      []*Volume
	CapAdds      []string
	Devices      []string
	StopTimeout  uint32
	AutoRemove   bool
	NetworkMode  string
	ExposedPorts meta.PortMappings
}

type Interface interface {
	InspectImage(image string) (*ImageInspectResult, error)
	PullImage(image string) (io.ReadCloser, error)
	ExportImage(image string) (io.ReadCloser, string, error)
	GetNetNS(containerID string) (string, error)
	RawClient() interface{}

	RunContainer(image string, config *ContainerConfig, name string) error
	RemoveContainer(container string) error
	StopContainer(container string, timeout *time.Duration) error
	KillContainer(container, signal string) error
}
