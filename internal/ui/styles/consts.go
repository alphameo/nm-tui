package styles

import "github.com/charmbracelet/lipgloss"

var (
	BorderOffset int = lipgloss.Width(Border.Left) * 2
	TabBarHeight int = BorderOffset + 1
)
