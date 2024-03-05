package dns

import (
	"fmt"
	"net"
	"time"

	"github.com/miekg/dns"
)

// udpHandler handles DNS requests over UDP and exchanges the DNS queries to upstream DNS over TLS
func (srv *Server) udpHandler() func(writer dns.ResponseWriter, msg *dns.Msg) {
	return func(writer dns.ResponseWriter, msg *dns.Msg) {
		// create a new DNS over TLS client
		c := new(dns.Client)
		c.Net = "tcp-tls"
		c.TLSConfig = srv.Client.TLSConfig
		c.Dialer = &net.Dialer{
			Timeout: 2 * time.Second,
		}
		// log the incoming request
		srv.logger.Infow("incoming request", "name", msg.Question[0].Name,
			"type", dns.Type(msg.Question[0].Qtype).String(),
		)

		// exchange the DNS query to upstream DNS over TLS
		rsp, _, err := c.Exchange(msg, fmt.Sprintf("%v:%v", srv.Client.Address, srv.Client.Port))
		if err != nil {
			srv.logger.Errorw("failed to communicate with upstream", "error", err)
			return
		}
		if err = writer.WriteMsg(rsp); err != nil {
			srv.logger.Errorw("failed to write response", "error", err)
		}
	}
}
