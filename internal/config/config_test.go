package config

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/calico32/kdl-go"
)

func TestDefaultConfig(t *testing.T) {
	t.Setenv("XDG_STATE_HOME", t.TempDir())

	cfg := DefaultConfig()

	if cfg.Colors == nil {
		t.Fatal("Colors is nil")
	}
	if cfg.Keys == nil {
		t.Fatal("Keys is nil")
	}
	if cfg.Logging == nil {
		t.Fatal("Logging is nil")
	}
	if cfg.Icons == nil {
		t.Fatal("Icons is nil")
	}
	if cfg.NotifCloseTime == nil {
		t.Fatal("NotifCloseTime is nil")
	}
	if got, want := *cfg.NotifCloseTime, 50; got != want {
		t.Errorf("NotifCloseTime = %d, want %d", got, want)
	}

	if got, want := *cfg.Colors.Text, ColorNone; got != want {
		t.Errorf("Colors.Text = %q, want %q", got, want)
	}
	if got, want := *cfg.Colors.Accent, ColorBlue; got != want {
		t.Errorf("Colors.Accent = %q, want %q", got, want)
	}
	if got, want := *cfg.Colors.Muted, ColorBrightBlack; got != want {
		t.Errorf("Colors.Muted = %q, want %q", got, want)
	}
	if got, want := *cfg.Colors.Error, ColorRed; got != want {
		t.Errorf("Colors.Error = %q, want %q", got, want)
	}
	if got, want := *cfg.Colors.Notif, ColorYellow; got != want {
		t.Errorf("Colors.Notif = %q, want %q", got, want)
	}

	if got, want := *cfg.Logging.Level, LogError; got != want {
		t.Errorf("Logging.Level = %q, want %q", got, want)
	}

	if cfg.Icons.NerdPreset == nil || *cfg.Icons.NerdPreset {
		t.Errorf("Icons.NerdPreset = %v, want false", cfg.Icons.NerdPreset)
	}
	if got, want := *cfg.Icons.BorderStyle, BorderASCII; got != want {
		t.Errorf("Icons.BorderStyle = %q, want %q", got, want)
	}
	if got, want := *cfg.Icons.SpinnerStyle, SpinnerLine; got != want {
		t.Errorf("Icons.SpinnerStyle = %q, want %q", got, want)
	}
}

func TestConfigMerge(t *testing.T) {
	t.Run("empty source produces no errors", func(t *testing.T) {
		cfg := DefaultConfig()
		if errs := cfg.merge(&Config{}); len(errs) != 0 {
			t.Errorf("unexpected errors: %v", errs)
		}
	})

	t.Run("nil child sections are handled", func(t *testing.T) {
		cfg := DefaultConfig()
		src := Config{
			Logging: &LogConfig{},
			Colors:  &ColorConfig{},
			Keys:    &KeyConfig{},
			Icons:   &IconConfig{},
		}
		if errs := cfg.merge(&src); len(errs) != 0 {
			t.Errorf("unexpected errors: %v", errs)
		}
	})

	t.Run("notification close time valid", func(t *testing.T) {
		cfg := DefaultConfig()
		src := Config{NotifCloseTime: new(250)}
		if errs := cfg.merge(&src); len(errs) != 0 {
			t.Fatalf("unexpected errors: %v", errs)
		}
		if got, want := *cfg.NotifCloseTime, 250; got != want {
			t.Errorf("NotifCloseTime = %d, want %d", got, want)
		}
	})

	t.Run("notification close time invalid", func(t *testing.T) {
		cfg := DefaultConfig()
		src := Config{NotifCloseTime: new(-1)}
		errs := cfg.merge(&src)
		if len(errs) != 1 {
			t.Fatalf("want 1 error, got %v", errs)
		}
		if !strings.Contains(errs[0].Error(), "notification_close_time") {
			t.Errorf("error does not mention notification_close_time: %v", errs[0])
		}
		if got, want := *cfg.NotifCloseTime, 50; got != want {
			t.Errorf("NotifCloseTime should stay default, got %d want %d", got, want)
		}
	})

	t.Run("nerd preset swaps icons to nerd defaults", func(t *testing.T) {
		cfg := DefaultConfig()
		src := Config{Icons: &IconConfig{
			NerdPreset:  new(true),
			BorderStyle: new(BorderSquare),
		}}
		if errs := cfg.merge(&src); len(errs) != 0 {
			t.Fatalf("unexpected errors: %v", errs)
		}
		if got, want := *cfg.Icons.BorderStyle, BorderSquare; got != want {
			t.Errorf("BorderStyle = %q, want %q", got, want)
		}
		if got, want := *cfg.Icons.SpinnerStyle, SpinnerMeter; got != want {
			t.Errorf("SpinnerStyle = %q, want nerd default %q", got, want)
		}
		if got, want := *cfg.Icons.ToggleOff, " "; got != want {
			t.Errorf("ToggleOff = %q, want nerd default %q", got, want)
		}
	})

	t.Run("icon merge errors are propagated", func(t *testing.T) {
		cfg := DefaultConfig()
		src := Config{Icons: &IconConfig{BorderStyle: new("bogus")}}
		errs := cfg.merge(&src)
		if len(errs) != 1 {
			t.Fatalf("want 1 error, got %v", errs)
		}
		if !strings.Contains(errs[0].Error(), "border_style") {
			t.Errorf("unexpected error: %v", errs[0])
		}
	})

	t.Run("logging merge errors are propagated", func(t *testing.T) {
		cfg := DefaultConfig()
		src := Config{Logging: &LogConfig{Level: new("verbose")}}
		errs := cfg.merge(&src)
		if len(errs) != 1 {
			t.Fatalf("want 1 error, got %v", errs)
		}
		if !strings.Contains(errs[0].Error(), "invalid log level") {
			t.Errorf("unexpected error: %v", errs[0])
		}
	})

	t.Run("colors merge errors are propagated", func(t *testing.T) {
		cfg := DefaultConfig()
		src := Config{Colors: &ColorConfig{Text: new("notacolor")}}
		errs := cfg.merge(&src)
		if len(errs) != 1 {
			t.Fatalf("want 1 error, got %v", errs)
		}
		if !strings.Contains(errs[0].Error(), "text color") {
			t.Errorf("unexpected error: %v", errs[0])
		}
	})

	t.Run("keys merge errors are propagated", func(t *testing.T) {
		cfg := DefaultConfig()
		src := Config{Keys: &KeyConfig{Toggle: keyBinding("notakey")}}
		errs := cfg.merge(&src)
		if len(errs) != 1 {
			t.Fatalf("want 1 error, got %v", errs)
		}
		if !strings.Contains(errs[0].Error(), "invalid key toggle") {
			t.Errorf("unexpected error: %v", errs[0])
		}
	})
}

