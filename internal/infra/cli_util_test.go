package infra

import (
	"errors"
	"os/exec"
	"testing"
)

func TestExtractStderr(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{"nil-error", nil, ""},
		{"plain-error", errors.New("boom"), ""},
		{"exit-error-with-stderr", &exec.ExitError{Stderr: []byte("some stderr")}, "some stderr"},
		{"exit-error-empty-stderr", &exec.ExitError{}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractStderr(tt.err); got != tt.want {
				t.Errorf("ExtractStderr(%v) = %q, want %q", tt.err, got, tt.want)
			}
		})
	}
}
