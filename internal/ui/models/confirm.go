package models

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/alphameo/nm-tui/internal/ui/styles"
)

type confirmKeyMap struct {
	accept  key.Binding
	decline key.Binding
}

type ConfirmModel struct {
	Question string
	Action   tea.Cmd

	help   help.Model
	styles lipgloss.Style
	keys   confirmKeyMap
}

func NewConfirmModel(style lipgloss.Style, keys confirmKeyMap) *ConfirmModel {
	help := help.New()
	help.Styles = styles.HelpStyles
	return &ConfirmModel{styles: style, keys: keys, help: help}
}

func (m *ConfirmModel) Init() tea.Cmd {
	return nil
}

func (m *ConfirmModel) Update(msg tea.Msg) (*ConfirmModel, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyPressMsg)
	if ok {
		switch {
		case key.Matches(keyMsg, m.keys.accept):
			return m, tea.Batch(m.Action, ClosePopupCmd())
		case key.Matches(keyMsg, m.keys.decline):
			return m, ClosePopupCmd()
		}
	}

	return m, nil
}

func (m *ConfirmModel) View() string {
	view := m.Question
	view = lipgloss.JoinVertical(
		lipgloss.Center,
		view,
		"",
		m.help.ShortHelpView([]key.Binding{m.keys.accept, m.keys.decline}),
	)
	return m.styles.Render(view)
}

type ConfirmMsg struct {
	question string
	cmd      tea.Cmd
}

func ConfirmCmd(question string, cmd tea.Cmd) tea.Cmd {
	return func() tea.Msg {
		return ConfirmMsg{
			question: question,
			cmd:      cmd,
		}
	}
}
