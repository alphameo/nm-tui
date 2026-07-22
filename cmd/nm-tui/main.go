package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/alphameo/nm-tui/internal/config"
	"github.com/alphameo/nm-tui/internal/infra/nmcli"
	"github.com/alphameo/nm-tui/internal/ui/models"
	"github.com/alphameo/nm-tui/internal/ui/styles"
)

func main() {
	loggerOpts := &slog.HandlerOptions{
		Level:     slog.LevelWarn,
		AddSource: true,
	}
	stdLogger := slog.New(slog.NewJSONHandler(os.Stderr, loggerOpts))
	slog.SetDefault(stdLogger)

	cfg, cfgErr := config.LoadOrDefaults()
	if cfgErr != nil {
		slog.Warn("errors in user config, falling back to defaults", "errors", cfgErr)
	}

	logPath := expandPath(*cfg.Logging.FilePath)
	logPathDir := filepath.Dir(logPath)
	if err := os.MkdirAll(logPathDir, 0o700); err != nil {
		err := fmt.Errorf("create log directory: %w", err)
		slog.Error(err.Error())
		panic(err)
	}
	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o600)
	if err != nil {
		err := fmt.Errorf("open log file: %w", err)
		slog.Error(err.Error())
		panic(err)
	}
	defer func() {
		_ = f.Close()
	}()

	logLevel, err := resolveLogLevel(*cfg.Logging.Level)
	if err != nil {
		err = fmt.Errorf("log level: %w", err)
		slog.Error(err.Error())
		panic(err)
	}
	loggerOpts.Level = logLevel
	fileLogger := slog.New(slog.NewJSONHandler(f, loggerOpts))

	slog.SetDefault(fileLogger)
	if cfgErr != nil {
		slog.Warn("errors in user config, falling back to defaults", "errors", cfgErr)
	}

	slog.Info("Style initialization")
	err = styles.Init(*cfg.Colors)
	if err != nil {
		slog.Error("errors during style initialization", "error", err.Error())
	}

	slog.Info("The program is running")
	defer slog.Info("Program is closed")

	nm := nmcli.New()
	p := tea.NewProgram(models.NewMainModel(nm, nm))
	if _, err := p.Run(); err != nil {
		slog.Error(err.Error())
	}
}

func resolveLogLevel(level string) (slog.Level, error) {
	logLevel := strings.ToLower(level)
	switch logLevel {
	case config.LogDebug:
		return slog.LevelDebug, nil
	case config.LogInfo:
		return slog.LevelInfo, nil
	case config.LogWarn:
		return slog.LevelWarn, nil
	case config.LogError:
		return slog.LevelError, nil
	}

	return slog.LevelError, fmt.Errorf("log level not resolved: %s", logLevel)
}

func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, _ := os.UserHomeDir()
		path = filepath.Join(home, path[2:])
	}
	return os.ExpandEnv(path)
}
