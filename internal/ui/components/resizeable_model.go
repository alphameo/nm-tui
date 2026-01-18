package components

import tea "github.com/charmbracelet/bubbletea"

type SizedModel interface {
	tea.Model
	Resize(width, height int)
	Width() int
	Height() int
}
