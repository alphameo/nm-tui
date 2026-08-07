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
		if cfg.FilePath == nil {
			t.Error("FilePath is nil")
		}
		if cfg.Level == nil {
			t.Error("Level is nil")
		}
	})

	t.Run("falls back to HOME", func(t *testing.T) {
		t.Setenv("XDG_STATE_HOME", "")
		home := t.TempDir()
		t.Setenv("HOME", home)

		cfg := DefaultLogConfig()
		if cfg.FilePath == nil {
			t.Error("FilePath is nil")
		}
		if cfg.Level == nil {
			t.Error("Level is nil")
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
