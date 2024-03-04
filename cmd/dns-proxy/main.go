package main

import (
	"fmt"
	"os"

	"github.com/13excite/dns-proxy/pkg/config"
	"github.com/13excite/dns-proxy/pkg/dns"
	"github.com/13excite/dns-proxy/pkg/logger"
)

func main() {
	// simple subcommand dispatcher with cobra or pflag
	if len(os.Args) < 2 {
		fmt.Println("expected 'tcp' or 'udp' subcommands")
		os.Exit(1)
	}
	switch os.Args[1] {
	case "tcp":
		c := &config.Config{}
		c.Defaults()

		err := logger.InitLogger(c)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// create a new server and listen on the TCP network address
		dnsServer := dns.NewServer("tcp", c)
		err = dnsServer.ListenAndServe()
		if err != nil {
			fmt.Println(err)
		}
	case "udp":
		fmt.Println("udp")
	default:
		fmt.Println("expected 'tcp' or 'udp' subcommands")
		os.Exit(1)
	}
}
