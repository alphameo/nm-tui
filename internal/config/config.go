package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/calico32/kdl-go"
)

const (
	AppName        = "nm-tui"
	ConfigFileName = "config.kdl"
	DefaultKeyword = "default"
)

type Config struct {
	Colors         *ColorConfig `kdl:"colors"`
	Keys           *KeyConfig   `kdl:"keys"`
	Logging        *LogConfig   `kdl:"logging"`
	Icons          *IconConfig  `kdl:"icons"`
	NotifCloseTime *int         `kdl:"notification_close_time"`
	RescanInterval *int         `kdl:"rescan_interval"`
}

func DefaultConfig() Config {
	return Config{
		Colors:         DefaultColorConfig(),
		Keys:           DefaultKeys(),
		Logging:        DefaultLogConfig(),
		Icons:          DefaultIconConfig(),
		NotifCloseTime: new(5),
		RescanInterval: new(10),
	}
}

func (c *Config) Merge(src *Config) []error {
	var errs []error

	if src.Logging != nil {
		errs = append(errs, c.Logging.Merge(src.Logging)...)
	}

	if src.Colors != nil {
		errs = append(errs, c.Colors.Merge(src.Colors)...)
	}

	if src.Keys != nil {
		errs = append(errs, c.Keys.Merge(src.Keys)...)
	}

	if src.Icons != nil {
		nerd := src.Icons.NerdPreset
		if nerd != nil && *nerd {
			c.Icons = DefaultNerdIconConfig()
		}
		errs = append(errs, c.Icons.Merge(src.Icons)...)
	}

	if src.NotifCloseTime != nil {
		time := *src.NotifCloseTime
		err := validatePositiveTime(time)
		if err != nil {
			err = fmt.Errorf("notification_close_time value: %w", err)
			errs = append(errs, err)
		} else {
			c.NotifCloseTime = src.NotifCloseTime
		}
	}

	if src.RescanInterval != nil {
		interval := *src.RescanInterval
		err := validateNonNegativeTime(interval)
		if err != nil {
			err = fmt.Errorf("rescan_interval value: %w", err)
			errs = append(errs, err)
		} else {
			c.RescanInterval = src.RescanInterval
		}
	}
	return errs
}

func Load() (*Config, error) {
	path, err := ResolveConfigPath()
	if err != nil {
		return nil, fmt.Errorf("resolve config path: %w", err)
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open config: %w", err)
	}
	defer func() {
		_ = f.Close()
	}()

	var cfg Config
	if err = kdl.Decode(f, &cfg); err != nil {
		return nil, fmt.Errorf("decode config: %w", err)
	}

	return &cfg, nil
}

func LoadOrDefaults() (Config, error) {
	cfg := DefaultConfig()
	userCfg, err := Load()
	if err != nil {
		return cfg, fmt.Errorf("user config loading failed: %w", err)
	}
	errs := cfg.Merge(userCfg)

	if len(errs) > 0 {
		err = fmt.Errorf("user config: %w", errors.Join(errs...))
	} else {
		err = nil
	}
	return cfg, err
}

func validatePositiveTime(time int) error {
	if time <= 0 {
		return fmt.Errorf("time <= 0: %d", time)
	}
	return nil
}

func validateNonNegativeTime(time int) error {
	if time < 0 {
		return fmt.Errorf("time < 0: %d", time)
	}
	return nil
}
