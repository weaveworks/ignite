package container

import (
	"context"
	"fmt"
	"github.com/firecracker-microvm/firecracker-go-sdk"
	models "github.com/firecracker-microvm/firecracker-go-sdk/client/models"
	"github.com/luxas/ignite/pkg/constants"
	"github.com/luxas/ignite/pkg/util"
	"github.com/pkg/errors"
	"net"
	"os"
	"os/signal"
	"path"
	"syscall"
)

/*
ip r list src 172.17.0.3

ip addr del "$IP" dev eth0

ip link add name br0 type bridge
ip tuntap add dev vm0 mode tap

ip link set br0 up
ip link set vm0 up

ip link set eth0 master br0
ip link set vm0 master br0
*/

// Array of container interfaces to ignore (not forward to VM)
var ignoreInterfaces = map[string]bool{
	"lo": true,
}

func NetworkSetup(dhcpIfaces *[]DHCPInterface) error {
	ifaces, err := net.Interfaces()
	if err != nil {
		return fmt.Errorf("cannot get local network interfaces: %v", err)
	}

	for _, iface := range ifaces {
		// Skip the interface if it's ignored
		if !ignoreInterfaces[iface.Name] {
			ipNet, err := takeAddress(&iface)
			if err != nil {
				return fmt.Errorf("parsing interface failed: %v", err)
			}

			dhcpIface, err := bridge(&iface)
			if err != nil {
				return fmt.Errorf("bridging interface %s failed: %v0", iface.Name, err)
			}

			// Gateway for now is just x.x.x.1 TODO: Better detection
			dhcpIface.GatewayIP = &net.IP{ipNet.IP[0], ipNet.IP[1], ipNet.IP[2], 1}
			dhcpIface.VMIPNet = ipNet

			*dhcpIfaces = append(*dhcpIfaces, dhcpIface)
		}
	}

	return nil
}

// bridge creates the TAP device and performs the bridging, returning the MAC address of the VM's adapter
func bridge(iface *net.Interface) (DHCPInterface, error) {
	tapName := constants.TAP_PREFIX + iface.Name
	bridgeName := constants.BRIDGE_PREFIX + iface.Name

	var d DHCPInterface

	if err := createTAPAdapter(tapName); err != nil {
		return d, err
	}

	if err := createBridge(bridgeName); err != nil {
		return d, err
	}

	if err := connectAdapterToBridge(tapName, bridgeName); err != nil {
		return d, err
	}

	if err := connectAdapterToBridge(iface.Name, bridgeName); err != nil {
		return d, err
	}

	d = DHCPInterface{
		VMTAP:  tapName,
		Bridge: bridgeName,
	}

	return d, nil
}

// takeAddress removes the first address of an interface and returns it
func takeAddress(iface *net.Interface) (*net.IPNet, error) {
	addrs, err := iface.Addrs()
	if err != nil {
		return nil, fmt.Errorf("interface %s has no address", iface.Name)
	}

	for _, addr := range addrs {
		var ip net.IP
		var mask net.IPMask

		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
			mask = v.Mask
		case *net.IPAddr:
			ip = v.IP
			mask = ip.DefaultMask()
		}

		if ip == nil {
			continue
		}

		ip = ip.To4()
		if ip == nil {
			continue
		}

		if _, err := util.ExecuteCommand("ip", "addr", "del", ip.String(), "dev", iface.Name); err != nil {
			return nil, errors.Wrapf(err, "failed to remove address from interface %s", iface.Name)
		}

		fmt.Printf("Found an deleted address: %s (%s)\n", ip.String(), fmt.Sprintf("%d.%d.%d.%d", mask[0], mask[1], mask[2], mask[3]))

		return &net.IPNet{
			IP:   ip,
			Mask: mask,
		}, nil
	}

	return nil, fmt.Errorf("interface %s has no valid addresses", iface.Name)
}

func createTAPAdapter(tapName string) error {
	if _, err := util.ExecuteCommand("ip", "tuntap", "add", "mode", "tap", tapName); err != nil {
		return err
	}

	return setLinkUp(tapName)
}

func createBridge(bridgeName string) error {
	if _, err := util.ExecuteCommand("ip", "link", "add", "name", bridgeName, "type", "bridge"); err != nil {
		return err
	}

	return setLinkUp(bridgeName)
}

func setLinkUp(adapterName string) error {
	_, err := util.ExecuteCommand("ip", "link", "set", adapterName, "up")
	return err
}

