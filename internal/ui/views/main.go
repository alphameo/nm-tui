// Package views provides various view models
package views

import (
	"fmt"
	"strings"

	"github.com/alphameo/nm-tui/internal/infra"
	"github.com/alphameo/nm-tui/internal/ui/components/label"
	"github.com/alphameo/nm-tui/internal/ui/components/overlay"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type sessionState uint

const (
	wifiView sessionState = iota
	timerView
	stateViewHeight int = 2
	borderOffset    int = 2
	tabBarHeight    int = 3
)

type MainModel struct {
	state        sessionState
	connections  ConnectionsModel
	popup        overlay.Model
	notification overlay.Model
	width        int
	height       int
}

func NewMainModel(networkManager infra.NetworkManager) MainModel {
	wifiTable := *NewConnectionsModel(networkManager)
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
	m := MainModel{
		connections:  wifiTable,
		popup:        popup,
		notification: notification,
	}
	m.connections.Resize(51, 20)
	return m
}

func (m MainModel) Init() tea.Cmd {
	return tea.Batch(
		m.connections.Init(),
		m.popup.Init(),
		m.notification.Init(),
	)
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Resize(msg.Width, msg.Height)
		return m, nil
	case PopupContentMsg:
		m.popup.Content = msg
		return m, m.popup.Content.Init()
	case PopupActivityMsg:
		m.popup.IsActive = bool(msg)
		return m, nil
	case NotificationTextMsg:
		m.notification.Content = label.New(string(msg))
		return m, nil
	case NotificationActivityMsg:
		m.notification.IsActive = bool(msg)
		return m, nil
	case tea.KeyMsg:
		return m, m.processKeyMsg(msg)
	}
	return m, m.processCommonMsg(msg)
}

func (m *MainModel) Resize(width, height int) {
	m.width = width
	m.height = height
	width -= borderOffset
	height -= borderOffset
	height -= stateViewHeight
	m.connections.Resize(width, height)
}

func (m MainModel) View() string {
	sb := strings.Builder{}
	fmt.Fprintf(&sb, "%s\n\n state: %v", m.connections.View(), m.state)
	view := sb.String()
	style := lipgloss.NewStyle().
		BorderStyle(BorderStyle).
		Width(m.width - borderOffset).
		Height(m.height - borderOffset)
	view = style.Render(view)

	if m.popup.IsActive {
		view = m.popup.Place(view, OverlayStyle)
	}
	if m.notification.IsActive {
		view = m.notification.Place(view, OverlayStyle)
	}
	return view
}

func (m *MainModel) processKeyMsg(keyMsg tea.KeyMsg) tea.Cmd {
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
	m.connections = upd.(ConnectionsModel)
	return cmd
}

func (m *MainModel) processCommonMsg(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	var upd tea.Model
	upd, cmd = m.connections.Update(msg)
	m.connections = upd.(ConnectionsModel)
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
