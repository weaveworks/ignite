package preflight

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"path"
	"path/filepath"
	"strings"

	"k8s.io/apimachinery/pkg/util/sets"

	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	"github.com/weaveworks/ignite/pkg/constants"
	"github.com/weaveworks/ignite/pkg/operations/lookup"
	"github.com/weaveworks/ignite/pkg/providers"
)

const (
	oldPathString  = "/"
	newPathString  = "-"
	noReplaceLimit = -1
)

type Checker interface {
	Check() error
	Name() string
	Type() string
}

type PortOpenChecker struct {
	port uint64
}

func (poc PortOpenChecker) Check() error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", poc.port))
	if err != nil {
		return fmt.Errorf("Port %d is in use", poc.port)
	}
	if err := listener.Close(); err != nil {
		return fmt.Errorf("Port %d is in used, failed to close it", poc.port)
	}
	return nil
}

func (poc PortOpenChecker) Name() string {
	return fmt.Sprintf("Port-%d", poc.port)
}

func (poc PortOpenChecker) Type() string {
	return "Port"
}

type ExistingFileChecker struct {
	filePath string
}

func (efc ExistingFileChecker) Check() error {
	if _, err := os.Stat(efc.filePath); os.IsNotExist(err) {
		return fmt.Errorf("File %s, does not exist", efc.filePath)
	}
	return nil
}

func (efc ExistingFileChecker) Name() string {
	return fmt.Sprintf("ExistingFile-%s", strings.Replace(efc.filePath, oldPathString, newPathString, noReplaceLimit))
}

func (efc ExistingFileChecker) Type() string {
	return "ExistingFile"
}

type AvailablePathChecker struct {
	path string
}

func (apc AvailablePathChecker) Check() error {
	if _, err := os.Stat(apc.path); !os.IsExist(err) {
		return fmt.Errorf("Path %s, already exist", apc.path)
	}
	return nil
}

func (apc AvailablePathChecker) Name() string {
	return fmt.Sprintf("AvailablePath-%s", strings.Replace(apc.path, oldPathString, newPathString, noReplaceLimit))
}

func (apc AvailablePathChecker) Type() string {
	return "AvailablePath"
}

func StartCmdChecks(vm *api.VM, ignoredPreflightErrors sets.String) error {
	checks := []Checker{}
	kernelUID, err := lookup.KernelUIDForVM(vm, providers.Client)
	if err != nil {
		return err
	}
	vmDir := filepath.Join(constants.VM_DIR, vm.GetUID().String())
	kernelDir := filepath.Join(constants.KERNEL_DIR, kernelUID.String())
	log.Println(vm.GetUID().String())
	checks = append(checks, ExistingFileChecker{filePath: path.Join(vmDir, constants.METADATA)})
	checks = append(checks, ExistingFileChecker{filePath: path.Join(kernelDir, constants.KERNEL_FILE)})
	checks = append(checks, ExistingFileChecker{filePath: "/dev/mapper/control"})
	checks = append(checks, ExistingFileChecker{filePath: "/dev/net/tun"})
	checks = append(checks, ExistingFileChecker{filePath: "/dev/kvm"})
	for _, port := range vm.Spec.Network.Ports {
		checks = append(checks, PortOpenChecker{port: port.HostPort})
	}
	return runChecks(checks, ignoredPreflightErrors)
}

func runChecks(checks []Checker, ignoredPreflightErrors sets.String) error {
	var errBuffer bytes.Buffer

	for _, check := range checks {
		name := check.Name()
		checkType := check.Type()

		err := check.Check()
		if isIgnoredPreflightError(ignoredPreflightErrors, checkType) {
			err = nil
		}

		if err != nil {
			errBuffer.WriteString(fmt.Sprintf("[ERROR %s]: %v\n", name, err))
		}
	}
	if errBuffer.Len() > 0 {
		return fmt.Errorf(errBuffer.String())
	}
	return nil
}

func isIgnoredPreflightError(ignoredPreflightError sets.String, checkType string) bool {
	return ignoredPreflightError.Has("all") || ignoredPreflightError.Has(strings.ToLower(checkType))
}
