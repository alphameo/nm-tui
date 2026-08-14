package config_test

import (
	"strings"
	"testing"

	"github.com/alphameo/nm-tui/internal/config"
)

func TestDefaultNerdIconConfig(t *testing.T) {
	assertNoNilFields(t, config.DefaultNerdIconConfig())
}

func TestDefaultNonNerdIconConfig(t *testing.T) {
	assertNoNilFields(t, config.DefaultNonNerdIconConfig())
}

func TestIconConfigMerge(t *testing.T) {
	t.Run("valid overrides applied", func(t *testing.T) {
		dst := config.DefaultIconConfig()
		src := &config.IconConfig{
			BorderStyle:      new(config.BorderRounded),
			SpinnerStyle:     new(config.SpinnerHamburger),
			InputCursorShape: new(config.CursorBlock),
			Check:            new("✓"),
		}
		errs := dst.Merge(src)
		if len(errs) != 0 {
			t.Fatalf("unexpected errors: %v", errs)
		}
		if *dst.BorderStyle != config.BorderRounded {
			t.Errorf("BorderStyle = %q, want %q", *dst.BorderStyle, config.BorderRounded)
		}
		if *dst.SpinnerStyle != config.SpinnerHamburger {
			t.Errorf("SpinnerStyle = %q, want %q", *dst.SpinnerStyle, config.SpinnerHamburger)
		}
		if *dst.InputCursorShape != config.CursorBlock {
			t.Errorf("InputCursorShape = %q, want %q", *dst.InputCursorShape, config.CursorBlock)
		}
		if *dst.Check != "✓" {
			t.Errorf("Check = %q, want ✓", *dst.Check)
		}
		if *dst.ToggleOff != "[ ]" {
			t.Errorf("ToggleOff = %q, want unchanged [ ]", *dst.ToggleOff)
		}
	})

	t.Run("invalid fields collected", func(t *testing.T) {
		dst := config.DefaultIconConfig()
		src := &config.IconConfig{
			BorderStyle:      new("dashed"),
			InputCursorShape: new("beam"),
			ToggleOff:        new(""),
			PwHiddenChar:     new("ab"),
		}
		errs := dst.Merge(src)
		if len(errs) != 4 {
			t.Fatalf("want 4 errors, got %v", errs)
		}
		for _, frag := range []string{"border_style", "cursor_shape", "empty toggle_off icon", "length of symbol password_symbol"} {
			found := false
			for _, err := range errs {
				if strings.Contains(err.Error(), frag) {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("no error contains %q in %v", frag, errs)
			}
		}
	})

	t.Run("invalid cursor shape collected", func(t *testing.T) {
		dst := config.DefaultIconConfig()
		src := &config.IconConfig{InputCursorShape: new("beam")}
		errs := dst.Merge(src)
		if len(errs) != 1 {
			t.Fatalf("want 1 error, got %v", errs)
		}
		if !strings.Contains(errs[0].Error(), "cursor_shape") {
			t.Errorf("unexpected error: %v", errs[0])
		}
		if *dst.InputCursorShape != config.CursorBar {
			t.Errorf("InputCursorShape = %q, want unchanged %q", *dst.InputCursorShape, config.CursorBar)
		}
	})

	t.Run("default keyword skips all fields", func(t *testing.T) {
		dst := config.DefaultIconConfig()
		src := &config.IconConfig{
			BorderStyle:  new(config.DefaultKeyword),
			ToggleOff:    new(config.DefaultKeyword),
			PwHiddenChar: new(config.DefaultKeyword),
		}
		if errs := dst.Merge(src); len(errs) != 0 {
			t.Fatalf("unexpected errors: %v", errs)
		}
		if *dst.BorderStyle != config.BorderASCII || *dst.ToggleOff != "[ ]" || *dst.PwHiddenChar != "*" {
			t.Errorf("defaults changed: border=%q off=%q pw=%q", *dst.BorderStyle, *dst.ToggleOff, *dst.PwHiddenChar)
		}
	})
}
