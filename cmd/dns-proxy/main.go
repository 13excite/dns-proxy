package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

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

	logger := zap.S().With("package", "cmd")

	// create the Group
	ctx := context.Background()
	group, _ := errgroup.WithContext(ctx)

	// create a new DNS server in the specified mode
	if *tcpMode {
		// create a new server and listen on the TCP network address
		dnsServer := dns.NewServer("tcp", c)
		group.Go(
			dnsServer.ListenAndServe,
		)
	}
	if *udpMode {
		dnsServer := dns.NewServer("udp", c)
		group.Go(
			dnsServer.ListenAndServe,
		)
	}

	err = group.Wait()
	if err != nil {
		logger.Errorw("waitgroup returned an error: %w", err)
	}
}
