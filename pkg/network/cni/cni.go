package cni

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/containernetworking/cni/libcni"
	cnitypes "github.com/containernetworking/cni/pkg/types"
	log "github.com/sirupsen/logrus"
	"github.com/weaveworks/ignite/pkg/network"
	"github.com/weaveworks/ignite/pkg/runtime"
)

// Disclaimer: This package is heavily influenced by
// https://github.com/kubernetes/kubernetes/blob/v1.15.0/pkg/kubelet/dockershim/network/cni/cni.go#L49

type cniNetworkPlugin struct {
	sync.RWMutex

	loNetwork      *cniNetwork
	defaultNetwork *cniNetwork

	runtime runtime.Interface
	confDir string
	binDirs []string
}

type cniNetwork struct {
	name          string
	NetworkConfig *libcni.NetworkConfigList
	CNIConfig     libcni.CNI
}

func GetCNINetworkPlugin(runtime runtime.Interface) (network.Plugin, error) {
	binDirs := []string{CNIBinDir}
	plugin := &cniNetworkPlugin{
		runtime:        runtime,
		defaultNetwork: nil,
		loNetwork:      getLoNetwork(binDirs),
		confDir:        CNIConfDir,
		binDirs:        binDirs,
	}

	return plugin, nil
}

func getLoNetwork(binDirs []string) *cniNetwork {
	loConfig, err := libcni.ConfListFromBytes([]byte(loopbackCNIConfig))
	if err != nil {
		// The hardcoded config above should always be valid and unit tests will
		// catch this
		panic(err)
	}

	return &cniNetwork{
		name:          "lo",
		NetworkConfig: loConfig,
		CNIConfig:     &libcni.CNIConfig{Path: binDirs},
	}
}

func getDefaultCNINetwork(confDir string, binDirs []string) (*cniNetwork, error) {
	files, err := libcni.ConfFiles(confDir, []string{".conf", ".conflist", ".json"})
	switch {
	case err != nil:
		return nil, err
	case len(files) == 0:
		return nil, fmt.Errorf("no networks found in %s", confDir)
	}

	sort.Strings(files)
	for _, confFile := range files {
		var confList *libcni.NetworkConfigList
		if strings.HasSuffix(confFile, ".conflist") {
			confList, err = libcni.ConfListFromFile(confFile)
			if err != nil {
				log.Infof("Error loading CNI config list file %s: %v", confFile, err)
				continue
			}
		} else {
			conf, err := libcni.ConfFromFile(confFile)
			if err != nil {
				log.Infof("Error loading CNI config file %s: %v", confFile, err)
				continue
			}

			// Ensure the config has a "type" so we know what plugin to run.
			// Also catches the case where somebody put a conflist into a conf file.
			if conf.Network.Type == "" {
				log.Infof("Error loading CNI config file %s: no 'type'; perhaps this is a .conflist?", confFile)
				continue
			}

			confList, err = libcni.ConfListFromConf(conf)
			if err != nil {
				log.Infof("Error converting CNI config file %s to list: %v", confFile, err)
				continue
			}
		}

		if len(confList.Plugins) == 0 {
			log.Infof("CNI config list %s has no networks, skipping", confFile)
			continue
		}

		log.Infof("Using CNI configuration file %s", confFile)

		network := &cniNetwork{
			name:          confList.Name,
			NetworkConfig: confList,
			CNIConfig:     &libcni.CNIConfig{Path: binDirs},
		}

		return network, nil
	}

	return nil, fmt.Errorf("no valid networks found in %s", confDir)
}

func (plugin *cniNetworkPlugin) syncNetworkConfig() error {
	network, err := getDefaultCNINetwork(plugin.confDir, plugin.binDirs)
	if err != nil {
		return fmt.Errorf("unable to get default CNI network: %v", err)
	}

	plugin.setDefaultNetwork(network)
	return nil
}

func (plugin *cniNetworkPlugin) getDefaultNetwork() *cniNetwork {
	plugin.RLock()
	defer plugin.RUnlock()
	return plugin.defaultNetwork
}

func (plugin *cniNetworkPlugin) setDefaultNetwork(n *cniNetwork) {
	plugin.Lock()
	defer plugin.Unlock()
	plugin.defaultNetwork = n
}

