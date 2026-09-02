package run

import (
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/alphameo/nm-tui/internal/config"
	"github.com/alphameo/nm-tui/internal/infra/logging"
	"github.com/alphameo/nm-tui/internal/infra/nm"
	"github.com/alphameo/nm-tui/internal/infra/portal"
	"github.com/alphameo/nm-tui/internal/ui/models"
)

func Run(version string) {
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Fprintf(os.Stdout, "nm-tui %s\n", version)
		os.Exit(0)
	}

	stdLogger := slog.New(slog.NewJSONHandler(
		os.Stderr,
		&slog.HandlerOptions{Level: slog.LevelWarn},
	))
	slog.SetDefault(stdLogger)

	cfg, cfgErr := config.LoadOrDefaults()
	if cfgErr != nil && !errors.Is(cfgErr, fs.ErrNotExist) {
		stdLogger.Warn("errors in user config, falling back to defaults", "errors", cfgErr)
	}

	logPath := *cfg.Logging.FilePath
	logPathDir := filepath.Dir(logPath)
	if err := os.MkdirAll(logPathDir, 0o700); err != nil {
		err = fmt.Errorf("create log directory: %w", err)
		stdLogger.Error(err.Error())
		panic(err)
	}
	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		err = fmt.Errorf("open log file: %w", err)
		stdLogger.Error(err.Error())
		panic(err)
	}
	defer func() {
		_ = f.Close()
	}()

	logLevel, err := resolveLogLevel(*cfg.Logging.Level)
	if err != nil {
		err = fmt.Errorf("log level: %w", err)
		stdLogger.Error(err.Error())
		panic(err)
	}
	fileLogger := slog.New(slog.NewJSONHandler(f, &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: logLevel <= slog.LevelDebug,
	}))
	slog.SetDefault(fileLogger)
	if cfgErr != nil {
		fileLogger.Warn("errors in user config, falling back to defaults", "errors", cfgErr)
	}

	fileLogger.Info("The program is running")
	defer fileLogger.Info("Program is closed")

	nm := nm.NewCLI()
	portalOpener := portal.New()
	networksMw := logging.NewNetworks(fileLogger, nm)
	deviceMw := logging.NewDevice(fileLogger, nm)
	portalMw := logging.NewPortal(fileLogger, portalOpener)
	model, err := models.NewMainModel(networksMw, deviceMw, portalMw, cfg)
	if err != nil {
		fileLogger.Error("error during model initialization", "errors", err.Error())
		return
	}

	p := tea.NewProgram(model)
	if _, err = p.Run(); err != nil {
		fileLogger.Error("runtime error", "error", err.Error())
		return
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
