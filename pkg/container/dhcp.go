package container

import (
	"fmt"
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

//func RunDHCP(gatewayIP, clientIP, subnetMask net.IP, clientMAC net.HardwareAddr, leaseDuration time.Duration, iface string) error {
func RunDHCP(gatewayIP, clientIP, subnetMask net.IP, leaseDuration time.Duration, iface string) error {
	/*sourceIP := net.ParseIP(*sourceIPFlag)
	if sourceIP == nil {
		return fmt.Errorf("--source-ip is invalid")
	}*/

	handler := &DHCPHandler{
		//sourceIP:      sourceIP.To4(),
		gatewayIP: gatewayIP.To4(),
		clientIP:  clientIP.To4(),
		//clientMAC:     clientMAC,
		dnsServer:     net.IP{1, 1, 1, 1},
		subnetMask:    subnetMask.To4(),
		leaseDuration: leaseDuration,
	}
	packetconn, err := conn.NewUDP4BoundListener(iface, ":67")
	if err != nil {
		return err
	}
	return dhcp.Serve(packetconn, handler)
}

type DHCPHandler struct {
	sourceIP  net.IP
	gatewayIP net.IP
	clientIP  net.IP
	//clientMAC     net.HardwareAddr
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
		//if requestingMAC == h.clientMAC.String() {
		fmt.Printf("Response: %s, Source %s, Client: %s, Options: %v, MAC: %s\n", respMsg.String(), h.gatewayIP.String(), h.clientIP.String(), optSlice, requestingMAC)
		return dhcp.ReplyPacket(p, respMsg, h.gatewayIP, h.clientIP, h.leaseDuration, optSlice)
		//}
	}
	return nil
}
