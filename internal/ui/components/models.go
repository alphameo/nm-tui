package components

import tea "charm.land/bubbletea/v2"

type Focusable interface {
	Focused() bool
	Focus() tea.Cmd
	Blur()
}
