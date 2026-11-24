GO ?= go
BINARY ?= cbor-tool
PKG ?= ./...
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS ?= -s -w -X main.version=$(VERSION)

.PHONY: all build test bench fuzz fuzz-ci cover cover-check clean lint security gosec govulncheck scan docker docker-test tools fmt fmt-check vet tidy check man completions

all: build test

build:
	$(GO) build -trimpath -ldflags "$(LDFLAGS)" -o bin/$(BINARY) ./cmd/cbor-tool

test:
	$(GO) test -race -shuffle=on $(PKG)

bench:
	$(GO) test -run '^$$' -bench . -benchmem $(PKG)

fuzz:
	@set -e; for pkg in $$($(GO) list ./...); do \
		for t in $$($(GO) test -list '^Fuzz' $$pkg | grep '^Fuzz'); do \
			echo "Fuzzing $$pkg/$$t"; \
			$(GO) test -run=$$t -fuzz=$$t -fuzztime=30s $$pkg; \
		done; \
	done

fuzz-ci:
	@set -e; for pkg in $$($(GO) list ./...); do \
		for t in $$($(GO) test -list '^Fuzz' $$pkg | grep '^Fuzz'); do \
			echo "Fuzzing $$pkg/$$t"; \
			$(GO) test -run=$$t -fuzz=$$t -fuzztime=10s $$pkg; \
		done; \
	done

cover:
	$(GO) test -coverprofile=coverage.out $(PKG)
	$(GO) tool cover -func=coverage.out

# Fail unless every statement in the repository is covered by a test.
cover-check: cover
	@total=$$($(GO) tool cover -func=coverage.out | tail -1 | awk '{print $$3}' | tr -d '%'); \
	echo "total coverage: $$total%"; \
	if [ "$$total" != "100.0" ] && [ "$$total" != "100" ]; then \
		echo "FAIL: coverage must be 100%"; exit 1; \
	fi

man:
	$(GO) run ./cmd/cbor-tool -man > docs/cbor-tool.1.tmp
	mv docs/cbor-tool.1.tmp docs/cbor-tool.1

completions:
	$(GO) run ./cmd/cbor-tool -completion bash > completions/cbor-tool.bash.tmp
	$(GO) run ./cmd/cbor-tool -completion zsh > completions/_cbor-tool.tmp
	$(GO) run ./cmd/cbor-tool -completion fish > completions/cbor-tool.fish.tmp
	mv completions/cbor-tool.bash.tmp completions/cbor-tool.bash
	mv completions/_cbor-tool.tmp completions/_cbor-tool
	mv completions/cbor-tool.fish.tmp completions/cbor-tool.fish

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
	rm -rf bin coverage.out docs/cbor-tool.1.tmp completions/*.tmp
	$(GO) clean -cache -testcache

lint:
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest run ./...

gosec:
	go run github.com/securego/gosec/v2/cmd/gosec@latest ./...

govulncheck:
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

security: gosec govulncheck

scan: security

check: fmt-check vet lint test cover-check security

tools:
	$(GO) install github.com/securego/gosec/v2/cmd/gosec@latest
	$(GO) install golang.org/x/vuln/cmd/govulncheck@latest
	$(GO) install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest

docker:
	docker build -t $(BINARY):latest .

docker-test:
	./scripts/docker-integration.sh
