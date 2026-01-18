package components

import (
	"github.com/alphameo/nm-tui/internal/infra"
	"github.com/alphameo/nm-tui/internal/ui/styles"
	"github.com/alphameo/nm-tui/internal/ui/tools/compositor"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	borderOffset int = 2
	tabBarHeight int = 3
)

type MainModel struct {
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
	popup.XAnchor = compositor.Center
	popup.YAnchor = compositor.Center
	popup.EscapeKeys = escKeys
	notification := *NewFloatingModel(nil, "Notification")
	notification.XAnchor = compositor.Center
	notification.YAnchor = compositor.Center
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

// UpdateMsg is fictive variable, which used to send as tea.Msg instead of nil to trigger main window re-render
var UpdateMsg = struct{}{}

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
	case tea.Cmd:
		return m, msg
	case tea.KeyMsg:
		return m, m.handleKeyMsg(msg)
	}
	return m, m.handleMsg(msg)
}

func (m *MainModel) Resize(width, height int) {
	m.width = width
	m.height = height

	m.tabs.Resize(width, height)
}

func (m MainModel) View() string {
	view := m.tabs.View()

	if m.popup.IsActive {
		view = m.popup.Place(view, styles.OverlayStyle)
	}
	if m.notification.IsActive {
		view = m.notification.Place(view, styles.OverlayStyle)
	}
	return view
}

func (m *MainModel) handleKeyMsg(keyMsg tea.KeyMsg) tea.Cmd {
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
	}
	upd, cmd := m.tabs.Update(keyMsg)
	m.tabs = upd.(TabsModel)
	return cmd
}

func (m *MainModel) handleMsg(msg tea.Msg) tea.Cmd {
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

func SetPopupContentCmd(content tea.Model, title string) tea.Cmd {
	return func() tea.Msg {
		return PopupContentMsg{content, title}
	}
}

func SetPopupActivityCmd(isActive bool) tea.Cmd {
	return func() tea.Msg {
		return PopupActivityMsg(isActive)
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
	return tea.Sequence(SetNotificationTextCmd(text), SetNotificationActivityCmd(true))
}
