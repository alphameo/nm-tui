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
	quit key.Binding
}

func (k mainKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.quit}
}

func (k mainKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.quit}}
}

type MainModel struct {
	tabs         TabsModel
	popup        FloatingModel
	notification FloatingModel

	keys mainKeyMap
	help help.Model

	width  int
	height int
}

var mainKeys = mainKeyMap{
	quit: key.NewBinding(
		key.WithKeys("q", "ctrl+q", "esc", "ctrl+c"),
		key.WithHelp("esc/q/^Q/^C", "quit"),
	),
}

var floatingKeys = floatingKeyMap{
	quit: key.NewBinding(
		key.WithKeys("ctrl+q", "esc", "ctrl+c"),
		key.WithHelp("esc/^Q/^C", "quit"),
	),
}

var tabsKeys = &tabsKeyMap{
	tabNext: key.NewBinding(
		key.WithKeys("]"),
		key.WithHelp("]", "next tab"),
	),
	tabPrev: key.NewBinding(
		key.WithKeys("["),
		key.WithHelp("[", "previous tab"),
	),
}

var toggleKeys = &toggleKeyMap{
	toggle: key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp(" ", "toggle"),
	),
}

var wifiStoredKeys = &WifiStoredKeyMap{
	edit: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "edit"),
	),
	connect: key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp(" ", "connect"),
	),
	disconnect: key.NewBinding(
		key.WithKeys("shift+"),
		key.WithHelp("shift+ ", "disconnect"),
	),
	update: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "rescan stored"),
	),
	delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete"),
	),
}

var wifiStoredInfoKeys = &wifiStoredInfoKeyMap{
	togglePasswordVisibility: key.NewBinding(
		key.WithKeys("ctrl+r"),
		key.WithHelp("^R", "toggle password visibility"),
	),
	up: key.NewBinding(
		key.WithKeys("ctrl+k"),
		key.WithHelp("^K", "up"),
	),
	down: key.NewBinding(
		key.WithKeys("ctrl+j"),
		key.WithHelp("^J", "down"),
	),
	submit: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "submit"),
	),
}

var wifiKeys = &wifiKeyMap{
	nextWindow: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "next window"),
	),
	firstWindow: key.NewBinding(
		key.WithKeys("1"),
		key.WithHelp("1", "first window"),
	),
	secondWindow: key.NewBinding(
		key.WithKeys("2"),
		key.WithHelp("2", "second window"),
	),
	rescan: key.NewBinding(
		key.WithKeys("ctrl+r"),
		key.WithHelp("^R", "rescan"),
	),
}

var wifiAvailableKeys = &wifiAvailableKeyMap{
	rescan: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "rescan"),
	),
	openConnector: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "open connector"),
	),
}

func NewMainModel(networkManager infra.NetworkManager) *MainModel {
	wifiTable := *NewTabsModel(networkManager, tabsKeys)
	popup := *NewFloatingModel(nil, "")
	popup.Width = 100
	popup.Height = 10
	popup.XAnchor = compositor.Center
	popup.YAnchor = compositor.Center
	popup.keys = floatingKeys
	notification := *NewFloatingModel(nil, "Notification")
	notification.XAnchor = compositor.Center
	notification.YAnchor = compositor.Center
	notification.Width = 100
	notification.Height = 10
	notification.keys = floatingKeys
	m := &MainModel{
		tabs:         wifiTable,
		popup:        popup,
		notification: notification,
		keys:         mainKeys,
		help:         help.New(),
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
	help := m.help.View(m.keys)

	view = lipgloss.JoinVertical(lipgloss.Center, view, help)
	return view
}

func (m *MainModel) handleKeyMsg(keyMsg tea.KeyMsg) tea.Cmd {
	if m.notification.IsActive {
		upd, cmd := m.notification.Update(keyMsg)
		m.notification = upd.(FloatingModel)
		return cmd
	}
	if m.popup.IsActive {
		upd, cmd := m.popup.Update(keyMsg)
		m.popup = upd.(FloatingModel)
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
