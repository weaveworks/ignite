package runtime

import (
	"io"
	"time"
)

type ImageInspectResult struct {
	ID          string
	RepoDigests []string
	Size        int64
}

type Interface interface {
	InspectImage(image string) (*ImageInspectResult, error)
	PullImage(image string) (io.ReadCloser, error)
	GetNetNS(containerID string) (string, error)
	RawClient() interface{}

	ExportImage(image string) (io.ReadCloser, string, error)
	RemoveContainer(container string) error
	StopContainer(container string, timeout *time.Duration) error
	KillContainer(container, signal string) error
}
