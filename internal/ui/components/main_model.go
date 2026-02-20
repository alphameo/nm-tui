package components

import (
	"fmt"
	"time"

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

	NotificationCloseTime int = 5
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

type Popup struct {
	content tea.Model
	active  bool
	title   string
}
type popupKeyMap struct {
	close key.Binding
}

func (k *popupKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.close}
}

func (k *popupKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.close}}
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

	keyMngr *keyMapManager
	help    help.Model

	width  int
	height int
}

func NewMainModel(networkManager infra.NetworkManager) *MainModel {
	keys := defaultKeyMap

	con := NewWifiConnector(networkManager, keys.wifiConnector)
	a := NewWifiAvailableModel(con, keys.wifiAvailable, networkManager)

	info := NewStoredInfoModel(keys.wifiStoredInfo, networkManager)
	s := NewWifiStoredModel(info, keys.wifiStored, networkManager)

	wifi := NewWifiModel(a, s, keys.wifi, networkManager)

	wifiTable := NewTabsModel([]TabModel{wifi}, keys.tabs, networkManager)

	p := &Popup{active: false}
	n := &Notification{}
	help := help.New()
	help.ShowAll = true

	return &MainModel{
		tabs:         *wifiTable,
		popup:        p,
		notification: n,
		keyMngr:      keys,
		help:         help,
	}
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
		return m, DeferedCloseNotificationCmd()
	case tea.Cmd:
		return m, msg
	case tea.KeyMsg:
		return m, m.handleKeyMsg(msg)
	}
	return m, m.handleMsg(msg)
}

func (m *MainModel) Resize(width, height int) {
	m.width = width
	helpHeight := lipgloss.Height(m.help.View(m.keyMngr))
	m.height = height - helpHeight

	m.tabs.Resize(m.width, m.height)
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

		title := fmt.Sprintf("[ %s ]", m.popup.title)

		popupView = compositor.Compose(
			title,
			popupView,
			compositor.Center,
			compositor.Begin,
			0,
			0,
		)
		view = compositor.Compose(
			popupView,
			view,
			compositor.Center,
			compositor.Center,
			0,
			0,
		)
	}
	if m.notification.active {
		notificationView := m.notification.message
		style := styles.NotifBorderedStyle
		notificationView = style.Render(notificationView)
		notificationView = compositor.Compose(
			m.notification.title,
			notificationView,
			compositor.Center,
			compositor.Begin,
			0,
			0,
		)
		view = compositor.Compose(
			notificationView,
			view,
			compositor.End,
			compositor.Begin,
			-1,
			1,
		)
	}

	help := m.help.View(m.keyMngr)
	if m.popup.active {
		help = m.help.View(m.keyMngr.popup)
	}
	view = lipgloss.JoinVertical(lipgloss.Center, view, help)
	return view
}

func (m *MainModel) handleKeyMsg(keyMsg tea.KeyMsg) tea.Cmd {
	if m.popup.active {
		if key.Matches(keyMsg, m.keyMngr.popup.close) {
			return SetPopupActivityCmd(false)
		}
		upd, cmd := m.popup.content.Update(keyMsg)
		m.popup.content = upd
		return cmd
	}
	if key.Matches(keyMsg, m.keyMngr.main.quit) {
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
	return tea.Sequence(
		SetNotificationTextCmd(text),
		SetNotificationActivityCmd(true),
	)
}

func DeferedCloseNotificationCmd() tea.Cmd {
	return func() tea.Msg {
		time.Sleep(time.Second * time.Duration(NotificationCloseTime))
		return NotificationActivityMsg(false)
	}
}
