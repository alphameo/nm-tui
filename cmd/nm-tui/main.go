package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/alphameo/nm-tui/internal/infra"
	"github.com/alphameo/nm-tui/internal/logger"
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

	logger.FilePath(logPath)
	logger.Level = logger.ErrorsLvl
	logger.Informln("The program is running")
	defer logger.Informln("Program is closed")

	nm := infra.NewNMCLI()
	p := tea.NewProgram(components.NewMainModel(nm), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		logger.Errln(err)
	}
}
