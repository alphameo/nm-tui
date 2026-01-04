package main

import (
	"github.com/alphameo/nm-tui/internal/infra"
	"github.com/alphameo/nm-tui/internal/logger"
	"github.com/alphameo/nm-tui/internal/ui/views"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	logger.FilePath("./log")
	logger.Level = logger.ErrorsLvl
	logger.Informln("The program is running")
	defer logger.Informln("Program is closed")
	p := tea.NewProgram(views.NewMainModel(infra.NewNMCLI()), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		logger.Errln(err)
	}
}
