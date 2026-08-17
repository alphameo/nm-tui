.PHONY: build build-dev run deps clean-build logs lint lint-fix format test

VERSION ?= $(shell git describe --tags --always 2>/dev/null || echo dev)

build:
	CGO_ENABLED=0 go build -ldflags "-X main.version=$(VERSION)" -o bin/nm-tui ./cmd/nm-tui/main.go

build-dev:
	CGO_ENABLED=0 go build -o bin/nm-tui ./cmd/nm-tui/main.go

run:
	go run ./cmd/nm-tui/main.go

deps:
	go mod tidy

clean-build:
	make deps
	make build

logs:
	cat ~/.local/state/nm-tui/nm-tui.log | tail -n 50

lint:
	golangci-lint run ./...

lint-fix:
	golangci-lint --fix run ./...

format:
	go fmt ./...

test:
	go test -v ./...
	go test -cover ./...
