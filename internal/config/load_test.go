package config

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/calico32/kdl-go"
)

func readTestdata(t *testing.T, name string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatalf("read testdata %s: %v", name, err)
	}
	return string(b)
}

func writeConfigFile(t *testing.T, src string) {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	cfgDir := filepath.Join(dir, appName)
	if err := os.MkdirAll(cfgDir, 0o750); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(cfgDir, configFileName)
	if err := os.WriteFile(path, []byte(src), 0o600); err != nil { //nolint:gosec // writes test config to a t.TempDir path
		t.Fatal(err)
	}
}

func assertConfigEqual(t *testing.T, name string, got, want any) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("%s mismatch:\n got: %#v\nwant: %#v", name, got, want)
	}
}

func TestLoadExampleConfig(t *testing.T) {
	src, err := os.ReadFile(filepath.Join("..", "..", "config.example.kdl"))
	if err != nil {
		t.Fatalf("read config.example.kdl: %v", err)
	}
	writeConfigFile(t, string(src))

	cfg, err := LoadOrDefaults()
	if err != nil {
		t.Fatalf("LoadOrDefaults() error: %v", err)
	}
	assertNoNilFields(t, cfg)

	// config.example.kdl mirrors the built-in defaults; the only explicit
	// deviation is the example log path. Update this expectation by hand
	// when you change config.example.kdl.
	want := DefaultConfig()
	want.Logging.FilePath = new("/home/user/.local/state/nm-tui/log")
	assertConfigEqual(t, "config.example.kdl decodes to defaults", cfg, want)
}

func TestLoadOrDefaultsOverrides(t *testing.T) {
	writeConfigFile(t, `
notification_close_time 300
colors {
    text "#111111"
    muted "default"
}
keys {
    toggle "t"
}
`)

	cfg, err := LoadOrDefaults()
	if err != nil {
		t.Fatalf("LoadOrDefaults() error: %v", err)
	}

	want := DefaultConfig()
	want.NotifCloseTime = new(300)
	want.Colors.Text = new("#111111")
	want.Keys.Toggle = keyBinding("t")
	assertConfigEqual(t, "overrides applied on top of defaults", cfg, want)
}

func TestLoadOrDefaultsInvalid(t *testing.T) {
	t.Setenv("XDG_STATE_HOME", t.TempDir())
	writeConfigFile(t, readTestdata(t, "invalid.config.kdl"))

	cfg, err := LoadOrDefaults()
	if err == nil {
		t.Fatal("LoadOrDefaults() expected error for invalid config")
	}

	wantFragments := []string{
		"notification_close_time",
		"cursor_shape",
		"accent color",
		"invalid log level",
		"border_style",
		"invalid key toggle",
	}
	for _, frag := range wantFragments {
		if !strings.Contains(err.Error(), frag) {
			t.Errorf("error does not contain %q: %v", frag, err)
		}
	}

	assertConfigEqual(t, "invalid values leave defaults intact", cfg, DefaultConfig())
}

func TestLoadOrDefaultsNoConfigFile(t *testing.T) {
	t.Setenv("XDG_STATE_HOME", t.TempDir())
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)

	cfg, err := LoadOrDefaults()
	if err == nil {
		t.Fatal("LoadOrDefaults() expected error when no config file exists")
	}
	if !strings.Contains(err.Error(), "user config loading failed") {
		t.Errorf("error does not mention user config loading failed: %v", err)
	}

	assertConfigEqual(t, "missing config keeps defaults", cfg, DefaultConfig())
}

func TestLoadMissingFile(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)

	_, err := Load()
	if err == nil {
		t.Fatal("Load() expected error for missing config file")
	}
	if !strings.Contains(err.Error(), "open config") {
		t.Errorf("error does not mention open config: %v", err)
	}
}

func TestLoadDecodeError(t *testing.T) {
	writeConfigFile(t, `notification_close_time "notanumber"`)

	_, err := Load()
	if err == nil {
		t.Fatal("Load() expected error for undecodable config")
	}
	if !strings.Contains(err.Error(), "decode config") {
		t.Errorf("error does not mention decode config: %v", err)
	}
}

func TestLoadPathResolutionError(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", "")
	t.Setenv("HOME", "")

	_, err := Load()
	if err == nil {
		t.Fatal("Load() expected error when config dir cannot be resolved")
	}
	if !strings.Contains(err.Error(), "resolve config path") {
		t.Errorf("error does not mention resolve config path: %v", err)
	}
}

func TestKeyBindingUnmarshalKDL(t *testing.T) {
	decode := func(t *testing.T, src string, kb *KeyBinding) error {
		t.Helper()
		doc, err := kdl.ParseString(src)
		if err != nil {
			t.Fatal(err)
		}
		name, _, _ := strings.Cut(src, " ")
		node := doc.GetNode(name)
		if node == nil {
			t.Fatalf("node %q not found", name)
		}
		return kb.UnmarshalKDL(node)
	}

	t.Run("single argument", func(t *testing.T) {
		var kb KeyBinding
		if err := decode(t, `key "space"`, &kb); err != nil {
			t.Fatalf("UnmarshalKDL() error: %v", err)
		}
		if want := (KeyBinding{"space"}); !reflect.DeepEqual(kb, want) {
			t.Errorf("got %v, want %v", kb, want)
		}
	})

	t.Run("multiple arguments", func(t *testing.T) {
		var kb KeyBinding
		if err := decode(t, `quit "esc" "ctrl+c" "q" "ctrl+q"`, &kb); err != nil {
			t.Fatalf("UnmarshalKDL() error: %v", err)
		}
		want := KeyBinding{"esc", "ctrl+c", "q", "ctrl+q"}
		if !reflect.DeepEqual(kb, want) {
			t.Errorf("got %v, want %v", kb, want)
		}
	})

	t.Run("empty node", func(t *testing.T) {
		var kb KeyBinding
		if err := decode(t, `key`, &kb); err != nil {
			t.Fatalf("UnmarshalKDL() error: %v", err)
		}
		if len(kb) != 0 {
			t.Errorf("got %v, want empty binding", kb)
		}
	})
}
