package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/calico32/kdl-go"
)

type Config struct {
	Colors ColorConfig `kdl:"colors"`
	Keys   Keys        `kdl:"keys"`
	Paths  PathConfig  `kdl:"paths"`
}

type ColorConfig struct {
	Text   string `kdl:"text,argument"`
	Accent string `kdl:"accent,argument"`
	Muted  string `kdl:"muted,argument"`
	Error  string `kdl:"error,argument"`
	Notif  string `kdl:"notif,argument"`
}

type PathConfig struct {
	CacheDir string `kdl:"cache_dir,argument"`
	LogFile  string `kdl:"log_file,argument"`
}

func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open config: %w", err)
	}
	defer f.Close()

	var cfg Config
	if err := kdl.Decode(f, &cfg); err != nil {
		return nil, fmt.Errorf("decode config: %w", err)
	}
	return &cfg, nil
}

func LoadOrDefaults() *Config {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil
	}
	cfg, err := Load(filepath.Join(configDir, "nm-tui", "config.kdl"))
	if err != nil {
		return nil
	}
	return cfg
}

func HelpFromKeys(keys []string) string {
	transformed := make([]string, len(keys))
	for i, key := range keys {
		transformed[i] = strings.ReplaceAll(key, "ctrl+", "^")
	}
	return strings.Join(transformed, "/")
}
