package infra_test

import (
	"errors"
	"fmt"
	"os/exec"
	"testing"

	"github.com/alphameo/nm-tui/internal/infra"
)

func TestExtractStderr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  error
		want string
	}{
		{"nil-error", nil, ""},
		{"plain-error", errors.New("boom"), ""},
		{"exit-error-with-stderr", &exec.ExitError{Stderr: []byte("some stderr")}, "some stderr"},
		{"exit-error-empty-stderr", &exec.ExitError{}, ""},
		{
			"wrapped-exit-error",
			fmt.Errorf("wrapped: %w", &exec.ExitError{Stderr: []byte("wrapped stderr")}),
			"wrapped stderr",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := infra.ExtractStderr(tt.err); got != tt.want {
				t.Errorf("ExtractStderr(%v) = %q, want %q", tt.err, got, tt.want)
			}
		})
	}
}
