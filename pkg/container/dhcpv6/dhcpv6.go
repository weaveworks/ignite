package dhcp6

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/ipv6"

	"github.com/insomniacslk/dhcp/dhcpv6"
	"github.com/insomniacslk/dhcp/dhcpv6/server6"
	"github.com/insomniacslk/dhcp/iana"
)

// ̣Borrowed from https://github.com/RedTeamPentesting/pretender
// under an MIT license

type Config struct {
	RelayIPv4      net.IP
	RelayIPv6      net.IP
	SOAHostname    string
	Interface      *net.Interface
	TTL            time.Duration
	LeaseLifetime  time.Duration
	RouterLifetime time.Duration
	LocalIPv6      net.IP
	RAPeriod       time.Duration

	NoDHCPv6DNSTakeover   bool
	NoDHCPv6              bool
	NoDNS                 bool
	NoRA                  bool
	NoMDNS                bool
	NoNetBIOS             bool
	NoLLMNR               bool
	NoLocalNameResolution bool
	NoIPv6LNR             bool

	StopAfter      time.Duration
	Verbose        bool
	NoColor        bool
	NoTimestamps   bool
	LogFileName    string
	NoHostInfo     bool
	HideIgnored    bool
	RedirectStderr bool
	ListInterfaces bool
}

// DHCPv6 default values.
const (
	// The valid lifetime for the IPv6 prefix in the option, expressed in units
	// of seconds.  A value of 0xFFFFFFFF represents infinity.
	dhcpv6DefaultValidLifetime = 60 * time.Second

	// The time at which the requesting router should contact the delegating
	// router from which the prefixes in the IA_PD were obtained to extend the
	// lifetimes of the prefixes delegated to the IA_PD; T1 is a time duration
	// relative to the current time expressed in units of seconds.
	dhcpv6T1 = 45 * time.Second

	// The time at which the requesting router should contact any available
	// delegating router to extend the lifetimes of the prefixes assigned to the
	// IA_PD; T2 is a time duration relative to the current time expressed in
	// units of seconds.
	dhcpv6T2 = 50 * time.Second
)

// dhcpv6LinkLocalPrefix is the 64-bit link local IPv6 prefix.
var dhcpv6LinkLocalPrefix = []byte{0xfe, 0x80, 0, 0, 0, 0, 0, 0}

// DHCPv6Handler holds the state for of the DHCPv6 handler method Handler().
type DHCPv6Handler struct {
	logger   *log.Logger
	serverID dhcpv6.Duid
	config   Config
}

// ListenUDPMulticast listens on a multicast group in a way that is supported on
// Unix for both IPv4 and IPv6.
func ListenUDPMulticast(iface *net.Interface, multicastGroup *net.UDPAddr) (net.PacketConn, error) {
	if multicastGroup.IP.To4() != nil {
		return net.ListenMulticastUDP("udp", iface, multicastGroup)
	}

	listenAddr := &net.UDPAddr{IP: multicastGroup.IP, Port: multicastGroup.Port}

	conn, err := net.ListenPacket("udp6", listenAddr.String())
	if err != nil {
		return nil, err
	}

	packetConn := ipv6.NewPacketConn(conn)

	err = packetConn.JoinGroup(iface, listenAddr)
	if err != nil {
		return nil, fmt.Errorf("join multicast group %s: %w", listenAddr.IP, err)
	}

	return conn, nil
}

// NewDHCPv6Handler returns a DHCPv6Handler.
func NewDHCPv6Handler(config Config, logger *log.Logger) *DHCPv6Handler {
	return &DHCPv6Handler{
		logger: logger,
		config: config,
		serverID: dhcpv6.Duid{
			Type:          dhcpv6.DUID_LL,
			HwType:        iana.HWTypeEthernet,
			LinkLayerAddr: config.Interface.HardwareAddr,
		},
	}
}

// Handler implements a server6.Handler.
func (h *DHCPv6Handler) Handler(conn net.PacketConn, peer net.Addr, m dhcpv6.DHCPv6) {
	err := h.handler(conn, peer, m)
	if err != nil {
		h.logger.Errorf(err.Error())
	}
}

func (h *DHCPv6Handler) handler(conn net.PacketConn, peerAddr net.Addr, m dhcpv6.DHCPv6) error {
	answer, err := h.createResponse(peerAddr, m)
	if errors.Is(err, errNoResponse) {
		return nil
	} else if err != nil {
		return err
	}

	_, err = conn.WriteTo(answer.ToBytes(), peerAddr)
	if err != nil {
		return fmt.Errorf("write to %s: %w", peerAddr, err)
	}

	return nil
}

