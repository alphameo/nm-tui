package config

import (
	"strings"
	"testing"
)

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
		if err := mergeBorderStyle(dst, new(DefaultKeyword)); err != nil {
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
		if err := mergeSpinnerStyle(dst, new(DefaultKeyword)); err != nil {
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
		if err := mergeCursorShape(dst, new(DefaultKeyword)); err != nil {
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
		if err := mergeIcon(dst, new(DefaultKeyword), "check"); err != nil {
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
		if err := mergeSymbol(dst, new(DefaultKeyword), "password_symbol"); err != nil {
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
