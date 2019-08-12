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
	var foundPlugin *network.PluginName
	for _, plugin := range plugins {
		if plugin.String() == val {
			foundPlugin = &plugin
			break
		}
	}
	if foundPlugin == nil {
		return fmt.Errorf("Invalid network mode %q, must be one of %v", val, plugins)
	}
	*nf.value = *foundPlugin
	return nil
}

func (nf *NetworkPluginFlag) String() string {
	if nf.value == nil {
		return ""
	}
	return nf.value.String()
}

func (nf *NetworkPluginFlag) Type() string {
	return "network-mode"
}

var _ pflag.Value = &NetworkPluginFlag{}

func NetworkPluginVar(fs *pflag.FlagSet, ptr *network.PluginName) {
	fs.Var(&NetworkPluginFlag{value: ptr}, "network-plugin", fmt.Sprintf("Networking mode to use. Available options are: %v", plugins))
}

// RegisterNetworkPluginFlag binds network.ActivePlugin to the --network-plugin flag
func RegisterNetworkPluginFlag(fs *pflag.FlagSet) {
	NetworkPluginVar(fs, &network.ActivePlugin)
}
