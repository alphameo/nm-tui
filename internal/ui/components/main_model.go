package components

import (
	"time"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/alphameo/nm-tui/internal/infra"
	"github.com/alphameo/nm-tui/internal/ui/components/tabview"
	"github.com/alphameo/nm-tui/internal/ui/styles"
	"github.com/alphameo/nm-tui/internal/ui/tools/compositor"
)

const notificationCloseTime time.Duration = 50 * time.Second

type mainKeyMap struct {
	quit key.Binding
}

func (k *mainKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.quit}
}

func (k *mainKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.quit}}
}

var mainKeys = &mainKeyMap{
	quit: key.NewBinding(
		key.WithKeys("q", "ctrl+q", "esc", "ctrl+c"),
		key.WithHelp("esc/q/^q/^c", "quit"),
	),
}

type PopupModel interface {
	Init() tea.Cmd
	PopupUpdate(msg tea.Msg) (PopupModel, tea.Cmd)
	View() string
}
type Popup struct {
	content PopupModel
	active  bool
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

var popupKeys = &popupKeyMap{
	close: key.NewBinding(
		key.WithKeys("ctrl+q", "esc", "ctrl+c"),
		key.WithHelp("esc/^q/^c", "close popup"),
	),
}

type Notification struct {
	message   string
	active    bool
	title     string
	closeTime time.Duration
	style     *lipgloss.Style
}

type MainModel struct {
	tabs         *tabview.Model
	popup        *Popup
	notification *Notification

	keyMngr *keyMapManager
	help    help.Model

	width  int
	height int
}

func NewMainModel(wifiManager infra.WifiManager, networkManager infra.NetworkManager) *MainModel {
	keys := defaultKeyMap

	conn := NewWifiConnector(keys.wifiConnector, wifiManager)
	a := NewWifiAvailableModel(conn, keys.wifiAvailable, wifiManager)

	info := NewWifiInfoModel(keys.wifiSavedInfo, wifiManager)
	s := NewWifiSavedModel(info, keys.wifiSaved, wifiManager)

	wifi := NewWifiModel(a, s, conn, keys.wifi, wifiManager)
	network := NewNetworkModel(networkManager, keys.network)

	wifiTable := tabview.New([]tabview.Tab{
		{Title: "Wi-Fi", Content: wifi},
		{Title: "Networking", Content: network},
	}, styles.TabViewStyles, keys.tabs)

	p := &Popup{
		active: false,
	}

	notifStyle := lipgloss.NewStyle().Inherit(styles.NotifBorderedStyle)
	n := &Notification{style: &notifStyle, closeTime: notificationCloseTime}

	help := help.New()
	help.ShowAll = true

	return &MainModel{
		tabs:         wifiTable,
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
		return m, m.popup.content.Init()
	case PopupActivityMsg:
		m.popup.active = bool(msg)
		return m, nil
	case NotificationTextMsg:
		m.notification.message = string(msg)
		return m, nil
	case NotificationActivityMsg:
		m.notification.active = bool(msg)
		return m, DeferedCloseNotificationCmd(m.notification.closeTime)
	case tea.Cmd:
		return m, msg
	case tea.KeyPressMsg:
		return m, m.handleKey(msg)
	}

	var cmd tea.Cmd
	m.tabs, cmd = m.tabs.Update(msg)
	if cmd != nil {
		return m, cmd
	}
	if m.popup.active {
		upd, cmd := m.popup.content.PopupUpdate(msg)
		m.popup.content = upd
		if cmd != nil {
			return m, cmd
		}
	}
	return m, nil
}

func (m *MainModel) handleKey(keyMsg tea.KeyPressMsg) tea.Cmd {
	if m.popup.active {
		if key.Matches(keyMsg, m.keyMngr.popup.close) {
			return SetPopupActivityCmd(false)
		}
		upd, cmd := m.popup.content.PopupUpdate(keyMsg)
		m.popup.content = upd
		return cmd
	}
	if key.Matches(keyMsg, m.keyMngr.main.quit) {
		return tea.Quit
	}
	var cmd tea.Cmd
	m.tabs, cmd = m.tabs.Update(keyMsg)
	return cmd
}

func (m *MainModel) Width() int {
	return m.width
}

func (m *MainModel) Height() int {
	return m.height
}

func (m *MainModel) Resize(width, height int) {
	m.width = width
	m.height = height
	helpHeight := lipgloss.Height(m.help.View(m.keyMngr))

	m.tabs.Resize(width, m.height-helpHeight)

	notifStyle := m.notification.style.Width(width / 2)
	m.notification.style = &notifStyle
}

func (m MainModel) View() tea.View {
	view := m.tabs.View()

	if m.popup.active {
		popupView := m.popup.content.View()
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
		notificationView = m.notification.style.Render(notificationView)
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
	v := tea.NewView(view)
	v.AltScreen = true
	return v
}

type (
	PopupContentMsg struct {
		model PopupModel
	}
	PopupActivityMsg bool
)

func SetPopupContentCmd(content PopupModel) tea.Cmd {
	return func() tea.Msg {
		return PopupContentMsg{content}
	}
}

func SetPopupActivityCmd(isActive bool) tea.Cmd {
	return func() tea.Msg {
		return PopupActivityMsg(isActive)
	}
}

func OpenPopup(content PopupModel) tea.Cmd {
	return tea.Sequence(
		SetPopupContentCmd(content),
		SetPopupActivityCmd(true),
	)
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
