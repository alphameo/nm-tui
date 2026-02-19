package components

import (
	"github.com/alphameo/nm-tui/internal/infra"
	"github.com/alphameo/nm-tui/internal/ui/styles"
	"github.com/alphameo/nm-tui/internal/ui/tools/compositor"
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
	quit       key.Binding
	closePopup key.Binding
}

func (k *mainKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.quit}
}

func (k *mainKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.quit}}
}

type Popup struct {
	content tea.Model
	active  bool
	title   string
}
type Notification struct {
	message string
	active  bool
	title   string
}

type MainModel struct {
	tabs         TabsModel
	popup        *Popup
	notification *Notification

	keys *mainKeyMap
	help help.Model

	width  int
	height int
}

func NewMainModel(networkManager infra.NetworkManager) *MainModel {
	keys := defaultKeyMap
	wifiTable := NewTabsModel(networkManager, keys)
	p := &Popup{active: false}
	// popup.Keys = keys.compositor
	// popup.Width = 100
	// popup.Height = 10
	// popup.XAnchor = compositor.Center
	// popup.YAnchor = compositor.Center
	// popup.ContentAlignHorizontal = lipgloss.Center
	// popup.ContentAlignVertical = lipgloss.Center
	n := &Notification{}
	// notification := compositor.New(nil)
	// notification.Keys = keys.compositor
	// notification.XAnchor = compositor.Center
	// notification.YAnchor = compositor.Center
	// notification.Width = 100
	// notification.Height = 10
	// notification.ContentAlignHorizontal = lipgloss.Center
	// notification.ContentAlignVertical = lipgloss.Center
	m := &MainModel{
		tabs:         *wifiTable,
		popup:        p,
		notification: n,
		keys:         keys.main,
		help:         help.New(),
	}
	return m
}

func (m MainModel) Init() tea.Cmd {
	return m.tabs.Init()
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
		m.popup.content = msg.model
		m.popup.title = msg.title
		return m, m.popup.content.Init()
	case PopupActivityMsg:
		m.popup.active = bool(msg)
		return m, nil
	case NotificationTextMsg:
		m.notification.message = string(msg)
		return m, nil
	case NotificationActivityMsg:
		m.notification.active = bool(msg)
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

	if m.popup.active {
		popupView := m.popup.content.View()
		style := styles.OverlayStyle.
			Align(lipgloss.Center).
			Width(100).
			Height(10)

		popupView = style.Render(popupView)

		popupView = compositor.Compose(m.popup.title, popupView, compositor.Center, compositor.Begin, 0, 0)
		view = compositor.Compose(popupView, view, compositor.Center, compositor.Center, 0, 0)
	}
	if m.notification.active {
		notificationView := m.notification.message
		style := styles.OverlayStyle.
			Align(lipgloss.Center)
		notificationView = style.Render(notificationView)
		notificationView = compositor.Compose(m.notification.title, notificationView, compositor.Center, compositor.Begin, 0, 0)
		view = compositor.Compose(notificationView, view, compositor.End, compositor.Begin, 0, 0)
	}
	help := m.help.View(m.keys)

	view = lipgloss.JoinVertical(lipgloss.Center, view, help)
	return view
}

func (m *MainModel) handleKeyMsg(keyMsg tea.KeyMsg) tea.Cmd {
	if m.popup.active {
		if key.Matches(keyMsg, m.keys.closePopup) {
			return SetPopupActivityCmd(false)
		}
		upd, cmd := m.popup.content.Update(keyMsg)
		m.popup.content = upd
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
	if m.popup.active {
		upd, cmd = m.popup.content.Update(msg)
		m.popup.content = upd
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
