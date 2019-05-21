package cmd

import (
	"fmt"
	"github.com/luxas/ignite/pkg/container"
	"github.com/luxas/ignite/pkg/errutils"
	"github.com/luxas/ignite/pkg/util"
	"github.com/miekg/dns"
	"github.com/spf13/cobra"
	"io"
	"os"
)

// NewContainerCmd runs the dhcp server and sets up routing inside Docker
func NewCmdContainer(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "container [id]",
		Hidden: true,
		Args:   cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			err := RunContainer(out, cmd, args)
			errutils.Check(err)
		},
	}

	//cmd.Flags().StringP("output", "o", "", "Output format; available options are 'yaml', 'json' and 'short'")
	return cmd
}

// RunBuild runs when the Container command is invoked
func RunContainer(out io.Writer, cmd *cobra.Command, args []string) error {
	// The VM to run in container mode
	id := args[0]

	// Load the metadata for the VM
	md, err := loadVMMetadata(id)
	if err != nil {
		return err
	}

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
		dhcpIface.Hostname = md.ID

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
	if err := md.setState(Running); err != nil {
		return fmt.Errorf("failed to update VM state: %v", err)
	}
	defer md.setState(Stopped)

	// Run the VM
	container.RunVM(md.ID, md.ObjectData.(*vmObjectData).KernelID, &dhcpIfaces)

	return nil
}
