VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -s -w -X github.com/Sigeth/opslab-cli/internal/cli.version=$(VERSION)

.PHONY: all build install test lint

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
