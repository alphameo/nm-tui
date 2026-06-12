package models

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

type PopupModel interface {
	Init() tea.Cmd
	UpdateAsPopup(msg tea.Msg) (PopupModel, tea.Cmd)
	View() string
}

type Popup struct {
	content PopupModel
	active  bool
}

type popupKeyMap struct {
	close key.Binding
}

func (k *popupKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.close}
}

func (k *popupKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.close}}
}

func popupKeys() *popupKeyMap {
	return &popupKeyMap{
		close: key.NewBinding(
			key.WithKeys("ctrl+q", "esc", "ctrl+c"),
			key.WithHelp("esc/^q/^c", "close popup"),
		),
	}
}
