package main

import (
	run "github.com/alphameo/nm-tui/internal"
	ver "github.com/alphameo/nm-tui/internal/version"
)

// Injects via `go build -ldflags "-X main.version=$(VERSION)"`.
var version = "dev"

func main() {
	run.Run(ver.Resolve(version))
}
