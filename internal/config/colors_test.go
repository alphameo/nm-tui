package config_test

import (
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

	tests := []struct {
		name      string
		src       *config.ColorConfig
		wantErr   int
		fragments []string
		check     func(t *testing.T, dst *config.ColorConfig)
	}{
		{name: "valid overrides applied", src: &config.ColorConfig{
			Text:  new("#111111"),
			Error: new(config.ColorBrightRed),
		}, check: func(t *testing.T, dst *config.ColorConfig) {
			if *dst.Text != "#111111" {
				t.Errorf("Text = %q, want #111111", *dst.Text)
			}
			if *dst.Error != config.ColorBrightRed {
				t.Errorf("Error = %q, want %q", *dst.Error, config.ColorBrightRed)
			}
			if *dst.Accent != config.ColorBlue {
				t.Errorf("Accent = %q, want unchanged %q", *dst.Accent, config.ColorBlue)
			}
		}},
		{name: "multiple invalid fields collected", src: &config.ColorConfig{
			Text:  new("bogus"),
			Error: new("also-bogus"),
			Muted: new(""),
		}, wantErr: 3, fragments: []string{"text color", "error color", "muted color"}},
		{name: "empty source produces no errors", src: &config.ColorConfig{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			dst := config.DefaultColorConfig()
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
