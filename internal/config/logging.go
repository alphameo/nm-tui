package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type LogConfig struct {
	Level    *string `kdl:"level"`
	FilePath *string `kdl:"file_path"`
}

func DefaultLogConfig() *LogConfig {
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

func (c *LogConfig) merge(src *LogConfig) []error {
	if src == nil {
		return nil
	}

	var errs []error

	if src.FilePath != nil && *src.FilePath != defaultKeyword {
		if *src.FilePath == "" {
			errs = append(errs, fmt.Errorf("empty log filepath"))
		} else {
			c.FilePath = src.FilePath
		}
	}

	if src.Level != nil && *src.Level != defaultKeyword {
		if !validLogLevel(*src.Level) {
			errs = append(errs, fmt.Errorf("invalid log level: %q", *src.Level))
		} else {
			c.Level = src.Level
		}
	}
	return errs
}

const (
	LogDebug = "debug"
	LogInfo  = "info"
	LogWarn  = "warn"
	LogError = "error"
)

func validLogLevel(level string) bool {
	switch strings.ToLower(level) {
	case LogDebug, LogInfo, LogWarn, LogError:
		return true
	}
	return false
}

func resolveConfigPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	path := filepath.Join(configDir, appName, configFileName)
	return path, nil
}
