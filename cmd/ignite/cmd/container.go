package cmd

import (
	"fmt"
	"github.com/luxas/ignite/pkg/container"
	"github.com/luxas/ignite/pkg/errutils"
	"github.com/luxas/ignite/pkg/metadata/vmmd"
	"github.com/luxas/ignite/pkg/util"
	"github.com/miekg/dns"
	"github.com/spf13/cobra"
	"io"
	"os"
)

type containerOptions struct {
	vm *vmmd.VMMetadata
}

// NewContainerCmd runs the DHCP server and sets up routing inside Docker
func NewCmdContainer(out io.Writer) *cobra.Command {
	co := &containerOptions{}

	cmd := &cobra.Command{
		Use:    "container [id]",
		Hidden: true,
		Args:   cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			errutils.Check(func() error {
				var err error
				if co.vm, err = matchSingleVM(args[0]); err != nil {
					return err
				}
				return RunContainer(co)
			}())
		},
	}

	return cmd
}

func RunContainer(co *containerOptions) error {
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
		dhcpIface.Hostname = co.vm.ID

		// Set the MAC address filter for the DHCP server
		dhcpIface.MACFilter = macAddresses[i]

		// Add the DNS servers from the container
		dhcpIface.SetDNSServers(clientConfig.Servers)

		go func() {
			fmt.Printf("Starting DHCP server for interface %s\n", dhcpIface.Bridge)
			if err := container.RunDHCP(dhcpIface); err != nil {
				fmt.Fprintf(os.Stderr, "%s DHCP server error: %v\n", dhcpIface.Bridge, err)
			}
		}()
	}

	// VM state handling
	if err := co.vm.SetState(vmmd.Running); err != nil {
		return fmt.Errorf("failed to update VM state: %v", err)
	}
	defer co.vm.SetState(vmmd.Stopped)

	// Run the VM
	if err := container.RunVM(co.vm, &dhcpIfaces); err != nil {
		return fmt.Errorf("runtime error for VM %q: %v", co.vm.ID, err)
	}

	return nil
}
