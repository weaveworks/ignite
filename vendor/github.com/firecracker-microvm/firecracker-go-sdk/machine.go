// Copyright 2018 Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//	http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package firecracker

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	models "github.com/firecracker-microvm/firecracker-go-sdk/client/models"
	log "github.com/sirupsen/logrus"
)

const (
	userAgent = "firecracker-go-sdk"
)

// ErrAlreadyStarted signifies that the Machine has already started and cannot
// be started again.
var ErrAlreadyStarted = errors.New("firecracker: machine already started")

// Config is a collection of user-configurable VMM settings
type Config struct {
	// SocketPath defines the file path where the Firecracker control socket
	// should be created.
	SocketPath string

	// LogFifo defines the file path where the Firecracker log named-pipe should
	// be located.
	LogFifo string

	// LogLevel defines the verbosity of Firecracker logging.  Valid values are
	// "Error", "Warning", "Info", and "Debug", and are case-sensitive.
	LogLevel string

	// MetricsFifo defines the file path where the Firecracker metrics
	// named-pipe should be located.
	MetricsFifo string

	// KernelImagePath defines the file path where the kernel image is located.
	// The kernel image must be an uncompressed ELF image.
	KernelImagePath string

	// KernelArgs defines the command-line arguments that should be passed to
	// the kernel.
	KernelArgs string

	// Drives specifies BlockDevices that should be made available to the
	// microVM.
	Drives []models.Drive

	// NetworkInterfaces specifies the tap devices that should be made available
	// to the microVM.
	NetworkInterfaces []NetworkInterface

	// FifoLogWriter is an io.Writer that is used to redirect the contents of the
	// fifo log to the writer.
	FifoLogWriter io.Writer

	// VsockDevices specifies the vsock devices that should be made available to
	// the microVM.
	VsockDevices []VsockDevice

	// Debug enables debug-level logging for the SDK.
	Debug bool

	// MachineCfg represents the firecracker microVM process configuration
	MachineCfg models.MachineConfiguration

	// DisableValidation allows for easier mock testing by disabling the
	// validation of configuration performed by the SDK.
	DisableValidation bool

	// EnableJailer will enable the jailer. By enabling the jailer, root level
	// permissions are required.
	EnableJailer bool

	// JailerCfg is configuration specific for the jailer process.
	JailerCfg JailerConfig
}

// Validate will ensure that the required fields are set and that
// the fields are valid values.
func (cfg *Config) Validate() error {
	if cfg.DisableValidation {
		return nil
	}

	if _, err := os.Stat(cfg.KernelImagePath); err != nil {
		return fmt.Errorf("failed to stat kernel image path, %q: %v", cfg.KernelImagePath, err)
	}

	rootPath := ""
	for _, drive := range cfg.Drives {
		if BoolValue(drive.IsRootDevice) {
			rootPath = StringValue(drive.PathOnHost)
			break
		}
	}

	if _, err := os.Stat(rootPath); err != nil {
		return fmt.Errorf("failed to stat host path, %q: %v", rootPath, err)
	}

	// Check the non-existence of some files:
	if _, err := os.Stat(cfg.SocketPath); err == nil {
		return fmt.Errorf("socket %s already exists", cfg.SocketPath)
	}

	if cfg.MachineCfg.VcpuCount == nil ||
		Int64Value(cfg.MachineCfg.VcpuCount) < 1 {
		return fmt.Errorf("machine needs a nonzero VcpuCount")
	}
	if cfg.MachineCfg.MemSizeMib == nil ||
		Int64Value(cfg.MachineCfg.MemSizeMib) < 1 {
		return fmt.Errorf("machine needs a nonzero amount of memory")
	}
	if cfg.MachineCfg.HtEnabled == nil {
		return fmt.Errorf("machine needs a setting for ht_enabled")
	}
	return nil
}

// Machine is the main object for manipulating Firecracker microVMs
type Machine struct {
	// Handlers holds the set of handlers that are run for validation and start
	Handlers Handlers

	cfg           Config
	client        *Client
	cmd           *exec.Cmd
	logger        *log.Entry
	machineConfig models.MachineConfiguration // The actual machine config as reported by Firecracker
	// startOnce ensures that the machine can only be started once
	startOnce sync.Once
	// exitCh is a channel which gets closed when the VMM exits
	exitCh chan struct{}
	// err records any error executing the VMM
	err error
}

