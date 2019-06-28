package container

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/firecracker-microvm/firecracker-go-sdk"
	models "github.com/firecracker-microvm/firecracker-go-sdk/client/models"
	"github.com/pkg/errors"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/metadata/vmmd"
	"github.com/weaveworks/ignite/pkg/util"
	"k8s.io/apimachinery/pkg/util/wait"
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

// Array of container interfaces to ignore (not forward to vm)
var ignoreInterfaces = map[string]bool{
	"lo": true,
}

func NetworkSetup(dhcpIfaces *[]DHCPInterface) error {
	interval := 1 * time.Second
	timeout := 1 * time.Minute
	return wait.PollImmediate(interval, timeout, func() (bool, error) {
		// this func returns true if it's done, and optionally an error
		retry, err := networkSetup(dhcpIfaces)
		if err == nil {
			// we're done here
			log.Printf("network setup done")
			return true, nil
		}
		if retry {
			// we got an error, but let's ignore it and try again
			log.Printf("retry, although error: %v", err)
			return false, nil
		}
		log.Printf("fatal error %v", err)
		// the error was fatal, return it
		return false, err
	})
}

func networkSetup(dhcpIfaces *[]DHCPInterface) (bool, error) {
	ifaces, err := net.Interfaces()
	if err != nil || ifaces == nil || len(ifaces) == 0 {
		return true, fmt.Errorf("cannot get local network interfaces: %v", err)
	}

	// interfacesCount counts the interfaces that are relevant to Ignite (in other words, not ignored)
	interfacesCount := 0
	for _, iface := range ifaces {
		// Skip the interface if it's ignored
		if !ignoreInterfaces[iface.Name] {
			// This is an interface we care about
			interfacesCount++

			ipNet, retry, err := takeAddress(&iface)
			if err != nil {
				return retry, fmt.Errorf("parsing interface failed: %v", err)
			}

			dhcpIface, err := bridge(&iface)
			if err != nil {
				return false, fmt.Errorf("bridging interface %s failed: %v0", iface.Name, err)
			}

			// Gateway for now is just x.x.x.1 TODO: Better detection
			dhcpIface.GatewayIP = &net.IP{ipNet.IP[0], ipNet.IP[1], ipNet.IP[2], 1}
			dhcpIface.VMIPNet = ipNet

			*dhcpIfaces = append(*dhcpIfaces, dhcpIface)
		}
	}

	// If there weren't any interfaces we cared about, retry the loop
	if interfacesCount == 0 {
		return true, fmt.Errorf("no active interfaces available yet")
	}

	return false, nil
}

// bridge creates the TAP device and performs the bridging, returning the MAC address of the vm's adapter
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
func takeAddress(iface *net.Interface) (*net.IPNet, bool, error) {
	addrs, err := iface.Addrs()
	if err != nil || addrs == nil || len(addrs) == 0 {
		// set the bool to true so the caller knows to retry
		return nil, true, fmt.Errorf("interface %s has no address", iface.Name)
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
			return nil, false, errors.Wrapf(err, "failed to remove address from interface %s", iface.Name)
		}

		fmt.Printf("Found an deleted address: %s (%s)\n", ip.String(), fmt.Sprintf("%d.%d.%d.%d", mask[0], mask[1], mask[2], mask[3]))

		return &net.IPNet{
			IP:   ip,
			Mask: mask,
		}, false, nil
	}

	return nil, false, fmt.Errorf("interface %s has no valid addresses", iface.Name)
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

func RunVM(md *vmmd.VMMetadata, dhcpIfaces *[]DHCPInterface) error {
	od := md.VMOD()
	drivePath := md.SnapshotDev()

	networkInterfaces := make([]firecracker.NetworkInterface, 0, len(*dhcpIfaces))
	for _, dhcpIface := range *dhcpIfaces {
		networkInterfaces = append(networkInterfaces, firecracker.NetworkInterface{
			MacAddress:  dhcpIface.MACFilter,
			HostDevName: dhcpIface.VMTAP,
		})
	}

	kernelCmd := od.KernelCmd
	if len(kernelCmd) == 0 {
		kernelCmd = constants.VM_KERNEL_ARGS
	}

	cfg := firecracker.Config{
		SocketPath:      constants.SOCKET_PATH,
		KernelImagePath: path.Join(constants.KERNEL_DIR, od.KernelID.String(), constants.KERNEL_FILE),
		KernelArgs:      kernelCmd,
		Drives: []models.Drive{{
			DriveID:      firecracker.String("1"),
			PathOnHost:   &drivePath,
			IsRootDevice: firecracker.Bool(true),
			IsReadOnly:   firecracker.Bool(false),
		}},
		NetworkInterfaces: networkInterfaces,
		MachineCfg: models.MachineConfiguration{
			VcpuCount:  od.VCPUs,
			MemSizeMib: od.Memory,
		},
		//JailerCfg: firecracker.JailerConfig{
		//	GID:      firecracker.Int(0),
		//	UID:      firecracker.Int(0),
		//	ID:       md.ID,
		//	NumaNode: firecracker.Int(0),
		//	ExecFile: "firecracker",
		//},

		// TODO: We could use /dev/null, but firecracker-go-sdk issues Mkfifo which collides with the existing device
		LogLevel:    constants.VM_LOG_LEVEL,
		LogFifo:     constants.LOG_FIFO,
		MetricsFifo: constants.METRICS_FIFO,
	}

	// Remove these FIFOs for now
	defer os.Remove(constants.LOG_FIFO)
	defer os.Remove(constants.METRICS_FIFO)

	ctx, vmmCancel := context.WithCancel(context.Background())
	defer vmmCancel()

	cmd := firecracker.VMCommandBuilder{}.
		WithBin("firecracker").
		WithSocketPath(constants.SOCKET_PATH).
		WithStdin(os.Stdin).
		WithStdout(os.Stdout).
		WithStderr(os.Stderr).
		Build(ctx)

	m, err := firecracker.NewMachine(ctx, cfg, firecracker.WithProcessRunner(cmd))
	if err != nil {
		return fmt.Errorf("failed to create machine: %s", err)
	}

	//defer os.Remove(cfg.SocketPath)

	//if opts.validMetadata != nil {
	//	m.EnableMetadata(opts.validMetadata)
	//}

	if err := m.Start(ctx); err != nil {
		return fmt.Errorf("failed to start machine: %v", err)
	}
	defer m.StopVMM()

	installSignalHandlers(ctx, m)

	// wait for the VMM to exit
	if err := m.Wait(ctx); err != nil {
		return fmt.Errorf("wait returned an error %s", err)
	}

	return nil
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
				time.Sleep(constants.STOP_TIMEOUT * time.Second)

				// There's no direct way of checking if a VM is running, so we test if we can send it another shutdown
				// request. If that fails, the VM is still running and we need to kill it.
				if err := m.Shutdown(ctx); err == nil {
					fmt.Println("Timeout exceeded, forcing shutdown") // TODO: Proper logging
					m.StopVMM()
				}
			case s == syscall.SIGQUIT:
				fmt.Println("Caught SIGTERM, forcing shutdown")
				m.StopVMM()
			}
		}
	}()
}