var errNoResponse = fmt.Errorf("no response")

// nolint:cyclop
func (h *DHCPv6Handler) createResponse(peerAddr net.Addr, m dhcpv6.DHCPv6) (*dhcpv6.Message, error) {
	msg, err := m.GetInnerMessage()
	if err != nil {
		return nil, fmt.Errorf("get inner message: %w", err)
	}

	peer := newPeerInfo(peerAddr, msg)

	// shouldRespond, reason := shouldRespondToDHCP(h.config, peer)
	// if !shouldRespond {
	// 	h.logger.Info("Ignoring request ", m.Type().String(), peer, reason)

	// 	return nil, errNoResponse
	// }

	var answer *dhcpv6.Message

	switch m.Type() {
	case dhcpv6.MessageTypeSolicit:
		answer, err = h.handleSolicit(msg, peer)
	case dhcpv6.MessageTypeRequest, dhcpv6.MessageTypeRebind, dhcpv6.MessageTypeRenew:
		answer, err = h.handleRequestRebindRenew(msg, peer)
	case dhcpv6.MessageTypeConfirm:
		answer, err = h.handleConfirm(msg, peer)
	case dhcpv6.MessageTypeRelease:
		answer, err = h.handleRelease(msg, peer)
	case dhcpv6.MessageTypeInformationRequest:
		h.logger.Debugf("ignoring %s from %s", msg.Type(), peer)

		return nil, errNoResponse
	default:
		h.logger.Debugf("unhandled DHCP message from %s:\n%s", peer, msg.Summary())

		return nil, errNoResponse
	}

	if err != nil {
		return nil, fmt.Errorf("configure response to %T from %s: %w", msg.Type(), peer, err)
	}

	if answer == nil {
		return nil, fmt.Errorf("answer to %T from %s was not configured", msg.Type(), peer)
	}

	return answer, nil
}

func (h *DHCPv6Handler) handleSolicit(msg *dhcpv6.Message, peer peerInfo) (*dhcpv6.Message, error) {
	iaNA, err := extractIANA(msg)
	if err != nil {
		return nil, fmt.Errorf("extract IANA: %w", err)
	}

	ip, opts, err := h.configureResponseOpts(iaNA, msg, peer)
	if err != nil {
		return nil, fmt.Errorf("configure response options: %w", err)
	}

	answer, err := dhcpv6.NewAdvertiseFromSolicit(msg, opts...)
	if err != nil {
		return nil, fmt.Errorf("create ADVERTISE: %w", err)
	}

	h.logger.Debug(msg.Type(), peer, ip)
	h.logger.Debug(answer)

	return answer, nil
}

func (h *DHCPv6Handler) handleRequestRebindRenew(msg *dhcpv6.Message, peer peerInfo) (*dhcpv6.Message, error) {
	iaNA, err := extractIANA(msg)
	if err != nil {
		return nil, fmt.Errorf("extract IANA: %w", err)
	}

	ip, opts, err := h.configureResponseOpts(iaNA, msg, peer)
	if err != nil {
		return nil, fmt.Errorf("configure response options: %w", err)
	}

	answer, err := dhcpv6.NewReplyFromMessage(msg, opts...)
	if err != nil {
		return nil, fmt.Errorf("create REPLY: %w", err)
	}

	h.logger.Debug(msg.Type(), peer, ip)

	return answer, nil
}

func (h *DHCPv6Handler) handleConfirm(msg *dhcpv6.Message, peer peerInfo) (*dhcpv6.Message, error) {
	answer, err := dhcpv6.NewReplyFromMessage(msg,
		dhcpv6.WithServerID(h.serverID),
		dhcpv6.WithDNS(h.config.LocalIPv6),
		dhcpv6.WithOption(&dhcpv6.OptStatusCode{
			StatusCode:    iana.StatusNotOnLink,
			StatusMessage: iana.StatusNotOnLink.String(),
		}))
	if err != nil {
		return nil, fmt.Errorf("create REPLY: %w", err)
	}

	h.logger.Debugf("rejecting %s from %s", msg.Type().String(), peer)

	return answer, nil
}

