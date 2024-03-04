package main

import (
	"fmt"

	"github.com/13excite/dns-proxy/pkg/config"
	"github.com/13excite/dns-proxy/pkg/dns"
	"github.com/13excite/dns-proxy/pkg/logger"
)

func main() {

	c := &config.Config{}
	c.Defaults()

	logger.InitLogger(c)

	dnsServer := dns.NewServer("tcp", c)
	err := dnsServer.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}
