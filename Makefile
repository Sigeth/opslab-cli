# Contract library, tested at its newest tag by contract-latest.
CORE_GO := github.com/Sigeth/opslab-core-go
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -s -w -X github.com/Sigeth/opslab-cli/internal/cli.version=$(VERSION)

.PHONY: all build install test lint contract-latest

all: lint test build

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags '$(LDFLAGS)' -o bin/opslab ./cmd/opslab

install:
	go install -ldflags '$(LDFLAGS)' ./cmd/opslab

test:
	go test -race ./...

lint:
	go vet ./...
	@unformatted=$$(gofmt -l cmd internal); if [ -n "$$unformatted" ]; then echo "$$unformatted"; exit 1; fi

## contract-latest: build and test against the newest opslab-core-go tag (CI),
## on a temporary copy of go.mod so the pinned version is not changed
contract-latest:
	@set -e; tmp=$$(mktemp -d); trap 'rm -rf "$$tmp"' EXIT; \
	cp go.mod go.sum "$$tmp/"; \
	export GOWORK=off GOFLAGS="-modfile=$$tmp/go.mod"; \
	go get $(CORE_GO)@latest; \
	echo "== $$(go list -m $(CORE_GO))"; \
	go build ./...; \
	go test ./...
