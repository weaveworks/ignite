package flag

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/weaveworks/ignite/pkg/network"
)

var plugins = network.ListPlugins()

type NetworkPluginFlag struct {
	value *network.PluginName
}

func (nf *NetworkPluginFlag) Set(val string) error {
	for _, plugin := range plugins {
		if plugin.String() == val {
			*nf.value = plugin
			return nil
		}
	}
	return fmt.Errorf("invalid network plugin %q, must be one of %v", val, plugins)
}

func (nf *NetworkPluginFlag) String() string {
	if nf.value == nil {
		return ""
	}
	return nf.value.String()
}

func (nf *NetworkPluginFlag) Type() string {
	return "plugin"
}

var _ pflag.Value = &NetworkPluginFlag{}

func NetworkPluginVar(fs *pflag.FlagSet, ptr *network.PluginName) {
	fs.Var(&NetworkPluginFlag{value: ptr}, "network-plugin", fmt.Sprintf("Network plugin to use. Available options are: %v (default %v)", plugins, network.PluginCNI))
}
