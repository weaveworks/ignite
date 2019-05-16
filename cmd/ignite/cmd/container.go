package cmd

import (
	"context"
	"fmt"
	"github.com/firecracker-microvm/firecracker-go-sdk"
	models "github.com/firecracker-microvm/firecracker-go-sdk/client/models"
	"github.com/luxas/ignite/pkg/constants"
	"github.com/luxas/ignite/pkg/errutils"
	"github.com/luxas/ignite/pkg/util"
	"github.com/spf13/cobra"
	"io"
	"os"
	"os/signal"
	"path"
	"syscall"
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

	md := &vmMetadata{
		ID: id,
	}

	if err := md.load(); err != nil {
		return err
	}

	md.State = Running

	runVM(md)

	md.State = Stopped

	return nil
}

func runVM(md *vmMetadata) {
	drivePath := path.Join(constants.VM_DIR, md.ID, constants.IMAGE_FS)

	fmt.Printf("%+v\n", md)
	const socketPath = "/tmp/firecracker.sock"
	cfg := firecracker.Config{
		SocketPath:      socketPath,
		KernelImagePath: path.Join(constants.KERNEL_DIR, md.KernelID, constants.KERNEL_FILE),
		KernelArgs:      "console=ttyS0 reboot=k panic=1 pci=off",
		Drives: []models.Drive{{
			DriveID:      firecracker.String("1"),
			PathOnHost:   &drivePath,
			IsRootDevice: firecracker.Bool(true),
			IsReadOnly:   firecracker.Bool(false),
			Partuuid:     "",
		}},
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

func createTAPAdapter(tapName string) error {
	_, err := util.ExecuteCommand("ip", "tuntap", "add", "mode", "tap", tapName)
	if err != nil {
		return err
	}
	_, err = util.ExecuteCommand("ip", "link", "set", tapName, "up")
	return err
}

func connectTAPToBridge(tapName, bridgeName string) error {
	_, err := util.ExecuteCommand("brctl", "addif", bridgeName, tapName)
	return err
}
