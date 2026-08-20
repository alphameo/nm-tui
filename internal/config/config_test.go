package config_test

import (
	"testing"

	"github.com/alphameo/nm-tui/internal/config"
)

func TestDefaultConfig(t *testing.T) {
	t.Parallel()

	assertNoNilFields(t, config.DefaultConfig())
}

func TestConfigMerge(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		src       *config.Config
		wantErr   int
		fragments []string
		check     func(t *testing.T, cfg *config.Config)
	}{
		{name: "empty source produces no errors", src: &config.Config{}},
		{name: "nil child sections are handled", src: &config.Config{
			Logging: &config.LogConfig{},
			Colors:  &config.ColorConfig{},
			Keys:    &config.KeyConfig{},
			Icons:   &config.IconConfig{},
		}},
		{
			name: "notification close time valid",
			src:  &config.Config{NotifCloseTime: new(250)},
			check: func(t *testing.T, cfg *config.Config) {
				if got, want := *cfg.NotifCloseTime, 250; got != want {
					t.Errorf("NotifCloseTime = %d, want %d", got, want)
				}
			},
		},
		{
			name:      "notification close time invalid",
			src:       &config.Config{NotifCloseTime: new(-1)},
			wantErr:   1,
			fragments: []string{"notification_close_time"},
			check: func(t *testing.T, cfg *config.Config) {
				if got, want := *cfg.NotifCloseTime, *config.DefaultConfig().NotifCloseTime; got != want {
					t.Errorf("NotifCloseTime should stay default, got %d want %d", got, want)
				}
			},
		},
		{
			name: "nerd preset swaps icons to nerd defaults",
			src: &config.Config{Icons: &config.IconConfig{
				NerdPreset:  new(true),
				BorderStyle: new(config.BorderSquare),
			}},
			check: func(t *testing.T, cfg *config.Config) {
				if got, want := *cfg.Icons.BorderStyle, config.BorderSquare; got != want {
					t.Errorf("BorderStyle = %q, want %q", got, want)
				}
				if got, want := *cfg.Icons.SpinnerStyle, config.SpinnerMeter; got != want {
					t.Errorf("SpinnerStyle = %q, want nerd default %q", got, want)
				}
				if got, want := *cfg.Icons.ToggleOff, " "; got != want {
					t.Errorf("ToggleOff = %q, want nerd default %q", got, want)
				}
			},
		},
		{
			name:      "icon merge errors are propagated",
			src:       &config.Config{Icons: &config.IconConfig{BorderStyle: new("bogus")}},
			wantErr:   1,
			fragments: []string{"border_style"},
		},
		{
			name:      "logging merge errors are propagated",
			src:       &config.Config{Logging: &config.LogConfig{Level: new("verbose")}},
			wantErr:   1,
			fragments: []string{"invalid log level"},
		},
		{
			name:      "colors merge errors are propagated",
			src:       &config.Config{Colors: &config.ColorConfig{Text: new("notacolor")}},
			wantErr:   1,
			fragments: []string{"text color"},
		},
		{
			name:      "keys merge errors are propagated",
			src:       &config.Config{Keys: &config.KeyConfig{Toggle: &config.KeyBinding{"notakey"}}},
			wantErr:   1,
			fragments: []string{"invalid key toggle"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cfg := config.DefaultConfig()
			errs := cfg.Merge(tt.src)
			if len(errs) != tt.wantErr {
				t.Fatalf("want %d errors, got %v", tt.wantErr, errs)
			}
			assertErrsContain(t, errs, tt.fragments...)
			if tt.check != nil {
				tt.check(t, &cfg)
			}
		})
	}
}
