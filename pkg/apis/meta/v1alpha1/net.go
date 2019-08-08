package v1alpha1

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/docker/go-connections/nat"
)

// PortMapping defines a port mapping between the VM and the host
type PortMapping struct {
	BindAddress net.IP   `json:"bindAddress,omitempty"`
	HostPort    uint64   `json:"hostPort"`
	VMPort      uint64   `json:"vmPort"`
	Protocol    Protocol `json:"protocol,omitempty"`
}

var _ fmt.Stringer = PortMapping{}

func (p PortMapping) String() string {
	var sb strings.Builder

	if p.BindAddress != nil {
		sb.WriteString(p.BindAddress.String())
	} else {
		sb.WriteString("0.0.0.0")
	}

	sb.WriteString(fmt.Sprintf(":%d->%d", p.HostPort, p.VMPort))

	if len(p.Protocol) > 0 {
		sb.WriteString(fmt.Sprintf("/%s", p.Protocol))
	}

	return sb.String()
}

// PortMappings represents a list of port mappings
type PortMappings []PortMapping

var _ fmt.Stringer = PortMappings{}

func ParsePortMappings(input []string) (PortMappings, error) {
	result := make(PortMappings, 0, len(input))

	_, bindings, err := nat.ParsePortSpecs(input)
	if err != nil {
		return nil, err
	}

	for port, bindings := range bindings {
		if len(bindings) > 1 {
			// TODO: For now only support mapping a VM port to a single host IP/port
			return nil, fmt.Errorf("only one host binding per VM binding supported for now, received %d", len(bindings))
		}

		binding := bindings[0]
		var err error
		var bindAddress net.IP
		var hostPort uint64
		var vmPort uint64
		var protocol Protocol

		if len(binding.HostIP) > 0 {
			if bindAddress = net.ParseIP(binding.HostIP); bindAddress == nil {
				return nil, fmt.Errorf("invalid bind address: %q", binding.HostIP)
			}
		}

		if hostPort, err = strconv.ParseUint(binding.HostPort, 10, 64); err != nil {
			return nil, fmt.Errorf("invalid host port: %q", binding.HostPort)
		}

		if vmPort, err = strconv.ParseUint(port.Port(), 10, 64); err != nil {
			return nil, fmt.Errorf("invalid VM port: %q", port.Port())
		}

		if protocol, err = protocolFromString(port.Proto()); err != nil {
			return nil, err
		}

		mapping := PortMapping{
			BindAddress: bindAddress,
			HostPort:    hostPort,
			VMPort:      vmPort,
			Protocol:    protocol,
		}

		for _, portMapping := range result {
			if portMapping.HostPort == mapping.HostPort && portMapping.Protocol == mapping.Protocol {
				return nil, fmt.Errorf("cannot use a port/protocol combination on the host twice")
			}
		}

		result = append(result, mapping)
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

// Protocol specifies a network port protocol
type Protocol string

const (
	ProtocolTCP Protocol = "tcp"
	ProtocolUDP Protocol = "udp"
)

var _ fmt.Stringer = Protocol("")

func protocolFromString(input string) (Protocol, error) {
	for _, protocol := range []Protocol{ProtocolTCP, ProtocolUDP} {
		if protocol.String() == input {
			return protocol, nil
		}
	}

	return "", fmt.Errorf("invalid protocol: %q", input)
}

func (p Protocol) String() string {
	return string(p)
}

func (p *Protocol) UnmarshalJSON(b []byte) (err error) {
	var s string
	if err = json.Unmarshal(b, &s); err != nil {
		return err
	}

	*p, err = protocolFromString(s)
	return
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
