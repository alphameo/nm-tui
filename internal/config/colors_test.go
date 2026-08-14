package config

import (
	"strings"
	"testing"
)

func TestDefaultColorConfig(t *testing.T) {
	assertNoNilFields(t, DefaultColorConfig())
}

func TestValidHex(t *testing.T) {
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
			if got := ValidHex(tt.in); got != tt.want {
				t.Errorf("ValidHex(%q) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

func TestValidCfgColor(t *testing.T) {
	valid := []string{
		ColorBlack, ColorRed, ColorGreen, ColorYellow,
		ColorBlue, ColorMagenta, ColorCyan, ColorWhite,
		ColorBrightBlack, ColorBrightRed, ColorBrightGreen, ColorBrightYellow,
		ColorBrightBlue, ColorBrightMagenta, ColorBrightCyan, ColorBrightWhite,
		ColorNone,
	}
	for _, c := range valid {
		if !ValidCfgColor(c) {
			t.Errorf("ValidCfgColor(%q) = false, want true", c)
		}
	}

	invalid := []string{
		"", "grey", "gray", "orange", "lightblue", "red2", "Red", "BLUE",
		"#ffffff", "bright", "nonexistent",
	}
	for _, c := range invalid {
		if ValidCfgColor(c) {
			t.Errorf("ValidCfgColor(%q) = true, want false", c)
		}
	}
}

func TestValidateColor(t *testing.T) {
	tests := []struct {
		in      string
		wantErr bool
	}{
		{ColorRed, false},
		{ColorBrightBlack, false},
		{ColorNone, false},
		{"#865fff", false},
		{"#ABCdef", false},
		{"#ffffff", false},
		{defaultKeyword, false},
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
			err := validateColor(tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateColor(%q) error = %v, wantErr %v", tt.in, err, tt.wantErr)
			}
		})
	}
}

func TestMergeColor(t *testing.T) {
	t.Run("nil source is a no-op", func(t *testing.T) {
		dst := new(ColorBlue)
		if err := mergeColor(dst, nil, "accent"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if *dst != ColorBlue {
			t.Errorf("dst = %q, want unchanged %q", *dst, ColorBlue)
		}
	})

	t.Run("valid color overrides", func(t *testing.T) {
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
		dst := new(ColorBlue)
		src := new(defaultKeyword)
		if err := mergeColor(dst, src, "accent"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if *dst != ColorBlue {
			t.Errorf("dst = %q, want unchanged %q", *dst, ColorBlue)
		}
	})

	t.Run("invalid color errors and does not override", func(t *testing.T) {
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

func TestColorConfigMerge(t *testing.T) {
	t.Run("valid overrides applied", func(t *testing.T) {
		dst := DefaultColorConfig()
		src := &ColorConfig{
			Text:  new("#111111"),
			Error: new(ColorBrightRed),
		}
		errs := dst.merge(src)
		if len(errs) != 0 {
			t.Fatalf("unexpected errors: %v", errs)
		}
		if *dst.Text != "#111111" {
			t.Errorf("Text = %q, want #111111", *dst.Text)
		}
		if *dst.Error != ColorBrightRed {
			t.Errorf("Error = %q, want %q", *dst.Error, ColorBrightRed)
		}
		if *dst.Accent != ColorBlue {
			t.Errorf("Accent = %q, want unchanged %q", *dst.Accent, ColorBlue)
		}
	})

	t.Run("multiple invalid fields collected", func(t *testing.T) {
		dst := DefaultColorConfig()
		src := &ColorConfig{
			Text:  new("bogus"),
			Error: new("also-bogus"),
			Muted: new(""),
		}
		errs := dst.merge(src)
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
		dst := DefaultColorConfig()
		if errs := dst.merge(&ColorConfig{}); len(errs) != 0 {
			t.Errorf("unexpected errors: %v", errs)
		}
	})
}
