package config

import (
	"strings"
	"testing"
)

func TestDefaultIconConfig(t *testing.T) {
	cfg := DefaultIconConfig()

	if cfg.NerdPreset == nil || *cfg.NerdPreset {
		t.Errorf("NerdPreset = %v, want false", cfg.NerdPreset)
	}
	if cfg.BorderStyle == nil || *cfg.BorderStyle != BorderASCII {
		t.Errorf("BorderStyle = %v, want %q", cfg.BorderStyle, BorderASCII)
	}
	if cfg.SpinnerStyle == nil || *cfg.SpinnerStyle != SpinnerLine {
		t.Errorf("SpinnerStyle = %v, want %q", cfg.SpinnerStyle, SpinnerLine)
	}
	if cfg.ToggleOff == nil || *cfg.ToggleOff != "[ ]" {
		t.Errorf("ToggleOff = %v, want [ ]", cfg.ToggleOff)
	}
	if cfg.ToggleOn == nil || *cfg.ToggleOn != "[x]" {
		t.Errorf("ToggleOn = %v, want [x]", cfg.ToggleOn)
	}
	if cfg.PwHiddenChar == nil || *cfg.PwHiddenChar != "*" {
		t.Errorf("PwHiddenChar = %v, want *", cfg.PwHiddenChar)
	}
}

func TestDefaultNerdIconConfig(t *testing.T) {
	cfg := DefaultNerdIconConfig()

	want := map[string]string{
		"BorderStyle":  BorderRounded,
		"SpinnerStyle": SpinnerMeter,
		"ToggleOff":    "\uf204 ",
		"ToggleOn":     "\uf205 ",
		"PwHiddenChar": "•",
		"Error":        "✗",
		"Check":        "\uf00c",
		"Signal":       "\uf012",
		"Connection":   "\U000f1616",
		"AccessPoint":  "\U000f0003",
		"Infra":        "🖳",
		"Mesh":         "\uf292",
		"AdHoc":        "\uf10b",
	}

	if cfg.NerdPreset == nil || !*cfg.NerdPreset {
		t.Errorf("NerdPreset = %v, want true", cfg.NerdPreset)
	}
	if got := *cfg.BorderStyle; got != want["BorderStyle"] {
		t.Errorf("BorderStyle = %q, want %q", got, want["BorderStyle"])
	}
	if got := *cfg.SpinnerStyle; got != want["SpinnerStyle"] {
		t.Errorf("SpinnerStyle = %q, want %q", got, want["SpinnerStyle"])
	}
	if got := *cfg.ToggleOff; got != want["ToggleOff"] {
		t.Errorf("ToggleOff = %q, want %q", got, want["ToggleOff"])
	}
	if got := *cfg.ToggleOn; got != want["ToggleOn"] {
		t.Errorf("ToggleOn = %q, want %q", got, want["ToggleOn"])
	}
	if got := *cfg.PwHiddenChar; got != want["PwHiddenChar"] {
		t.Errorf("PwHiddenChar = %q, want %q", got, want["PwHiddenChar"])
	}
	if got := *cfg.Error; got != want["Error"] {
		t.Errorf("Error = %q, want %q", got, want["Error"])
	}
	if got := *cfg.Check; got != want["Check"] {
		t.Errorf("Check = %q, want %q", got, want["Check"])
	}
	if got := *cfg.Signal; got != want["Signal"] {
		t.Errorf("Signal = %q, want %q", got, want["Signal"])
	}
	if got := *cfg.Connection; got != want["Connection"] {
		t.Errorf("Connection = %q, want %q", got, want["Connection"])
	}
	if got := *cfg.AccessPoint; got != want["AccessPoint"] {
		t.Errorf("AccessPoint = %q, want %q", got, want["AccessPoint"])
	}
	if got := *cfg.Infra; got != want["Infra"] {
		t.Errorf("Infra = %q, want %q", got, want["Infra"])
	}
	if got := *cfg.Mesh; got != want["Mesh"] {
		t.Errorf("Mesh = %q, want %q", got, want["Mesh"])
	}
	if got := *cfg.AdHoc; got != want["AdHoc"] {
		t.Errorf("AdHoc = %q, want %q", got, want["AdHoc"])
	}
}

