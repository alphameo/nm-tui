// Package logging wraps the infra managers and emits a structured log
// record for every call.
package logging

import (
	"errors"
	"log/slog"
	"os/exec"
	"time"
)

// middleware carries the logger and the operation prefix shared by all
// middleware implementations.
type middleware struct {
	logger *slog.Logger
	prefix string
}

func (mw middleware) call(name string, fn func() error) error {
	return logCall(mw.logger, mw.prefix+"."+name, fn)
}

func callResult[T any](mw middleware, name string, fn func() (T, error)) (T, error) {
	return logCallResult(mw.logger, mw.prefix+"."+name, fn)
}

func logCall(logger *slog.Logger, operation string, fn func() error) error {
	start := time.Now()
	err := fn()
	log(logger, operation, start, err)
	return err
}

func logCallResult[T any](logger *slog.Logger, operation string, fn func() (T, error)) (T, error) {
	start := time.Now()
	res, err := fn()
	log(logger, operation, start, err)
	return res, err
}

func log(logger *slog.Logger, operation string, start time.Time, err error) {
	if err == nil {
		logger.Debug("manager call",
			"operation", operation,
			"duration", time.Since(start),
		)
		return
	}
	logger.Error("manager call",
		"operation", operation,
		"duration", time.Since(start),
		"error", err,
		"exit_code", exitCode(err),
	)
}

func exitCode(err error) int {
	if exitErr, ok := errors.AsType[*exec.ExitError](err); ok {
		return exitErr.ExitCode()
	}
	return -1
}
