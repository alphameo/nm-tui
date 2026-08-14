package config_test

import (
	"strings"
	"testing"

	"github.com/alphameo/nm-tui/internal/config"
)

func TestDefaultKeys(t *testing.T) {
	t.Parallel()

	assertNoNilFields(t, config.DefaultKeys())
}

func TestKeyConfigMerge(t *testing.T) {
	t.Parallel()

	t.Run("nil source is a no-op", func(t *testing.T) {
		t.Parallel()

		dst := config.DefaultKeys()
		if errs := dst.Merge(nil); len(errs) != 0 {
			t.Fatalf("unexpected errors: %v", errs)
		}
	})

	t.Run("empty source is a no-op", func(t *testing.T) {
		t.Parallel()

		dst := config.DefaultKeys()
		if errs := dst.Merge(&config.KeyConfig{}); len(errs) != 0 {
			t.Fatalf("unexpected errors: %v", errs)
		}
		assertKeyBinding(t, "toggle", dst.Toggle, "space")
	})

	t.Run("partial override merges only set fields", func(t *testing.T) {
		t.Parallel()

		dst := config.DefaultKeys()
		src := &config.KeyConfig{
			Toggle: &config.KeyBinding{"t"},
			Main:   &config.MainKeys{Quit: &config.KeyBinding{"q"}},
		}
		if errs := dst.Merge(src); len(errs) != 0 {
			t.Fatalf("unexpected errors: %v", errs)
		}
		assertKeyBinding(t, "toggle", dst.Toggle, "t")
		assertKeyBinding(t, "main.quit", dst.Main.Quit, "q")
		assertKeyBinding(t, "main.next_tab", dst.Main.TabNext, "]")
		assertKeyBinding(t, "rescan", dst.Rescan, "r")
	})

	t.Run("invalid key propagates error", func(t *testing.T) {
		t.Parallel()

		dst := config.DefaultKeys()
		src := &config.KeyConfig{Toggle: &config.KeyBinding{"notakey"}}
		errs := dst.Merge(src)
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
	t.Parallel()

	dst := &config.MainKeys{TabNext: &config.KeyBinding{"]"}}
	if errs := dst.Merge(nil); len(errs) != 0 {
		t.Fatalf("nil source: unexpected errors: %v", errs)
	}

	src := &config.MainKeys{TabNext: &config.KeyBinding{"n"}, Quit: &config.KeyBinding{"q"}}
	if errs := dst.Merge(src); len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	assertKeyBinding(t, "next_tab", dst.TabNext, "n")
	assertKeyBinding(t, "quit", dst.Quit, "q")
	assertNilKeyBinding(t, "prev_tab", dst.TabPrev)
}

func TestDialogKeysMerge(t *testing.T) {
	t.Parallel()

	dst := &config.DialogKeys{TogglePWVisibility: &config.KeyBinding{"ctrl+j"}}
	if errs := dst.Merge(nil); len(errs) != 0 {
		t.Fatalf("nil source: unexpected errors: %v", errs)
	}

	src := &config.DialogKeys{TogglePWVisibility: &config.KeyBinding{"j"}, Close: &config.KeyBinding{"esc"}}
	if errs := dst.Merge(src); len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	assertKeyBinding(t, "focus_down", dst.TogglePWVisibility, "j")
	assertKeyBinding(t, "close", dst.Close, "esc")
	assertNilKeyBinding(t, "focus_up", dst.Accept)
}

func TestWifiKeysMerge(t *testing.T) {
	t.Parallel()

	dst := &config.NetworksKeys{CreateProfile: &config.KeyBinding{"a", "c"}}
	if errs := dst.Merge(nil); len(errs) != 0 {
		t.Fatalf("nil source: unexpected errors: %v", errs)
	}

	src := &config.NetworksKeys{CreateProfile: &config.KeyBinding{"p"}}
	if errs := dst.Merge(src); len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	assertKeyBinding(t, "create_profile", dst.CreateProfile, "p")
	assertNilKeyBinding(t, "create_hotspot", dst.CreateHotspot)
}

func TestWifiAvailableKeysMerge(t *testing.T) {
	t.Parallel()

	dst := &config.AvailableNetworksKeys{Connect: &config.KeyBinding{"enter"}}
	if errs := dst.Merge(nil); len(errs) != 0 {
		t.Fatalf("nil source: unexpected errors: %v", errs)
	}

	src := &config.AvailableNetworksKeys{Connect: &config.KeyBinding{"space"}}
	if errs := dst.Merge(src); len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	assertKeyBinding(t, "connect", dst.Connect, "space")
}

func TestWifiSavedKeysMerge(t *testing.T) {
	t.Parallel()

	dst := &config.NetworkProfilesKeys{Delete: &config.KeyBinding{"d", "delete"}}
	if errs := dst.Merge(nil); len(errs) != 0 {
		t.Fatalf("nil source: unexpected errors: %v", errs)
	}

	src := &config.NetworkProfilesKeys{Delete: &config.KeyBinding{"x"}, Disconnect: &config.KeyBinding{"ctrl+space"}}
	if errs := dst.Merge(src); len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	assertKeyBinding(t, "delete", dst.Delete, "x")
	assertKeyBinding(t, "disconnect", dst.Disconnect, "ctrl+space")
	assertNilKeyBinding(t, "edit", dst.Edit)
}

func TestMergeKeyList(t *testing.T) {
	t.Parallel()

	t.Run("nil source is a no-op", func(t *testing.T) {
		t.Parallel()

		dst := &config.KeyBinding{"a"}
		if errs := config.MergeKeyList(&dst, nil, "toggle"); len(errs) != 0 {
			t.Fatalf("unexpected errors: %v", errs)
		}
		assertKeyBinding(t, "dst", dst, "a")
	})

	t.Run("valid source overrides destination", func(t *testing.T) {
		t.Parallel()

		dst := &config.KeyBinding{"a"}
		src := &config.KeyBinding{"b", "c"}
		if errs := config.MergeKeyList(&dst, src, "toggle"); len(errs) != 0 {
			t.Fatalf("unexpected errors: %v", errs)
		}
		assertKeyBinding(t, "dst", dst, "b", "c")
	})

	t.Run("invalid source errors and keeps destination", func(t *testing.T) {
		t.Parallel()

		dst := &config.KeyBinding{"a"}
		src := &config.KeyBinding{"b", "notakey"}
		errs := config.MergeKeyList(&dst, src, "toggle")
		if len(errs) != 1 {
			t.Fatalf("want 1 error, got %v", errs)
		}
		if !strings.Contains(errs[0].Error(), `invalid key toggle: "notakey"`) {
			t.Errorf("unexpected error: %v", errs[0])
		}
		assertKeyBinding(t, "dst", dst, "a")
	})

	t.Run("multiple invalid keys produce multiple errors", func(t *testing.T) {
		t.Parallel()

		dst := &config.KeyBinding{"a"}
		src := &config.KeyBinding{"bad1", "bad2"}
		errs := config.MergeKeyList(&dst, src, "main.quit")
		if len(errs) != 2 {
			t.Fatalf("want 2 errors, got %v", errs)
		}
	})

	t.Run("nil destination is populated", func(t *testing.T) {
		t.Parallel()

		var dst *config.KeyBinding
		src := &config.KeyBinding{"x"}
		if errs := config.MergeKeyList(&dst, src, "toggle"); len(errs) != 0 {
			t.Fatalf("unexpected errors: %v", errs)
		}
		assertKeyBinding(t, "dst", dst, "x")
	})
}
