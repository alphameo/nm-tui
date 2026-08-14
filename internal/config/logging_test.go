package config

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestDefaultLogConfig(t *testing.T) {
	t.Run("with XDG_STATE_HOME", func(t *testing.T) {
		dir := t.TempDir()
		t.Setenv("XDG_STATE_HOME", dir)

		cfg := DefaultLogConfig()
		assertNoNilFields(t, cfg)

		if got, want := *cfg.Level, LogError; got != want {
			t.Errorf("Level = %q, want %q", got, want)
		}
		wantPath := filepath.Join(dir, appName, "log")
		if got := *cfg.FilePath; got != wantPath {
			t.Errorf("FilePath = %q, want %q", got, wantPath)
		}
	})

	t.Run("falls back to HOME", func(t *testing.T) {
		t.Setenv("XDG_STATE_HOME", "")
		home := t.TempDir()
		t.Setenv("HOME", home)

		cfg := DefaultLogConfig()
		assertNoNilFields(t, cfg)

		if got, want := *cfg.Level, LogError; got != want {
			t.Errorf("Level = %q, want %q", got, want)
		}
		wantPath := filepath.Join(home, ".local", "state", appName, "log")
		if got := *cfg.FilePath; got != wantPath {
			t.Errorf("FilePath = %q, want %q", got, wantPath)
		}
	})
}

func TestValidLogLevel(t *testing.T) {
	for _, lvl := range []string{LogDebug, LogInfo, LogWarn, LogError} {
		if !validLogLevel(lvl) {
			t.Errorf("validLogLevel(%q) = false, want true", lvl)
		}
	}
	for _, lvl := range []string{"DEBUG", "Info", "WARN", "Error"} {
		if !validLogLevel(lvl) {
			t.Errorf("validLogLevel(%q) = false, want true (case-insensitive)", lvl)
		}
	}
	for _, lvl := range []string{"", "verbose", "trace", "fatal", "info ", "debug\n"} {
		if validLogLevel(lvl) {
			t.Errorf("validLogLevel(%q) = true, want false", lvl)
		}
	}
}

func TestLogConfigMerge(t *testing.T) {
	t.Run("valid level and path applied", func(t *testing.T) {
		dst := DefaultLogConfig()
		src := &LogConfig{Level: new("debug"), FilePath: new("/tmp/x.log")}
		if errs := dst.merge(src); len(errs) != 0 {
			t.Fatalf("unexpected errors: %v", errs)
		}
		if *dst.Level != "debug" {
			t.Errorf("Level = %q, want debug", *dst.Level)
		}
		if *dst.FilePath != "/tmp/x.log" {
			t.Errorf("FilePath = %q, want /tmp/x.log", *dst.FilePath)
		}
	})

	t.Run("tilde path is expanded", func(t *testing.T) {
		home := t.TempDir()
		t.Setenv("HOME", home)

		dst := DefaultLogConfig()
		src := &LogConfig{FilePath: new("~/nm-tui.log")}
		if errs := dst.merge(src); len(errs) != 0 {
			t.Fatalf("unexpected errors: %v", errs)
		}
		want := filepath.Join(home, "nm-tui.log")
		if *dst.FilePath != want {
			t.Errorf("FilePath = %q, want %q", *dst.FilePath, want)
		}
	})

	t.Run("nil source no-op", func(t *testing.T) {
		dst := DefaultLogConfig()
		if errs := dst.merge(nil); len(errs) != 0 {
			t.Fatalf("unexpected errors: %v", errs)
		}
	})

	t.Run("empty source no-op", func(t *testing.T) {
		dst := DefaultLogConfig()
		if errs := dst.merge(&LogConfig{}); len(errs) != 0 {
			t.Fatalf("unexpected errors: %v", errs)
		}
	})

	t.Run("invalid level errors and keeps default", func(t *testing.T) {
		dst := DefaultLogConfig()
		src := &LogConfig{Level: new("verbose")}
		errs := dst.merge(src)
		if len(errs) != 1 {
			t.Fatalf("want 1 error, got %v", errs)
		}
		if !strings.Contains(errs[0].Error(), `invalid log level: "verbose"`) {
			t.Errorf("unexpected error: %v", errs[0])
		}
		if *dst.Level != LogError {
			t.Errorf("Level = %q, want unchanged %q", *dst.Level, LogError)
		}
	})

	t.Run("empty file path errors and keeps default", func(t *testing.T) {
		dst := DefaultLogConfig()
		src := &LogConfig{FilePath: new("")}
		errs := dst.merge(src)
		if len(errs) != 1 {
			t.Fatalf("want 1 error, got %v", errs)
		}
		if !strings.Contains(errs[0].Error(), "empty log filepath") {
			t.Errorf("unexpected error: %v", errs[0])
		}
	})

	t.Run("default keyword skips both", func(t *testing.T) {
		dst := DefaultLogConfig()
		src := &LogConfig{Level: new(defaultKeyword), FilePath: new(defaultKeyword)}
		if errs := dst.merge(src); len(errs) != 0 {
			t.Fatalf("unexpected errors: %v", errs)
		}
		if *dst.Level != LogError {
			t.Errorf("Level = %q, want unchanged %q", *dst.Level, LogError)
		}
	})

	t.Run("multiple errors collected", func(t *testing.T) {
		dst := DefaultLogConfig()
		src := &LogConfig{Level: new("verbose"), FilePath: new("")}
		errs := dst.merge(src)
		if len(errs) != 2 {
			t.Fatalf("want 2 errors, got %v", errs)
		}
	})
}

func TestExpandPath(t *testing.T) {
	t.Run("tilde expands to home", func(t *testing.T) {
		home := t.TempDir()
		t.Setenv("HOME", home)

		got := expandPath("~/nm-tui/log")
		want := filepath.Join(home, "nm-tui", "log")
		if got != want {
			t.Errorf("expandPath(%q) = %q, want %q", "~/nm-tui/log", got, want)
		}
	})

	t.Run("bare tilde unchanged", func(t *testing.T) {
		if got := expandPath("~"); got != "~" {
			t.Errorf("expandPath(%q) = %q, want %q", "~", got, "~")
		}
	})

	t.Run("env vars expanded", func(t *testing.T) {
		dir := t.TempDir()
		t.Setenv("NM_TUI_TEST_DIR", dir)

		got := expandPath("$NM_TUI_TEST_DIR/log")
		want := filepath.Join(dir, "log")
		if got != want {
			t.Errorf("expandPath($NM_TUI_TEST_DIR/log) = %q, want %q", got, want)
		}
	})

	t.Run("absolute path unchanged", func(t *testing.T) {
		p := "/tmp/nm-tui.log"
		if got := expandPath(p); got != p {
			t.Errorf("expandPath(%q) = %q, want %q", p, got, p)
		}
	})
}

func TestResolveConfigPath(t *testing.T) {
	t.Run("with XDG_CONFIG_HOME", func(t *testing.T) {
		dir := t.TempDir()
		t.Setenv("XDG_CONFIG_HOME", dir)

		got, err := resolveConfigPath()
		if err != nil {
			t.Fatalf("resolveConfigPath() error: %v", err)
		}
		want := filepath.Join(dir, appName, configFileName)
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("error when no config dir", func(t *testing.T) {
		t.Setenv("XDG_CONFIG_HOME", "")
		t.Setenv("HOME", "")

		if _, err := resolveConfigPath(); err == nil {
			t.Fatal("resolveConfigPath() expected error")
		}
	})
}
