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

func (c *ColorConfig) tryMerge(src *ColorConfig) []error {
	var errs []error
	v, err := resolveColor(src.Text)
	if err != nil {
		err := fmt.Errorf("text color: %w", err)
		errs = append(errs, err)
	} else {
		c.Text = v
	}

	v, err = resolveColor(src.Accent)
	if err != nil {
		err := fmt.Errorf("accent color: %w", err)
		errs = append(errs, err)
	} else {
		c.Accent = v
	}

	v, err = resolveColor(src.Error)
	if err != nil {
		err := fmt.Errorf("error color: %w", err)
		errs = append(errs, err)
	} else {
		c.Error = v
	}

	v, err = resolveColor(src.Muted)
	if err != nil {
		err := fmt.Errorf("muted color: %w", err)
		errs = append(errs, err)
	} else {
		c.Muted = v
	}

	v, err = resolveColor(src.Notif)
	if err != nil {
		err := fmt.Errorf("notif color: %w", err)
		errs = append(errs, err)
	} else {
		c.Notif = v
	}
	return errs
}

func (c *LogConfig) tryMerge(src *LogConfig) []error {
	var errs []error
	if src.FilePath == "" {
		err := fmt.Errorf("invalid log filepath: %s", src.FilePath)
		errs = append(errs, err)
	} else {
		c.FilePath = src.FilePath
	}

	if !validLogLevel(src.Level) {
		err := fmt.Errorf("invalid log level: %s", src.Level)
		errs = append(errs, err)
	} else {
		c.Level = src.Level
	}
	return errs
}

func (c *Config) tryMerge(src *Config) []error {
	var errs []error

	logErrs := c.Logging.tryMerge(&src.Logging)
	errs = append(errs, logErrs...)

	colorErrs := c.Colors.tryMerge(&src.Colors)
	errs = append(errs, colorErrs...)

	return errs
}

func validLogLevel(s string) bool {
	switch strings.ToLower(s) {
	case "debug", "info", "warn", "error":
		return true
	}
	return false
}

func validHex(color string) bool {
	if len(color) != 7 || color[0] != '#' {
		return false
	}
	_, err := strconv.ParseUint(color[1:], 16, 64)
	return err == nil
}
func resolveColor(color string) (string, error) {
	if validHex(color) {
		return color, nil
	}
	return "", fmt.Errorf("color not resolved: %s", color)
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
	cfg := DefaultConfig()
	userCfg, err := Load()
	if err != nil {
		slog.Warn("user config loading failed", "err", err.Error())
		return cfg
	}
	errs := cfg.tryMerge(userCfg)
	for _, err := range errs {
		slog.Warn("error inside user config, fallback to default values", "error", err.Error())
	}
	return cfg
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
