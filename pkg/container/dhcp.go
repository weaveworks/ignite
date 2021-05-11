package container

import (
	"fmt"
	"net"
	"time"

	dhcp "github.com/krolaw/dhcp4"
	"github.com/krolaw/dhcp4/conn"
	"github.com/miekg/dns"
	log "github.com/sirupsen/logrus"
	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	"github.com/weaveworks/ignite/pkg/constants"
)

var leaseDuration, _ = time.ParseDuration(constants.DHCP_INFINITE_LEASE) // Infinite lease time

// StartDHCPServers starts multiple DHCP servers for the VM, one per interface
// It returns the IP addresses that the API object may post in .status, and a potential error
func StartDHCPServers(vm *api.VM, dhcpIfaces []DHCPInterface) error {

	// Fetch the DNS servers given to the container
	clientConfig, err := dns.ClientConfigFromFile("/etc/resolv.conf")
	if err != nil {
		return fmt.Errorf("failed to get DNS configuration: %v", err)
	}

	for i := range dhcpIfaces {
		dhcpIface := &dhcpIfaces[i]
		// Set the VM hostname to the VM ID
		dhcpIface.Hostname = vm.GetUID().String()

		// Add the DNS servers from the container
		dhcpIface.SetDNSServers(clientConfig.Servers)

		go func() {
			log.Infof("Starting DHCP server for interface %q (%s)\n", dhcpIface.Bridge, dhcpIface.VMIPNet.IP)
			if err := dhcpIface.StartBlockingServer(); err != nil {
				log.Errorf("%q DHCP server error: %v\n", dhcpIface.Bridge, err)
			}
		}()
	}

	return nil
}

type DHCPInterface struct {
	VMIPNet    *net.IPNet
	GatewayIP  *net.IP
	VMTAP      string
	Bridge     string
	Hostname   string
	MACFilter  string
	dnsServers []byte
}

// StartBlockingServer starts a blocking DHCP server on port 67
func (i *DHCPInterface) StartBlockingServer() error {
	packetConn, err := conn.NewUDP4BoundListener(i.Bridge, ":67")
	if err != nil {
		return err
	}

	return dhcp.Serve(packetConn, i)
}

// ServeDHCP responds to a DHCP request
func (i *DHCPInterface) ServeDHCP(p dhcp.Packet, msgType dhcp.MessageType, options dhcp.Options) dhcp.Packet {
	var respMsg dhcp.MessageType

	switch msgType {
	case dhcp.Discover:
		respMsg = dhcp.Offer
	case dhcp.Request:
		respMsg = dhcp.ACK
	}

	//fmt.Printf("Packet %v, Request: %s, Options: %v, Response: %v\n", p, msgType.String(), options, respMsg.String())
	if respMsg != 0 {
		requestingMAC := p.CHAddr().String()
		if requestingMAC == i.MACFilter {
			opts := dhcp.Options{
				dhcp.OptionSubnetMask:       []byte(i.VMIPNet.Mask),
				dhcp.OptionRouter:           []byte(*i.GatewayIP),
				dhcp.OptionDomainNameServer: i.dnsServers,
				dhcp.OptionHostName:         []byte(i.Hostname),
			}

			optSlice := opts.SelectOrderOrAll(options[dhcp.OptionParameterRequestList])
			//fmt.Printf("Response: %s, Source %s, Client: %s, Options: %v, MAC: %s\n", respMsg.String(), i.GatewayIP.String(), i.VMIPNet.IP.String(), optSlice, requestingMAC)
			return dhcp.ReplyPacket(p, respMsg, *i.GatewayIP, i.VMIPNet.IP, leaseDuration, optSlice)
		}
	}

	return nil
}

// Parse the DNS servers for the DHCP server
func (i *DHCPInterface) SetDNSServers(dns []string) {
	for _, server := range dns {
		i.dnsServers = append(i.dnsServers, []byte(net.ParseIP(server).To4())...)
	}
}
