package dhcpv6

import (
	"fmt"
	"net"

	"github.com/insomniacslk/dhcp/rfc1035label"
	"github.com/u-root/uio/uio"
)

// NTPSuboptionSrvAddr is NTP_SUBOPTION_SRV_ADDR according to RFC 5908.
type NTPSuboptionSrvAddr net.IP

// Code returns the suboption code.
func (n *NTPSuboptionSrvAddr) Code() OptionCode {
	return NTPSuboptionSrvAddrCode
}

// ToBytes returns the byte serialization of the suboption.
func (n *NTPSuboptionSrvAddr) ToBytes() []byte {
	buf := uio.NewBigEndianBuffer(nil)
	buf.Write16(uint16(NTPSuboptionSrvAddrCode))
	buf.Write16(uint16(net.IPv6len))
	buf.WriteBytes(net.IP(*n).To16())
	return buf.Data()
}

func (n *NTPSuboptionSrvAddr) String() string {
	return fmt.Sprintf("Server Address: %s", net.IP(*n).String())
}

// NTPSuboptionMCAddr is NTP_SUBOPTION_MC_ADDR according to RFC 5908.
type NTPSuboptionMCAddr net.IP

// Code returns the suboption code.
func (n *NTPSuboptionMCAddr) Code() OptionCode {
	return NTPSuboptionMCAddrCode
}

// ToBytes returns the byte serialization of the suboption.
func (n *NTPSuboptionMCAddr) ToBytes() []byte {
	buf := uio.NewBigEndianBuffer(nil)
	buf.Write16(uint16(NTPSuboptionMCAddrCode))
	buf.Write16(uint16(net.IPv6len))
	buf.WriteBytes(net.IP(*n).To16())
	return buf.Data()
}

func (n *NTPSuboptionMCAddr) String() string {
	return fmt.Sprintf("Multicast Address: %s", net.IP(*n).String())
}

// NTPSuboptionSrvFQDN is NTP_SUBOPTION_SRV_FQDN according to RFC 5908.
type NTPSuboptionSrvFQDN rfc1035label.Labels

// Code returns the suboption code.
func (n *NTPSuboptionSrvFQDN) Code() OptionCode {
	return NTPSuboptionSrvFQDNCode
}

// ToBytes returns the byte serialization of the suboption.
func (n *NTPSuboptionSrvFQDN) ToBytes() []byte {
	buf := uio.NewBigEndianBuffer(nil)
	buf.Write16(uint16(NTPSuboptionSrvFQDNCode))
	l := rfc1035label.Labels(*n)
	buf.Write16(uint16(l.Length()))
	buf.WriteBytes(l.ToBytes())
	return buf.Data()
}

func (n *NTPSuboptionSrvFQDN) String() string {
	l := rfc1035label.Labels(*n)
	return fmt.Sprintf("Server FQDN: %s", l.String())
}

// NTPSuboptionSrvAddr is the value of NTP_SUBOPTION_SRV_ADDR according to RFC 5908.
const (
	NTPSuboptionSrvAddrCode = OptionCode(1)
	NTPSuboptionMCAddrCode  = OptionCode(2)
	NTPSuboptionSrvFQDNCode = OptionCode(3)
)

// parseNTPSuboption implements the OptionParser interface.
func parseNTPSuboption(code OptionCode, data []byte) (Option, error) {
	//var o Options
	buf := uio.NewBigEndianBuffer(data)
	length := len(data)
	data, err := buf.ReadN(length)
	if err != nil {
		return nil, fmt.Errorf("failed to read %d bytes for suboption: %w", length, err)
	}
	switch code {
	case NTPSuboptionSrvAddrCode, NTPSuboptionMCAddrCode:
		if length != net.IPv6len {
			return nil, fmt.Errorf("invalid suboption length, want %d, got %d", net.IPv6len, length)
		}
		var so Option
		switch code {
		case NTPSuboptionSrvAddrCode:
			sos := NTPSuboptionSrvAddr(data)
			so = &sos
		case NTPSuboptionMCAddrCode:
			som := NTPSuboptionMCAddr(data)
			so = &som
		}
		return so, nil
	case NTPSuboptionSrvFQDNCode:
		l, err := rfc1035label.FromBytes(data)
		if err != nil {
			return nil, fmt.Errorf("failed to parse rfc1035 labels: %w", err)
		}
		// TODO according to rfc3315, this label must not be compressed.
		// Need to add support for compression detection to the
		// `rfc1035label` package in order to do that.
		so := NTPSuboptionSrvFQDN(*l)
		return &so, nil
	default:
		gopt := OptionGeneric{OptionCode: code, OptionData: data}
		return &gopt, nil
	}
}

// ParseOptNTPServer parses a sequence of bytes into an OptNTPServer object.
func ParseOptNTPServer(data []byte) (*OptNTPServer, error) {
	var so Options
	if err := so.FromBytesWithParser(data, parseNTPSuboption); err != nil {
		return nil, err
	}
	return &OptNTPServer{
		Suboptions: so,
	}, nil
}

// OptNTPServer is an option NTP server as defined by RFC 5908.
type OptNTPServer struct {
	Suboptions Options
}

// Code returns the option code
func (op *OptNTPServer) Code() OptionCode {
	return OptionNTPServer
}

// ToBytes returns the option serialized to bytes.
func (op *OptNTPServer) ToBytes() []byte {
	buf := uio.NewBigEndianBuffer(nil)
	for _, so := range op.Suboptions {
		buf.WriteBytes(so.ToBytes())
	}
	return buf.Data()
}

func (op *OptNTPServer) String() string {
	return fmt.Sprintf("NTP: %v", op.Suboptions)
}
