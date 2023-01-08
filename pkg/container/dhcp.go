package container

import (
	"context"
	"fmt"
	"net"
	"time"

	dhcp "github.com/krolaw/dhcp4"
	"github.com/krolaw/dhcp4/conn"
	"github.com/miekg/dns"
	log "github.com/sirupsen/logrus"
	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	"github.com/weaveworks/ignite/pkg/constants"
	internalv6 "github.com/weaveworks/ignite/pkg/container/dhcpv6"
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
			log.Infof("Starting DHCPv4 server for interface %q (%s)\n", dhcpIface.Bridge, dhcpIface.VMIPNet.IP)
			if err := dhcpIface.StartBlockingServerV4(); err != nil {
				log.Errorf("%q DHCPv4 server error: %v\n", dhcpIface.Bridge, err)
			}
		}()

		go func() {
			log.Infof("Starting DHCPv6 server for interface %q (%s)\n", dhcpIface.Bridge, dhcpIface.VMIPNet.IP)
			if err := dhcpIface.StartBlockingServerV6(); err != nil {
				log.Errorf("%q DHCPv6 server error: %v\n", dhcpIface.Bridge, err)
			}
		}()

	}

	return nil
}

type DHCPInterface struct {
	VMIPNet      *net.IPNet
	GatewayIP    *net.IP
	VMTAP        string
	Bridge       string
	Hostname     string
	MACFilter    string
	dnsv4Servers []byte
	dnsv6Servers []byte
}

// StartBlockingServerV4 starts a blocking DHCPv4 server on port 67
func (i *DHCPInterface) StartBlockingServerV4() error {
	packetConn, err := conn.NewUDP4BoundListener(i.Bridge, ":67")
	if err != nil {
		return err
	}

	return dhcp.Serve(packetConn, i)
}

// StartBlockingServerV6 starts a blocking DHCPv6 server on port 547
func (i *DHCPInterface) StartBlockingServerV6() error {

	interf, err := net.InterfaceByName(i.Bridge)

	if err != nil {
		return err
	}

	c := internalv6.Config{Interface: interf}
	return internalv6.RunDHCPv6Server(context.TODO(), log.StandardLogger(), c)
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

	log.Debugf("Packet %v, Request: %s, Options: %v, Response: %v\n", p, msgType.String(), options, respMsg.String())
	if respMsg != 0 {
		requestingMAC := p.CHAddr().String()
		if requestingMAC == i.MACFilter {
			opts := dhcp.Options{
				dhcp.OptionSubnetMask:       []byte(i.VMIPNet.Mask),
				dhcp.OptionRouter:           []byte(*i.GatewayIP),
				dhcp.OptionDomainNameServer: i.dnsv4Servers,
				dhcp.OptionHostName:         []byte(i.Hostname),
			}

			optSlice := opts.SelectOrderOrAll(options[dhcp.OptionParameterRequestList])
			log.Debugf("Response: %s, Source %s, Client: %s, Options: %v, MAC: %s\n", respMsg.String(), i.GatewayIP.String(), i.VMIPNet.IP.String(), optSlice, requestingMAC)
			return dhcp.ReplyPacket(p, respMsg, *i.GatewayIP, i.VMIPNet.IP, leaseDuration, optSlice)
		}
	}

	return nil
}

// Parse the DNS servers for the DHCP server
func (i *DHCPInterface) SetDNSServers(dns []string) {
	for _, server := range dns {

		if ipv4 := net.ParseIP(server).To4(); ipv4 != nil {
			i.dnsv4Servers = append(i.dnsv4Servers, ipv4...)
		} else if ipv6 := net.ParseIP(server).To16(); ipv6 != nil {
			i.dnsv6Servers = append(i.dnsv4Servers, ipv6...)
		} else {
			log.Errorf("failed to parse dns %v as either ipv4 or ipv6", server)
		}

	}
}
