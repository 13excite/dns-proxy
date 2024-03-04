package dns

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"

	"github.com/13excite/dns-proxy/pkg/config"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"go.uber.org/zap"
)

type Server struct {
	Addr   string
	Port   int
	Net    string // "tcp" or "ud"
	Client *Client
	logger *zap.SugaredLogger
}

func NewServer(netPrefix string, config *config.Config) *Server {
	return &Server{
		Addr:   config.Addr,
		Port:   config.Port,
		Net:    netPrefix,
		Client: DefaultClinet(),
		logger: zap.S().With("package", "dns-server"),
	}
}

// ListenAndServe listens on the TCP network address srv.Addr and then
// calls Serve to handle requests on incoming connections
func (srv *Server) ListenAndServe() error {
	srv.logger.Infow("Server is starting")
	ip := net.ParseIP(srv.Addr)
	if ip == nil {
		srv.logger.Errorw("failed to parse address", "address", srv.Addr)
		return fmt.Errorf("failed to parse address: %s", srv.Addr)
	}
	switch srv.Net {
	case "tcp":
		listner, err := net.ListenTCP("tcp", &net.TCPAddr{IP: ip, Port: srv.Port})
		if err != nil {
			return err
		}
		srv.logger.Infow("started", "address", srv.Addr, "port", srv.Port, "network", srv.Net)
		return srv.serveTCP(listner)

	case "udp":
		fmt.Println("udp")
		return nil
	}
	return fmt.Errorf("unsupported network type: %s", srv.Net)
}

func (srv *Server) serveTCP(listner *net.TCPListener) error {
	defer listner.Close()

	for {
		downStreamConn, err := listner.Accept()
		if err != nil {
			// if the error is a timeout, just continue to the next iteration
			if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
				srv.logger.Debugw("timeout", "error", err)
				continue
			}
			return err
		}

		go func(downStreamConn net.Conn) {
			defer downStreamConn.Close()

			tbuff := make([]byte, 4096)
			size, err := downStreamConn.Read(tbuff)
			if err != nil {
				srv.logger.Errorw("failed to read from connection", "error", err)
				return
			}

			err = srv.validateDNSReq(tbuff, size)
			if err != nil {
				return
			}

			dnsConn, err := srv.Client.Dial()
			if err != nil {
				fmt.Println(err)
				return
			}
			defer dnsConn.Close()

			_, err = dnsConn.Write(tbuff[:size])
			if err != nil {
				srv.logger.Errorw("failed to write to the clinet connection", "error", err)
				return
			}

			// straight forward copy response  from dnsConn to downStreamConn
			_, err = io.Copy(downStreamConn, dnsConn)
			if err != nil {
				fmt.Println(err)
				return
			}

		}(downStreamConn)
	}
}

func (srv *Server) validateDNSReq(b []byte, size int) error {
	DNSMsgLength := binary.BigEndian.Uint16(b[:2])
	if int(DNSMsgLength) != size-2 {
		srv.logger.Errorw("Size incorrect", "expected", DNSMsgLength, "actual", size)
		return fmt.Errorf("Size incorrect: %d, %d", DNSMsgLength, size)
	}
	srv.logger.Debugw("Size correct", "size", DNSMsgLength)

	// read the DNS packet with gopacket
	packet := gopacket.NewPacket(b[2:], layers.LayerTypeDNS, gopacket.Default)
	dnsPacket := packet.Layer(layers.LayerTypeDNS)
	// cast the packet to a DNS layer and check if it's valid
	tcp, ok := dnsPacket.(*layers.DNS)
	if !ok {
		srv.logger.Errorw("cannot cast dns packet to layers.DNS")
		return fmt.Errorf("cannot cast dns packet to layers.DNS")
	}

	srv.logger.Infow("incoming request", "name", string(tcp.Questions[0].Name), "type", tcp.Questions[0].Type)

	return nil
}
