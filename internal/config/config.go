package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/calico32/kdl-go"
)

func init() {
	initKeyConfig()
}

const (
	appName        = "nm-tui"
	configFileName = "config.kdl"
	defaultKeyword = "default"
)

type Config struct {
	Colors  *ColorConfig `kdl:"colors"`
	Keys    *KeyConfig   `kdl:"keys"`
	Logging *LogConfig   `kdl:"logging"`
	Icons   *IconConfig  `kdl:"icons"`
}

func DefaultConfig() Config {
	return Config{
		Colors:  DefaultColorConfig(),
		Keys:    DefaultKeys(),
		Logging: DefaultLogConfig(),
		Icons:   DefaultIconConfig(),
	}
}

func (c *Config) merge(src *Config) []error {
	var errs []error

	if src.Logging != nil {
		errs = append(errs, c.Logging.merge(src.Logging)...)
	}

	if src.Colors != nil {
		errs = append(errs, c.Colors.merge(src.Colors)...)
	}

	if src.Keys != nil {
		errs = append(errs, c.Keys.merge(src.Keys)...)
	}

	if src.Icons != nil {
		nerd := src.Icons.NerdPreset
		if nerd != nil && *nerd {
			c.Icons = DefaultNerdIconConfig()
		}
		errs = append(errs, c.Icons.merge(src.Icons)...)
	}

	return errs
}

func Load() (*Config, error) {
	path, err := resolveConfigPath()
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
	if err := kdl.Decode(f, &cfg); err != nil {
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
	errs := cfg.merge(userCfg)

	if len(errs) > 0 {
		err = fmt.Errorf("user config: %w", errors.Join(errs...))
	} else {
		err = nil
	}
	return cfg, err
}
