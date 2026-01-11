package views

import (
	"fmt"
	"strings"

	"github.com/alphameo/nm-tui/internal/infra"
	"github.com/alphameo/nm-tui/internal/ui/components/overlay"
	"github.com/alphameo/nm-tui/internal/ui/styles"
	tea "github.com/charmbracelet/bubbletea"
)

type sessionState uint

const (
	wifiView sessionState = iota
	timerView
	stateViewHeight int = 1
	borderOffset    int = 2
	tabBarHeight    int = 3
)

type MainModel struct {
	state        sessionState
	tabs         TabsModel
	popup        FloatingModel
	notification FloatingModel
	width        int
	height       int
}

func NewMainModel(networkManager infra.NetworkManager) MainModel {
	wifiTable := *NewConnectionsModel(networkManager)
	escKeys := []string{"ctrl+q", "esc", "ctrl+c"}
	popup := *NewFloatingModel(nil, "")
	popup.Width = 100
	popup.Height = 10
	popup.XAnchor = overlay.Center
	popup.YAnchor = overlay.Center
	popup.EscapeKeys = escKeys
	notification := *NewFloatingModel(nil, "Notification")
	notification.XAnchor = overlay.Center
	notification.YAnchor = overlay.Center
	notification.Width = 100
	notification.Height = 10
	notification.EscapeKeys = escKeys
	m := MainModel{
		tabs:         wifiTable,
		popup:        popup,
		notification: notification,
	}
	m.tabs.Resize(51, 20)
	return m
}

func (m MainModel) Init() tea.Cmd {
	return tea.Batch(
		m.tabs.Init(),
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
		m.popup.Content = msg.model
		m.popup.Title = msg.title
		return m, m.popup.Content.Init()
	case PopupActivityMsg:
		m.popup.IsActive = bool(msg)
		return m, nil
	case NotificationTextMsg:
		m.notification.Content = NewTextModel(string(msg))
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

	height -= stateViewHeight

	m.tabs.Resize(width, height)
}

func (m MainModel) View() string {
	sb := strings.Builder{}
	fmt.Fprintf(&sb, "%s\n state: %v", m.tabs.View(), m.state)
	view := sb.String()

	if m.popup.IsActive {
		view = m.popup.Place(view, styles.OverlayStyle)
	}
	if m.notification.IsActive {
		view = m.notification.Place(view, styles.OverlayStyle)
	}
	return view
}

func (m *MainModel) processKeyMsg(keyMsg tea.KeyMsg) tea.Cmd {
	if m.notification.IsActive {
		upd, cmd := m.notification.Update(keyMsg)
		m.notification = upd.(FloatingModel)
		return cmd
	} else if m.popup.IsActive {
		upd, cmd := m.popup.Update(keyMsg)
		m.popup = upd.(FloatingModel)
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
	upd, cmd := m.tabs.Update(keyMsg)
	m.tabs = upd.(TabsModel)
	return cmd
}

func (m *MainModel) processCommonMsg(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	var upd tea.Model
	upd, cmd = m.tabs.Update(msg)
	m.tabs = upd.(TabsModel)
	if cmd != nil {
		return cmd
	}
	if m.notification.IsActive {
		upd, cmd = m.notification.Update(msg)
		m.notification = upd.(FloatingModel)
		if cmd != nil {
			return cmd
		}
	}
	if m.popup.IsActive {
		upd, cmd = m.popup.Update(msg)
		m.popup = upd.(FloatingModel)
		if cmd != nil {
			return cmd
		}
	}
	return nil
}

type (
	PopupContentMsg struct {
		model tea.Model
		title string
	}
	PopupActivityMsg bool
)

func SetPopupContent(content tea.Model, title string) tea.Cmd {
	return func() tea.Msg {
		return PopupContentMsg{content, title}
	}
}

func SetPopupActivity(isActive bool) tea.Cmd {
	return func() tea.Msg {
		return PopupActivityMsg(isActive)
	}
}

type (
	NotificationTextMsg     string
	NotificationActivityMsg bool
)

func SetNotificationText(text string) tea.Cmd {
	return func() tea.Msg {
		return NotificationTextMsg(text)
	}
}

func SetNotificationActivity(isActive bool) tea.Cmd {
	return func() tea.Msg {
		return NotificationActivityMsg(isActive)
	}
}

func Notify(text string) tea.Cmd {
	return tea.Batch(SetNotificationActivity(true), SetNotificationText(text))
}
