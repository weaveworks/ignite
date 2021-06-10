package cni

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path"
	"strings"
	"sync"

	gocni "github.com/containerd/go-cni"
	"github.com/containernetworking/plugins/pkg/ip"
	"github.com/containernetworking/plugins/pkg/utils"
	"github.com/coreos/go-iptables/iptables"
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

	// defaultCNIConfFilename is the vanity filename of Ignite's default CNI configuration file
	defaultCNIConfFilename = "10-ignite.conflist"
	// defaultNetworkName names the "docker-bridge"-like CNI plugin-chain installed when no other CNI configuration is present.
	// This value appears in iptables comments created by CNI.
	defaultNetworkName = "ignite-cni-bridge"
	// defaultBridgeName is the default bridge device name used in the defaultCNIConf
	defaultBridgeName = "ignite0"
	// defaultSubnet is the default subnet used in the defaultCNIConf -- this value is set to not collide with common container networking subnets:
	// - 172.{17..31}.0.0/16 and 192.168.({1..15}*16).0/20 are defaults used by `docker network create`.
	// - 10.32.0.0/12 is used with some weavenet CNI installs. (https://github.com/weaveworks/weave/blob/master/site/kubernetes/kube-addon.md)
	// - 10.{42,43}.0.0/16 are used in Rancher CNI installs. (https://rancher.com/docs/rancher/v1.6/en/faqs/troubleshooting/#the-default-subnet-10420016-used-by-rancher-is-already-used-in-my-network-and-prohibiting-the-managed-network-how-do-i-change-the-subnet)
	// - 10.96.0.0/12 is the default kubeadm CNI pod network. (https://godoc.org/k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1beta1#pkg-constants)
	// - 10.244.0.0/16 is the default Flannel CNI pod network. (https://coreos.com/flannel/docs/latest/kubernetes.html, https://github.com/coreos/flannel/blob/b30e689/Documentation/kube-flannel.yml#L125-L131)
	// Avoiding collisions with docker is necessary so ignite CNI networking can function on the same machine as dockerd without routing conflicts.
	// Using the same subnet as another CNI implementation is less consequential. If the other CNI implementation is configured as the default, ignite vm's will just use that network.
	// It's still best to pick a unique, right-sized subnet to avoid confusion and make documentation and issue threads easier to search for.
	// Since a large host could potentially start thousands to tens-of-thousands of firecracker vm's, perhaps a /18, /17, or /16 is appropriate.
	defaultSubnet = "10.61.0.0/16"
)

// defaultCNIConf is a CNI configuration chain that enables VMs to access the internet (docker-bridge style)
var defaultCNIConf = fmt.Sprintf(`{
	"cniVersion": "0.4.0",
	"name": "%s",
	"plugins": [
		{
			"type": "bridge",
			"bridge": "%s",
			"isGateway": true,
			"isDefaultGateway": true,
			"promiscMode": true,
			"ipMasq": true,
			"ipam": {
				"type": "host-local",
				"subnet": "%s"
			}
		},
		{
			"type": "portmap",
			"capabilities": {
				"portMappings": true
			}
		},
		{
			"type": "firewall"
		}
	]
}
`, defaultNetworkName, defaultBridgeName, defaultSubnet)

type cniNetworkPlugin struct {
	cni       gocni.CNI
	cniConfig *gocni.ConfigResult
	runtime   runtime.Interface
	once      *sync.Once
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

	pms := make([]gocni.PortMapping, 0, len(portMappings))
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
		if err = ioutil.WriteFile(path.Join(CNIConfDir, defaultCNIConfFilename), []byte(defaultCNIConf), constants.DATA_DIR_FILE_PERM); err != nil {
			return
		}
	}

	plugin.once.Do(func() {
		if err = plugin.cni.Load(gocni.WithLoNetwork, gocni.WithDefaultConf); err != nil {
			log.Errorf("failed to load cni configuration: %v", err)
		}
	})

	plugin.cniConfig = plugin.cni.GetConfig()

	return
}

func cniToIgniteResult(r *gocni.CNIResult) *network.Result {
	result := &network.Result{}
	for _, iface := range r.Interfaces {
		for _, i := range iface.IPConfigs {
			result.Addresses = append(result.Addresses, network.Address{
				IP:      i.IP,
				Gateway: i.Gateway,
			})
		}
	}

	return result
}

func (plugin *cniNetworkPlugin) RemoveContainerNetwork(containerID string, portMappings ...meta.PortMapping) (err error) {
	if err = plugin.initialize(); err != nil {
		return err
	}

	cleanupErr := plugin.cleanupBridges(containerID)
	if cleanupErr != nil {
		defer util.DeferErr(&err, func() error {
			return cleanupErr
		})
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

	pms := make([]gocni.PortMapping, 0, len(portMappings))
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

	return plugin.cni.Remove(context.Background(), containerID, netnsPath, gocni.WithCapabilityPortMap(pms))
}

// cleanupBridges makes the defaultNetworkName CNI network config not leak iptables rules
// It could possibly help with rule cleanup for other CNI network configs as well
func (plugin *cniNetworkPlugin) cleanupBridges(containerID string) error {
	// Get the amount of combinations between an IP mask, and an iptables chain, with the specified container ID
	result, err := getIPChains(containerID)
	if err != nil {
		return err
	}

	var teardownErrs []error
	for _, net := range plugin.cniConfig.Networks {
		var hasBridge bool
		for _, plugin := range net.Config.Plugins {
			if plugin.Network.Type == "bridge" {
				hasBridge = true
			}
		}

		if hasBridge {
			log.Debugf("Teardown IPMasq for container %q on CNI network %q which contains a bridge", containerID, net.Config.Name)
			comment := utils.FormatComment(net.Config.Name, containerID)
			for _, t := range result {
				if err = ip.TeardownIPMasq(t.ip, t.chain, comment); err != nil {
					teardownErrs = append(teardownErrs, err)
				}
			}
		}
	}

	if len(teardownErrs) == 1 {
		return teardownErrs[0]
	}
	if len(teardownErrs) > 0 {
		return fmt.Errorf("Errors occured cleaning up bridges: %v", teardownErrs)
	}

	return nil
}

type ipChain struct {
	ip    *net.IPNet
	chain string
}

func getIPChains(containerID string) (result []*ipChain, err error) {
	ipt, err := iptables.NewWithProtocol(iptables.ProtocolIPv4)
	if err != nil {
		return
	}

	rawStats, err := ipt.Stats("nat", "POSTROUTING")
	if err != nil {
		return
	}

	quotedContainerID := fmt.Sprintf("id: %q", containerID)
	const statOptionsIndex = 9
	for _, rawStat := range rawStats {
		// stat.Options has a comment that looks like:
		//   /* name: "ignite-cni-bridge" id: "ignite-9a10b07d7c0d4ce9" */
		if strings.Contains(rawStat[statOptionsIndex], quotedContainerID) {
			// only parse the IP's for the rules we need
			// ( avoids https://github.com/coreos/go-iptables/issues/70 )
			var stat iptables.Stat
			stat, err = ipt.ParseStat(rawStat)
			if err != nil {
				return
			}
			result = append(result, &ipChain{
				ip:    stat.Source,
				chain: stat.Target,
			})
		}
	}

	return
}
