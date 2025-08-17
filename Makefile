# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=go-zed-test-toggle
BINARY_UNIX=$(BINARY_NAME)_unix
BINARY_WINDOWS=$(BINARY_NAME).exe

# Build variables
VERSION?=dev
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-s -w -X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.BuildTime=$(BUILD_TIME)"

# Installation variables
PREFIX?=/usr/local
BINDIR=$(PREFIX)/bin

all: build

build:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) -v

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_UNIX) -v

build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BINARY_WINDOWS) -v

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	rm -f $(BINARY_WINDOWS)

run:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) -v
	./$(BINARY_NAME)

deps:
	$(GOMOD) download
	$(GOMOD) tidy

install: build
	@echo "Installing $(BINARY_NAME) to $(BINDIR)"
	@install -d $(BINDIR)
	@install -m 755 $(BINARY_NAME) $(BINDIR)

uninstall:
	@echo "Removing $(BINARY_NAME) from $(BINDIR)"
	@rm -f $(BINDIR)/$(BINARY_NAME)

# Cross compilation
build-all: build build-linux build-windows

.PHONY: all build build-linux build-windows test clean run deps install uninstall build-all
