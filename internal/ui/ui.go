// Package ui contains Model, which represents main window of TUI
package ui

import (
	"fmt"
	"strings"

	"github.com/alphameo/nm-tui/internal/ui/components/label"
	"github.com/alphameo/nm-tui/internal/ui/components/overlay"
	"github.com/alphameo/nm-tui/internal/ui/connections"
	"github.com/alphameo/nm-tui/internal/ui/controls"
	"github.com/alphameo/nm-tui/internal/ui/styles"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type sessionState uint

const (
	wifiView sessionState = iota
	timerView
	stateViewHeight int = 2
)

type Model struct {
	state        sessionState
	connections  connections.Model
	popup        overlay.Model
	notification overlay.Model
	width        int
	height       int
}

func New() Model {
	wifiTable := *connections.New()
	escKeys := []string{"ctrl+q", "esc", "ctrl+c"}
	popup := *overlay.New(nil)
	popup.Width = 100
	popup.Height = 10
	popup.XAnchor = overlay.Center
	popup.YAnchor = overlay.Center
	popup.EscapeKeys = escKeys
	notification := *overlay.New(nil)
	notification.XAnchor = overlay.Center
	notification.YAnchor = overlay.Center
	notification.Width = 100
	notification.Height = 10
	notification.EscapeKeys = escKeys
	m := Model{
		connections:  wifiTable,
		timer:        timer,
		popup:        popup,
		notification: notification,
	}
	m.connections.Resize(51, 20)
	return m
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.connections.Init(),
		m.popup.Init(),
		m.notification.Init(),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Resize(msg.Width, msg.Height)
		return m, nil
	case controls.PopupContentMsg:
		m.popup.Content = msg
		return m, m.popup.Content.Init()
	case controls.PopupActivityMsg:
		m.popup.IsActive = bool(msg)
		return m, nil
	case controls.NotificationTextMsg:
		m.notification.Content = label.New(string(msg))
		return m, nil
	case controls.NotificationActivityMsg:
		m.notification.IsActive = bool(msg)
		return m, nil
	case tea.KeyMsg:
		return m, m.processKeyMsg(msg)
	}
	return m, m.processCommonMsg(msg)
}

func (m *Model) Resize(width, height int) {
	m.width = width
	m.height = height
	width -= styles.BorderOffset
	height -= styles.BorderOffset
	height -= stateViewHeight
	m.connections.Resize(width, height)
}

func (m Model) View() string {
	sb := strings.Builder{}
	fmt.Fprintf(&sb, "%s\n\n state: %v", m.connections.View(), m.state)
	view := sb.String()
	style := lipgloss.NewStyle().
		BorderStyle(styles.BorderStyle).
		Width(m.width - styles.BorderOffset).
		Height(m.height - styles.BorderOffset)
	view = style.Render(view)

	if m.popup.IsActive {
		view = m.popup.Place(view, styles.OverlayStyle)
	}
	if m.notification.IsActive {
		view = m.notification.Place(view, styles.OverlayStyle)
	}
	return view
}

func (m *Model) processKeyMsg(keyMsg tea.KeyMsg) tea.Cmd {
	if m.notification.IsActive {
		upd, cmd := m.notification.Update(keyMsg)
		m.notification = upd.(overlay.Model)
		return cmd
	} else if m.popup.IsActive {
		upd, cmd := m.popup.Update(keyMsg)
		m.popup = upd.(overlay.Model)
		return cmd
	}
	switch keyMsg.String() {
	case "q", "ctrl+q", "esc", "ctrl+c":
		return tea.Quit
	case "s":
		if m.state == wifiView {
			m.state = timerView
		} else {
			m.state = wifiView
		}
		return nil
	}
	upd, cmd := m.connections.Update(keyMsg)
	m.connections = upd.(connections.Model)
	return cmd
}

func (m *Model) processCommonMsg(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	var upd tea.Model
	upd, cmd = m.connections.Update(msg)
	m.connections = upd.(connections.Model)
	if cmd != nil {
		return cmd
	}
	if m.notification.IsActive {
		upd, cmd = m.notification.Update(msg)
		m.notification = upd.(overlay.Model)
		if cmd != nil {
			return cmd
		}
	}
	if m.popup.IsActive {
		upd, cmd = m.popup.Update(msg)
		m.popup = upd.(overlay.Model)
		if cmd != nil {
			return cmd
		}
	}
	return nil
}
