package config_test

import (
	"testing"

	"github.com/alphameo/nm-tui/internal/config"
)

func TestDefaultNerdIconConfig(t *testing.T) {
	t.Parallel()

	assertNoNilFields(t, config.DefaultNerdIconConfig())
}

func TestDefaultNonNerdIconConfig(t *testing.T) {
	t.Parallel()

	assertNoNilFields(t, config.DefaultNonNerdIconConfig())
}

func TestIconConfigMerge(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		src       *config.IconConfig
		wantErr   int
		fragments []string
		check     func(t *testing.T, dst *config.IconConfig)
	}{
		{name: "valid overrides applied", src: &config.IconConfig{
			BorderStyle:      new(config.BorderRounded),
			SpinnerStyle:     new(config.SpinnerHamburger),
			InputCursorShape: new(config.CursorBlock),
			Check:            new("✓"),
		}, check: assertIconOverridesApplied},
		{
			name: "invalid fields collected",
			src: &config.IconConfig{
				BorderStyle:      new("dashed"),
				InputCursorShape: new("beam"),
				ToggleOff:        new(""),
				PwHiddenChar:     new("ab"),
			},
			wantErr: 4,
			fragments: []string{
				"border_style",
				"cursor_shape",
				"empty toggle_off icon",
				"length of symbol password_symbol",
			},
		},
		{
			name:      "invalid cursor shape collected",
			src:       &config.IconConfig{InputCursorShape: new("beam")},
			wantErr:   1,
			fragments: []string{"cursor_shape"},
			check: func(t *testing.T, dst *config.IconConfig) {
				if *dst.InputCursorShape != config.CursorBar {
					t.Errorf("InputCursorShape = %q, want unchanged %q", *dst.InputCursorShape, config.CursorBar)
				}
			},
		},
		{name: "default keyword skips all fields", src: &config.IconConfig{
			BorderStyle:  new(config.DefaultKeyword),
			ToggleOff:    new(config.DefaultKeyword),
			PwHiddenChar: new(config.DefaultKeyword),
		}, check: func(t *testing.T, dst *config.IconConfig) {
			if *dst.BorderStyle != config.BorderASCII || *dst.ToggleOff != "[ ]" || *dst.PwHiddenChar != "*" {
				t.Errorf(
					"defaults changed: border=%q off=%q pw=%q",
					*dst.BorderStyle,
					*dst.ToggleOff,
					*dst.PwHiddenChar,
				)
			}
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			dst := config.DefaultIconConfig()
			errs := dst.Merge(tt.src)
			if len(errs) != tt.wantErr {
				t.Fatalf("want %d errors, got %v", tt.wantErr, errs)
			}
			assertErrsContain(t, errs, tt.fragments...)
			if tt.check != nil {
				tt.check(t, dst)
			}
		})
	}
}

func assertIconOverridesApplied(t *testing.T, dst *config.IconConfig) {
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
}
