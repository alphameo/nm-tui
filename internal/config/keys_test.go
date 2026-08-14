package config

import (
	"fmt"
	"strings"
	"testing"
)

func TestDefaultKeys(t *testing.T) {
	assertNoNilFields(t, DefaultKeys())
}

func TestValidKeyName(t *testing.T) {
	for name := range validKey {
		t.Run("named_"+name, func(t *testing.T) {
			if !validKeyName(name) {
				t.Errorf("validKeyName(%q) = false, want true", name)
			}
		})
	}

	for name := range validModifier {
		t.Run("modifier_"+name, func(t *testing.T) {
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

func TestMergeKeyList(t *testing.T) {
	t.Run("nil source is a no-op", func(t *testing.T) {
		dst := keyBinding("a")
		if errs := mergeKeyList(&dst, nil, "toggle"); len(errs) != 0 {
			t.Fatalf("unexpected errors: %v", errs)
		}
		assertKeyBinding(t, "dst", dst, "a")
	})

	t.Run("valid source overrides destination", func(t *testing.T) {
		dst := keyBinding("a")
		src := keyBinding("b", "c")
		if errs := mergeKeyList(&dst, src, "toggle"); len(errs) != 0 {
			t.Fatalf("unexpected errors: %v", errs)
		}
		assertKeyBinding(t, "dst", dst, "b", "c")
	})

	t.Run("invalid source errors and keeps destination", func(t *testing.T) {
		dst := keyBinding("a")
		src := keyBinding("b", "notakey")
		errs := mergeKeyList(&dst, src, "toggle")
		if len(errs) != 1 {
			t.Fatalf("want 1 error, got %v", errs)
		}
		if !strings.Contains(errs[0].Error(), `invalid key toggle: "notakey"`) {
			t.Errorf("unexpected error: %v", errs[0])
		}
		assertKeyBinding(t, "dst", dst, "a")
	})

	t.Run("multiple invalid keys produce multiple errors", func(t *testing.T) {
		dst := keyBinding("a")
		src := keyBinding("bad1", "bad2")
		errs := mergeKeyList(&dst, src, "main.quit")
		if len(errs) != 2 {
			t.Fatalf("want 2 errors, got %v", errs)
		}
	})

	t.Run("nil destination is populated", func(t *testing.T) {
		var dst *KeyBinding
		src := keyBinding("x")
		if errs := mergeKeyList(&dst, src, "toggle"); len(errs) != 0 {
			t.Fatalf("unexpected errors: %v", errs)
		}
		assertKeyBinding(t, "dst", dst, "x")
	})
}

func TestKeyConfigMerge(t *testing.T) {
	t.Run("nil source is a no-op", func(t *testing.T) {
		dst := DefaultKeys()
		if errs := dst.merge(nil); len(errs) != 0 {
			t.Fatalf("unexpected errors: %v", errs)
		}
	})

	t.Run("empty source is a no-op", func(t *testing.T) {
		dst := DefaultKeys()
		if errs := dst.merge(&KeyConfig{}); len(errs) != 0 {
			t.Fatalf("unexpected errors: %v", errs)
		}
		assertKeyBinding(t, "toggle", dst.Toggle, "space")
	})

	t.Run("partial override merges only set fields", func(t *testing.T) {
		dst := DefaultKeys()
		src := &KeyConfig{
			Toggle: keyBinding("t"),
			Main:   &MainKeys{Quit: keyBinding("q")},
		}
		if errs := dst.merge(src); len(errs) != 0 {
			t.Fatalf("unexpected errors: %v", errs)
		}
		assertKeyBinding(t, "toggle", dst.Toggle, "t")
		assertKeyBinding(t, "main.quit", dst.Main.Quit, "q")
		assertKeyBinding(t, "main.next_tab", dst.Main.TabNext, "]")
		assertKeyBinding(t, "rescan", dst.Rescan, "r")
	})

	t.Run("invalid key propagates error", func(t *testing.T) {
		dst := DefaultKeys()
		src := &KeyConfig{Toggle: keyBinding("notakey")}
		errs := dst.merge(src)
		if len(errs) != 1 {
			t.Fatalf("want 1 error, got %v", errs)
		}
		if !strings.Contains(errs[0].Error(), "invalid key toggle") {
			t.Errorf("unexpected error: %v", errs[0])
		}
		assertKeyBinding(t, "toggle", dst.Toggle, "space")
	})
}

func TestMainKeysMerge(t *testing.T) {
	dst := &MainKeys{TabNext: keyBinding("]")}
	if errs := dst.merge(nil); len(errs) != 0 {
		t.Fatalf("nil source: unexpected errors: %v", errs)
	}

	src := &MainKeys{TabNext: keyBinding("n"), Quit: keyBinding("q")}
	if errs := dst.merge(src); len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	assertKeyBinding(t, "next_tab", dst.TabNext, "n")
	assertKeyBinding(t, "quit", dst.Quit, "q")
	assertNilKeyBinding(t, "prev_tab", dst.TabPrev)
}

func TestDialogKeysMerge(t *testing.T) {
	dst := &DialogKeys{TogglePWVisibility: keyBinding("ctrl+j")}
	if errs := dst.merge(nil); len(errs) != 0 {
		t.Fatalf("nil source: unexpected errors: %v", errs)
	}

	src := &DialogKeys{TogglePWVisibility: keyBinding("j"), Close: keyBinding("esc")}
	if errs := dst.merge(src); len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	assertKeyBinding(t, "focus_down", dst.TogglePWVisibility, "j")
	assertKeyBinding(t, "close", dst.Close, "esc")
	assertNilKeyBinding(t, "focus_up", dst.Accept)
}

func TestWifiKeysMerge(t *testing.T) {
	dst := &NetworksKeys{CreateProfile: keyBinding("a", "c")}
	if errs := dst.merge(nil); len(errs) != 0 {
		t.Fatalf("nil source: unexpected errors: %v", errs)
	}

	src := &NetworksKeys{CreateProfile: keyBinding("p")}
	if errs := dst.merge(src); len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	assertKeyBinding(t, "create_profile", dst.CreateProfile, "p")
	assertNilKeyBinding(t, "create_hotspot", dst.CreateHotspot)
}

func TestWifiAvailableKeysMerge(t *testing.T) {
	dst := &AvailableNetworksKeys{Connect: keyBinding("enter")}
	if errs := dst.merge(nil); len(errs) != 0 {
		t.Fatalf("nil source: unexpected errors: %v", errs)
	}

	src := &AvailableNetworksKeys{Connect: keyBinding("space")}
	if errs := dst.merge(src); len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	assertKeyBinding(t, "connect", dst.Connect, "space")
}

func TestWifiSavedKeysMerge(t *testing.T) {
	dst := &NetworkProfilesKeys{Delete: keyBinding("d", "delete")}
	if errs := dst.merge(nil); len(errs) != 0 {
		t.Fatalf("nil source: unexpected errors: %v", errs)
	}

	src := &NetworkProfilesKeys{Delete: keyBinding("x"), Disconnect: keyBinding("ctrl+space")}
	if errs := dst.merge(src); len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	assertKeyBinding(t, "delete", dst.Delete, "x")
	assertKeyBinding(t, "disconnect", dst.Disconnect, "ctrl+space")
	assertNilKeyBinding(t, "edit", dst.Edit)
}
