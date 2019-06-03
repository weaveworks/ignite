package run

import (
	"fmt"
	"github.com/luxas/ignite/pkg/container"
	"github.com/luxas/ignite/pkg/metadata/vmmd"
	"github.com/luxas/ignite/pkg/util"
	"github.com/miekg/dns"
	"os"
)

type ContainerOptions struct {
	VM *vmmd.VMMetadata
}

func Container(co *ContainerOptions) error {
	var dhcpIfaces []container.DHCPInterface

	// New networking setup
	if err := container.NetworkSetup(&dhcpIfaces); err != nil {
		return fmt.Errorf("network setup failed: %v", err)
	}

	// Generate the MAC addresses for the VM's adapters
	macAddresses := make([]string, 0, len(dhcpIfaces))
	if err := util.NewMAC(&macAddresses); err != nil {
		return fmt.Errorf("failed to generate MAC addresses: %v", err)
	}

	// Fetch the DNS servers given to the container
	clientConfig, err := dns.ClientConfigFromFile("/etc/resolv.conf")
	if err != nil {
		return fmt.Errorf("failed to get DNS configuration: %v", err)
	}

	for i := range dhcpIfaces {
		dhcpIface := &dhcpIfaces[i]
		// Set the VM hostname to the VM ID
		dhcpIface.Hostname = co.VM.ID

		// Set the MAC address filter for the DHCP server
		dhcpIface.MACFilter = macAddresses[i]

		// Add the DNS servers from the container
		dhcpIface.SetDNSServers(clientConfig.Servers)

		co.VM.AddIPAddress(dhcpIface.VMIPNet.IP)

		go func() {
			fmt.Printf("Starting DHCP server for interface %s (%s)\n", dhcpIface.Bridge, dhcpIface.VMIPNet.IP)
			if err := container.RunDHCP(dhcpIface); err != nil {
				fmt.Fprintf(os.Stderr, "%s DHCP server error: %v\n", dhcpIface.Bridge, err)
			}
		}()
	}

	// Remove the IP addresses post-run
	defer co.VM.ClearIPAddresses()

	// Remove the snapshot overlay post-run, which also removes the detached backing loop devices
	defer co.VM.RemoveSnapshot()

	// VM state handling
	if err := co.VM.SetState(vmmd.Running); err != nil {
		return fmt.Errorf("failed to update VM state: %v", err)
	}
	defer co.VM.SetState(vmmd.Stopped)

	// Run the VM
	if err := container.RunVM(co.VM, &dhcpIfaces); err != nil {
		return fmt.Errorf("runtime error for VM %q: %v", co.VM.ID, err)
	}

	return nil
}
