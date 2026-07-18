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
		Level:     slog.LevelError,
		AddSource: true,
	}
	stdLogger := slog.New(slog.NewJSONHandler(os.Stderr, loggerOpts))
	slog.SetDefault(stdLogger)

	cfg := config.LoadOrDefaults()
	styles.Init(cfg.Colors)

	logPathDir := filepath.Dir(cfg.Logging.FilePath)
	if err := os.MkdirAll(logPathDir, 0o700); err != nil {
		panic(fmt.Errorf("create log directory: %w", err))
	}
	f, err := os.OpenFile(cfg.Logging.FilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o600)
	if err != nil {
		panic(fmt.Errorf("open log file: %w", err))
	}

	defer func() {
		_ = f.Close()
	}()

	loggerOpts.Level = ResolveLogLevel(cfg.Logging.Level)
	fileLogger := slog.New(slog.NewJSONHandler(f, loggerOpts))

	slog.SetDefault(fileLogger)
	slog.Info("The program is running")
	defer slog.Info("Program is closed")

	nm := nmcli.New()
	p := tea.NewProgram(models.NewMainModel(nm, nm))
	if _, err := p.Run(); err != nil {
		slog.Error(err.Error())
	}
}

func ResolveLogLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelError
	}
}
