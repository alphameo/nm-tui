package components

import (
	"github.com/alphameo/nm-tui/internal/infra"
	"github.com/alphameo/nm-tui/internal/ui/components/floating"
	"github.com/alphameo/nm-tui/internal/ui/components/text"
	"github.com/alphameo/nm-tui/internal/ui/styles"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	BorderOffset int = 2
	TabBarHeight int = 3
)

type mainKeyMap struct {
	quit key.Binding
}

func (k *mainKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.quit}
}

func (k *mainKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.quit}}
}

type MainModel struct {
	tabs         TabsModel
	popup        floating.Model
	notification floating.Model

	keys *mainKeyMap
	help help.Model

	width  int
	height int
}

func NewMainModel(networkManager infra.NetworkManager) *MainModel {
	keys := defaultKeyMap
	wifiTable := NewTabsModel(networkManager, keys)
	popup := floating.New(nil)
	popup.Title = "Popup"
	popup.Keys = keys.floating
	popup.Width = 100
	popup.Height = 10
	popup.XAnchor = floating.Center
	popup.YAnchor = floating.Center
	popup.ContentAlignHorizontal = lipgloss.Center
	popup.ContentAlignVertical = lipgloss.Center
	notification := floating.New(nil)
	notification.Keys = keys.floating
	notification.XAnchor = floating.Center
	notification.YAnchor = floating.Center
	notification.Width = 100
	notification.Height = 10
	notification.ContentAlignHorizontal = lipgloss.Center
	notification.ContentAlignVertical = lipgloss.Center
	m := &MainModel{
		tabs:         *wifiTable,
		popup:        *popup,
		notification: *notification,
		keys:         keys.main,
		help:         help.New(),
	}
	return m
}

func (m MainModel) Init() tea.Cmd {
	return tea.Batch(
		m.tabs.Init(),
		m.popup.Init(),
		m.notification.Init(),
	)
}

// NilMsg is a fictive struct, which used to send as tea.Msg instead of nil to trigger main window re-render
type NilMsg struct{}

// NilCmd is a function, which returns fictive Msg to trigger Model Update
var NilCmd = func() tea.Msg {
	return NilMsg{}
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
		m.notification.Content = text.New(string(msg))
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
	help := m.help.View(m.keys)

	view = lipgloss.JoinVertical(lipgloss.Center, view, help)
	return view
}

func (m *MainModel) handleKeyMsg(keyMsg tea.KeyMsg) tea.Cmd {
	if m.notification.IsActive {
		upd, cmd := m.notification.Update(keyMsg)
		m.notification = upd.(floating.Model)
		return cmd
	}
	if m.popup.IsActive {
		upd, cmd := m.popup.Update(keyMsg)
		m.popup = upd.(floating.Model)
		return cmd
	}
	if key.Matches(keyMsg, m.keys.quit) {
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
		m.notification = upd.(floating.Model)
		if cmd != nil {
			return cmd
		}
	}
	if m.popup.IsActive {
		upd, cmd = m.popup.Update(msg)
		m.popup = upd.(floating.Model)
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
