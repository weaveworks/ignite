package cmd

import (
	"fmt"
	"github.com/luxas/ignite/pkg/container"
	"github.com/luxas/ignite/pkg/errutils"
	"github.com/luxas/ignite/pkg/filter"
	"github.com/luxas/ignite/pkg/metadata"
	"github.com/luxas/ignite/pkg/util"
	"github.com/miekg/dns"
	"github.com/spf13/cobra"
	"io"
	"os"
)

// NewContainerCmd runs the DHCP server and sets up routing inside Docker
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
	// Load all VM metadata as Filterable objects
	mdf, err := metadata.LoadVMMetadataFilterable()
	if err != nil {
		return err
	}

	// Create an IDNameFilter to match a single VM
	d, err := filter.NewFilterer(metadata.NewVMFilter(args[0])).Single(mdf)
	if err != nil {
		return err
	}

	// Convert the result Filterable to a VMMetadata
	md, err := metadata.ToVMMetadata(d)
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
	if err := md.SetState(metadata.VMRunning); err != nil {
		return fmt.Errorf("failed to update VM state: %v", err)
	}
	defer md.SetState(metadata.VMStopped)

	// Run the VM
	container.RunVM(md.ID, md.GetKernelID(), &dhcpIfaces)

	return nil
}
