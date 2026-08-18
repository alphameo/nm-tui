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

type (
	OpenPopupMsg struct {
		model PopupModel
	}
	ClosePopupMsg struct{}
)

func OpenPopupCmd(content PopupModel) tea.Cmd {
	return func() tea.Msg {
		return OpenPopupMsg{model: content}
	}
}

func ClosePopupCmd() tea.Cmd {
	return func() tea.Msg {
		return ClosePopupMsg{}
	}
}

type (
	openConnectorMsg      string
	openHotspotCreatorMsg struct{}
	openProfileCreatorMsg struct{}
	openProfileEditorMsg  string
)

func OpenConnectorCmd(ssid string) tea.Cmd {
	return func() tea.Msg {
		return openConnectorMsg(ssid)
	}
}

func OpenHotspotCreatorCmd() tea.Cmd {
	return func() tea.Msg {
		return openHotspotCreatorMsg{}
	}
}

func OpenProfileCreatorCmd() tea.Cmd {
	return func() tea.Msg {
		return openProfileCreatorMsg{}
	}
}

func OpenProfileEditorCmd(name string) tea.Cmd {
	return func() tea.Msg {
		return openProfileEditorMsg(name)
	}
}
