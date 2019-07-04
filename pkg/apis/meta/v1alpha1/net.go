package v1alpha1

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

// PortMapping defines a port mapping between the VM and the host
type PortMapping struct {
	HostPort uint64 `json:"hostPort"`
	VMPort   uint64 `json:"vmPort"`
}

var _ fmt.Stringer = PortMapping{}

func (p PortMapping) String() string {
	return fmt.Sprintf("0.0.0.0:%d->%d", p.HostPort, p.VMPort)
}

// PortMappings represents a list of port mappings
type PortMappings []PortMapping

var _ fmt.Stringer = PortMappings{}

func ParsePortMappings(input []string) (PortMappings, error) {
	result := make(PortMappings, 0, len(input))

	for _, portMapping := range input {
		ports := strings.Split(portMapping, ":")
		if len(ports) != 2 {
			return nil, fmt.Errorf("port mappings must be of form <host port>:<VM port>")
		}

		hostPort, err := strconv.ParseUint(ports[0], 10, 64)
		if err != nil {
			return nil, err
		}

		vmPort, err := strconv.ParseUint(ports[1], 10, 64)
		if err != nil {
			return nil, err
		}

		for _, portMapping := range result {
			if portMapping.HostPort == hostPort {
				return nil, fmt.Errorf("cannot use a port on the host twice")
			}
		}

		result = append(result, PortMapping{
			HostPort: hostPort,
			VMPort:   vmPort,
		})
	}

	return result, nil
}

func (p PortMappings) String() string {
	var sb strings.Builder
	var index int

	for _, portMapping := range p {
		sb.WriteString(portMapping.String())

		index++
		if index < len(p) {
			sb.WriteString(", ")
		}
	}

	return sb.String()
}

// IPAddresses represents a list of VM IP addresses
type IPAddresses []net.IP

var _ fmt.Stringer = IPAddresses{}

func (i IPAddresses) String() string {
	var sb strings.Builder
	var index int

	for _, ip := range i {
		sb.WriteString(ip.String())

		index++
		if index < len(i) {
			sb.WriteString(", ")
		}
	}

	return sb.String()
}
