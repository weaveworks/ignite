package providers

import (
	log "github.com/sirupsen/logrus"
	"github.com/weaveworks/libgitops/pkg/storage"
	"github.com/weaveworks/ignite/pkg/client"
	"github.com/weaveworks/ignite/pkg/network"
	"github.com/weaveworks/ignite/pkg/runtime"
)

// NetworkPluginName binds to the global flag to select the network plugin
// The default network plugin is "cni"
var NetworkPluginName = network.PluginCNI

// NetworkPlugin provides the chosen network plugin that should be used
// This should be set after parsing user input on what network plugin to use
var NetworkPlugin network.Plugin

// RuntimeName binds to the global flag to select the container runtime
// The default runtime is "containerd"
var RuntimeName = runtime.RuntimeContainerd

// Runtime provides the chosen container runtime for retrieving OCI images and running VM containers
// This should be set after parsing user input on what runtime to use
var Runtime runtime.Interface

// Client is the default client that can be easily used
var Client *client.Client

// Storage is the default storage implementation
var Storage storage.Storage

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
