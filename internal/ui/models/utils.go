package models

import (
	"errors"
	"fmt"
	"time"

	tea "charm.land/bubbletea/v2"
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

var ErrPasswordFmt error = errors.New("wrong password format")

func passwordValidator(input string) error {
	if len(input) < 8 {
		return fmt.Errorf("%w: length < 8", ErrPasswordFmt)
	}
	return nil
}