func TestValidateTime(t *testing.T) {
	tests := []struct {
		name string
		in   int
		want bool
	}{
		{"positive", 1, true},
		{"large positive", 10000, true},
		{"zero", 0, false},
		{"negative", -5, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateTime(tt.in)
			if (err == nil) != tt.want {
				t.Errorf("validateTime(%d) error = %v, want ok=%v", tt.in, err, tt.want)
			}
		})
	}
}

const fullConfigKDL = `
notification_close_time 250

colors {
    text "red"
    accent "#865fff"
    muted "default"
    error "bright_red"
    notification "default"
}

icons {
    nerd_preset false
    border_style "thick_square"
    spinner_style "jump"
    input_cursor_shape "underline"
    toggle_off "[ ]"
    toggle_on "[x]"
    password_hidden_character "*"
    error "!"
    check "v"
    connection "con"
    signal "sig"
    access_point "ap"
    infra "infr"
    mesh "#"
    ad_hoc "ah"
}

logging {
    level "debug"
    file_path "/tmp/nm-tui-test.log"
}

keys {
    toggle "space"
    rescan "r"
    rescan_focused "ctrl+r"
    focus_1 "1"
    focus_2 "2"
    focus_3 "3"
    focus_4 "4"
    focus_5 "5"
    focus_6 "6"
    focus_7 "7"
    focus_8 "8"
    focus_9 "9"
    focus_10 "10"
    main {
        next_tab "]"
        prev_tab "["
        focus_next "tab"
        focus_prev "shift+tab"
        quit "esc" "ctrl+c" "q" "ctrl+q"
    }
    dialog {
        focus_down "ctrl+j"
        focus_up "ctrl+k"
        toggle_pw_visibility "ctrl+p"
        accept "ctrl+enter"
        close "esc" "ctrl+q" "ctrl+c"
    }
    wifi {
        create_profile "a" "c"
        open_network_login "l"
        enable_hotspot "ctrl+h"
        create_hotspot "h"
    }
    wifi_available {
        connect "enter"
    }
    wifi_saved {
        edit "enter"
        connect "space"
        disconnect "ctrl+space"
        delete "d" "delete"
    }
}
`

