package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/alphameo/nm-tui/internal/infra"
	"github.com/alphameo/nm-tui/internal/ui/components"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	logPath := os.ExpandEnv("$HOME") + "/.cache/nm-tui/log"
	_, err := os.Stat(logPath)
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
	slog.Debug("yappi")

	slog.Info("The program is running")
	slog.SetDefault(logger)
	defer slog.Info("Program is closed")

	nm := infra.NewNMCLI()
	p := tea.NewProgram(components.NewMainModel(nm), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		slog.Error(err.Error())
	}
}