func connectAdapterToBridge(adapterName, bridgeName string) error {
	_, err := util.ExecuteCommand("ip", "link", "set", adapterName, "master", bridgeName)
	return err
}

func RunVM(vmID string, kernelID string, dhcpIfaces *[]DHCPInterface) {
	drivePath := path.Join(constants.VM_DIR, vmID, constants.IMAGE_FS)

	networkInterfaces := make([]firecracker.NetworkInterface, 0, len(*dhcpIfaces))
	for _, dhcpIface := range *dhcpIfaces {
		networkInterfaces = append(networkInterfaces, firecracker.NetworkInterface{
			MacAddress:  dhcpIface.MACFilter,
			HostDevName: dhcpIface.VMTAP,
		})
	}

	const socketPath = "/tmp/firecracker.sock"
	cfg := firecracker.Config{
		SocketPath:      socketPath,
		KernelImagePath: path.Join(constants.KERNEL_DIR, kernelID, constants.KERNEL_FILE),
		KernelArgs:      "console=ttyS0 reboot=k panic=1 pci=off",
		Drives: []models.Drive{{
			DriveID:      firecracker.String("1"),
			PathOnHost:   &drivePath,
			IsRootDevice: firecracker.Bool(true),
			IsReadOnly:   firecracker.Bool(false),
			Partuuid:     "",
		}},
		NetworkInterfaces: networkInterfaces,
		MachineCfg: models.MachineConfiguration{
			VcpuCount: 1,
		},
		//JailerCfg: firecracker.JailerConfig{
		//	GID:      firecracker.Int(0),
		//	UID:      firecracker.Int(0),
		//	ID:       md.ID,
		//	NumaNode: firecracker.Int(0),
		//	ExecFile: "firecracker",
		//},
	}

	// stdout will be directed to this file
	//stdoutPath := "/tmp/stdout.log"
	//stdout, err := os.OpenFile(stdoutPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	//if err != nil {
	//	panic(fmt.Errorf("failed to create stdout file: %v", err))
	//}

	// stderr will be directed to this file
	//stderrPath := "/tmp/stderr.log"
	//stderr, err := os.OpenFile(stderrPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	//if err != nil {
	//	panic(fmt.Errorf("failed to create stderr file: %v", err))
	//}

	ctx, vmmCancel := context.WithCancel(context.Background())
	defer vmmCancel()
	// build our custom command that contains our two files to
	// write to during process execution
	cmd := firecracker.VMCommandBuilder{}.
		WithBin("firecracker").
		WithSocketPath(socketPath).
		WithStdin(os.Stdin).
		WithStdout(os.Stdout).
		WithStderr(os.Stderr).
		Build(ctx)

	m, err := firecracker.NewMachine(ctx, cfg, firecracker.WithProcessRunner(cmd))
	if err != nil {
		panic(fmt.Errorf("Failed creating machine: %s", err))
	}

	//if opts.validMetadata != nil {
	//	m.EnableMetadata(opts.validMetadata)
	//}

	if err := m.Start(ctx); err != nil {
		panic(fmt.Errorf("Failed to start machine: %v", err))
	}
	defer m.StopVMM()

	installSignalHandlers(ctx, m)

	// wait for the VMM to exit
	if err := m.Wait(ctx); err != nil {
		panic(fmt.Errorf("Wait returned an error %s", err))
	}
	fmt.Println("Start machine was happy")

	//m, err := firecracker.NewMachine(ctx, cfg, firecracker.WithProcessRunner(cmd))
	//if err != nil {
	//	panic(fmt.Errorf("failed to create new machine: %v", err))
	//}
	//
	//defer os.Remove(cfg.SocketPath)
	//
	//fmt.Println("Starting machine...")
	//if err := m.Start(ctx); err != nil {
	//	panic(fmt.Errorf("failed to initialize machine: %v", err))
	//}
	//
	//// wait for VMM to execute
	//if err := m.Wait(ctx); err != nil {
	//	panic(err)
	//}
}

// Install custom signal handlers:
func installSignalHandlers(ctx context.Context, m *firecracker.Machine) {
	go func() {
		// Clear some default handlers installed by the firecracker SDK:
		signal.Reset(os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

		for {
			switch s := <-c; {
			case s == syscall.SIGTERM || s == os.Interrupt:
				fmt.Println("Caught SIGINT, requesting clean shutdown")
				m.Shutdown(ctx)
			case s == syscall.SIGQUIT:
				fmt.Println("Caught SIGTERM, forcing shutdown")
				m.StopVMM()
			}
		}
	}()
}
