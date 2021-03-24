package container

import (
	"fmt"
	"net"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
	"k8s.io/apimachinery/pkg/util/wait"
)

// Array of container interfaces to ignore (not forward to vm)
var ignoreInterfaces = map[string]struct{}{
	"lo": {},
}

type netInterface struct {
	VMIPNet   *net.IPNet
	GatewayIP *net.IP
	VMTAP     string
	MAC       string
}

func SetupContainerNetworking() ([]netInterface, error) {
	var dhcpIfaces []netInterface
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

func networkSetup(dhcpIfaces *[]netInterface) (bool, error) {
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

		ipNet, gw, _, err := takeAddress(&iface)
		if err != nil {
			// Log the problem, but don't quit the function here as there might be other good interfaces
			log.Errorf("Parsing interface %q failed: %v", iface.Name, err)
			// Try with the next interface
			continue
		}

		link, err := netlink.LinkByName(iface.Name)
		if err != nil {
			log.Errorf("Failed to get details of interface %q: %v", iface.Name, err)
			// Try with the next interface
			continue
		}
		log.Infof("interface %q is %q, index %d", iface.Name, link.Type(), link.Attrs().Index)
		//if link.Type() == "macvtap"
		/*_, err = createMacvtapDevice(link)
		if err != nil {
			log.Errorf("Failed to create macvtap device: %v", err)
			// Try with the next interface
			continue
		}*/

		dhcpIface := &netInterface{
			VMIPNet:   ipNet,
			VMTAP:     iface.Name,
			MAC:       iface.HardwareAddr.String(),
			GatewayIP: gw, // important! can be nil
		} //Bridge:    iface.Name, // listen for DHCP here

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

// takeAddress removes the first address of an interface and returns it and the appropriate gateway
func takeAddress(iface *net.Interface) (*net.IPNet, *net.IP, bool, error) {
	addrs, err := iface.Addrs()
	if err != nil || addrs == nil || len(addrs) == 0 {
		// set the bool to true so the caller knows to retry
		return nil, nil, true, fmt.Errorf("interface %q has no address", iface.Name)
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
			return nil, nil, false, fmt.Errorf("failed to get interface %q by name: %v", iface.Name, err)
		}

		var gw *net.IP
		gwString := "<nil>"
		routes, err := netlink.RouteList(link, netlink.FAMILY_ALL)
		if err != nil {
			return nil, nil, false, fmt.Errorf("failed to get default gateway for interface %q: %v", iface.Name, err)
		}
		for _, rt := range routes {
			if rt.Gw != nil {
				gw = &rt.Gw
				gwString = gw.String()
				break
			}
		}

		/* Not deleting address for macvtap testing
		delAddr := &netlink.Addr{
			IPNet: &net.IPNet{
				IP:   ip,
				Mask: mask,
			},
		}
		if err = netlink.AddrDel(link, delAddr); err != nil {
			return nil, nil, false, fmt.Errorf("failed to remove address %q from interface %q: %v", delAddr, iface.Name, err)
		}
		*/

		log.Infof("Moving IP address %s (%s) with gateway %s from container to VM", ip.String(), maskString(mask), gwString)

		return &net.IPNet{
			IP:   ip,
			Mask: mask,
		}, gw, false, nil
	}

	return nil, nil, false, fmt.Errorf("interface %s has no valid addresses", iface.Name)
}

func maskString(mask net.IPMask) string {
	if len(mask) < 4 {
		return "<nil>"
	}

	return fmt.Sprintf("%d.%d.%d.%d", mask[0], mask[1], mask[2], mask[3])
}

/*func createMacvtapDevice(link netlink.Link) (string, error) {
	/*if util.FileExists(devPath) { // don't re-create it if already exists
		return devPath, nil
	}*

	filename := fmt.Sprintf("/sys/devices/virtual/net/%s/macvtap/tap%d/dev", link.Attrs().Name, link.Attrs().Index)
	file, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("Failed to open sys device %q: %v", filename, err)
	}
	var buf [128]byte
	n, err := file.Read(buf[:])
	if err != nil {
		return "", fmt.Errorf("Failed to read from sys device %q: %v", filename, err)
	}
	log.Infof("interface %q is %q", link.Attrs().Name, buf[:n])
	var maj, min uint32
	count, err := fmt.Sscanf(string(buf[:n]), "%d:%d", &maj, &min)
	if err != nil {
		return "", fmt.Errorf("Failed to parse sys device %q: %v", filename, err)
	}
	if count != 2 {
		return "", fmt.Errorf("Failed to extract major/minor from sys device %q", filename)
	}

	devPath := fmt.Sprintf("/dev/net/%s", link.Attrs().Name)
	if err := unix.Mknod(devPath, 0644|syscall.S_IFCHR, int(unix.Mkdev(maj, min))); err != nil {
		log.Errorf("Failed to mknod device %q: %v", devPath, err)
		return devPath, nil
	}

	return devPath, nil
}

// bridge creates the TAP device and performs the bridging, returning the base configuration for a DHCP server
func bridge(iface *net.Interface) (*netInterface, error) {
	tapName := constants.TAP_PREFIX + iface.Name
	bridgeName := constants.BRIDGE_PREFIX + iface.Name

	eth, err := netlink.LinkByIndex(iface.Index)
	if err != nil {
		return nil, errors.Wrap(err, "LinkByIndex")
	}

	tuntap, err := createTAPAdapter(tapName)
	if err != nil {
		return nil, errors.Wrap(err, "createTAPAdapter")
	}

	bridge, err := createBridge(bridgeName)
	if err != nil {
		return nil, errors.Wrap(err, "createBridge")
	}

	if err := setMaster(bridge, tuntap, eth); err != nil {
		return nil, err
	}

	return &netInterface{
		VMTAP:  tapName,
		Bridge: bridgeName,
	}, nil
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

	// Assign a specific mac to the bridge - if we don't do this it will adopt
	// the lowest address of an attached device, hence change over time.
	mac, err := randomMAC()
	if err != nil {
		return nil, errors.Wrap(err, "creating random MAC")
	}
	la.HardwareAddr = mac

	// Disable MAC address age tracking. This causes issues in the container,
	// the bridge is unable to resolve MACs from outside resulting in it never
	// establishing the internal routes. This "optimization" is only really useful
	// with more than two interfaces attached to the bridge anyways, so we're not
	// taking any performance hit by disabling it here.
	ageingTime := uint32(0)
	bridge := &netlink.Bridge{LinkAttrs: la, AgeingTime: &ageingTime}
	return bridge, addLink(bridge)
}

// addLink creates the given link and brings it up
func addLink(link netlink.Link) (err error) {
	if err = netlink.LinkAdd(link); err == nil {
		err = netlink.LinkSetUp(link)
	}

	return
}

func randomMAC() (net.HardwareAddr, error) {
	mac := make([]byte, 6)
	if _, err := rand.Read(mac); err != nil {
		return nil, err
	}

	// In the first byte of the MAC, the 'multicast' bit should be
	// clear and 'locally administered' bit should be set.
	mac[0] = (mac[0] & 0xFE) | 0x02

	return net.HardwareAddr(mac), nil
}

func setMaster(master netlink.Link, links ...netlink.Link) error {
	masterIndex := master.Attrs().Index

	for _, link := range links {
		if err := netlink.LinkSetMasterByIndex(link, masterIndex); err != nil {
			return errors.Wrapf(err, "setMaster %s %s", master.Attrs().Name, link.Attrs().Name)
		}
	}

	return nil
}*/
