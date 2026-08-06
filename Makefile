.PHONY: build test clean all logs

build:
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

test:
	make test-config

test-config:
	go test ./internal/config/ -v
	go test ./internal/config/ -cover
