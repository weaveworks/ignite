package runtime

import (
	"io"
	"net"
	"time"

	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
)

type ImageInspectResult struct {
	ID   *meta.OCIContentID
	Size int64
}

type ContainerInspectResult struct {
	ID        string
	Image     string
	Status    string
	IPAddress net.IP
	PID       uint32
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
}

type Interface interface {
	PullImage(image meta.OCIImageRef) error
	InspectImage(image meta.OCIImageRef) (*ImageInspectResult, error)
	ExportImage(image meta.OCIImageRef) (io.ReadCloser, func() error, error)

	InspectContainer(container string) (*ContainerInspectResult, error)
	AttachContainer(container string) error
	RunContainer(image meta.OCIImageRef, config *ContainerConfig, name string) (string, error)
	StopContainer(container string, timeout *time.Duration) error
	KillContainer(container, signal string) error
	RemoveContainer(container string) error
	ContainerLogs(container string) (io.ReadCloser, error)

	RawClient() interface{}
}
