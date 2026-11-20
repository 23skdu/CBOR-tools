GO ?= go
BINARY ?= cbor-tool
PKG ?= ./...
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS ?= -s -w -X main.version=$(VERSION)

.PHONY: all build test fuzz cover clean lint security gosec govulncheck scan docker tools

all: build test

build:
	$(GO) build -trimpath -ldflags "$(LDFLAGS)" -o bin/$(BINARY) ./cmd/cbor-tool

test:
	$(GO) test -race -shuffle=on $(PKG)

fuzz:
	$(GO) test -fuzz=. -fuzztime=30s $(PKG)

cover:
	$(GO) test -coverprofile=coverage.out $(PKG)
	$(GO) tool cover -func=coverage.out

clean:
	rm -rf bin coverage.out
	$(GO) clean -cache -testcache

lint:
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest run ./...

gosec:
	go run github.com/securego/gosec/v2/cmd/gosec@latest ./...

govulncheck:
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

security: gosec govulncheck

scan: security

tools:
	$(GO) install github.com/securego/gosec/v2/cmd/gosec@latest
	$(GO) install golang.org/x/vuln/cmd/govulncheck@latest
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

docker:
	docker build -t $(BINARY):latest .