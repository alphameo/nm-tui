.PHONY: build build-dev run deps clean-build logs lint lint-fix test test-config test-infra test-compositor test-tabview

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
	cat ~/.cache/nm-tui/log | tail -n 50

lint:
	golangci-lint run ./...

lint-fix:
	golangci-lint --fix run ./...

test:
	go test -v ./...
	go test -cover ./...

test-config:
	go test ./internal/config/ -v
	go test ./internal/config/ -cover

test-infra:
	go test ./internal/infra/ -v
	go test ./internal/infra/ -cover

test-compositor:
	go test ./internal/ui/tools/compositor/ -v
	go test ./internal/ui/tools/compositor/ -cover

test-tabview:
	go test ./internal/ui/models/tabview/ -v
	go test ./internal/ui/models/tabview/ -cover
