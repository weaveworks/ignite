package container

import (
	"fmt"
	"net"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
	"github.com/weaveworks/ignite/pkg/constants"
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
var ignoreInterfaces = map[string]struct{}{
	"lo": {},
}

func SetupContainerNetworking() ([]DHCPInterface, error) {
	var dhcpIfaces []DHCPInterface
	interval := 1 * time.Second
	timeout := 1 * time.Minute

	err := wait.PollImmediate(interval, timeout, func() (bool, error) {
		// This func returns true if it's done, and optionally an error
		retry, err := networkSetup(&dhcpIfaces)
		if err == nil {
			// We're done here
			return true, nil
		}
		if retry {
			// We got an error, but let's ignore it and try again
			log.Warnf("Got an error while trying to set up networking, but retrying: %v", err)
			return false, nil
		}
		// The error was fatal, return it
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
		if _, ok := ignoreInterfaces[iface.Name]; ok {
			continue
		}

		// Try to transfer the address from the container to the DHCP server
		ipNet, _, err := takeAddress(&iface)
		if err != nil {
			// Log the problem, but don't quit the function here as there might be other good interfaces
			log.Errorf("Parsing interface %q failed: %v", iface.Name, err)
			// Try with the next interface
			continue
		}

		// Bridge the Firecracker TAP interface with the container veth interface
		dhcpIface, err := bridge(&iface)
		if err != nil {
			// Log the problem, but don't quit the function here as there might be other good interfaces
			// Don't set shouldRetry here as there is no point really with retrying with this interface
			// that seems broken/unsupported in some way.
			log.Errorf("Bridging interface %q failed: %v", iface.Name, err)
			// Try with the next interface
			continue
		}

		// Gateway for now is just x.x.x.1 TODO: Better detection
		dhcpIface.GatewayIP = &net.IP{ipNet.IP[0], ipNet.IP[1], ipNet.IP[2], 1}
		dhcpIface.VMIPNet = ipNet

		*dhcpIfaces = append(*dhcpIfaces, *dhcpIface)

		// This is an interface we care about
		interfacesCount++
	}

	// If there weren't any interfaces that were valid or active yet, retry the loop
	if interfacesCount == 0 {
		return true, fmt.Errorf("no active or valid interfaces available yet")
	}

	return false, nil
}

// bridge creates the TAP device and performs the bridging, returning the base configuration for a DHCP server
func bridge(iface *net.Interface) (d *DHCPInterface, err error) {
	tapName := constants.TAP_PREFIX + iface.Name
	bridgeName := constants.BRIDGE_PREFIX + iface.Name

	var bridge *netlink.Bridge
	var eth, tuntap netlink.Link

	if eth, err = netlink.LinkByIndex(iface.Index); err != nil {
		return
	}

	if tuntap, err = createTAPAdapter(tapName); err != nil {
		return
	}

	if bridge, err = createBridge(bridgeName); err != nil {
		return
	}

	if err = setMaster(bridge, tuntap, eth); err != nil {
		return
	}

	d = &DHCPInterface{
		VMTAP:  tapName,
		Bridge: bridgeName,
	}

	return
}

// takeAddress removes the first address of an interface and returns it
func takeAddress(iface *net.Interface) (*net.IPNet, bool, error) {
	addrs, err := iface.Addrs()
	if err != nil || addrs == nil || len(addrs) == 0 {
		// set the bool to true so the caller knows to retry
		return nil, true, fmt.Errorf("interface %q has no address", iface.Name)
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

		link, err := netlink.LinkByName(iface.Name)
		if err != nil {
			return nil, false, fmt.Errorf("failed to get interface %q by name: %v", iface.Name, err)
		}

		delAddr, err := netlink.ParseAddr(addr.String())
		if err != nil {
			return nil, false, fmt.Errorf("failed to parse address from stringified IP %q: %v", addr.String(), err)
		}

		if err = netlink.AddrDel(link, delAddr); err != nil {
			return nil, false, fmt.Errorf("failed to remove address from interface %q: %v", iface.Name, err)
		}

		log.Infof("Moving IP address %s (%s) from container to VM", ip.String(), maskString(mask))

		return &net.IPNet{
			IP:   ip,
			Mask: mask,
		}, false, nil
	}

	return nil, false, fmt.Errorf("interface %s has no valid addresses", iface.Name)
}

// createTAPAdapter creates a new TAP device with the given name
func createTAPAdapter(tapName string) (*netlink.Tuntap, error) {
	la := netlink.NewLinkAttrs()
	la.Name = tapName
	tuntap := &netlink.Tuntap{
		LinkAttrs: la,
		Mode:      netlink.TUNTAP_MODE_TAP,
	}
	return tuntap, addLink(tuntap)
}

// createBridge creates a new bridge device with the given name
func createBridge(bridgeName string) (*netlink.Bridge, error) {
	la := netlink.NewLinkAttrs()
	la.Name = bridgeName
	bridge := &netlink.Bridge{LinkAttrs: la}
	return bridge, addLink(bridge)
}

// addLink creates the given link and brings it up
func addLink(link netlink.Link) (err error) {
	if err = netlink.LinkAdd(link); err == nil {
		err = netlink.LinkSetUp(link)
	}

	return
}

// This is a MAC address persistence workaround, netlink.LinkSetMaster{,ByIndex}()
// has a bug that arbitrarily changes the MAC addresses of the bridge and virtual
// device to be bound to it. TODO: Remove when fixed upstream
func setMaster(master netlink.Link, links ...netlink.Link) error {
	masterIndex := master.Attrs().Index
	masterMAC, err := getMAC(master)
	if err != nil {
		return err
	}

	for _, link := range links {
		mac, err := getMAC(link)
		if err != nil {
			return err
		}

		if err = netlink.LinkSetMasterByIndex(link, masterIndex); err != nil {
			return err
		}

		if err = netlink.LinkSetHardwareAddr(link, mac); err != nil {
			return err
		}
	}

	return netlink.LinkSetHardwareAddr(master, masterMAC)
}

// getMAC fetches the generated MAC address for the given link
func getMAC(link netlink.Link) (addr net.HardwareAddr, err error) {
	if link, err = netlink.LinkByIndex(link.Attrs().Index); err != nil {
		return
	}

	addr = link.Attrs().HardwareAddr
	return
}

func maskString(mask net.IPMask) string {
	if len(mask) < 4 {
		return "<nil>"
	}

	return fmt.Sprintf("%d.%d.%d.%d", mask[0], mask[1], mask[2], mask[3])
}