// Logger returns a logrus logger appropriate for logging hypervisor messages
func (m *Machine) Logger() *log.Entry {
	return m.logger.WithField("subsystem", userAgent)
}

// NetworkInterface represents a Firecracker microVM's network interface.
type NetworkInterface struct {
	// MacAddress defines the MAC address that should be assigned to the network
	// interface inside the microVM.
	MacAddress string
	// HostDevName defines the file path of the tap device on the host.
	HostDevName string
	// AllowMMDS makes the Firecracker MMDS available on this network interface.
	AllowMMDS bool

	// InRateLimiter limits the incoming bytes.
	InRateLimiter *models.RateLimiter
	// OutRateLimiter limits the outgoing bytes.
	OutRateLimiter *models.RateLimiter
}

// VsockDevice represents a vsock connection between the host and the guest
// microVM.
type VsockDevice struct {
	// Path defines the filesystem path of the vsock device on the host.
	Path string
	// CID defines the 32-bit Context Identifier for the vsock device.  See
	// the vsock(7) manual page for more information.
	CID uint32
}

// SocketPath returns the filesystem path to the socket used for VMM
// communication
func (m *Machine) socketPath() string {
	return m.cfg.SocketPath
}

// LogFile returns the filesystem path of the VMM log
func (m *Machine) LogFile() string {
	return m.cfg.LogFifo
}

// LogLevel returns the VMM log level.
func (m *Machine) LogLevel() string {
	return m.cfg.LogLevel
}

// NewMachine initializes a new Machine instance and performs validation of the
// provided Config.
func NewMachine(ctx context.Context, cfg Config, opts ...Opt) (*Machine, error) {
	m := &Machine{
		exitCh: make(chan struct{}),
	}

	m.Handlers = defaultHandlers

	if cfg.EnableJailer {
		m.Handlers.Validation = m.Handlers.Validation.Append(JailerConfigValidationHandler)
		if err := jail(ctx, m, &cfg); err != nil {
			return nil, err
		}
	} else {
		m.Handlers.Validation = m.Handlers.Validation.Append(ConfigValidationHandler)
		m.cmd = defaultFirecrackerVMMCommandBuilder.
			WithSocketPath(cfg.SocketPath).
			Build(ctx)
	}

	for _, opt := range opts {
		opt(m)
	}

	if m.logger == nil {
		logger := log.New()
		if cfg.Debug {
			logger.SetLevel(log.DebugLevel)
		}

		m.logger = log.NewEntry(logger)
	}

	if m.client == nil {
		m.client = NewClient(cfg.SocketPath, m.logger, cfg.Debug)
	}

	m.machineConfig = cfg.MachineCfg
	m.cfg = cfg

	m.logger.Debug("Called NewMachine()")
	return m, nil
}

// Start will iterate through the handler list and call each handler. If an
// error occurred during handler execution, that error will be returned. If the
// handlers succeed, then this will start the VMM instance.
// Start may only be called once per Machine.  Subsequent calls will return
// ErrAlreadyStarted.
func (m *Machine) Start(ctx context.Context) error {
	m.logger.Debug("Called Machine.Start()")
	alreadyStarted := true
	m.startOnce.Do(func() {
		m.logger.Debug("Marking Machine as Started")
		alreadyStarted = false
	})
	if alreadyStarted {
		return ErrAlreadyStarted
	}

	if err := m.Handlers.Run(ctx, m); err != nil {
		return err
	}

	return m.startInstance(ctx)
}

// Shutdown requests a clean shutdown of the VM by sending CtrlAltDelete on the virtual keyboard
func (m *Machine) Shutdown(ctx context.Context) error {
	m.logger.Debug("Called machine.Shutdown()")
	return m.sendCtrlAltDel(ctx)
}

// Wait will wait until the firecracker process has finished.  Wait is safe to
// call concurrently, and will deliver the same error to all callers, subject to
// each caller's context cancellation.
func (m *Machine) Wait(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-m.exitCh:
		return m.err
	}
}

