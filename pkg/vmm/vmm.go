package vmm

/*
import (
	"context"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	firecracker "github.com/firecracker-microvm/firecracker-go-sdk"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// VMM represents a virtual machine monitor
type VMM struct {
	binary        string
	name          string
	rootDrivePath string
	cfg           firecracker.Config
	metadata      interface{}
	fifoLogFile   string
	copyFiles     []string
	cleanupFns    []func() error
}

// Run a VMM with a given set of options
func (vmm *VMM) Run(ctx context.Context) error {
	vmmlogger := log.New()
	vmmlogger.SetLevel(log.GetLevel())

	if err := vmm.createRuntimeDir(); err != nil {
		return err
	}

	if err := vmm.copyFilesFromHost(); err != nil {
		return err
	}

	logWriter, err := vmm.handleFifos(createFifoFile)
	if err != nil {
		return err
	}
	defer vmm.Cleanup()
	vmm.cfg.FifoLogWriter = logWriter

	vmmCtx, vmmCancel := context.WithCancel(ctx)
	defer vmmCancel()

	machineOpts := []firecracker.Opt{
		firecracker.WithLogger(log.NewEntry(vmmlogger)),
	}

	if len(vmm.binary) != 0 {
		if err := verifyFileIsExecutable(vmm.binary); err != nil {
			return err
		}

		cmd := firecracker.VMCommandBuilder{}.
			WithBin(vmm.binary).
			WithSocketPath(vmm.cfg.SocketPath).
			WithStdin(os.Stdin).
			WithStdout(os.Stdout).
			WithStderr(os.Stderr).
			Build(ctx)

		machineOpts = append(machineOpts, firecracker.WithProcessRunner(cmd))
	}

	m, err := firecracker.NewMachine(vmmCtx, vmm.cfg, machineOpts...)
	if err != nil {
		return errors.Errorf("Failed creating machine: %s", err)
	}

	if vmm.metadata != nil {
		m.EnableMetadata(vmm.metadata)
	}

	log.Printf("Booting VMM now")
	if err := m.Start(vmmCtx); err != nil {
		return errors.Errorf("Failed to start machine: %v", err)
	}
	defer m.StopVMM()

	// wait for the VMM to exit
	if err := m.Wait(vmmCtx); err != nil {
		return errors.Errorf("Wait returned an error %s", err)
	}
	log.Printf("VMM has stopped successfully")
	return nil
}

// handleFifos will see if any fifos need to be generated and if a fifo log
// file should be created.
func (vmm *VMM) handleFifos(createFifoFn func(string) (*os.File, error)) (io.Writer, error) {
	// these booleans are used to check whether or not the fifo queue or metrics
	// fifo queue needs to be generated. If any which need to be generated, then
	// we know we need to create a temporary directory. Otherwise, a temporary
	// directory does not need to be created.
	generateFifoFilename := len(vmm.cfg.LogFifo) == 0
	generateMetricFifoFilename := len(vmm.cfg.MetricsFifo) == 0
	var err error
	var fifo io.WriteCloser

	if len(vmm.fifoLogFile) > 0 {
		if fifo, err = createFifoFn(vmm.fifoLogFile); err != nil {
			return nil, errors.Wrap(err, errUnableToCreateFifoLogFile.Error())
		}
		vmm.addCleanupFn(func() error {
			return fifo.Close()
		})
	}

	if generateFifoFilename || generateMetricFifoFilename {
		dir, err := ioutil.TempDir(os.TempDir(), "fcfifo")
		if err != nil {
			return nil, errors.Errorf("Fail to create temporary directory: %v", err)
		}
		vmm.addCleanupFn(func() error {
			return os.RemoveAll(dir)
		})

		if generateFifoFilename {
			vmm.cfg.LogFifo = filepath.Join(dir, "fc_fifo")
		}
		if generateMetricFifoFilename {
			vmm.cfg.MetricsFifo = filepath.Join(dir, "fc_metrics_fifo")
		}
	}
	return fifo, nil
}

// createRuntimeDir creates the runtime directory for the VM, and registers the deferred cleanup func
func (vmm *VMM) createRuntimeDir() error {
	vmdir := filepath.Join(RuntimeDir, vmm.name)
	// After execution, we can clean it up. This function is run when vmm.Cleanup() is called
	vmm.addCleanupFn(func() error {
		return os.RemoveAll(vmdir)
	})
	return os.MkdirAll(vmdir, 0755)
}

func (vmm *VMM) addCleanupFn(c func() error) {
	vmm.cleanupFns = append(vmm.cleanupFns, c)
}

// Cleanup removes temporarily used
func (vmm *VMM) Cleanup() {
	for _, closer := range vmm.cleanupFns {
		if err := closer(); err != nil {
			log.Error(err)
		}
	}
}

func createFifoFile(fifoPath string) (*os.File, error) {
	return os.OpenFile(fifoPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
}
*/
