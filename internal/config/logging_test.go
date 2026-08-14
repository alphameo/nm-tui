package config_test

import (
	"path/filepath"
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

//nolint:tparallel // uses t.Setenv, which forbids parallel execution
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

	tests := []struct {
		name      string
		src       *config.LogConfig
		wantErr   int
		fragments []string
		keepLevel bool
	}{
		{name: "nil source no-op", src: nil},
		{name: "empty source no-op", src: &config.LogConfig{}},
		{
			name:      "default keyword skips both",
			src:       &config.LogConfig{Level: new(config.DefaultKeyword), FilePath: new(config.DefaultKeyword)},
			keepLevel: true,
		},
		{
			name:      "invalid level errors and keeps default",
			src:       &config.LogConfig{Level: new("verbose")},
			wantErr:   1,
			fragments: []string{`invalid log level: "verbose"`},
			keepLevel: true,
		},
		{
			name:      "empty file path errors and keeps default",
			src:       &config.LogConfig{FilePath: new("")},
			wantErr:   1,
			fragments: []string{"empty log filepath"},
		},
		{
			name:    "multiple errors collected",
			src:     &config.LogConfig{Level: new("verbose"), FilePath: new("")},
			wantErr: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			dst := config.DefaultLogConfig()
			errs := dst.Merge(tt.src)
			if len(errs) != tt.wantErr {
				t.Fatalf("want %d errors, got %v", tt.wantErr, errs)
			}
			assertErrsContain(t, errs, tt.fragments...)
			if tt.keepLevel && *dst.Level != config.LogError {
				t.Errorf("Level = %q, want unchanged %q", *dst.Level, config.LogError)
			}
		})
	}
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
