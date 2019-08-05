package cmdutil

import (
	"fmt"

	"github.com/spf13/pflag"
	api "github.com/weaveworks/ignite/pkg/apis/ignite"
)

type NetworkModeFlag struct {
	value *api.NetworkMode
}

func (nf *NetworkModeFlag) Set(val string) error {
	*nf.value = api.NetworkMode(val)
	return nil
}

func (nf *NetworkModeFlag) String() string {
	if nf.value == nil {
		return ""
	}
	return nf.value.String()
}

func (nf *NetworkModeFlag) Type() string {
	return "network-mode"
}

var _ pflag.Value = &NetworkModeFlag{}

func NetworkModeVar(fs *pflag.FlagSet, ptr *api.NetworkMode) {
	fs.Var(&NetworkModeFlag{value: ptr}, "net", fmt.Sprintf("Networking mode to use. Available options are: %v", api.GetNetworkModes()))
}
