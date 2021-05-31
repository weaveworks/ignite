package providers

import (
	log "github.com/sirupsen/logrus"
	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	"github.com/weaveworks/ignite/pkg/client"
	"github.com/weaveworks/ignite/pkg/network"
	"github.com/weaveworks/ignite/pkg/runtime"
	"github.com/weaveworks/libgitops/pkg/storage"
)

// IDPrefix is used for vm, container, and snapshot file/device names
// It's set by the ComponentConfig and Flag override logic and should default to `constants.IGNITE_PREFIX`
var IDPrefix string

// NetworkPluginName binds to the global flag to select the network plugin
// The default network plugin is "cni"
var NetworkPluginName network.PluginName

// NetworkPlugin provides the chosen network plugin that should be used
// This should be set after parsing user input on what network plugin to use
var NetworkPlugin network.Plugin

// RuntimeName binds to the global flag to select the container runtime
// The default runtime is "containerd"
var RuntimeName runtime.Name

// Runtime provides the chosen container runtime for retrieving OCI images and running VM containers
// This should be set after parsing user input on what runtime to use
var Runtime runtime.Interface

// Client is the default client that can be easily used
var Client *client.Client

// Storage is the default storage implementation
var Storage storage.Storage

var ComponentConfig *api.Configuration

// RegistryConfigDir is the container runtime registry configuration directory.
// This is used during operations like image import for loading registry
// configurations.
var RegistryConfigDir string

type ProviderInitFunc func() error

// Populate initializes all given providers
func Populate(providers []ProviderInitFunc) error {
	log.Trace("Populating providers...")
	for _, init := range providers {
		if err := init(); err != nil {
			return err
		}
	}

	return nil
}
