package models

import (
	tea "charm.land/bubbletea/v2"
)

type PopupModel[T any] interface {
	Init() tea.Cmd
	Update(msg tea.Msg) (T, tea.Cmd)
	View() string
}

type Popup[T PopupModel[T]] struct {
	content T
}

func (p Popup[T]) Init() tea.Cmd {
	return p.content.Init()
}

func (p Popup[T]) Update(msg tea.Msg) (Popup[T], tea.Cmd) {
	var cmd tea.Cmd
	p.content, cmd = p.content.Update(msg)
	return p, cmd
}

func (p Popup[T]) View() string {
	return p.content.View()
}

type (
	OpenPopupMsg struct {
		kind popupKind
	}
	ClosePopupMsg struct{}
)

func OpenPopupCmd(kind popupKind) tea.Cmd {
	return func() tea.Msg {
		return OpenPopupMsg{kind: kind}
	}
}

func ClosePopupCmd() tea.Cmd {
	return func() tea.Msg {
		return ClosePopupMsg{}
	}
}

type (
	openConnectorMsg      struct{ ssid string }
	openHotspotCreatorMsg struct{}
	openProfileCreatorMsg struct{}
	openHelpMsg           struct{}
	openProfileEditorMsg  struct{ deviceID string }
	openDeviceInfoMsg     struct{ deviceName string }
)

func OpenConnectorCmd(ssid string) tea.Cmd {
	return func() tea.Msg {
		return openConnectorMsg{ssid: ssid}
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

func OpenHelpCmd() tea.Cmd {
	return func() tea.Msg {
		return openHelpMsg{}
	}
}

func OpenProfileEditorCmd(name string) tea.Cmd {
	return func() tea.Msg {
		return openProfileEditorMsg{deviceID: name}
	}
}

func OpenDeviceInfoCmd(name string) tea.Cmd {
	return func() tea.Msg {
		return openDeviceInfoMsg{deviceName: name}
	}
}