func (h *DHCPv6Handler) handleRelease(msg *dhcpv6.Message, peer peerInfo) (*dhcpv6.Message, error) {
	iaNAs, err := extractIANAs(msg)
	if err != nil {
		return nil, err
	}

	opts := []dhcpv6.Modifier{
		dhcpv6.WithOption(&dhcpv6.OptStatusCode{
			StatusCode:    iana.StatusSuccess,
			StatusMessage: iana.StatusSuccess.String(),
		}),
		dhcpv6.WithServerID(h.serverID),
	}

	// send status NoBinding for each address
	for _, iaNA := range iaNAs {
		opts = append(opts, dhcpv6.WithOption(&dhcpv6.OptIANA{
			IaId: iaNA.IaId,
			Options: dhcpv6.IdentityOptions{
				Options: []dhcpv6.Option{
					&dhcpv6.OptStatusCode{
						StatusCode:    iana.StatusNoBinding,
						StatusMessage: iana.StatusNoBinding.String(),
					},
				},
			},
		}))
	}

	answer, err := dhcpv6.NewReplyFromMessage(msg, opts...)
	if err != nil {
		return nil, fmt.Errorf("create REPLY: %w", err)
	}

	h.logger.Debugf("aggreeing to RELEASE from %s", peer)

	return answer, nil
}

// configureResponseOpts returns the IP that should be assigned based on the
// request IA_NA and the modifiers to configure the response with that IP and
// the DNS server configured in the DHCPv6Handler.
func (h *DHCPv6Handler) configureResponseOpts(requestIANA *dhcpv6.OptIANA,
	msg *dhcpv6.Message, peer peerInfo,
) (net.IP, []dhcpv6.Modifier, error) {
	cid := msg.GetOneOption(dhcpv6.OptionClientID)
	if cid == nil {
		return nil, nil, fmt.Errorf("no client ID option from DHCPv6 message")
	}

	duid, err := dhcpv6.DuidFromBytes(cid.ToBytes())
	if err != nil {
		return nil, nil, fmt.Errorf("deserialize DUI")
	}

	var leasedIP net.IP

	if duid.LinkLayerAddr == nil {
		h.logger.Debugf("DUID does not contain link layer address")

		randomIP, err := generateDeterministicRandomAddress(peer.IP)
		if err != nil {
			h.logger.Debugf("could not generate deterministic address (using SLAAC IP instead): %v", err)

			leasedIP = peer.IP
		} else {
			leasedIP = randomIP
		}
	} else {
		if h.logger != nil {
			go h.logger.Debug(peer.IP, duid.LinkLayerAddr)
		}

		leasedIP = append(leasedIP, dhcpv6LinkLocalPrefix...)
		leasedIP = append(leasedIP, 0, 0)
		leasedIP = append(leasedIP, duid.LinkLayerAddr...)
	}

	// if the IP has the first bit after the prefix set, Windows won't route
	// queries via this IP and use the regular self-generated link-local address
	// instead.
	leasedIP[8] |= 0b10000000

	return leasedIP, []dhcpv6.Modifier{
		dhcpv6.WithServerID(h.serverID),
		dhcpv6.WithDNS(h.config.LocalIPv6),
		dhcpv6.WithOption(&dhcpv6.OptIANA{
			IaId: requestIANA.IaId,
			T1:   dhcpv6T1,
			T2:   dhcpv6T2,
			Options: dhcpv6.IdentityOptions{
				Options: []dhcpv6.Option{
					&dhcpv6.OptIAAddress{
						IPv6Addr:          leasedIP,
						PreferredLifetime: h.config.LeaseLifetime,
						ValidLifetime:     h.config.LeaseLifetime,
					},
				},
			},
		}),
	}, nil
}

func generateDeterministicRandomAddress(peer net.IP) (net.IP, error) {
	if len(peer) != net.IPv6len {
		return nil, fmt.Errorf("invalid length of IPv6 address: %d bytes", len(peer))
	}

	prefixLength := net.IPv6len / 2 // nolint:gomnd

	seed := binary.LittleEndian.Uint64(peer[prefixLength:])

	deterministicAddress := make([]byte, prefixLength)

	n, err := rand.New(rand.NewSource(int64(seed))).Read(deterministicAddress) // nolint:gosec
	if err != nil {
		return nil, err
	}

	if n != prefixLength {
		return nil, fmt.Errorf("read %d random bytes instead of %d", n, prefixLength)
	}

	var newIP net.IP
	newIP = append(newIP, dhcpv6LinkLocalPrefix...)
	newIP = append(newIP, deterministicAddress...)

	return newIP, nil
}

