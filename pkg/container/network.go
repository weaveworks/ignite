package container

import (
	"fmt"
	"net"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/util"

	"k8s.io/apimachinery/pkg/util/wait"
)

/*
ip r list src 172.17.0.3

ip addr del "$IP" dev eth0

ip link add name br0 type bridge
ip tuntap add dev vm0 mode tap

ip link set br0 up
ip link set vm0 up

ip link set eth0 master br0
ip link set vm0 master br0
*/

// Array of container interfaces to ignore (not forward to vm)
var ignoreInterfaces = map[string]bool{
	"lo": true,
}

func SetupContainerNetworking() ([]DHCPInterface, error) {
	var dhcpIfaces []DHCPInterface
	interval := 1 * time.Second
	timeout := 1 * time.Minute
	err := wait.PollImmediate(interval, timeout, func() (bool, error) {
		// this func returns true if it's done, and optionally an error
		retry, err := networkSetup(&dhcpIfaces)
		if err == nil {
			// we're done here
			return true, nil
		}
		if retry {
			// we got an error, but let's ignore it and try again
			return false, nil
		}
		// the error was fatal, return it
		return false, err
	})
	if err != nil {
		return nil, err
	}
	return dhcpIfaces, nil
}

func networkSetup(dhcpIfaces *[]DHCPInterface) (bool, error) {
	ifaces, err := net.Interfaces()
	if err != nil || ifaces == nil || len(ifaces) == 0 {
		return true, fmt.Errorf("cannot get local network interfaces: %v", err)
	}

	// interfacesCount counts the interfaces that are relevant to Ignite (in other words, not ignored)
	interfacesCount := 0
	for _, iface := range ifaces {
		// Skip the interface if it's ignored
		if !ignoreInterfaces[iface.Name] {
			// This is an interface we care about
			interfacesCount++

			ipNet, retry, err := takeAddress(&iface)
			if err != nil {
				return retry, fmt.Errorf("parsing interface failed: %v", err)
			}

			dhcpIface, err := bridge(&iface)
			if err != nil {
				return false, fmt.Errorf("bridging interface %s failed: %v0", iface.Name, err)
			}

			// Gateway for now is just x.x.x.1 TODO: Better detection
			dhcpIface.GatewayIP = &net.IP{ipNet.IP[0], ipNet.IP[1], ipNet.IP[2], 1}
			dhcpIface.VMIPNet = ipNet

			*dhcpIfaces = append(*dhcpIfaces, *dhcpIface)
		}
	}

	// If there weren't any interfaces we cared about, retry the loop
	if interfacesCount == 0 {
		return true, fmt.Errorf("no active interfaces available yet")
	}

	return false, nil
}

// bridge creates the TAP device and performs the bridging, returning the MAC address of the vm's adapter
func bridge(iface *net.Interface) (*DHCPInterface, error) {
	tapName := constants.TAP_PREFIX + iface.Name
	bridgeName := constants.BRIDGE_PREFIX + iface.Name

	if err := createTAPAdapter(tapName); err != nil {
		return nil, err
	}

	if err := createBridge(bridgeName); err != nil {
		return nil, err
	}

	if err := connectAdapterToBridge(tapName, bridgeName); err != nil {
		return nil, err
	}

	if err := connectAdapterToBridge(iface.Name, bridgeName); err != nil {
		return nil, err
	}

	return &DHCPInterface{
		VMTAP:  tapName,
		Bridge: bridgeName,
	}, nil
}

// takeAddress removes the first address of an interface and returns it
func takeAddress(iface *net.Interface) (*net.IPNet, bool, error) {
	addrs, err := iface.Addrs()
	if err != nil || addrs == nil || len(addrs) == 0 {
		// set the bool to true so the caller knows to retry
		return nil, true, fmt.Errorf("interface %s has no address", iface.Name)
	}

	for _, addr := range addrs {
		var ip net.IP
		var mask net.IPMask

		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
			mask = v.Mask
		case *net.IPAddr:
			ip = v.IP
			mask = ip.DefaultMask()
		}

		if ip == nil {
			continue
		}

		ip = ip.To4()
		if ip == nil {
			continue
		}

		if _, err := util.ExecuteCommand("ip", "addr", "del", ip.String(), "dev", iface.Name); err != nil {
			return nil, false, errors.Wrapf(err, "failed to remove address from interface %s", iface.Name)
		}

		log.Infof("Moving IP address %s (%s) from container to VM\n", ip.String(), fmt.Sprintf("%d.%d.%d.%d", mask[0], mask[1], mask[2], mask[3]))

		return &net.IPNet{
			IP:   ip,
			Mask: mask,
		}, false, nil
	}

	return nil, false, fmt.Errorf("interface %s has no valid addresses", iface.Name)
}

func createTAPAdapter(tapName string) error {
	if _, err := util.ExecuteCommand("ip", "tuntap", "add", "mode", "tap", tapName); err != nil {
		return err
	}

	return setLinkUp(tapName)
}

func createBridge(bridgeName string) error {
	if _, err := util.ExecuteCommand("ip", "link", "add", "name", bridgeName, "type", "bridge"); err != nil {
		return err
	}

	return setLinkUp(bridgeName)
}

func setLinkUp(adapterName string) error {
	_, err := util.ExecuteCommand("ip", "link", "set", adapterName, "up")
	return err
}

func connectAdapterToBridge(adapterName, bridgeName string) error {
	_, err := util.ExecuteCommand("ip", "link", "set", adapterName, "master", bridgeName)
	return err
}
