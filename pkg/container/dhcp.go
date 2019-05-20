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

var leaseDuration, _ = time.ParseDuration(constants.DHCP_INFINITE_LEASE) // Infinite lease time

type DHCPInterface struct {
	VMIPNet    *net.IPNet
	GatewayIP  *net.IP
	VMTAP      string
	Bridge     string
	dnsServers []byte
}

func RunDHCP(dhcpIface *DHCPInterface) error {
	packetConn, err := conn.NewUDP4BoundListener(dhcpIface.Bridge, ":67")
	if err != nil {
		return err
	}
	return dhcp.Serve(packetConn, dhcpIface)
}

func (i *DHCPInterface) ServeDHCP(p dhcp.Packet, msgType dhcp.MessageType, options dhcp.Options) dhcp.Packet {
	var respMsg dhcp.MessageType
	switch msgType {
	case dhcp.Discover:
		respMsg = dhcp.Offer
	case dhcp.Request:
		respMsg = dhcp.ACK
	}
	fmt.Printf("Packet %v, Request: %s, Options: %v, Response: %v\n", p, msgType.String(), options, respMsg.String())
	if respMsg != 0 {
		opts := dhcp.Options{
			dhcp.OptionSubnetMask:       []byte(i.VMIPNet.Mask),
			dhcp.OptionRouter:           []byte(*i.GatewayIP),
			dhcp.OptionDomainNameServer: i.dnsServers,
		}
		optSlice := opts.SelectOrderOrAll(options[dhcp.OptionParameterRequestList])
		requestingMAC := p.CHAddr().String()
		fmt.Printf("Response: %s, Source %s, Client: %s, Options: %v, MAC: %s\n", respMsg.String(), i.GatewayIP.String(), i.VMIPNet.IP.String(), optSlice, requestingMAC)
		return dhcp.ReplyPacket(p, respMsg, *i.GatewayIP, i.VMIPNet.IP, leaseDuration, optSlice)
	}
	return nil
}

// Parse the DNS servers for the DHCP server
func (i *DHCPInterface) SetDNSServers(dns []string) {
	for _, server := range dns {
		i.dnsServers = append(i.dnsServers, []byte(net.ParseIP(server).To4())...)
	}
}
