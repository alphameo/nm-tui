package infra

import (
	"errors"
	"os/exec"
)

func ExtractStderr(err error) string {
	var stderr string
	if exitErr, ok := errors.AsType[*exec.ExitError](err); ok {
		stderr = string(exitErr.Stderr)
	}
	return stderr
}
