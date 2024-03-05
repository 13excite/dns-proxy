package dns

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"github.com/13excite/dns-proxy/pkg/config"
	"go.uber.org/zap"
)

// Client is a DNS over TLS client
type Client struct {
	Address   string // DNS over TLS address
	Port      int    // DNS over TLS port
	TLSConfig *tls.Config
	Dialer    *net.Dialer
	logger    *zap.SugaredLogger
}

func DefaultClinet() *Client {
	return &Client{
		TLSConfig: &tls.Config{
			ServerName: "one.one.one.one",
			MinVersion: tls.VersionTLS12,
			// think about algorithm supported by TLS
			// CipherSuites: []uint16{tls.TLS_AES_256_GCM_SHA384},
		},
		Dialer: &net.Dialer{
			Timeout: 2 * time.Second,
		},
		Address: "1.1.1.1",
		Port:    853,
		logger:  zap.S().With("package", "dns-client"),
	}
}

func NewClient(conf *config.Config) *Client {
	return &Client{
		TLSConfig: &tls.Config{
			ServerName: conf.Upstream.Hostname,
			MinVersion: tls.VersionTLS12,
		},
		Dialer: &net.Dialer{
			Timeout: time.Duration(conf.Upstream.Timeout) * time.Second,
		},
		Address: conf.Upstream.Address,
		Port:    conf.Upstream.Port,
		logger:  zap.S().With("package", "dns-client"),
	}
}

// Dial connects to the DNS over TLS
// TODO: need to implement mock tls server for testing
func (c *Client) Dial() (conn *tls.Conn, err error) {
	c.logger.Debugw("dialing", "address", c.Address, "port", c.Port)
	conn, err = tls.DialWithDialer(c.Dialer, "tcp", fmt.Sprintf("%s:%d", c.Address, c.Port), c.TLSConfig)
	if err != nil {
		c.logger.Errorw("failed to dial", "error", err)
		return nil, err
	}
	return conn, nil
}
