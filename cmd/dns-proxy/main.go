package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/13excite/dns-proxy/pkg/config"
	"github.com/13excite/dns-proxy/pkg/dns"
	"github.com/13excite/dns-proxy/pkg/logger"
)

func main() {
	tcpMode := flag.Bool("tcp", false, "tcp mode")
	udpMode := flag.Bool("udp", false, "udp mode")
	configPath := flag.String("config", "", "path to the cfg file (optional)")
	flag.Parse()

	// check if the user has specified a mode
	if !*tcpMode && !*udpMode {
		fmt.Println("You must specify a mode: tcp or udp")
		flag.Usage()
		os.Exit(1)
	}
	c := &config.Config{}
	c.Defaults()

	if *configPath != "" {
		c.ReadConfigFile(*configPath)
	}
	// initialize the logger
	err := logger.InitLogger(c)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// create a new DNS server in the specified mode
	if *tcpMode {
		// create a new server and listen on the TCP network address
		dnsServer := dns.NewServer("tcp", c)
		err = dnsServer.ListenAndServe()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	if *udpMode {
		fmt.Println("udp")
		dnsServer := dns.NewServer("udp", c)
		err = dnsServer.ListenAndServe()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