func (m *Machine) addVsocks(ctx context.Context, vsocks ...VsockDevice) error {
	for _, dev := range m.cfg.VsockDevices {
		if err := m.addVsock(ctx, dev); err != nil {
			return err
		}
	}
	return nil
}

func (m *Machine) createNetworkInterfaces(ctx context.Context, ifaces ...NetworkInterface) error {
	for id, iface := range ifaces {
		if err := m.createNetworkInterface(ctx, iface, id+1); err != nil {
			return err
		}
		m.logger.Debugf("createNetworkInterface returned for %s", iface.HostDevName)
	}

	return nil
}

func (m *Machine) attachDrives(ctx context.Context, drives ...models.Drive) error {
	for _, dev := range drives {
		if err := m.attachDrive(ctx, dev); err != nil {
			m.logger.Errorf("While attaching drive %s, got error %s", StringValue(dev.PathOnHost), err)
			return err
		}
		m.logger.Debugf("attachDrive returned for %s", StringValue(dev.PathOnHost))
	}

	return nil
}

// startVMM starts the firecracker vmm process and configures logging.
func (m *Machine) startVMM(ctx context.Context) error {
	m.logger.Printf("Called startVMM(), setting up a VMM on %s", m.cfg.SocketPath)

	errCh := make(chan error)

	err := m.cmd.Start()
	if err != nil {
		m.logger.Errorf("Failed to start VMM: %s", err)
		close(m.exitCh)
		return err
	}
	m.logger.Debugf("VMM started socket path is %s", m.cfg.SocketPath)

	go func() {
		if err := m.cmd.Wait(); err != nil {
			m.logger.Warnf("firecracker exited: %s", err.Error())
		} else {
			m.logger.Printf("firecracker exited: status=0")
		}

		os.Remove(m.cfg.SocketPath)
		os.Remove(m.cfg.LogFifo)
		os.Remove(m.cfg.MetricsFifo)
		errCh <- err

		// Notify subscribers that there will be no more values.
		// When err is nil, two reads are performed (waitForSocket and close exitCh goroutine),
		// second one never ends as it tries to read from empty channel.
		close(errCh)
	}()

	// Set up a signal handler and pass INT, QUIT, and TERM through to firecracker
	sigchan := make(chan os.Signal)
	signal.Notify(sigchan, os.Interrupt,
		syscall.SIGQUIT,
		syscall.SIGTERM,
		syscall.SIGHUP,
		syscall.SIGABRT)
	m.logger.Debugf("Setting up signal handler")
	go func() {
		sig := <-sigchan
		m.logger.Printf("Caught signal %s", sig)
		m.cmd.Process.Signal(sig)
	}()

	// Wait for firecracker to initialize:
	err = m.waitForSocket(3*time.Second, errCh)
	if err != nil {
		msg := fmt.Sprintf("Firecracker did not create API socket %s: %s", m.cfg.SocketPath, err)
		err = errors.New(msg)
		close(m.exitCh)
		return err
	}
	go func() {
		select {
		case <-ctx.Done():
			m.err = ctx.Err()
		case err := <-errCh:
			m.err = err
		}

		close(m.exitCh)
	}()

	m.logger.Debugf("returning from startVMM()")
	return nil
}

//StopVMM stops the current VMM.
func (m *Machine) StopVMM() error {
	return m.stopVMM()
}

func (m *Machine) stopVMM() error {
	if m.cmd != nil && m.cmd.Process != nil {
		log.Debug("stopVMM(): sending sigterm to firecracker")
		return m.cmd.Process.Signal(syscall.SIGTERM)
	}
	log.Debug("stopVMM(): no firecracker process running, not sending a signal")

	// don't return an error if the process isn't even running
	return nil
}

// createFifos sets up the firecracker logging and metrics FIFOs
func createFifos(logFifo, metricsFifo string) error {
	log.Debugf("Creating FIFO %s", logFifo)
	if err := syscall.Mkfifo(logFifo, 0700); err != nil {
		return fmt.Errorf("Failed to create log fifo: %v", err)
	}

	log.Debugf("Creating metric FIFO %s", metricsFifo)
	if err := syscall.Mkfifo(metricsFifo, 0700); err != nil {
		return fmt.Errorf("Failed to create metric fifo: %v", err)
	}
	return nil
}

