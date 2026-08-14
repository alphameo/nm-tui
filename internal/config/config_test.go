package config

import (
	"strings"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	assertNoNilFields(t, DefaultConfig())
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
