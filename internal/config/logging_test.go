package config_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/alphameo/nm-tui/internal/config"
)

func TestDefaultLogConfig(t *testing.T) {
	t.Run("with XDG_STATE_HOME", func(t *testing.T) {
		dir := t.TempDir()
		t.Setenv("XDG_STATE_HOME", dir)

		cfg := config.DefaultLogConfig()
		assertNoNilFields(t, cfg)

		if got, want := *cfg.Level, config.LogError; got != want {
			t.Errorf("Level = %q, want %q", got, want)
		}
		wantPath := filepath.Join(dir, config.AppName, "log")
		if got := *cfg.FilePath; got != wantPath {
			t.Errorf("FilePath = %q, want %q", got, wantPath)
		}
	})

	t.Run("falls back to HOME", func(t *testing.T) {
		t.Setenv("XDG_STATE_HOME", "")
		home := t.TempDir()
		t.Setenv("HOME", home)

		cfg := config.DefaultLogConfig()
		assertNoNilFields(t, cfg)

		if got, want := *cfg.Level, config.LogError; got != want {
			t.Errorf("Level = %q, want %q", got, want)
		}
		wantPath := filepath.Join(home, ".local", "state", config.AppName, "log")
		if got := *cfg.FilePath; got != wantPath {
			t.Errorf("FilePath = %q, want %q", got, wantPath)
		}
	})
}

func TestLogConfigMerge(t *testing.T) {
	t.Run("valid level and path applied", func(t *testing.T) {
		t.Parallel()

		dst := config.DefaultLogConfig()
		src := &config.LogConfig{Level: new("debug"), FilePath: new("/tmp/x.log")}
		if errs := dst.Merge(src); len(errs) != 0 {
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

		dst := config.DefaultLogConfig()
		src := &config.LogConfig{FilePath: new("~/nm-tui.log")}
		if errs := dst.Merge(src); len(errs) != 0 {
			t.Fatalf("unexpected errors: %v", errs)
		}
		want := filepath.Join(home, "nm-tui.log")
		if *dst.FilePath != want {
			t.Errorf("FilePath = %q, want %q", *dst.FilePath, want)
		}
	})

	t.Run("nil source no-op", func(t *testing.T) {
		t.Parallel()

		dst := config.DefaultLogConfig()
		if errs := dst.Merge(nil); len(errs) != 0 {
			t.Fatalf("unexpected errors: %v", errs)
		}
	})

	t.Run("empty source no-op", func(t *testing.T) {
		t.Parallel()

		dst := config.DefaultLogConfig()
		if errs := dst.Merge(&config.LogConfig{}); len(errs) != 0 {
			t.Fatalf("unexpected errors: %v", errs)
		}
	})

	t.Run("invalid level errors and keeps default", func(t *testing.T) {
		t.Parallel()

		dst := config.DefaultLogConfig()
		src := &config.LogConfig{Level: new("verbose")}
		errs := dst.Merge(src)
		if len(errs) != 1 {
			t.Fatalf("want 1 error, got %v", errs)
		}
		if !strings.Contains(errs[0].Error(), `invalid log level: "verbose"`) {
			t.Errorf("unexpected error: %v", errs[0])
		}
		if *dst.Level != config.LogError {
			t.Errorf("Level = %q, want unchanged %q", *dst.Level, config.LogError)
		}
	})

	t.Run("empty file path errors and keeps default", func(t *testing.T) {
		t.Parallel()

		dst := config.DefaultLogConfig()
		src := &config.LogConfig{FilePath: new("")}
		errs := dst.Merge(src)
		if len(errs) != 1 {
			t.Fatalf("want 1 error, got %v", errs)
		}
		if !strings.Contains(errs[0].Error(), "empty log filepath") {
			t.Errorf("unexpected error: %v", errs[0])
		}
	})

	t.Run("default keyword skips both", func(t *testing.T) {
		t.Parallel()

		dst := config.DefaultLogConfig()
		src := &config.LogConfig{Level: new(config.DefaultKeyword), FilePath: new(config.DefaultKeyword)}
		if errs := dst.Merge(src); len(errs) != 0 {
			t.Fatalf("unexpected errors: %v", errs)
		}
		if *dst.Level != config.LogError {
			t.Errorf("Level = %q, want unchanged %q", *dst.Level, config.LogError)
		}
	})

	t.Run("multiple errors collected", func(t *testing.T) {
		t.Parallel()

		dst := config.DefaultLogConfig()
		src := &config.LogConfig{Level: new("verbose"), FilePath: new("")}
		errs := dst.Merge(src)
		if len(errs) != 2 {
			t.Fatalf("want 2 errors, got %v", errs)
		}
	})
}

//nolint:tparallel // uses t.Setenv, which forbids parallel execution
func TestExpandPath(t *testing.T) {
	t.Run("tilde expands to home", func(t *testing.T) {
		home := t.TempDir()
		t.Setenv("HOME", home)

		got := config.ExpandPath("~/nm-tui/log")
		want := filepath.Join(home, "nm-tui", "log")
		if got != want {
			t.Errorf("expandPath(%q) = %q, want %q", "~/nm-tui/log", got, want)
		}
	})

	t.Run("bare tilde unchanged", func(t *testing.T) {
		t.Parallel()

		if got := config.ExpandPath("~"); got != "~" {
			t.Errorf("expandPath(%q) = %q, want %q", "~", got, "~")
		}
	})

	t.Run("env vars expanded", func(t *testing.T) {
		dir := t.TempDir()
		t.Setenv("NM_TUI_TEST_DIR", dir)

		got := config.ExpandPath("$NM_TUI_TEST_DIR/log")
		want := filepath.Join(dir, "log")
		if got != want {
			t.Errorf("expandPath($NM_TUI_TEST_DIR/log) = %q, want %q", got, want)
		}
	})

	t.Run("absolute path unchanged", func(t *testing.T) {
		t.Parallel()

		p := "/tmp/nm-tui.log"
		if got := config.ExpandPath(p); got != p {
			t.Errorf("expandPath(%q) = %q, want %q", p, got, p)
		}
	})
}

func TestResolveConfigPath(t *testing.T) {
	t.Run("with XDG_CONFIG_HOME", func(t *testing.T) {
		dir := t.TempDir()
		t.Setenv("XDG_CONFIG_HOME", dir)

		got, err := config.ResolveConfigPath()
		if err != nil {
			t.Fatalf("resolveConfigPath() error: %v", err)
		}
		want := filepath.Join(dir, config.AppName, config.ConfigFileName)
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("error when no config dir", func(t *testing.T) {
		t.Setenv("XDG_CONFIG_HOME", "")
		t.Setenv("HOME", "")

		if _, err := config.ResolveConfigPath(); err == nil {
			t.Fatal("resolveConfigPath() expected error")
		}
	})
}