func (m *Machine) setupLogging(ctx context.Context) error {
	if len(m.cfg.LogFifo) == 0 || len(m.cfg.MetricsFifo) == 0 {
		// No logging configured
		m.logger.Printf("VMM logging and metrics disabled.")
		return nil
	}

	l := models.Logger{
		LogFifo:       String(m.cfg.LogFifo),
		Level:         String(m.cfg.LogLevel),
		MetricsFifo:   String(m.cfg.MetricsFifo),
		ShowLevel:     Bool(true),
		ShowLogOrigin: Bool(false),
		Options:       []string{},
	}

	_, err := m.client.PutLogger(ctx, &l)
	if err != nil {
		return err
	}

	m.logger.Debugf("Configured VMM logging to %s, metrics to %s",
		m.cfg.LogFifo,
		m.cfg.MetricsFifo,
	)

	if m.cfg.FifoLogWriter != nil {
		if err := captureFifoToFile(m.logger, m.cfg.LogFifo, m.cfg.FifoLogWriter); err != nil {
			return err
		}
	}

	return nil
}

func captureFifoToFile(logger *log.Entry, fifoPath string, fifo io.Writer) error {
	// create the fifo pipe which will be used
	// to write its contents to a file.
	fifoPipe, err := os.OpenFile(fifoPath, os.O_RDONLY, 0600)
	if err != nil {
		return fmt.Errorf("Failed to open fifo path at %q: %v", fifoPath, err)
	}

	if err := syscall.Unlink(fifoPath); err != nil {
		logger.Warnf("Failed to unlink %s", fifoPath)
	}

	logger.Debugf("Capturing %q to writer", fifoPath)

	// Uses a go routine to do a non-blocking io.Copy. The fifo
	// file should be closed when the appication has finished, since
	// the forked firecracker application will be closed resulting
	// in the pipe to return an io.EOF
	go func() {
		defer fifoPipe.Close()

		if _, err := io.Copy(fifo, fifoPipe); err != nil {
			logger.Warnf("io.Copy failed to copy contents of fifo pipe: %v", err)
		}
	}()

	return nil
}

func (m *Machine) createMachine(ctx context.Context) error {
	resp, err := m.client.PutMachineConfiguration(ctx, &m.cfg.MachineCfg)
	if err != nil {
		m.logger.Errorf("PutMachineConfiguration returned %s", resp.Error())
		return err
	}

	m.logger.Debug("PutMachineConfiguration returned")
	err = m.refreshMachineConfiguration()
	if err != nil {
		log.Errorf("Unable to inspect Firecracker MachineConfiguration. Continuing anyway. %s", err)
	}
	m.logger.Debug("createMachine returning")
	return err
}

func (m *Machine) createBootSource(ctx context.Context, imagePath, kernelArgs string) error {
	bsrc := models.BootSource{
		KernelImagePath: &imagePath,
		BootArgs:        kernelArgs,
	}

	resp, err := m.client.PutGuestBootSource(ctx, &bsrc)
	if err == nil {
		m.logger.Printf("PutGuestBootSource: %s", resp.Error())
	}

	return err
}

func (m *Machine) createNetworkInterface(ctx context.Context, iface NetworkInterface, iid int) error {
	ifaceID := strconv.Itoa(iid)
	m.logger.Printf("Attaching NIC %s (hwaddr %s) at index %s", iface.HostDevName, iface.MacAddress, ifaceID)

	ifaceCfg := models.NetworkInterface{
		IfaceID:           &ifaceID,
		GuestMac:          iface.MacAddress,
		HostDevName:       String(iface.HostDevName),
		AllowMmdsRequests: iface.AllowMMDS,
	}

	if iface.InRateLimiter != nil {
		ifaceCfg.RxRateLimiter = iface.InRateLimiter
	}

	if iface.OutRateLimiter != nil {
		ifaceCfg.TxRateLimiter = iface.OutRateLimiter
	}

	if iface.InRateLimiter != nil {
		ifaceCfg.RxRateLimiter = iface.InRateLimiter
	}

	if iface.OutRateLimiter != nil {
		ifaceCfg.TxRateLimiter = iface.OutRateLimiter
	}

	resp, err := m.client.PutGuestNetworkInterfaceByID(ctx, ifaceID, &ifaceCfg)
	if err == nil {
		m.logger.Printf("PutGuestNetworkInterfaceByID: %s", resp.Error())
	}

	return err
}

