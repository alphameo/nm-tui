package config

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/calico32/kdl-go"
)

const (
	appName        = "nm-tui"
	configFileName = "config.kdl"
)

type Config struct {
	Colors  ColorConfig `kdl:"colors"`
	Keys    Keys        `kdl:"keys"`
	Logging LogConfig   `kdl:"logging"`
}

type ColorConfig struct {
	Text   string `kdl:"text"`
	Accent string `kdl:"accent"`
	Muted  string `kdl:"muted"`
	Error  string `kdl:"error"`
	Notif  string `kdl:"notif"`
}

type LogConfig struct {
	Level    string `kdl:"level"`
	FilePath string `kdl:"file_path"`
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

func DefaultLogConfig() LogConfig {
	stateDir := os.Getenv("XDG_STATE_HOME")
	if stateDir == "" {
		home, _ := os.UserHomeDir()
		stateDir = filepath.Join(home, ".local", "state")
	}
	logPath := filepath.Join(stateDir, appName, "log")
	return LogConfig{
		Level:    "error",
		FilePath: logPath,
	}
}

func DefaultConfig() Config {
	return Config{
		Colors:  DefaultColorConfig(),
		Keys:    DefaultKeys(),
		Logging: DefaultLogConfig(),
	}
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

	cfg.Logging.FilePath = expandPath(cfg.Logging.FilePath)
	return &cfg, nil
}

func LoadOrDefaults() Config {
	cfg, err := Load()
	if err != nil {
		slog.Error("config loading failed", "err", err.Error())
		return DefaultConfig()
	}
	return *cfg
}

func resolveConfigPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	path := filepath.Join(configDir, appName, configFileName)
	return path, nil
}

func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, _ := os.UserHomeDir()
		path = filepath.Join(home, path[2:])
	}
	return os.ExpandEnv(path)
}
