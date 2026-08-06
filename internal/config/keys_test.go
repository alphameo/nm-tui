package config

import (
	"fmt"
	"strings"
	"testing"
)

func TestDefaultKeys(t *testing.T) {
	k := DefaultKeys()

	assertKeyBinding(t, "toggle", k.Toggle, "space")
	assertKeyBinding(t, "rescan", k.Rescan, "r")
	assertKeyBinding(t, "rescan_focused", k.RescanFocused, "ctrl+r")

	foci := []*KeyBinding{
		k.Focus1, k.Focus2, k.Focus3, k.Focus4, k.Focus5,
		k.Focus6, k.Focus7, k.Focus8, k.Focus9, k.Focus10,
	}
	for i, fb := range foci {
		assertKeyBinding(t, fmt.Sprintf("focus_%d", i+1), fb, fmt.Sprintf("%d", i+1))
	}

	if k.Main == nil {
		t.Fatal("Main is nil")
	}
	assertKeyBinding(t, "main.next_tab", k.Main.NextTab, "]")
	assertKeyBinding(t, "main.prev_tab", k.Main.PrevTab, "[")
	assertKeyBinding(t, "main.focus_next", k.Main.FocusNext, "tab")
	assertKeyBinding(t, "main.focus_prev", k.Main.FocusPrev, "shift+tab")
	assertKeyBinding(t, "main.quit", k.Main.Quit, "esc", "ctrl+c", "q", "ctrl+q")

	if k.Dialog == nil {
		t.Fatal("Dialog is nil")
	}
	assertKeyBinding(t, "dialog.focus_down", k.Dialog.FocusDown, "ctrl+j")
	assertKeyBinding(t, "dialog.focus_up", k.Dialog.FocusUp, "ctrl+k")
	assertKeyBinding(t, "dialog.toggle_pw_visibility", k.Dialog.TogglePWVisibility, "ctrl+p")
	assertKeyBinding(t, "dialog.accept", k.Dialog.Accept, "ctrl+enter")
	assertKeyBinding(t, "dialog.close", k.Dialog.Close, "esc", "ctrl+q", "ctrl+c")

	if k.Wifi == nil {
		t.Fatal("Wifi is nil")
	}
	assertKeyBinding(t, "wifi.create_profile", k.Wifi.CreateProfile, "a", "c")
	assertKeyBinding(t, "wifi.open_network_login", k.Wifi.OpenCaptivePortal, "l")
	assertKeyBinding(t, "wifi.enable_hotspot", k.Wifi.EnableHotspot, "ctrl+h")
	assertKeyBinding(t, "wifi.create_hotspot", k.Wifi.CreateHotspot, "h")

	if k.WifiAvailable == nil {
		t.Fatal("WifiAvailable is nil")
	}
	assertKeyBinding(t, "wifi_available.connect", k.WifiAvailable.Connect, "enter")

	if k.WifiSaved == nil {
		t.Fatal("WifiSaved is nil")
	}
	assertKeyBinding(t, "wifi_saved.edit", k.WifiSaved.Edit, "enter")
	assertKeyBinding(t, "wifi_saved.connect", k.WifiSaved.Connect, "space")
	assertKeyBinding(t, "wifi_saved.disconnect", k.WifiSaved.Disconnect, "ctrl+space")
	assertKeyBinding(t, "wifi_saved.delete", k.WifiSaved.Delete, "d", "delete")
}

func TestValidKeyName(t *testing.T) {
	for name := range validKey {
		name := name
		t.Run("named_"+name, func(t *testing.T) {
			if !validKeyName(name) {
				t.Errorf("validKeyName(%q) = false, want true", name)
			}
		})
	}

	for name := range validModifier {
		name := name
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
		assertKeyBinding(t, "main.next_tab", dst.Main.NextTab, "]")
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
	dst := &MainKeys{NextTab: keyBinding("]")}
	if errs := dst.merge(nil); len(errs) != 0 {
		t.Fatalf("nil source: unexpected errors: %v", errs)
	}

	src := &MainKeys{NextTab: keyBinding("n"), Quit: keyBinding("q")}
	if errs := dst.merge(src); len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	assertKeyBinding(t, "next_tab", dst.NextTab, "n")
	assertKeyBinding(t, "quit", dst.Quit, "q")
	assertNilKeyBinding(t, "prev_tab", dst.PrevTab)
}

func TestDialogKeysMerge(t *testing.T) {
	dst := &DialogKeys{FocusDown: keyBinding("ctrl+j")}
	if errs := dst.merge(nil); len(errs) != 0 {
		t.Fatalf("nil source: unexpected errors: %v", errs)
	}

	src := &DialogKeys{FocusDown: keyBinding("j"), Close: keyBinding("esc")}
	if errs := dst.merge(src); len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	assertKeyBinding(t, "focus_down", dst.FocusDown, "j")
	assertKeyBinding(t, "close", dst.Close, "esc")
	assertNilKeyBinding(t, "focus_up", dst.FocusUp)
}

func TestWifiKeysMerge(t *testing.T) {
	dst := &WifiKeys{CreateProfile: keyBinding("a", "c")}
	if errs := dst.merge(nil); len(errs) != 0 {
		t.Fatalf("nil source: unexpected errors: %v", errs)
	}

	src := &WifiKeys{CreateProfile: keyBinding("p")}
	if errs := dst.merge(src); len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	assertKeyBinding(t, "create_profile", dst.CreateProfile, "p")
	assertNilKeyBinding(t, "create_hotspot", dst.CreateHotspot)
}

func TestWifiAvailableKeysMerge(t *testing.T) {
	dst := &WifiAvailableKeys{Connect: keyBinding("enter")}
	if errs := dst.merge(nil); len(errs) != 0 {
		t.Fatalf("nil source: unexpected errors: %v", errs)
	}

	src := &WifiAvailableKeys{Connect: keyBinding("space")}
	if errs := dst.merge(src); len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	assertKeyBinding(t, "connect", dst.Connect, "space")
}

func TestWifiSavedKeysMerge(t *testing.T) {
	dst := &WifiSavedKeys{Delete: keyBinding("d", "delete")}
	if errs := dst.merge(nil); len(errs) != 0 {
		t.Fatalf("nil source: unexpected errors: %v", errs)
	}

	src := &WifiSavedKeys{Delete: keyBinding("x"), Disconnect: keyBinding("ctrl+space")}
	if errs := dst.merge(src); len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	assertKeyBinding(t, "delete", dst.Delete, "x")
	assertKeyBinding(t, "disconnect", dst.Disconnect, "ctrl+space")
	assertNilKeyBinding(t, "edit", dst.Edit)
}
