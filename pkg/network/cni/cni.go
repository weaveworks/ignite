package cni

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sync"

	gocni "github.com/containerd/go-cni"
	log "github.com/sirupsen/logrus"
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/network"
	"github.com/weaveworks/ignite/pkg/runtime"
	"github.com/weaveworks/ignite/pkg/util"
)

const (
	// TODO: CNIBinDir and CNIConfDir should maybe be globally configurable?

	// CNIBinDir describes the directory where the CNI binaries are stored
	CNIBinDir = "/opt/cni/bin"
	// CNIConfDir describes the directory where the CNI plugin's configuration is stored
	CNIConfDir = "/etc/cni/net.d"
	// netNSPathFmt gives the path to the a process network namespace, given the pid
	netNSPathFmt = "/proc/%d/ns/net"
	// igniteCNIConfName is the filename of Ignite's CNI configuration file
	igniteCNIConfName = "10-ignite.conf"
)

// igniteCNIConf is a base CNI configuration that will enable VMs to access the internet connection (docker-bridge style)
var igniteCNIConf = `{
	"cniVersion": "0.4.0",
	"name": "ignite-containerd-bridge",
	"plugins": [
		{
			"type": "bridge",
			"bridge": "cni0",
			"isGateway": true,
			"isDefaultGateway": true,
			"promiscMode": true,
			"ipMasq": true,
			"ipam": {
				"type": "host-local",
				"subnet": "172.18.0.0/16"
			}
		},
		{
			"type": "portmap",
			"capabilities": {
				"portMappings": true
			}
		}
	]
}
`

type cniNetworkPlugin struct {
	cni     gocni.CNI
	runtime runtime.Interface
	once    *sync.Once
}

func GetCNINetworkPlugin(runtime runtime.Interface) (network.Plugin, error) {
	// If the CNI configuration directory doesn't exist, create it
	if !util.DirExists(CNIConfDir) {
		if err := os.MkdirAll(CNIConfDir, constants.DATA_DIR_PERM); err != nil {
			return nil, err
		}
	}

	binDirs := []string{CNIBinDir}
	cniInstance, err := gocni.New(gocni.WithMinNetworkCount(2),
		gocni.WithPluginConfDir(CNIConfDir),
		gocni.WithPluginDir(binDirs))
	if err != nil {
		return nil, err
	}

	return &cniNetworkPlugin{
		runtime: runtime,
		cni:     cniInstance,
		once:    &sync.Once{},
	}, nil
}

func (plugin *cniNetworkPlugin) Name() network.PluginName {
	return network.PluginCNI
}

func (plugin *cniNetworkPlugin) PrepareContainerSpec(container *runtime.ContainerConfig) error {
	// No need for the container runtime to set up networking, as this plugin will do it
	container.NetworkMode = "none"
	return nil
}

func (plugin *cniNetworkPlugin) SetupContainerNetwork(containerid string, portMappings ...meta.PortMapping) (*network.Result, error) {
	if err := plugin.initialize(); err != nil {
		return nil, err
	}

	c, err := plugin.runtime.InspectContainer(containerid)
	if err != nil {
		return nil, fmt.Errorf("CNI failed to retrieve network namespace path: %v", err)
	}

	pms := []gocni.PortMapping{}
	for _, pm := range portMappings {
		hostIP := ""
		if pm.BindAddress != nil {
			hostIP = pm.BindAddress.String()
		}
		pms = append(pms, gocni.PortMapping{
			HostPort:      int32(pm.HostPort),
			ContainerPort: int32(pm.VMPort),
			Protocol:      pm.Protocol.String(),
			HostIP:        hostIP,
		})
	}

	netnsPath := fmt.Sprintf(netNSPathFmt, c.PID)
	result, err := plugin.cni.Setup(context.Background(), containerid, netnsPath, gocni.WithCapabilityPortMap(pms))
	if err != nil {
		log.Errorf("failed to setup network for namespace %q: %v", containerid, err)
		return nil, err
	}

	return cniToIgniteResult(result), nil
}

func (plugin *cniNetworkPlugin) initialize() (err error) {
	// If there's no existing CNI configuration, write ignite's example config to the CNI directory
	if util.DirEmpty(CNIConfDir) {
		if err = ioutil.WriteFile(path.Join(CNIConfDir, igniteCNIConfName), []byte(igniteCNIConf), constants.DATA_DIR_FILE_PERM); err != nil {
			return
		}
	}

	plugin.once.Do(func() {
		if err = plugin.cni.Load(gocni.WithLoNetwork, gocni.WithDefaultConf); err != nil {
			log.Errorf("failed to load cni configuration: %v", err)
		}
	})
	return
}

func cniToIgniteResult(r *gocni.CNIResult) *network.Result {
	result := &network.Result{}
	for _, iface := range r.Interfaces {
		for _, ip := range iface.IPConfigs {
			result.Addresses = append(result.Addresses, network.Address{
				IP:      ip.IP,
				Gateway: ip.Gateway,
			})
		}
	}
	return result
}

func (plugin *cniNetworkPlugin) RemoveContainerNetwork(containerID string) error {
	if err := plugin.initialize(); err != nil {
		return err
	}

	// Lack of namespace should not be fatal on teardown
	c, err := plugin.runtime.InspectContainer(containerID)
	if err != nil {
		log.Infof("CNI failed to retrieve network namespace path: %v", err)
		return nil
	}

	netnsPath := fmt.Sprintf(netNSPathFmt, c.PID)
	if c.PID == 0 {
		log.Info("CNI failed to retrieve network namespace path, PID was 0")
		return nil
	}

	return plugin.cni.Remove(context.Background(), containerID, netnsPath)
}
