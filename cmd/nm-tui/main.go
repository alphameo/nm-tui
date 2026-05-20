package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	tea "charm.land/bubbletea/v2"
	"github.com/alphameo/nm-tui/internal/infra/nmcli"
	"github.com/alphameo/nm-tui/internal/ui/components"
)

func main() {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		panic(fmt.Errorf("failed to get cache directory: %w", err))
	} else {
		cacheDir = filepath.Join(cacheDir, "nm-tui")
	}

	if err := os.MkdirAll(cacheDir, 0700); err != nil {
		panic(fmt.Errorf("failed to create cache directory: %w", err))
	}

	logPath := filepath.Join(cacheDir, "log")
	_, err = os.Stat(logPath)
	if errors.Is(err, os.ErrNotExist) {
		_, err = os.Create(logPath)
		fmt.Println(err)
		if err != nil {
			os.Exit(1)
		}
	}
	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o600)
	if err != nil {
		panic(err)
	}

	opts := &slog.HandlerOptions{
		Level:     slog.LevelError,
		AddSource: true,
	}
	logger := slog.New(slog.NewJSONHandler(f, opts))

	slog.SetDefault(logger)
	slog.Info("The program is running")
	defer slog.Info("Program is closed")

	nm := nmcli.New()
	p := tea.NewProgram(components.NewMainModel(nm, nm))
	if _, err := p.Run(); err != nil {
		slog.Error(err.Error())
	}
}
