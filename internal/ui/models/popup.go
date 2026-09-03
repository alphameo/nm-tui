package models

import (
	tea "charm.land/bubbletea/v2"
)

type PopupModel[T any] interface {
	Init() tea.Cmd
	Update(msg tea.Msg) (PopupModel[T], tea.Cmd)
	View() string
}

type Popup[T any] struct {
	content PopupModel[T]
	active  bool
}

func (p Popup[T]) Init() tea.Cmd {
	return p.content.Init()
}

func (p Popup[T]) Update(msg tea.Msg) (Popup[T], tea.Cmd) {
	if !p.active {
		return p, nil
	}

	var cmd tea.Cmd
	p.content, cmd = p.content.Update(msg)
	return p, cmd
}

func (p Popup[T]) View() string {
	return p.content.View()
}

type (
	OpenPopupMsg[T any] struct {
		model PopupModel[T]
	}
	ClosePopupMsg struct{}
)

func OpenPopupCmd[T any](content PopupModel[T]) tea.Cmd {
	return func() tea.Msg {
		return OpenPopupMsg[T]{model: content}
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
	openDeviceInfoMsg     string
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

func OpenDeviceInfoCmd(name string) tea.Cmd {
	return func() tea.Msg {
		return openDeviceInfoMsg(name)
	}
}