func writeConfigFile(t *testing.T, src string) {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)
	cfgDir := filepath.Join(dir, appName)
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(cfgDir, configFileName)
	if err := os.WriteFile(path, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestLoad(t *testing.T) {
	writeConfigFile(t, fullConfigKDL)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if cfg.NotifCloseTime == nil || *cfg.NotifCloseTime != 250 {
		t.Errorf("NotifCloseTime = %v, want 250", cfg.NotifCloseTime)
	}
	if cfg.Icons.InputCursorShape == nil || *cfg.Icons.InputCursorShape != "underline" {
		t.Errorf("Icons.InputCursorShape = %v, want underline", cfg.Icons.InputCursorShape)
	}

	if cfg.Colors == nil || cfg.Colors.Text == nil || *cfg.Colors.Text != "red" {
		t.Errorf("Colors.Text = %v, want red", cfg.Colors.Text)
	}
	if cfg.Colors.Accent == nil || *cfg.Colors.Accent != "#865fff" {
		t.Errorf("Colors.Accent = %v, want #865fff", cfg.Colors.Accent)
	}
	if cfg.Colors.Muted == nil || *cfg.Colors.Muted != "default" {
		t.Errorf("Colors.Muted = %v, want default", cfg.Colors.Muted)
	}

	if cfg.Icons == nil || cfg.Icons.NerdPreset == nil || *cfg.Icons.NerdPreset {
		t.Errorf("Icons.NerdPreset = %v, want false", cfg.Icons.NerdPreset)
	}
	if cfg.Icons.BorderStyle == nil || *cfg.Icons.BorderStyle != "thick_square" {
		t.Errorf("Icons.BorderStyle = %v, want thick_square", cfg.Icons.BorderStyle)
	}

	if cfg.Logging == nil || cfg.Logging.Level == nil || *cfg.Logging.Level != "debug" {
		t.Errorf("Logging.Level = %v, want debug", cfg.Logging.Level)
	}
	if cfg.Logging.FilePath == nil || *cfg.Logging.FilePath != "/tmp/nm-tui-test.log" {
		t.Errorf("Logging.FilePath = %v, want /tmp/nm-tui-test.log", cfg.Logging.FilePath)
	}

	if cfg.Keys == nil || cfg.Keys.Main == nil {
		t.Fatal("keys.main not decoded")
	}
	assertKeyBinding(t, "keys.toggle", cfg.Keys.Toggle, "space")
	assertKeyBinding(t, "keys.rescan_focused", cfg.Keys.RescanFocused, "ctrl+r")
	assertKeyBinding(t, "keys.main.next_tab", cfg.Keys.Main.NextTab, "]")
	assertKeyBinding(t, "keys.main.quit", cfg.Keys.Main.Quit, "esc", "ctrl+c", "q", "ctrl+q")
	assertKeyBinding(t, "keys.dialog.close", cfg.Keys.Dialog.Close, "esc", "ctrl+q", "ctrl+c")
	assertKeyBinding(t, "keys.wifi.create_profile", cfg.Keys.Wifi.CreateProfile, "a", "c")
	assertKeyBinding(t, "keys.wifi_available.connect", cfg.Keys.WifiAvailable.Connect, "enter")
	assertKeyBinding(t, "keys.wifi_saved.delete", cfg.Keys.WifiSaved.Delete, "d", "delete")
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

func TestLoadOrDefaultsValid(t *testing.T) {
	writeConfigFile(t, fullConfigKDL)

	cfg, err := LoadOrDefaults()
	if err != nil {
		t.Fatalf("LoadOrDefaults() error: %v", err)
	}

	if *cfg.NotifCloseTime != 250 {
		t.Errorf("NotifCloseTime = %d, want 250", *cfg.NotifCloseTime)
	}
	if *cfg.Icons.InputCursorShape != "underline" {
		t.Errorf("Icons.InputCursorShape = %q, want underline", *cfg.Icons.InputCursorShape)
	}

	if *cfg.Colors.Text != "red" {
		t.Errorf("Colors.Text = %q, want red", *cfg.Colors.Text)
	}
	if *cfg.Colors.Accent != "#865fff" {
		t.Errorf("Colors.Accent = %q, want #865fff", *cfg.Colors.Accent)
	}
	if *cfg.Colors.Muted != ColorBrightBlack {
		t.Errorf("Colors.Muted = %q, want default %q", *cfg.Colors.Muted, ColorBrightBlack)
	}
	if *cfg.Colors.Error != "bright_red" {
		t.Errorf("Colors.Error = %q, want bright_red", *cfg.Colors.Error)
	}
	if *cfg.Colors.Notif != ColorYellow {
		t.Errorf("Colors.Notif = %q, want default %q", *cfg.Colors.Notif, ColorYellow)
	}

	if *cfg.Icons.BorderStyle != "thick_square" {
		t.Errorf("Icons.BorderStyle = %q, want thick_square", *cfg.Icons.BorderStyle)
	}
	if *cfg.Icons.SpinnerStyle != "jump" {
		t.Errorf("Icons.SpinnerStyle = %q, want jump", *cfg.Icons.SpinnerStyle)
	}
	if *cfg.Icons.NerdPreset {
		t.Errorf("Icons.NerdPreset = true, want false")
	}

	if *cfg.Logging.Level != "debug" {
		t.Errorf("Logging.Level = %q, want debug", *cfg.Logging.Level)
	}
	if *cfg.Logging.FilePath != "/tmp/nm-tui-test.log" {
		t.Errorf("Logging.FilePath = %q, want /tmp/nm-tui-test.log", *cfg.Logging.FilePath)
	}

	assertKeyBinding(t, "keys.main.quit", cfg.Keys.Main.Quit, "esc", "ctrl+c", "q", "ctrl+q")
}

const invalidConfigKDL = `
notification_close_time -5
colors {
    accent "notacolor"
}
logging {
    level "verbose"
}
icons {
    border_style "dashed"
    input_cursor_shape "squiggle"
}
keys {
    toggle "notakey"
}
`

func TestLoadOrDefaultsInvalid(t *testing.T) {
	t.Setenv("XDG_STATE_HOME", t.TempDir())
	writeConfigFile(t, invalidConfigKDL)

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

	if *cfg.NotifCloseTime != 50 {
		t.Errorf("NotifCloseTime = %d, want default 50", *cfg.NotifCloseTime)
	}
	if *cfg.Icons.InputCursorShape != CursorBar {
		t.Errorf("Icons.InputCursorShape = %q, want default %q", *cfg.Icons.InputCursorShape, CursorBar)
	}
	if *cfg.Colors.Accent != ColorBlue {
		t.Errorf("Colors.Accent = %q, want default %q", *cfg.Colors.Accent, ColorBlue)
	}
	if *cfg.Logging.Level != LogError {
		t.Errorf("Logging.Level = %q, want default %q", *cfg.Logging.Level, LogError)
	}
	if *cfg.Icons.BorderStyle != BorderASCII {
		t.Errorf("Icons.BorderStyle = %q, want default %q", *cfg.Icons.BorderStyle, BorderASCII)
	}
	assertKeyBinding(t, "keys.toggle", cfg.Keys.Toggle, "space")
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

	if *cfg.NotifCloseTime != 50 {
		t.Errorf("NotifCloseTime = %d, want default 50", *cfg.NotifCloseTime)
	}
	if *cfg.Icons.InputCursorShape != CursorBar {
		t.Errorf("Icons.InputCursorShape = %q, want default %q", *cfg.Icons.InputCursorShape, CursorBar)
	}
	if *cfg.Colors.Text != ColorNone {
		t.Errorf("Colors.Text = %q, want default %q", *cfg.Colors.Text, ColorNone)
	}
	assertKeyBinding(t, "keys.main.quit", cfg.Keys.Main.Quit, "esc", "ctrl+c", "q", "ctrl+q")
}

func TestKeyBindingUnmarshalKDL(t *testing.T) {
	t.Run("single argument", func(t *testing.T) {
		doc, err := kdl.ParseString(`key "space"`)
		if err != nil {
			t.Fatal(err)
		}
		node := doc.GetNode("key")
		if node == nil {
			t.Fatal("node not found")
		}
		var kb KeyBinding
		if err := kb.UnmarshalKDL(node); err != nil {
			t.Fatalf("UnmarshalKDL() error: %v", err)
		}
		if want := (KeyBinding{"space"}); !reflect.DeepEqual(kb, want) {
			t.Errorf("got %v, want %v", kb, want)
		}
	})

	t.Run("multiple arguments", func(t *testing.T) {
		doc, err := kdl.ParseString(`quit "esc" "ctrl+c" "q" "ctrl+q"`)
		if err != nil {
			t.Fatal(err)
		}
		node := doc.GetNode("quit")
		if node == nil {
			t.Fatal("node not found")
		}
		var kb KeyBinding
		if err := kb.UnmarshalKDL(node); err != nil {
			t.Fatalf("UnmarshalKDL() error: %v", err)
		}
		want := KeyBinding{"esc", "ctrl+c", "q", "ctrl+q"}
		if !reflect.DeepEqual(kb, want) {
			t.Errorf("got %v, want %v", kb, want)
		}
	})

	t.Run("empty node", func(t *testing.T) {
		doc, err := kdl.ParseString(`key`)
		if err != nil {
			t.Fatal(err)
		}
		node := doc.GetNode("key")
		if node == nil {
			t.Fatal("node not found")
		}
		var kb KeyBinding
		if err := kb.UnmarshalKDL(node); err != nil {
			t.Fatalf("UnmarshalKDL() error: %v", err)
		}
		if len(kb) != 0 {
			t.Errorf("got %v, want empty binding", kb)
		}
	})
}
