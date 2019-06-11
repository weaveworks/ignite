package run

import (
	"fmt"
	"os"

	"github.com/weaveworks/ignite/cmd/ignite/run/runutil"

	"github.com/miekg/dns"
	"github.com/weaveworks/ignite/pkg/container"
	"github.com/weaveworks/ignite/pkg/metadata/vmmd"
	"github.com/weaveworks/ignite/pkg/util"
)

type containerOptions struct {
	vm *vmmd.VMMetadata
}

func NewContainerOptions(l *runutil.ResLoader, vmMatch string) (*containerOptions, error) {
	co := &containerOptions{}

	if allVMS, err := l.VMs(); err == nil {
		if co.vm, err = allVMS.MatchSingle(vmMatch); err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}

	return co, nil
}

func Container(co *containerOptions) error {
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
		dhcpIface.Hostname = co.vm.ID.String()

		// Set the MAC address filter for the DHCP server
		dhcpIface.MACFilter = macAddresses[i]

		// Add the DNS servers from the container
		dhcpIface.SetDNSServers(clientConfig.Servers)

		co.vm.AddIPAddress(dhcpIface.VMIPNet.IP)

		go func() {
			fmt.Printf("Starting DHCP server for interface %s (%s)\n", dhcpIface.Bridge, dhcpIface.VMIPNet.IP)
			if err := container.RunDHCP(dhcpIface); err != nil {
				fmt.Fprintf(os.Stderr, "%s DHCP server error: %v\n", dhcpIface.Bridge, err)
			}
		}()
	}

	// VM state handling
	if err := co.vm.SetState(vmmd.Running); err != nil {
		return fmt.Errorf("failed to update VM state: %v", err)
	}
	defer co.vm.SetState(vmmd.Stopped) // Performs a save, all other metadata-modifying defers need to be after this

	// Remove the snapshot overlay post-run, which also removes the detached backing loop devices
	defer co.vm.RemoveSnapshot()

	// Remove the IP addresses post-run
	defer co.vm.ClearIPAddresses()

	// Remove the port mappings post-run
	defer co.vm.ClearPortMappings()

	// Run the vm
	if err := container.RunVM(co.vm, &dhcpIfaces); err != nil {
		return fmt.Errorf("runtime error for VM %q: %v", co.vm.ID, err)
	}

	return nil
}
