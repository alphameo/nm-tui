package config

import (
	"fmt"
	"testing"
)

func TestValidKeyName(t *testing.T) {
	t.Parallel()

	for name := range validKey {
		t.Run("named_"+name, func(t *testing.T) {
			t.Parallel()

			if !validKeyName(name) {
				t.Errorf("validKeyName(%q) = false, want true", name)
			}
		})
	}

	for name := range validModifier {
		t.Run("modifier_"+name, func(t *testing.T) {
			t.Parallel()

			combo := name + "+enter"
			if !validKeyName(combo) {
				t.Errorf("validKeyName(%q) = false, want true", combo)
			}
		})
	}

	for _, f := range []int{1, 12, 24, 63} {
		name := fmt.Sprintf("f%d", f)
		if !validKeyName(name) {
			t.Errorf("validKeyName(%q) = false, want true", name)
		}
	}

	singleChars := []string{"a", "Z", "0", "9", "!", " ", "-", "é", "中"}
	for _, c := range singleChars {
		if !validKeyName(c) {
			t.Errorf("validKeyName(%q) = false, want true", c)
		}
	}

	nonPrintableSingleRunes := []string{"\x00", "\x01", "\x7f"}
	for _, c := range nonPrintableSingleRunes {
		if validKeyName(c) {
			t.Errorf("validKeyName(%q) = true, want false", c)
		}
	}

	caseVariants := []string{"CTRL+R", "Enter", "SHIFT+TAB", "F5", "Ctrl+Shift+X"}
	for _, c := range caseVariants {
		if !validKeyName(c) {
			t.Errorf("validKeyName(%q) = false, want true", c)
		}
	}
}

func TestValidKeyNameInvalid(t *testing.T) {
	t.Parallel()

	tests := []string{
		"",
		"notakey",
		"f64",
		"f0",
		"ctrl+notakey",
		"badmod+enter",
		"ctrl+shift",
		"ctrl++",
		"super+duper+enter",
		"enter+",
	}
	for _, in := range tests {
		if validKeyName(in) {
			t.Errorf("validKeyName(%q) = true, want false", in)
		}
	}
}
