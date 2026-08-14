package config_test

import (
	"strings"
	"testing"

	"github.com/alphameo/nm-tui/internal/config"
)

func TestDefaultConfig(t *testing.T) {
	assertNoNilFields(t, config.DefaultConfig())
}

func TestConfigMerge(t *testing.T) {
	t.Run("empty source produces no errors", func(t *testing.T) {
		cfg := config.DefaultConfig()
		if errs := cfg.Merge(&config.Config{}); len(errs) != 0 {
			t.Errorf("unexpected errors: %v", errs)
		}
	})

	t.Run("nil child sections are handled", func(t *testing.T) {
		cfg := config.DefaultConfig()
		src := config.Config{
			Logging: &config.LogConfig{},
			Colors:  &config.ColorConfig{},
			Keys:    &config.KeyConfig{},
			Icons:   &config.IconConfig{},
		}
		if errs := cfg.Merge(&src); len(errs) != 0 {
			t.Errorf("unexpected errors: %v", errs)
		}
	})

	t.Run("notification close time valid", func(t *testing.T) {
		cfg := config.DefaultConfig()
		src := config.Config{NotifCloseTime: new(250)}
		if errs := cfg.Merge(&src); len(errs) != 0 {
			t.Fatalf("unexpected errors: %v", errs)
		}
		if got, want := *cfg.NotifCloseTime, 250; got != want {
			t.Errorf("NotifCloseTime = %d, want %d", got, want)
		}
	})

	t.Run("notification close time invalid", func(t *testing.T) {
		cfg := config.DefaultConfig()
		src := config.Config{NotifCloseTime: new(-1)}
		errs := cfg.Merge(&src)
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
		cfg := config.DefaultConfig()
		src := config.Config{Icons: &config.IconConfig{
			NerdPreset:  new(true),
			BorderStyle: new(config.BorderSquare),
		}}
		if errs := cfg.Merge(&src); len(errs) != 0 {
			t.Fatalf("unexpected errors: %v", errs)
		}
		if got, want := *cfg.Icons.BorderStyle, config.BorderSquare; got != want {
			t.Errorf("BorderStyle = %q, want %q", got, want)
		}
		if got, want := *cfg.Icons.SpinnerStyle, config.SpinnerMeter; got != want {
			t.Errorf("SpinnerStyle = %q, want nerd default %q", got, want)
		}
		if got, want := *cfg.Icons.ToggleOff, " "; got != want {
			t.Errorf("ToggleOff = %q, want nerd default %q", got, want)
		}
	})

	t.Run("icon merge errors are propagated", func(t *testing.T) {
		cfg := config.DefaultConfig()
		src := config.Config{Icons: &config.IconConfig{BorderStyle: new("bogus")}}
		errs := cfg.Merge(&src)
		if len(errs) != 1 {
			t.Fatalf("want 1 error, got %v", errs)
		}
		if !strings.Contains(errs[0].Error(), "border_style") {
			t.Errorf("unexpected error: %v", errs[0])
		}
	})

	t.Run("logging merge errors are propagated", func(t *testing.T) {
		cfg := config.DefaultConfig()
		src := config.Config{Logging: &config.LogConfig{Level: new("verbose")}}
		errs := cfg.Merge(&src)
		if len(errs) != 1 {
			t.Fatalf("want 1 error, got %v", errs)
		}
		if !strings.Contains(errs[0].Error(), "invalid log level") {
			t.Errorf("unexpected error: %v", errs[0])
		}
	})

	t.Run("colors merge errors are propagated", func(t *testing.T) {
		cfg := config.DefaultConfig()
		src := config.Config{Colors: &config.ColorConfig{Text: new("notacolor")}}
		errs := cfg.Merge(&src)
		if len(errs) != 1 {
			t.Fatalf("want 1 error, got %v", errs)
		}
		if !strings.Contains(errs[0].Error(), "text color") {
			t.Errorf("unexpected error: %v", errs[0])
		}
	})

	t.Run("keys merge errors are propagated", func(t *testing.T) {
		cfg := config.DefaultConfig()
		src := config.Config{Keys: &config.KeyConfig{Toggle: &config.KeyBinding{"notakey"}}}
		errs := cfg.Merge(&src)
		if len(errs) != 1 {
			t.Fatalf("want 1 error, got %v", errs)
		}
		if !strings.Contains(errs[0].Error(), "invalid key toggle") {
			t.Errorf("unexpected error: %v", errs[0])
		}
	})
}
