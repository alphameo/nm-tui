package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
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
		Text:   "none",
		Accent: "blue",
		Muted:  "bright_black",
		Error:  "red",
		Notif:  "yellow",
	}
}

func DefaultLogConfig() LogConfig {
	stateDir := os.Getenv("XDG_STATE_HOME")
	if stateDir == "" {
		home, _ := os.UserHomeDir()
		stateDir = filepath.Join(home, ".local", "state")
	}
	logPath := filepath.Join(stateDir, appName, "log")
	level := LogError
	return &LogConfig{
		Level:    &level,
		FilePath: &logPath,
	}
}

func DefaultConfig() Config {
	return Config{
		Colors:  DefaultColorConfig(),
		Keys:    DefaultKeys(),
		Logging: DefaultLogConfig(),
	}
}

func (c *ColorConfig) merge(src *ColorConfig) []error {
	var errs []error
	err := validateColor(src.Text)
	if err != nil {
		err := fmt.Errorf("text color: %w", err)
		errs = append(errs, err)
	} else {
		c.Text = src.Text
	}

	err = validateColor(src.Accent)
	if err != nil {
		err := fmt.Errorf("accent color: %w", err)
		errs = append(errs, err)
	} else {
		c.Accent = src.Accent
	}

	err = validateColor(src.Error)
	if err != nil {
		err := fmt.Errorf("error color: %w", err)
		errs = append(errs, err)
	} else {
		c.Error = src.Error
	}

	err = validateColor(src.Muted)
	if err != nil {
		err := fmt.Errorf("muted color: %w", err)
		errs = append(errs, err)
	} else {
		c.Muted = src.Muted
	}

	err = validateColor(src.Notif)
	if err != nil {
		err := fmt.Errorf("notif color: %w", err)
		errs = append(errs, err)
	} else {
		c.Notif = src.Notif
	}
	return errs
}

func (c *LogConfig) merge(src *LogConfig) []error {
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

func (c *Config) merge(src *Config) []error {
	var errs []error

	logErrs := c.Logging.merge(&src.Logging)
	errs = append(errs, logErrs...)

	colorErrs := c.Colors.merge(&src.Colors)
	errs = append(errs, colorErrs...)

	return errs
}

const (
	LogDebug = "debug"
	LogInfo  = "info"
	LogWarn  = "warn"
	LogError = "error"
)

func validLogLevel(s string) bool {
	switch strings.ToLower(s) {
	case LogDebug, LogInfo, LogWarn, LogError:
		return true
	}
	return false
}

const (
	CBlack         = "black"
	CRed           = "red"
	CGreen         = "green"
	CYellow        = "yellow"
	CBlue          = "blue"
	CMagenta       = "magenta"
	CCyan          = "cyan"
	CWhite         = "white"
	CBrightBlack   = "bright_black"
	CBrightRed     = "bright_red"
	CBrightGreen   = "bright_green"
	CBrightYellow  = "bright_yellow"
	CBrightBlue    = "bright_blue"
	CBrightMagenta = "bright_magenta"
	CBrightCyan    = "bright_cyan"
	CBrightWhite   = "bright_white"
	CNone          = "none"
)

func ValidCfgColor(color string) bool {
	switch color {
	case CBlack, CRed, CGreen, CYellow, CBlue, CMagenta, CCyan, CWhite,
		CBrightBlack, CBrightRed, CBrightGreen, CBrightYellow,
		CBrightBlue, CBrightMagenta, CBrightCyan, CBrightWhite,
		CNone:
		return true
	default:
		return false
	}
}

func ValidHex(color string) bool {
	if len(color) != 7 || color[0] != '#' {
		return false
	}
	_, err := strconv.ParseUint(color[1:], 16, 64)
	return err == nil
}

func validateColor(color string) error {
	if ValidHex(color) {
		return nil
	}
	if ValidCfgColor(color) {
		return nil
	}
	return fmt.Errorf("unknown color: %s", color)
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

func resolveConfigPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	path := filepath.Join(configDir, appName, configFileName)
	return path, nil
}
