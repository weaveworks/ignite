package providers

import (
	"github.com/weaveworks/ignite/pkg/network/cni"
)

var NetworkPlugin cni.NetworkPlugin

func SetCNINetworkPlugin() (err error) {
	NetworkPlugin, err = cni.GetCNINetworkPlugin(Runtime)
	return
}