// attachDrive attaches a secondary block device
func (m *Machine) attachDrive(ctx context.Context, dev models.Drive) error {
	hostPath := StringValue(dev.PathOnHost)
	log.Infof("Attaching drive %s, slot %s, root %t.", hostPath, StringValue(dev.DriveID), BoolValue(dev.IsRootDevice))
	respNoContent, err := m.client.PutGuestDriveByID(ctx, StringValue(dev.DriveID), &dev)
	if err == nil {
		m.logger.Printf("Attached drive %s: %s", hostPath, respNoContent.Error())
	} else {
		m.logger.Errorf("Attach drive failed: %s: %s", hostPath, err)
	}
	return err
}

// addVsock adds a vsock to the instance
func (m *Machine) addVsock(ctx context.Context, dev VsockDevice) error {
	vsockCfg := models.Vsock{
		GuestCid: int64(dev.CID),
		ID:       &dev.Path,
	}

	resp, _, err := m.client.PutGuestVsockByID(ctx, dev.Path, &vsockCfg)
	if err != nil {
		return err
	}
	m.logger.Debugf("Attach vsock %s successful: %s", dev.Path, resp.Error())
	return nil
}

func (m *Machine) startInstance(ctx context.Context) error {
	action := models.InstanceActionInfoActionTypeInstanceStart
	info := models.InstanceActionInfo{
		ActionType: &action,
	}

	resp, err := m.client.CreateSyncAction(ctx, &info)
	if err == nil {
		m.logger.Printf("startInstance successful: %s", resp.Error())
	} else {
		m.logger.Errorf("Starting instance: %s", err)
	}
	return err
}

func (m *Machine) sendCtrlAltDel(ctx context.Context) error {
	action := models.InstanceActionInfoActionTypeSendCtrlAltDel
	info := models.InstanceActionInfo{
		ActionType: &action,
	}

	resp, err := m.client.CreateSyncAction(ctx, &info)
	if err == nil {
		m.logger.Printf("Sent instance shutdown request: %s", resp.Error())
	} else {
		m.logger.Errorf("Unable to send CtrlAltDel: %s", err)
	}
	return err
}

// SetMetadata sets the machine's metadata for MDDS
func (m *Machine) SetMetadata(ctx context.Context, metadata interface{}) error {
	if _, err := m.client.PutMmds(ctx, metadata); err != nil {
		m.logger.Errorf("Setting metadata: %s", err)
		return err
	}

	m.logger.Printf("SetMetadata successful")
	return nil
}

// UpdateGuestDrive will modify the current guest drive of ID index with the new
// parameters of the partialDrive.
func (m *Machine) UpdateGuestDrive(ctx context.Context, driveID, pathOnHost string, opts ...PatchGuestDriveByIDOpt) error {
	if _, err := m.client.PatchGuestDriveByID(ctx, driveID, pathOnHost, opts...); err != nil {
		m.logger.Errorf("PatchGuestDrive failed: %v", err)
		return err
	}

	m.logger.Printf("PatchGuestDrive successful")
	return nil
}

// refreshMachineConfiguration synchronizes our cached representation of the machine configuration
// with that reported by the Firecracker API
func (m *Machine) refreshMachineConfiguration() error {
	resp, err := m.client.GetMachineConfiguration()
	if err != nil {
		return err
	}

	m.logger.Infof("refreshMachineConfiguration: %s", resp.Error())
	m.machineConfig = *resp.Payload
	return nil
}

// waitForSocket waits for the given file to exist
func (m *Machine) waitForSocket(timeout time.Duration, exitchan chan error) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ticker := time.NewTicker(10 * time.Millisecond)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-exitchan:
			return err
		case <-ticker.C:
			if _, err := os.Stat(m.cfg.SocketPath); err != nil {
				continue
			}

			// Send test HTTP request to make sure socket is available
			if _, err := m.client.GetMachineConfiguration(); err != nil {
				continue
			}

			return nil
		}
	}
}
