package infra

import (
	"os/exec"
)

func ExtractStderr(err error) string {
	var stderr string
	if exitErr, ok := err.(*exec.ExitError); ok {
		stderr = string(exitErr.Stderr)
	}
	return stderr
}
