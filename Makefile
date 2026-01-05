.PHONY: build test clean all

build:
	go build ./cmd/nm-tui/main.go

run:
	go run ./cmd/nm-tui/main.go