func extractIANA(innerMessage *dhcpv6.Message) (*dhcpv6.OptIANA, error) {
	iaNAOpt := innerMessage.GetOneOption(dhcpv6.OptionIANA)
	if iaNAOpt == nil {
		return nil, fmt.Errorf("message does not contain IANA:\n%s", innerMessage.Summary())
	}

	iaNA, ok := iaNAOpt.(*dhcpv6.OptIANA)
	if !ok {
		return nil, fmt.Errorf("unexpected type for IANA option: %T", iaNAOpt)
	}

	return iaNA, nil
}

func extractIANAs(innerMessage *dhcpv6.Message) ([]*dhcpv6.OptIANA, error) {
	iaNAOpts := innerMessage.GetOption(dhcpv6.OptionIANA)
	if iaNAOpts == nil {
		return nil, fmt.Errorf("message does not contain IANAs:\n%s", innerMessage.Summary())
	}

	iaNAs := make([]*dhcpv6.OptIANA, 0, len(iaNAOpts))

	for i, iaNAOpt := range iaNAOpts {
		iaNA, ok := iaNAOpt.(*dhcpv6.OptIANA)
		if !ok {
			return nil, fmt.Errorf("unexpected type for IANA option %d: %T", i, iaNAOpt)
		}

		iaNAs = append(iaNAs, iaNA)
	}

	return iaNAs, nil
}

// RunDHCPv6Server starts a DHCPv6 server which assigns a DNS server.
func RunDHCPv6Server(ctx context.Context, logger *log.Logger, config Config) error {
	listenAddr := &net.UDPAddr{
		IP:   dhcpv6.AllDHCPRelayAgentsAndServers,
		Port: dhcpv6.DefaultServerPort,
		Zone: config.Interface.Name,
	}

	dhcvpv6Handler := NewDHCPv6Handler(config, logger)

	conn, err := ListenUDPMulticast(config.Interface, listenAddr)
	if err != nil {
		return err
	}

	server, err := server6.NewServer(config.Interface.Name, nil, dhcvpv6Handler.Handler,
		server6.WithConn(conn))
	if err != nil {
		return fmt.Errorf("starting DHCPv6 server: %w", err)
	}

	go func() {
		<-ctx.Done()

		_ = server.Close()
	}()

	logger.Infof("listening via UDP on %s", listenAddr)

	err = server.Serve()

	// if the server is stopped via ctx, we suppress the resulting errors that
	// result from server.Close closing the connection.
	if ctx.Err() != nil {
		return nil
	}

	return err
}

type peerInfo struct {
	IP        net.IP
	Hostnames []string
}

func newPeerInfo(addr net.Addr, innerMessage *dhcpv6.Message) peerInfo {
	p := peerInfo{
		IP: addrToIP(addr),
	}

	fqdnOpt := innerMessage.GetOneOption(dhcpv6.OptionFQDN)
	if fqdnOpt == nil {
		return p
	}

	fqdn, ok := fqdnOpt.(*dhcpv6.OptFQDN)
	if !ok {
		return p
	}

	// workaround, because the DHCP library seems not be be able to decode
	// simple length+label values
	if len(fqdn.DomainName.Labels) == 0 {
		rawLabel := fqdn.DomainName.ToBytes()
		if len(rawLabel) > 2 && int(rawLabel[0]) == (len(rawLabel)-1) {
			fqdn.DomainName.Labels = append(fqdn.DomainName.Labels, string(rawLabel[1:]))
		}
	}

	p.Hostnames = fqdn.DomainName.Labels

	return p
}

// String returns the string representation of a peerInfo.
func (p peerInfo) String() string {
	if len(p.Hostnames) > 0 {
		return p.IP.String() + " (" + strings.Join(p.Hostnames, ", ") + ")"
	}

	return p.IP.String()
}

func addrToIP(addr net.Addr) net.IP {
	udpAddr, ok := addr.(*net.UDPAddr)
	if ok {
		return udpAddr.IP
	}

	addrString := addr.String()

	for strings.Contains(addrString, "/") || strings.Contains(addrString, "%") {
		addrString = strings.SplitN(addrString, "/", 2)[0] // nolint:gomnd
		addrString = strings.SplitN(addrString, "%", 2)[0] // nolint:gomnd
	}

	splitAddr, _, err := net.SplitHostPort(addrString)
	if err == nil {
		addrString = splitAddr
	}

	return net.ParseIP(addrString)
}
