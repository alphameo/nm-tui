// Package text provides simple text component
package text

import tea "github.com/charmbracelet/bubbletea"

type TextModel string

func New(label string) TextModel {
	return TextModel(label)
}

func (m TextModel) Init() tea.Cmd {
	return nil
}

func (m TextModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m TextModel) View() string {
	return string(m)
}
