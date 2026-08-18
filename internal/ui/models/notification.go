package models

import (
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type Notification struct {
	message   string
	active    bool
	title     string
	closeTime time.Duration
	style     lipgloss.Style
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
	return tea.Tick(t, func(time.Time) tea.Msg {
		return NotificationActivityMsg(false)
	})
}
