package config_test

import (
	"strings"
	"testing"

	"github.com/alphameo/nm-tui/internal/config"
)

func TestDefaultColorConfig(t *testing.T) {
	t.Parallel()

	assertNoNilFields(t, config.DefaultColorConfig())
}

func TestValidHex(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in   string
		want bool
	}{
		{"#000000", true},
		{"#ffffff", true},
		{"#FFFFFF", true},
		{"#865fff", true},
		{"#123456", true},
		{"#1a2B3c", true},
		{"#fffff", false},
		{"#ffffffff", false},
		{"865fff", false},
		{"#gggggg", false},
		{"#zzz123", false},
		{"", false},
		{"#", false},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			t.Parallel()

			if got := config.ValidHex(tt.in); got != tt.want {
				t.Errorf("ValidHex(%q) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

func TestValidCfgColor(t *testing.T) {
	t.Parallel()

	valid := []string{
		config.ColorBlack, config.ColorRed, config.ColorGreen, config.ColorYellow,
		config.ColorBlue, config.ColorMagenta, config.ColorCyan, config.ColorWhite,
		config.ColorBrightBlack, config.ColorBrightRed, config.ColorBrightGreen, config.ColorBrightYellow,
		config.ColorBrightBlue, config.ColorBrightMagenta, config.ColorBrightCyan, config.ColorBrightWhite,
		config.ColorNone,
	}
	for _, c := range valid {
		if !config.ValidCfgColor(c) {
			t.Errorf("ValidCfgColor(%q) = false, want true", c)
		}
	}

	invalid := []string{
		"", "grey", "gray", "orange", "lightblue", "red2", "Red", "BLUE",
		"#ffffff", "bright", "nonexistent",
	}
	for _, c := range invalid {
		if config.ValidCfgColor(c) {
			t.Errorf("ValidCfgColor(%q) = true, want false", c)
		}
	}
}

func TestValidateColor(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in      string
		wantErr bool
	}{
		{config.ColorRed, false},
		{config.ColorBrightBlack, false},
		{config.ColorNone, false},
		{"#865fff", false},
		{"#ABCdef", false},
		{"#ffffff", false},
		{config.DefaultKeyword, false},
		{"DEFAULT", false},
		{"Red", false},
		{"RED", false},
		{"#865ff", true},
		{"#865ffz", true},
		{"notacolor", true},
		{"", true},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			t.Parallel()

			err := config.ValidateColor(tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateColor(%q) error = %v, wantErr %v", tt.in, err, tt.wantErr)
			}
		})
	}
}

func TestColorConfigMerge(t *testing.T) {
	t.Parallel()

	t.Run("valid overrides applied", func(t *testing.T) {
		t.Parallel()

		dst := config.DefaultColorConfig()
		src := &config.ColorConfig{
			Text:  new("#111111"),
			Error: new(config.ColorBrightRed),
		}
		errs := dst.Merge(src)
		if len(errs) != 0 {
			t.Fatalf("unexpected errors: %v", errs)
		}
		if *dst.Text != "#111111" {
			t.Errorf("Text = %q, want #111111", *dst.Text)
		}
		if *dst.Error != config.ColorBrightRed {
			t.Errorf("Error = %q, want %q", *dst.Error, config.ColorBrightRed)
		}
		if *dst.Accent != config.ColorBlue {
			t.Errorf("Accent = %q, want unchanged %q", *dst.Accent, config.ColorBlue)
		}
	})

	t.Run("multiple invalid fields collected", func(t *testing.T) {
		t.Parallel()

		dst := config.DefaultColorConfig()
		src := &config.ColorConfig{
			Text:  new("bogus"),
			Error: new("also-bogus"),
			Muted: new(""),
		}
		errs := dst.Merge(src)
		if len(errs) != 3 {
			t.Fatalf("want 3 errors, got %v", errs)
		}
		for _, frag := range []string{"text color", "error color", "muted color"} {
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

	t.Run("empty source produces no errors", func(t *testing.T) {
		t.Parallel()

		dst := config.DefaultColorConfig()
		if errs := dst.Merge(&config.ColorConfig{}); len(errs) != 0 {
			t.Errorf("unexpected errors: %v", errs)
		}
	})
}
