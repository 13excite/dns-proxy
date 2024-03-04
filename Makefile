.DEFAULT_GOAL := help

SHELL := /bin/bash

# constant variables
PROJECT_NAME 	= dns-proxy
BINARY_NAME 	= dns-proxy
GIT_COMMIT 		= $(shell git rev-parse HEAD)
BINARY_TAR_DIR 	= $(BINARY_NAME)-$(GIT_COMMIT)
BINARY_TAR_FILE	= $(BINARY_TAR_DIR).tar.gz
BUILD_VERSION 	= $(shell cat VERSION.txt)
BUILD_DATE 		= $(shell date -u '+%Y-%m-%d_%H:%M:%S')

# Terminal colors config
NO_COLOR=\033[0m
OK_COLOR=\033[32;01m

# golangci-lint config
golangci_lint_version=latest
vols=-v `pwd`:/app -w /app
run_lint=docker run --rm $(vols) golangci/golangci-lint:$(golangci_lint_version)

# LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD)"

SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

.PHONY: lint fmt help test build

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

## build: compiles the binary
build:
	@echo 'compiling binary...'
	# @cd cmd/ && GOARCH=amd64 GOOS=linux go build -ldflags "-X main.buildTimestamp=$(BUILD_DATE) -X main.gitHash=$(GIT_COMMIT) -X main.buildVersion=$(BUILD_VERSION)" -o ../$(BINARY_NAME)
	@cd cmd/${PROJECT_NAME} && GOARCH=arm64 GOOS=linux go build -ldflags "-X main.buildTimestamp=$(BUILD_DATE) -X main.gitHash=$(GIT_COMMIT) -X main.buildVersion=$(BUILD_VERSION)" -o ../../$(BINARY_NAME)