func (plugin *cniNetworkPlugin) checkInitialized() error {
	if plugin.getDefaultNetwork() == nil {
		// Sync the network configuration if the plugin is not initialized
		if err := plugin.syncNetworkConfig(); err != nil {
			return err
		}
	}

	return nil
}

func (plugin *cniNetworkPlugin) Name() string {
	return CNIPluginName
}

func (plugin *cniNetworkPlugin) Status() error {
	// Can't set up pods if we don't have any CNI network configs yet
	return plugin.checkInitialized()
}

func (plugin *cniNetworkPlugin) SetupContainerNetwork(containerid string) error {
	if err := plugin.checkInitialized(); err != nil {
		return err
	}

	netnsPath, err := plugin.runtime.ContainerNetNS(containerid)
	if err != nil {
		return fmt.Errorf("CNI failed to retrieve network namespace path: %v", err)
	}

	if _, err = plugin.addToNetwork(plugin.loNetwork, containerid, netnsPath); err != nil {
		return err
	}

	_, err = plugin.addToNetwork(plugin.getDefaultNetwork(), containerid, netnsPath)
	return err
}

func (plugin *cniNetworkPlugin) RemoveContainerNetwork(containerid string) error {
	if err := plugin.checkInitialized(); err != nil {
		return err
	}

	// Lack of namespace should not be fatal on teardown
	netnsPath, err := plugin.runtime.ContainerNetNS(containerid)
	if err != nil {
		log.Infof("CNI failed to retrieve network namespace path: %v", err)
	}

	return plugin.deleteFromNetwork(plugin.getDefaultNetwork(), containerid, netnsPath, nil)
}

func (plugin *cniNetworkPlugin) addToNetwork(network *cniNetwork, containerID string, netnsPath string) (cnitypes.Result, error) {
	rt, err := plugin.buildCNIRuntimeConf(containerID, netnsPath)
	if err != nil {
		return nil, fmt.Errorf("Error adding network when building cni runtime conf: %v", err)
	}

	netConf, cniNet := network.NetworkConfig, network.CNIConfig
	log.Debugf("Adding %s to network %s/%s netns %q", containerID, netConf.Plugins[0].Network.Type, netConf.Name, netnsPath)
	res, err := cniNet.AddNetworkList(context.Background(), netConf, rt)
	if err != nil {
		return nil, fmt.Errorf("Error adding %s to network %s/%s: %v", containerID, netConf.Plugins[0].Network.Type, netConf.Name, err)
	}

	log.Debugf("Added %s to network %s: %v", containerID, netConf.Name, res)
	return res, nil
}

func (plugin *cniNetworkPlugin) deleteFromNetwork(network *cniNetwork, containerID string, netnsPath string, annotations map[string]string) error {
	rt, err := plugin.buildCNIRuntimeConf(containerID, netnsPath)
	if err != nil {
		return fmt.Errorf("Error deleting network when building cni runtime conf: %v", err)
	}

	netConf, cniNet := network.NetworkConfig, network.CNIConfig
	log.Debugf("Deleting %s from network %s/%s netns %q", containerID, netConf.Plugins[0].Network.Type, netConf.Name, netnsPath)
	err = cniNet.DelNetworkList(context.Background(), netConf, rt)
	// The pod may not get deleted successfully at the first time.
	// Ignore "no such file or directory" error in case the network has already been deleted in previous attempts.
	if err != nil && !strings.Contains(err.Error(), "no such file or directory") {
		return fmt.Errorf("Error deleting %s from network %s/%s: %v", containerID, netConf.Plugins[0].Network.Type, netConf.Name, err)
	}

	log.Debugf("Deleted %s from network %s/%s", containerID, netConf.Plugins[0].Network.Type, netConf.Name)
	return nil
}

func (plugin *cniNetworkPlugin) buildCNIRuntimeConf(containerID string, netnsPath string) (*libcni.RuntimeConf, error) {
	rt := &libcni.RuntimeConf{
		ContainerID: containerID,
		NetNS:       netnsPath,
		IfName:      DefaultInterfaceName,
	}

	return rt, nil
}
