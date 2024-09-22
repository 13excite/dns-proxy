.DEFAULT_GOAL := help

SHELL := /bin/bash

# constant variables
PROJECT_NAME	= dns-proxy
BINARY_NAME	= dns-proxy
GIT_COMMIT	= $(shell git rev-parse HEAD)
BINARY_TAR_DIR	= $(BINARY_NAME)-$(GIT_COMMIT)
BINARY_TAR_FILE	= $(BINARY_TAR_DIR).tar.gz
BUILD_VERSION	= $(shell cat VERSION.txt)
BUILD_DATE	= $(shell date -u '+%Y-%m-%d_%H:%M:%S')

# Terminal colors config
NO_COLOR=\033[0m
OK_COLOR=\033[32;01m

# golangci-lint config
golangci_lint_version=v1.60.3
vols=-v `pwd`:/app -w /app
run_lint=docker run --rm $(vols) golangci/golangci-lint:$(golangci_lint_version)

# LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

.PHONY: lint fmt help test build build-docker docker-run-tcp docker-run-udp docker-stop clean

## help: prints this help message
help:
	@echo "Usage:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## fmt: runs gofmt on all source files
fmt:
	@gofmt -l -w $(SRC)

## lint: runs golangci-lint on all source files
lint:
	@printf "$(OK_COLOR)==> Running golang-ci-linter via Docker$(NO_COLOR)\n"
	@$(run_lint) golangci-lint run --timeout=5m --verbose

## test: runs go test on all source files
test:
	@printf "$(OK_COLOR)==> Running tests$(NO_COLOR)\n"
	@go test -v -count=1 -covermode=atomic -coverpkg=./... -coverprofile=coverage.txt ./...
	@go tool cover -func coverage.txt


## build: compiles the binary for linux without architecture
build:
	@echo 'compiling binary...'
	@cd cmd/${PROJECT_NAME} && GOOS=linux go build -o ../../$(BINARY_NAME)

## build-docker: builds the docker image with $PROJECT_NAME $BUILD_VERSION
build-docker:
	@echo 'building docker image...'
	@docker build -t $(PROJECT_NAME):$(BUILD_VERSION) .

## docker-run-tcp: runs the docker image with tcp protocol
docker-run-tcp:
	@echo 'running docker image...'
	@docker run -d -p 53:1153 $(PROJECT_NAME):$(BUILD_VERSION)

## docker-run-udp: runs the docker image with udp protocol
docker-run-udp:
	@echo 'running docker image...'
	@docker run -d -p 53:1153/udp $(PROJECT_NAME):$(BUILD_VERSION) /usr/local/bin/dns-proxy -udp

## docker-stop: stops the running docker image
docker-stop:
	@echo 'stopping docker image...'
	@docker stop $(shell docker ps -q --filter ancestor=$(PROJECT_NAME):$(BUILD_VERSION))

## clean: try to remove binary
clean:
	@rm ./dns-proxy

