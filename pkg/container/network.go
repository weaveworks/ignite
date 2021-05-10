package container

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
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
var ignoreInterfaces = map[string]struct{}{
	"lo": {},
}

var (
	maxIntfsVar     = "IGNITE_INTFS"
	defaultMaxIntfs = 1
)

type FCInterface struct {
	VMTAP     string
	MACFilter string
}

func SetupContainerNetworking(args map[string]string) ([]FCInterface, []DHCPInterface, error) {
	var dhcpIfaces []DHCPInterface
	var allIfaces []FCInterface

	interval := 1 * time.Second
	timeout := 1 * time.Minute
	maxIntfs, err := strconv.Atoi(args[maxIntfsVar])
	if err != nil {
		maxIntfs = defaultMaxIntfs
	}

	err = wait.PollImmediate(interval, timeout, func() (bool, error) {
		// This func returns true if it's done, and optionally an error
		retry, err := networkSetup(&allIfaces, &dhcpIfaces, maxIntfs)
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
		return nil, nil, err
	}

	return allIfaces, dhcpIfaces, nil
}

func isIgnored(link net.Interface, allIfaces *[]FCInterface) bool {
	// ignore if in explicit ignore list
	if _, ok := ignoreInterfaces[link.Name]; ok {
		return true
	}

	// ignore if _not_ a veth
	eth, err := netlink.LinkByIndex(link.Index)
	if err != nil {
		return true
	}

	_, ok := eth.(*netlink.Veth)

	return !ok
}

func networkSetup(allIfaces *[]FCInterface, dhcpIfaces *[]DHCPInterface, intfNum int) (bool, error) {
	ifaces, err := net.Interfaces()
	if err != nil || ifaces == nil || len(ifaces) == 0 {
		return true, fmt.Errorf("cannot get local network interfaces: %v", err)
	}

	// interfacesCount counts the interfaces that are relevant to Ignite (in other words, not ignored)
	interfacesCount := 0
	for _, iface := range ifaces {
		// Skip the interface if it's ignored
		if isIgnored(iface, allIfaces) {
			continue
		}

		// Try to transfer the address from the container to the DHCP server
		ipNet, gw, noIPs, err := takeAddress(&iface)

		// If interface has no IPs configured, setup tc redirect
		if noIPs {
			tcInterface, err := addTcRedirect(&iface)
			if err != nil {
				log.Errorf("Failed to setup tc redirect %v", err)
				continue
			}

			*allIfaces = append(*allIfaces, *tcInterface)
			ignoreInterfaces[iface.Name] = struct{}{}

			continue
		}
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

		dhcpIface.VMIPNet = ipNet
		dhcpIface.GatewayIP = gw

		*dhcpIfaces = append(*dhcpIfaces, *dhcpIface)

		*allIfaces = append(*allIfaces, FCInterface{
			VMTAP:     dhcpIface.VMTAP,
			MACFilter: dhcpIface.MACFilter,
		})
		ignoreInterfaces[iface.Name] = struct{}{}

		// This is an interface we care about
		interfacesCount++
	}

	// If there weren't any interfaces that were valid or active yet, retry the loop
	if interfacesCount < intfNum {
		return true, fmt.Errorf("not enough active or valid interfaces available yet")
	}

	return false, nil
}

// addTcRedirect sets up tc redirect betweeb veth and tap https://github.com/awslabs/tc-redirect-tap/blob/master/internal/netlink.go
// on WSL2 this requires `CONFIG_NET_CLS_U32=y`
func addTcRedirect(iface *net.Interface) (*FCInterface, error) {

	eth, err := netlink.LinkByIndex(iface.Index)
	if err != nil {
		return nil, err
	}

	tapName := constants.TAP_PREFIX + iface.Name
	tuntap, err := createTAPAdapter(tapName)
	if err != nil {
		return nil, err
	}

	err = addIngressQdisc(eth)
	if err != nil {
		return nil, err
	}
	err = addIngressQdisc(tuntap)
	if err != nil {
		return nil, err
	}

	err = addRedirectFilter(eth, tuntap)
	if err != nil {
		return nil, err
	}

	log.Printf("Adding tc filter to tuntap %s", tuntap.Attrs().Name)
	err = addRedirectFilter(tuntap, eth)
	if err != nil {
		return nil, err
	}

	return &FCInterface{
		VMTAP:     tapName,
		MACFilter: eth.Attrs().HardwareAddr.String(),
	}, nil
}

// tc qdisc add dev $SRC_IFACE ingress
func addIngressQdisc(link netlink.Link) error {
	qdisc := &netlink.Ingress{
		QdiscAttrs: netlink.QdiscAttrs{
			LinkIndex: link.Attrs().Index,
			Parent:    netlink.HANDLE_INGRESS,
		},
	}

	if err := netlink.QdiscAdd(qdisc); err != nil {
		if !os.IsExist(err) {
			return err
		}
	}
	return nil
}

// tc filter add dev $SRC_IFACE parent ffff:
// protocol all
// u32 match u32 0 0
// action mirred egress mirror dev $DST_IFACE
func addRedirectFilter(linkSrc, linkDest netlink.Link) error {
	filter := &netlink.U32{
		FilterAttrs: netlink.FilterAttrs{
			LinkIndex: linkSrc.Attrs().Index,
			Parent:    netlink.MakeHandle(0xffff, 0),
			Protocol:  syscall.ETH_P_ALL,
		},
		Actions: []netlink.Action{
			&netlink.MirredAction{
				ActionAttrs: netlink.ActionAttrs{
					Action: netlink.TC_ACT_STOLEN,
				},
				MirredAction: netlink.TCA_EGRESS_MIRROR,
				Ifindex:      linkDest.Attrs().Index,
			},
		},
	}
	return netlink.FilterAdd(filter)
}

// bridge creates the TAP device and performs the bridging, returning the base configuration for a DHCP server
func bridge(iface *net.Interface) (*DHCPInterface, error) {
	tapName := constants.TAP_PREFIX + iface.Name
	bridgeName := constants.BRIDGE_PREFIX + iface.Name

	eth, err := netlink.LinkByIndex(iface.Index)
	if err != nil {
		return nil, err
	}

	tuntap, err := createTAPAdapter(tapName)
	if err != nil {
		return nil, err
	}

	bridge, err := createBridge(bridgeName)
	if err != nil {
		return nil, err
	}

	if err := setMaster(bridge, tuntap, eth); err != nil {
		return nil, err
	}

	// Generate the MAC addresses for the VM's adapters
	macAddress := make([]string, 0, 1)
	if err := util.NewMAC(&macAddress); err != nil {
		return nil, fmt.Errorf("failed to generate MAC addresses: %v", err)
	}

	return &DHCPInterface{
		VMTAP:     tapName,
		Bridge:    bridgeName,
		MACFilter: macAddress[0],
	}, nil
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

		delAddr := &netlink.Addr{
			IPNet: &net.IPNet{
				IP:   ip,
				Mask: mask,
			},
		}
		if err = netlink.AddrDel(link, delAddr); err != nil {
			return nil, nil, false, fmt.Errorf("failed to remove address %q from interface %q: %v", delAddr, iface.Name, err)
		}

		log.Infof("Moving IP address %s (%s) with gateway %s from container to VM", ip.String(), maskString(mask), gwString)

		return &net.IPNet{
			IP:   ip,
			Mask: mask,
		}, gw, false, nil
	}

	return nil, nil, true, fmt.Errorf("interface %s has no valid addresses", iface.Name)
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
	// The attributes of the netlink.Link passed to this function do not contain HardwareAddr
	// as it is expected to be generated by the networking subsystem. Thus, "reload" the Link
	// by querying it to retrieve the generated attributes after the link has been created.
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
