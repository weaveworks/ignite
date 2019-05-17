package container

import (
	"fmt"
	"github.com/luxas/ignite/pkg/constants"
	"net"
	"time"

	dhcp "github.com/krolaw/dhcp4"
	"github.com/krolaw/dhcp4/conn"
)

//var (
//	sourceIPFlag   = flag.String("source-ip", "", "The source IP by the DHCP server")
//	gatewayIPFlag  = flag.String("gateway-ip", "", "The IP of the default gateway")
//	clientIPFlag   = flag.String("client-ip", "", "The IP for the one lease this DHCP server serves for the MAC address specified")
//	clientMACFlag  = flag.String("client-mac", "", "The MAC address that should be recognized and get the IP specified")
//	ifaceFlag      = flag.String("iface", "br0", "The bridge to listen on")
//	subnetMaskFlag = flag.String("subnet-mask", "255.255.0.0", "The subnet for network")
//)

type DHCPInterface struct {
	VMIPNet   *net.IPNet
	GatewayIP *net.IP
	VMTAP     string
	Bridge    string
}

func RunDHCP(dhcpIface *DHCPInterface) error {
	leaseDuration, _ := time.ParseDuration(constants.DHCP_INFINITE_LEASE) // Infinite lease time

	handler := &DHCPHandler{
		gatewayIP:     dhcpIface.GatewayIP.To4(),
		clientIP:      dhcpIface.VMIPNet.IP.To4(),
		dnsServer:     net.IP{1, 1, 1, 1},
		subnetMask:    maskToIP(dhcpIface.VMIPNet.Mask),
		leaseDuration: leaseDuration,
	}

	packetconn, err := conn.NewUDP4BoundListener(dhcpIface.Bridge, ":67")
	if err != nil {
		return err
	}
	return dhcp.Serve(packetconn, handler)
}

type DHCPHandler struct {
	gatewayIP     net.IP
	clientIP      net.IP
	dnsServer     net.IP
	subnetMask    net.IP
	leaseDuration time.Duration
}

func (h *DHCPHandler) ServeDHCP(p dhcp.Packet, msgType dhcp.MessageType, options dhcp.Options) dhcp.Packet {
	var respMsg dhcp.MessageType
	switch msgType {
	case dhcp.Discover: // Just answer Discover calls for that specific MAC address
		respMsg = dhcp.Offer
	case dhcp.Request:
		respMsg = dhcp.ACK
	}
	fmt.Printf("Packet %v, Request: %s, Options: %v, Response: %v\n", p, msgType.String(), options, respMsg.String())
	if respMsg != 0 {
		opts := dhcp.Options{
			dhcp.OptionSubnetMask:       []byte(h.subnetMask),
			dhcp.OptionRouter:           []byte(h.gatewayIP),
			dhcp.OptionDomainNameServer: []byte(h.dnsServer),
		}
		optSlice := opts.SelectOrderOrAll(options[dhcp.OptionParameterRequestList])
		requestingMAC := p.CHAddr().String()
		fmt.Printf("Response: %s, Source %s, Client: %s, Options: %v, MAC: %s\n", respMsg.String(), h.gatewayIP.String(), h.clientIP.String(), optSlice, requestingMAC)
		return dhcp.ReplyPacket(p, respMsg, h.gatewayIP, h.clientIP, h.leaseDuration, optSlice)
	}
	return nil
}

func maskToIP(m net.IPMask) net.IP {
	return net.ParseIP(m.String()).To4()
}
