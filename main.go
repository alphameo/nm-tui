package main

import (
	run "github.com/alphameo/nm-tui/internal"
)

// Injects via `go build -ldflags "-X main.version=$(VERSION)"`.
var version = "dev"

func main() {
	run.Run(version)
}
