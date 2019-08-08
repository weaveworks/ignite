package v1alpha1

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"
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

var errInvalidPortMappingFormat = fmt.Errorf("port mappings must be of form [<bind address>:]<host port>:<VM port>[/<protocol>]")

func ParsePortMappings(input []string) (PortMappings, error) {
	result := make(PortMappings, 0, len(input))

	for _, portMapping := range input {
		ports := strings.Split(portMapping, ":")
		if len(ports) > 3 || len(ports) < 2 {
			return nil, errInvalidPortMappingFormat
		}

		var bindAddress net.IP
		var protocol Protocol
		offset := 0

		if len(ports) == 3 {
			offset = 1

			if bindAddress = net.ParseIP(ports[0]); bindAddress == nil {
				return nil, errInvalidPortMappingFormat
			}
		}

		hostPort, err := strconv.ParseUint(ports[0+offset], 10, 64)
		if err != nil {
			return nil, err
		}

		proto := strings.Split(ports[1+offset], "/")

		if len(proto) > 2 {
			return nil, errInvalidPortMappingFormat
		}

		if len(proto) == 2 {
			if protocol, err = protocolFromString(proto[1]); err != nil {
				return nil, err
			}
		}

		vmPort, err := strconv.ParseUint(proto[0], 10, 64)
		if err != nil {
			return nil, err
		}

		for _, portMapping := range result {
			if portMapping.HostPort == hostPort && portMapping.Protocol == protocol {
				return nil, fmt.Errorf("cannot use a port/protocol combination on the host twice")
			}
		}

		// TODO: Check for duplicate VM ports

		result = append(result, PortMapping{
			BindAddress: bindAddress,
			HostPort:    hostPort,
			VMPort:      vmPort,
			Protocol:    protocol,
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

func (p Protocol) MarshalJSON() ([]byte, error) {
	return json.Marshal(p)
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