func TestDefaultNonNerdIconConfig(t *testing.T) {
	cfg := DefaultNonNerdIconConfig()

	if cfg.NerdPreset == nil || *cfg.NerdPreset {
		t.Errorf("NerdPreset = %v, want false", cfg.NerdPreset)
	}
	want := map[string]string{
		"BorderStyle":  BorderASCII,
		"SpinnerStyle": SpinnerLine,
		"ToggleOff":    "[ ]",
		"ToggleOn":     "[x]",
		"PwHiddenChar": "*",
		"Error":        "!",
		"Check":        "v",
		"Signal":       "sig",
		"Connection":   "con",
		"AccessPoint":  "ap",
		"Infra":        "infr",
		"Mesh":         "#",
		"AdHoc":        "ah",
	}

	if got := *cfg.ToggleOff; got != want["ToggleOff"] {
		t.Errorf("ToggleOff = %q, want %q", got, want["ToggleOff"])
	}
	if got := *cfg.ToggleOn; got != want["ToggleOn"] {
		t.Errorf("ToggleOn = %q, want %q", got, want["ToggleOn"])
	}
	if got := *cfg.PwHiddenChar; got != want["PwHiddenChar"] {
		t.Errorf("PwHiddenChar = %q, want %q", got, want["PwHiddenChar"])
	}
	if got := *cfg.Error; got != want["Error"] {
		t.Errorf("Error = %q, want %q", got, want["Error"])
	}
	if got := *cfg.Check; got != want["Check"] {
		t.Errorf("Check = %q, want %q", got, want["Check"])
	}
	if got := *cfg.Signal; got != want["Signal"] {
		t.Errorf("Signal = %q, want %q", got, want["Signal"])
	}
	if got := *cfg.Connection; got != want["Connection"] {
		t.Errorf("Connection = %q, want %q", got, want["Connection"])
	}
	if got := *cfg.AccessPoint; got != want["AccessPoint"] {
		t.Errorf("AccessPoint = %q, want %q", got, want["AccessPoint"])
	}
	if got := *cfg.Infra; got != want["Infra"] {
		t.Errorf("Infra = %q, want %q", got, want["Infra"])
	}
	if got := *cfg.Mesh; got != want["Mesh"] {
		t.Errorf("Mesh = %q, want %q", got, want["Mesh"])
	}
	if got := *cfg.AdHoc; got != want["AdHoc"] {
		t.Errorf("AdHoc = %q, want %q", got, want["AdHoc"])
	}
}

func TestValidateBorderStyle(t *testing.T) {
	valid := []string{
		BorderASCII, BorderMarkdown,
		BorderRounded,
		BorderSquare, BorderThickSquare, BorderDoubleSquare,
		BorderBlock, BorderOuterHalfBlock, BorderInnerHalfBlock,
	}
	for _, b := range valid {
		if err := validateBorderStyle(b); err != nil {
			t.Errorf("validateBorderStyle(%q) error: %v", b, err)
		}
	}
	for _, b := range []string{"", "dashed", "double", "none", "ROUNDED"} {
		if err := validateBorderStyle(b); err == nil {
			t.Errorf("validateBorderStyle(%q) = nil, want error", b)
		}
	}
}

func TestValidateSpinnerStyle(t *testing.T) {
	valid := []string{
		SpinnerLine, SpinnerEllipsis,
		SpinnerDot, SpinnerMiniDot,
		SpinnerJump, SpinnerPulse, SpinnerPoints, SpinnerMeter, SpinnerHamburger,
	}
	for _, s := range valid {
		if err := validateSpinnerStyle(s); err != nil {
			t.Errorf("validateSpinnerStyle(%q) error: %v", s, err)
		}
	}
	for _, s := range []string{"", "dots", "bounce", "LINE"} {
		if err := validateSpinnerStyle(s); err == nil {
			t.Errorf("validateSpinnerStyle(%q) = nil, want error", s)
		}
	}
}

func TestMergeBorderStyle(t *testing.T) {
	t.Run("valid override", func(t *testing.T) {
		dst := new(BorderASCII)
		if err := mergeBorderStyle(dst, new(BorderSquare)); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if *dst != BorderSquare {
			t.Errorf("dst = %q, want %q", *dst, BorderSquare)
		}
	})

	t.Run("nil source no-op", func(t *testing.T) {
		dst := new(BorderASCII)
		if err := mergeBorderStyle(dst, nil); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if *dst != BorderASCII {
			t.Errorf("dst = %q, want unchanged", *dst)
		}
	})

	t.Run("default keyword no-op", func(t *testing.T) {
		dst := new(BorderASCII)
		if err := mergeBorderStyle(dst, new(defaultKeyword)); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if *dst != BorderASCII {
			t.Errorf("dst = %q, want unchanged", *dst)
		}
	})

	t.Run("invalid errors and does not override", func(t *testing.T) {
		dst := new(BorderASCII)
		err := mergeBorderStyle(dst, new("dashed"))
		if err == nil {
			t.Fatal("expected error")
		}
		if !strings.Contains(err.Error(), "border style") {
			t.Errorf("unexpected error: %v", err)
		}
		if *dst != BorderASCII {
			t.Errorf("dst = %q, want unchanged", *dst)
		}
	})
}

