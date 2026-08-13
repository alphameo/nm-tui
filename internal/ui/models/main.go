package models

import (
	"fmt"
	"time"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/alphameo/nm-tui/internal/config"
	"github.com/alphameo/nm-tui/internal/infra"
	"github.com/alphameo/nm-tui/internal/ui/models/tabview"
	"github.com/alphameo/nm-tui/internal/ui/styles"
	"github.com/alphameo/nm-tui/internal/ui/tools/compositor"
)

type mainConfig struct {
	notificationCloseTime time.Duration
}

var mainCfg = mainConfig{
	notificationCloseTime: 50 * time.Second,
}

type mainKeyMap struct {
	quit       key.Binding
	closePopup key.Binding
	help       key.Binding
}

type MainModel struct {
	tabs         tabview.Model
	popup        Popup
	notification Notification

	connector      *ConnectorModel
	profileCreator *ProfileCreatorModel
	hotspotCreator *HotspotCreatorModel
	profileEditor  *ProfileEditorModel

	keys *mainKeyMap
	help *HelpModel

	width  int
	height int
}

func NewMainModel(networksManager infra.NetworksManager, connectivityManager infra.ConnectivityManager, portalOpener infra.CaptivePortalOpener, cfg config.Config) (*MainModel, error) {
	err := styles.Init(cfg)
	if err != nil {
		return nil, fmt.Errorf("style initialization: %w", err)
	}

	keys := initKeys(*cfg.Keys)

	mainCfg.notificationCloseTime = time.Duration(*cfg.NotifCloseTime) * time.Second

	connector := NewConnectorModel(keys.connector, networksManager)
	profileCreator := NewProfileCreatorModel(keys.profileCreator, networksManager)
	hotspotCreator := NewHotspotCreatorModel(keys.hotspotCreator, networksManager)
	profileEditor := NewProfileEditorModel(keys.profileEditor, networksManager)

	available := NewAvailableNetworksModel(keys.availableNetworks, networksManager)
	profiles := NewNetworkProfilesModel(keys.networkProfiles, networksManager)

	wifi := NewNetworksModel(available, profiles, keys.networks, networksManager, portalOpener)
	connectivity := NewConnectivityModel(keys.connectivity, connectivityManager)

	wifiTable := tabview.New([]tabview.Tab{
		{Title: wifi.Title(), Content: wifi},
		{Title: connectivity.Title(), Content: connectivity},
	})
	wifiTable.SetStyles(styles.TabViewStyles)
	wifiTable.Keys = keys.tabs

	p := Popup{
		active: false,
	}

	notifStyle := lipgloss.NewStyle().Inherit(styles.NotifBorderedStyle)
	n := Notification{style: notifStyle, closeTime: mainCfg.notificationCloseTime}

	return &MainModel{
		tabs:         wifiTable,
		popup:        p,
		notification: n,

		connector:      connector,
		profileCreator: profileCreator,
		hotspotCreator: hotspotCreator,
		profileEditor:  profileEditor,

		keys: &keys.main,
		help: NewHelpModel(keys),
	}, nil
}

func (m MainModel) Init() tea.Cmd {
	var cmds []tea.Cmd
	cmds = append(cmds, m.tabs.Init())

	// NOTE: Request base text color for directing it into colorscheme,
	// preventing wide nerd icons narrowing
	if styles.TextColor == lipgloss.Color("") {
		cmds = append(cmds, tea.RequestForegroundColor)
	}

	return tea.Batch(cmds...)
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Resize(msg.Width, msg.Height)
		return m, nil
	case tea.ForegroundColorMsg:
		styles.TextColor = msg
		styles.UpdateColorscheme()
	case OpenPopupMsg:
		m.popup.content = msg.model
		m.popup.active = true
		return m, m.popup.content.Init()
	case ClosePopupMsg:
		m.popup.content = nil
		m.popup.active = false
		return m, nil
	case openConnectorMsg:
		return m, tea.Batch(
			m.connector.setNew(string(msg)),
			OpenPopupCmd(m.connector),
		)
	case openHotspotCreatorMsg:
		return m, tea.Batch(
			m.hotspotCreator.Reset(),
			OpenPopupCmd(m.hotspotCreator),
		)
	case openProfileCreatorMsg:
		return m, tea.Batch(
			m.profileCreator.Reset(),
			OpenPopupCmd(m.profileCreator),
		)
	case openProfileEditorMsg:
		return m, tea.Batch(
			m.profileEditor.setNew(string(msg)),
			OpenPopupCmd(m.profileEditor),
		)
	case NotificationTextMsg:
		m.notification.message = string(msg)
		return m, nil
	case NotificationActivityMsg:
		m.notification.active = bool(msg)
		return m, DeferedCloseNotificationCmd(m.notification.closeTime)
	case tea.Cmd:
		return m, msg
	case tea.KeyPressMsg:
		var cmd tea.Cmd
		if m.popup.active {
			switch {
			case key.Matches(msg, m.keys.closePopup):
				return m, ClosePopupCmd()
			}
			m.popup, cmd = m.popup.Update(msg)
			return m, cmd
		}
		switch {
		case key.Matches(msg, m.keys.quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.help):
			return m, OpenPopupCmd(m.help)
		}
		m.tabs, cmd = m.tabs.Update(msg)
		return m, cmd
	}

	var cmd tea.Cmd
	var cmds []tea.Cmd

	m.tabs, cmd = m.tabs.Update(msg)
	cmds = append(cmds, cmd)

	m.popup, cmd = m.popup.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m MainModel) View() tea.View {
	view := m.tabs.View()

	if m.popup.active {
		popupView := m.popup.View()
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

	help := m.help.ShortView()
	view = lipgloss.JoinVertical(lipgloss.Center, view, help)
	v := tea.NewView(view)
	v.AltScreen = true
	return v
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
	helpHeight := lipgloss.Height(m.help.ShortView())

	m.tabs.Resize(width, m.height-helpHeight)
	m.help.Resize(width/2, height/2)

	notifStyle := m.notification.style.Width(width / 2)
	m.notification.style = notifStyle
}
