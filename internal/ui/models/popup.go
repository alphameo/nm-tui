package models

import (
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
