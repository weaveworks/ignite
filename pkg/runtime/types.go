package runtime

import "io"

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
}
