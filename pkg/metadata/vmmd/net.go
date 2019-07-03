package vmmd

import (
	"fmt"
	"strconv"
	"strings"

	api "github.com/weaveworks/ignite/pkg/apis/ignite/v1alpha1"
)

type PortMappings []api.PortMapping

var _ fmt.Stringer = &PortMappings{}

func (pm *PortMappings) String() string {
	var sb strings.Builder
	var index int

	for _, portMapping := range *pm {
		sb.WriteString(fmt.Sprintf("0.0.0.0:%d->%d", portMapping.HostPort, portMapping.VMPort))

		index++
		if index < len(*pm) {
			sb.WriteString(", ")
		}
	}

	return sb.String()
}

func (md *VM) NewPortMappings(input []string) error {
	result := PortMappings{}

	for _, portMapping := range input {
		ports := strings.Split(portMapping, ":")
		if len(ports) != 2 {
			return fmt.Errorf("port mappings must be of form <host port>:<VM port>")
		}

		hostPort, err := strconv.ParseUint(ports[0], 10, 64)
		if err != nil {
			return err
		}

		vmPort, err := strconv.ParseUint(ports[1], 10, 64)
		if err != nil {
			return err
		}

		for _, portMapping := range result {
			if portMapping.HostPort == hostPort {
				return fmt.Errorf("cannot use a port on the host twice")
			}
		}

		result[hostPort] = api.PortMapping{
			HostPort: hostPort,
			VMPort:   vmPort,
		}
	}

	md.Spec.Ports = result
	return nil
}

func (md *VM) ClearPortMappings() {
	md.Spec.Ports = nil
}
