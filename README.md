# Simple DNS Proxy

Project contains a simple DNS proxy that listens to conventional DNS
and sends it over TLS.

## Table of Contents

- [Implementation Details](#implementation-details)
- [Getting Started](#getting-started)
- [Requirements](#requirements)

## Implementation Details

The DNS proxy listens for DNS requests on a specified port (by default it's
1153 udp or tcp) and forwards these requests securely over TLS to a configured
upstream DNS server (cloudflare by default).

## Getting Started

Repo contains the `Makefile` with the next commands:

```make
Usage:
  help             prints this help message
  fmt              runs gofmt on all source files
  lint             runs golangci-lint on all source files
  test             runs go test on all source files
  build            compiles the binary for linux without architecture
  build-docker     builds the docker image with $PROJECT_NAME $BUILD_VERSION
  docker-run-tcp   runs the docker image with tcp protocol
  docker-run-udp   runs the docker image with udp protocol
  docker-stop      stops the running docker image
  clean            try to remove binary
```

Binary contains the next arguments:

```bash
  -config string
     path to the cfg file (optional)
  -tcp
     tcp mode
  -udp
     udp mode
```

In order to build the container and start the service you can use the
following commands:

```bash

# build and run container with TCP listner
make build-docker && make docker-run-tcp

# build and run container with UDP listner
make build-docker && make docker-run-udp

# test cmd how it works for TCP
dig -p 53 +short +tcp  goole.com @localhost

# test cmd how it works for UDP
dig -p 53 +short  google.com @localhost
```

## Requirements

- **Docker**  is installed
