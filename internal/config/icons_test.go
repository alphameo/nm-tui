package config

import (
	"strings"
	"testing"
)

func TestDefaultIconConfig(t *testing.T) {
	assertNoNilFields(t, DefaultIconConfig())
}

func TestDefaultNerdIconConfig(t *testing.T) {
	assertNoNilFields(t, DefaultNerdIconConfig())
}

func TestDefaultNonNerdIconConfig(t *testing.T) {
	assertNoNilFields(t, DefaultNonNerdIconConfig())
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
		if !strings.Contains(err.Error(), "border_style") {
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
		if !strings.Contains(err.Error(), "spinner_style") {
			t.Errorf("unexpected error: %v", err)
		}
		if *dst != SpinnerLine {
			t.Errorf("dst = %q, want unchanged", *dst)
		}
	})
}

func TestValidateCursorShape(t *testing.T) {
	for _, c := range []string{CursorBar, CursorUnderline, CursorBlock} {
		if err := validateCursorShape(c); err != nil {
			t.Errorf("validateCursorShape(%q) error: %v", c, err)
		}
	}
	for _, c := range []string{"", "beam", "block_bar", "BAR"} {
		if err := validateCursorShape(c); err == nil {
			t.Errorf("validateCursorShape(%q) = nil, want error", c)
		}
	}
}

func TestMergeCursorShape(t *testing.T) {
	t.Run("valid override", func(t *testing.T) {
		dst := new(CursorBar)
		if err := mergeCursorShape(dst, new(CursorUnderline)); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if *dst != CursorUnderline {
			t.Errorf("dst = %q, want %q", *dst, CursorUnderline)
		}
	})

	t.Run("nil source no-op", func(t *testing.T) {
		dst := new(CursorBar)
		if err := mergeCursorShape(dst, nil); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if *dst != CursorBar {
			t.Errorf("dst = %q, want unchanged", *dst)
		}
	})

	t.Run("default keyword no-op", func(t *testing.T) {
		dst := new(CursorBar)
		if err := mergeCursorShape(dst, new(defaultKeyword)); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if *dst != CursorBar {
			t.Errorf("dst = %q, want unchanged", *dst)
		}
	})

	t.Run("invalid errors and does not override", func(t *testing.T) {
		dst := new(CursorBar)
		err := mergeCursorShape(dst, new("beam"))
		if err == nil {
			t.Fatal("expected error")
		}
		if !strings.Contains(err.Error(), "cursor_shape") {
			t.Errorf("unexpected error: %v", err)
		}
		if *dst != CursorBar {
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
			BorderStyle:      new(BorderRounded),
			SpinnerStyle:     new(SpinnerHamburger),
			InputCursorShape: new(CursorBlock),
			Check:            new("✓"),
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
		if *dst.InputCursorShape != CursorBlock {
			t.Errorf("InputCursorShape = %q, want %q", *dst.InputCursorShape, CursorBlock)
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
			BorderStyle:      new("dashed"),
			InputCursorShape: new("beam"),
			ToggleOff:        new(""),
			PwHiddenChar:     new("ab"),
		}
		errs := dst.merge(src)
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
		dst := DefaultIconConfig()
		src := &IconConfig{InputCursorShape: new("beam")}
		errs := dst.merge(src)
		if len(errs) != 1 {
			t.Fatalf("want 1 error, got %v", errs)
		}
		if !strings.Contains(errs[0].Error(), "cursor_shape") {
			t.Errorf("unexpected error: %v", errs[0])
		}
		if *dst.InputCursorShape != CursorBar {
			t.Errorf("InputCursorShape = %q, want unchanged %q", *dst.InputCursorShape, CursorBar)
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
