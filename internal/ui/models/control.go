package models

import (
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/alphameo/nm-tui/internal/infra"
)

type Focusable interface {
	Focused() bool
	Focus() tea.Cmd
	Blur()
}

// NilMsg is a fictive struct, which used to send as tea.Msg instead of nil to trigger main window re-render
type NilMsg struct{}

// NilCmd is a function, which returns fictive Msg to trigger Model Update
var NilCmd = func() tea.Msg {
	return NilMsg{}
}

type (
	PopupContentMsg struct {
		model PopupModel
	}
	PopupActivityMsg bool
)

func SetPopupContentCmd(content PopupModel) tea.Cmd {
	return func() tea.Msg {
		return PopupContentMsg{content}
	}
}

func SetPopupActivityCmd(isActive bool) tea.Cmd {
	return func() tea.Msg {
		return PopupActivityMsg(isActive)
	}
}

func OpenPopup(content PopupModel) tea.Cmd {
	return tea.Sequence(
		SetPopupContentCmd(content),
		SetPopupActivityCmd(true),
	)
}

type (
	openConnectorMsg      string
	openHotspotCreatorMsg struct{}
	openProfileCreatorMsg struct{}
	openProfileEditorMsg  infra.NetworkInfo
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

func OpenProfileEditorCmd(info infra.NetworkInfo) tea.Cmd {
	return func() tea.Msg {
		return openProfileEditorMsg(info)
	}
}

type (
	NotificationTextMsg     string
	NotificationActivityMsg bool
)

func SetNotificationTextCmd(text string) tea.Cmd {
	return func() tea.Msg {
		return NotificationTextMsg(text)
	}
}

func SetNotificationActivityCmd(isActive bool) tea.Cmd {
	return func() tea.Msg {
		return NotificationActivityMsg(isActive)
	}
}

func NotifyCmd(text string) tea.Cmd {
	return tea.Sequence(
		SetNotificationTextCmd(text),
		SetNotificationActivityCmd(true),
	)
}

func DeferedCloseNotificationCmd(t time.Duration) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(t)
		return NotificationActivityMsg(false)
	}
}