func TestMergeSpinnerStyle(t *testing.T) {
	t.Run("valid override", func(t *testing.T) {
		dst := new(SpinnerLine)
		if err := mergeSpinnerStyle(dst, new(SpinnerDot)); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if *dst != SpinnerDot {
			t.Errorf("dst = %q, want %q", *dst, SpinnerDot)
		}
	})

	t.Run("nil source no-op", func(t *testing.T) {
		dst := new(SpinnerLine)
		if err := mergeSpinnerStyle(dst, nil); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("default keyword no-op", func(t *testing.T) {
		dst := new(SpinnerLine)
		if err := mergeSpinnerStyle(dst, new(defaultKeyword)); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if *dst != SpinnerLine {
			t.Errorf("dst = %q, want unchanged", *dst)
		}
	})

	t.Run("invalid errors and does not override", func(t *testing.T) {
		dst := new(SpinnerLine)
		err := mergeSpinnerStyle(dst, new("bounce"))
		if err == nil {
			t.Fatal("expected error")
		}
		if !strings.Contains(err.Error(), "spinner style") {
			t.Errorf("unexpected error: %v", err)
		}
		if *dst != SpinnerLine {
			t.Errorf("dst = %q, want unchanged", *dst)
		}
	})
}

func TestMergeIcon(t *testing.T) {
	t.Run("valid icon", func(t *testing.T) {
		dst := new("old")
		if err := mergeIcon(dst, new("new"), "check"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if *dst != "new" {
			t.Errorf("dst = %q, want new", *dst)
		}
	})

	t.Run("nil source no-op", func(t *testing.T) {
		dst := new("old")
		if err := mergeIcon(dst, nil, "check"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("default keyword no-op", func(t *testing.T) {
		dst := new("old")
		if err := mergeIcon(dst, new(defaultKeyword), "check"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if *dst != "old" {
			t.Errorf("dst = %q, want unchanged", *dst)
		}
	})

	t.Run("empty icon errors", func(t *testing.T) {
		err := mergeIcon(new("old"), new(""), "toggle_off")
		if err == nil {
			t.Fatal("expected error")
		}
		if !strings.Contains(err.Error(), "empty toggle_off icon") {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

func TestMergeSymbol(t *testing.T) {
	t.Run("single rune symbol", func(t *testing.T) {
		dst := new("x")
		if err := mergeSymbol(dst, new("•"), "password_symbol"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if *dst != "•" {
			t.Errorf("dst = %q, want •", *dst)
		}
	})

	t.Run("nil source no-op", func(t *testing.T) {
		dst := new("x")
		if err := mergeSymbol(dst, nil, "password_symbol"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("default keyword no-op", func(t *testing.T) {
		dst := new("x")
		if err := mergeSymbol(dst, new(defaultKeyword), "password_symbol"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if *dst != "x" {
			t.Errorf("dst = %q, want unchanged", *dst)
		}
	})

	t.Run("empty symbol errors", func(t *testing.T) {
		err := mergeSymbol(new("x"), new(""), "password_symbol")
		if err == nil {
			t.Fatal("expected error")
		}
		if !strings.Contains(err.Error(), "empty password_symbol symbol") {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("multi rune symbol errors", func(t *testing.T) {
		err := mergeSymbol(new("x"), new("ab"), "password_symbol")
		if err == nil {
			t.Fatal("expected error")
		}
		if !strings.Contains(err.Error(), "length of symbol password_symbol > 1") {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

func TestIconConfigMerge(t *testing.T) {
	t.Run("valid overrides applied", func(t *testing.T) {
		dst := DefaultIconConfig()
		src := &IconConfig{
			BorderStyle:  new(BorderRounded),
			SpinnerStyle: new(SpinnerHamburger),
			Check:        new("✓"),
		}
		errs := dst.merge(src)
		if len(errs) != 0 {
			t.Fatalf("unexpected errors: %v", errs)
		}
		if *dst.BorderStyle != BorderRounded {
			t.Errorf("BorderStyle = %q, want %q", *dst.BorderStyle, BorderRounded)
		}
		if *dst.SpinnerStyle != SpinnerHamburger {
			t.Errorf("SpinnerStyle = %q, want %q", *dst.SpinnerStyle, SpinnerHamburger)
		}
		if *dst.Check != "✓" {
			t.Errorf("Check = %q, want ✓", *dst.Check)
		}
		if *dst.ToggleOff != "[ ]" {
			t.Errorf("ToggleOff = %q, want unchanged [ ]", *dst.ToggleOff)
		}
	})

	t.Run("invalid fields collected", func(t *testing.T) {
		dst := DefaultIconConfig()
		src := &IconConfig{
			BorderStyle:  new("dashed"),
			ToggleOff:    new(""),
			PwHiddenChar: new("ab"),
		}
		errs := dst.merge(src)
		if len(errs) != 3 {
			t.Fatalf("want 3 errors, got %v", errs)
		}
		for _, frag := range []string{"border style", "empty toggle_off icon", "length of symbol password_symbol"} {
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

	t.Run("default keyword skips all fields", func(t *testing.T) {
		dst := DefaultIconConfig()
		src := &IconConfig{
			BorderStyle:  new(defaultKeyword),
			ToggleOff:    new(defaultKeyword),
			PwHiddenChar: new(defaultKeyword),
		}
		if errs := dst.merge(src); len(errs) != 0 {
			t.Fatalf("unexpected errors: %v", errs)
		}
		if *dst.BorderStyle != BorderASCII || *dst.ToggleOff != "[ ]" || *dst.PwHiddenChar != "*" {
			t.Errorf("defaults changed: border=%q off=%q pw=%q", *dst.BorderStyle, *dst.ToggleOff, *dst.PwHiddenChar)
		}
	})
}
