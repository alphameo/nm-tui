package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

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

	cacheDir, err := os.UserCacheDir()
	if err != nil {
		panic(fmt.Errorf("get cache directory: %w", err))
	} else {
		cacheDir = filepath.Join(cacheDir, "nm-tui")
	}
	if err := os.MkdirAll(cacheDir, 0o700); err != nil {
		panic(fmt.Errorf("create cache directory: %w", err))
	}
	logPath := filepath.Join(cacheDir, "log")
	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o600)
	if err != nil {
		panic(fmt.Errorf("open log file: %w", err))
	}

	defer func() {
		_ = f.Close()
	}()

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
