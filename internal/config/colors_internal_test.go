package config

import (
	"strings"
	"testing"
)

func TestMergeColor(t *testing.T) {
	t.Parallel()

	t.Run("nil source is a no-op", func(t *testing.T) {
		t.Parallel()

		dst := new(ColorBlue)
		if err := mergeColor(dst, nil, "accent"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if *dst != ColorBlue {
			t.Errorf("dst = %q, want unchanged %q", *dst, ColorBlue)
		}
	})

	t.Run("valid color overrides", func(t *testing.T) {
		t.Parallel()

		dst := new(ColorBlue)
		src := new(ColorGreen)
		if err := mergeColor(dst, src, "accent"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if *dst != ColorGreen {
			t.Errorf("dst = %q, want %q", *dst, ColorGreen)
		}
	})

	t.Run("default keyword keeps destination", func(t *testing.T) {
		t.Parallel()

		dst := new(ColorBlue)
		src := new(DefaultKeyword)
		if err := mergeColor(dst, src, "accent"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if *dst != ColorBlue {
			t.Errorf("dst = %q, want unchanged %q", *dst, ColorBlue)
		}
	})

	t.Run("invalid color errors and does not override", func(t *testing.T) {
		t.Parallel()

		dst := new(ColorBlue)
		src := new("notacolor")
		err := mergeColor(dst, src, "accent")
		if err == nil {
			t.Fatal("expected error")
		}
		if !strings.Contains(err.Error(), "accent color") {
			t.Errorf("error does not mention accent color: %v", err)
		}
		if *dst != ColorBlue {
			t.Errorf("dst = %q, want unchanged %q", *dst, ColorBlue)
		}
	})
}
