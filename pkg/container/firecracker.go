package container

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/firecracker-microvm/firecracker-go-sdk"
	models "github.com/firecracker-microvm/firecracker-go-sdk/client/models"
	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	"github.com/weaveworks/ignite/pkg/client"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/operations/lookup"
)

// ExecuteFirecracker executes the firecracker process using the Go SDK
func ExecuteFirecracker(vm *api.VM, dhcpIfaces []DHCPInterface) error {
	drivePath := vm.SnapshotDev()

	networkInterfaces := make([]firecracker.NetworkInterface, 0, len(dhcpIfaces))
	for _, dhcpIface := range dhcpIfaces {
		networkInterfaces = append(networkInterfaces, firecracker.NetworkInterface{
			MacAddress:  dhcpIface.MACFilter,
			HostDevName: dhcpIface.VMTAP,
		})
	}

	vCPUCount := int64(vm.Spec.CPUs)
	memSizeMib := int64(vm.Spec.Memory.MBytes())

	cvmLine := vm.Spec.Kernel.CvmLine
	if len(cvmLine) == 0 {
		// if for some reason cvmline would be unpopulated, set it to the default
		cvmLine = constants.VM_DEFAULT_KERNEL_ARGS
	}

	kernelUID, err := lookup.KernelUIDForVM(vm, client.DefaultClient)
	if err != nil {
		return err
	}

	firecrackerSocketPath := path.Join(vm.ObjectPath(), constants.FIRECRACKER_API_SOCKET)
	logSocketPath := path.Join(vm.ObjectPath(), constants.LOG_FIFO)
	metricsSocketPath := path.Join(vm.ObjectPath(), constants.METRICS_FIFO)
	cfg := firecracker.Config{
		SocketPath:      firecrackerSocketPath,
		KernelImagePath: path.Join(constants.KERNEL_DIR, kernelUID.String(), constants.KERNEL_FILE),
		KernelArgs:      cvmLine,
		Drives: []models.Drive{{
			DriveID:      firecracker.String("1"),
			PathOnHost:   &drivePath,
			IsRootDevice: firecracker.Bool(true),
			IsReadOnly:   firecracker.Bool(false),
		}},
		NetworkInterfaces: networkInterfaces,
		MachineCfg: models.MachineConfiguration{
			VcpuCount:  &vCPUCount,
			MemSizeMib: &memSizeMib,
			HtEnabled:  boolPtr(true),
		},
		//JailerCfg: firecracker.JailerConfig{
		//	GID:      firecracker.Int(0),
		//	UID:      firecracker.Int(0),
		//	ID:       vm.ID,
		//	NumaNode: firecracker.Int(0),
		//	ExecFile: "firecracker",
		//},

		// TODO: We could use /dev/null, but firecracker-go-sdk issues Mkfifo which collides with the existing device
		LogLevel:    constants.VM_LOG_LEVEL,
		LogFifo:     logSocketPath,
		MetricsFifo: metricsSocketPath,
	}

	// Remove these FIFOs for now
	defer os.Remove(logSocketPath)
	defer os.Remove(metricsSocketPath)

	ctx, vmmCancel := context.WithCancel(context.Background())
	defer vmmCancel()

	cvm := firecracker.VMCommandBuilder{}.
		WithBin("firecracker").
		WithSocketPath(firecrackerSocketPath).
		WithStdin(os.Stdin).
		WithStdout(os.Stdout).
		WithStderr(os.Stderr).
		Build(ctx)

	m, err := firecracker.NewMachine(ctx, cfg, firecracker.WithProcessRunner(cvm))
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

func boolPtr(val bool) *bool {
	return &val
}
