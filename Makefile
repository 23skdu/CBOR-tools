GO ?= go
BINARY ?= cbor-tool
PKG ?= ./...
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS ?= -s -w -X main.version=$(VERSION)

.PHONY: all build test fuzz cover clean lint security gosec govulncheck scan docker tools fmt fmt-check vet tidy check

all: build test

build:
	$(GO) build -trimpath -ldflags "$(LDFLAGS)" -o bin/$(BINARY) ./cmd/cbor-tool

test:
	$(GO) test -race -shuffle=on $(PKG)

fuzz:
	@set -e; for t in $$($(GO) test -list '^Fuzz' ./internal/cborx/ | grep '^Fuzz'); do \
		echo "Fuzzing $$t"; \
		$(GO) test -run=$$t -fuzz=$$t -fuzztime=30s ./internal/cborx/; \
	done

fuzz-ci:
	@set -e; for t in $$( $(GO) test -list '^Fuzz' ./internal/cborx/ | grep '^Fuzz'); do \
		echo "Fuzzing $$t"; \
		$(GO) test -run=$$t -fuzz=$$t -fuzztime=10s ./internal/cborx/; \
	done

cover:
	$(GO) test -coverprofile=coverage.out $(PKG)
	$(GO) tool cover -func=coverage.out

fmt:
	gofmt -w .

fmt-check:
	@test -z "$$(gofmt -l .)" || (echo "gofmt needed:"; gofmt -l .; exit 1)

vet:
	$(GO) vet ./...

tidy:
	$(GO) mod tidy
	@git diff --exit-code go.mod go.sum

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

check: fmt-check vet lint test security

tools:
	$(GO) install github.com/securego/gosec/v2/cmd/gosec@latest
	$(GO) install golang.org/x/vuln/cmd/govulncheck@latest
	$(GO) install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest

docker:
	docker build -t $(BINARY):latest .
