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

func (p Popup) Init() tea.Cmd {
	return p.content.Init()
}

func (p Popup) Update(msg tea.Msg) (Popup, tea.Cmd) {
	if !p.active {
		return p, nil
	}

	var cmd tea.Cmd
	p.content, cmd = p.content.UpdateAsPopup(msg)
	return p, cmd
}

func (p Popup) View() string {
	return p.content.View()
}
