package tabview

import tea "charm.land/bubbletea/v2"

type TabModel interface {
	Init() tea.Cmd

	// TabModel updates content.
	// Just wrap Update(msg) method
	//
	// Example:
	// func (m Model) TabUpdate(msg tea.Msg) (tabview.TabModel, tea.Cmd) {
	//	return m.Update(msg)
	//}
	UpdateAsTab(msg tea.Msg) (TabModel, tea.Cmd)

	View() string

	Resize(width, height int)
	Width() int
	Height() int
}
