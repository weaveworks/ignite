package containerruntime

type Interface interface {
	GetNetNS(containerID string) (string, error)
	RawClient() interface{}
}
