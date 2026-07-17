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
	Text   string `kdl:"text"`
	Accent string `kdl:"accent"`
	Muted  string `kdl:"muted"`
	Error  string `kdl:"error"`
	Notif  string `kdl:"notif"`
}

func DefaultColorConfig() ColorConfig {
	return ColorConfig{
		Text:   "#cbcbcb",
		Accent: "#865fff",
		Muted:  "#595959",
		Error:  "#ff0000",
		Notif:  "#e4bf7a",
	}
}

type PathConfig struct {
	CacheDir string `kdl:"cache_dir"`
	LogFile  string `kdl:"log_file"`
}

const (
	configDirName = "nm-tui"
	configName    = "config.kdl"
)

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

func LoadOrDefaults() *Config {
	cfg, err := Load()
	if err != nil {
		return nil
	}
	return cfg
}

func resolveConfigPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	path := filepath.Join(configDir, configDirName, configName)
	return path, nil
}

func HelpFromKeys(keys []string) string {
	transformed := make([]string, len(keys))
	for i, key := range keys {
		transformed[i] = strings.ReplaceAll(key, "ctrl+", "^")
	}
	return strings.Join(transformed, "/")
}
