# Simple DNS Proxy

Project contains a simple DNS proxy that listens to conventional DNS
and sends it over TLS.

## Table of Contents

- [Implementation Details](#implementation-details)
- [Getting Started](#getting-started)
- [Security Concerns](#security-concerns)
- [Integration in a Distributed, Microservices-Oriented, Containerized Architecture](#integration-in-a-distributed-microservices-oriented-containerized-architecture)
- [Future Improvements](#future-improvements)
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

# test cmd how it works for TCP
dig -p 53 +short  google.com @localhost
```

## Security Concerns

Deploying this proxy in a prod leds several security concerns:

- **Custom service**: It's recommended to take one of the products that provides this
functionality, or use one of the already written libraries so as not to write
everything from scratch, for example `github.com/miekg/dns`. This will probably reduce
the number of critical problems.

- **TLS Configuration**: Needs to set-up correct TLS configuration is crucial to prevent
MitM attacks and unauthorized access to DNS queries.

- **Access Control**: Needs to resctrict access to the proxy to authorized clients or
networks.It can help prevent abuse. It's also necessary to limit the ability to resolve
by domain and request types.

- **Monitoring:**: Need to implement logging, tracing and monitoring mechanisms
for detecting and responding to suspicious activities.

- **CI checks**: It's necessary to implement `CI` with checking the code, used libraries and docker images for vulnerabilities

## Integration in a Distributed, Microservices-Oriented, Containerized Architecture

To integrate this solution in such an architecture:

- Need to deploy the DNS proxy as a microservice within containers.
- Recommend to use orchestration tools like Kubernetes for managing scalability and availability,
not self-maded bash wrappers for managing docker.
- Need to implement logging, tracing and monitoring and also health checks.
- Need to implement service discovery mechanisms.
- It's necessary to configure network policies to restrict `IN/OUT` traffic
and enforce security between services.

## Future Improvements

Some potential improvements to the project include:

- Needs to implement support of metrics and tracing.
- Needs to implement caching mechanisms to improve performance.
- Try to add support for IPv6.
- It's require to add health checks and automatic failover mechanisms for high availability.
- Enhancing logging capabilities.

## Requirements

- **Docker**  is installed
