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
	ready bool

	tabs         tabview.Model
	popup        Popup
	notification Notification

	networks *NetworksModel
	device   *DeviceModel

	connector      *ConnectorModel
	profileCreator *ProfileCreatorModel
	hotspotCreator *HotspotCreatorModel
	profileEditor  *ProfileEditorModel

	keys  *mainKeyMap
	help  *HelpModel
	style lipgloss.Style
}

func NewMainModel(
	networksManager infra.NetworksManager,
	deviceManager infra.DeviceManager,
	portalOpener infra.CaptivePortalOpener,
	cfg config.Config,
) (*MainModel, error) {
	err := styles.Init(cfg)
	if err != nil {
		return nil, fmt.Errorf("style initialization: %w", err)
	}

	keys := initKeys(*cfg.Keys)

	mainCfg.notificationCloseTime = time.Duration(*cfg.NotifCloseTime) * time.Second

	connector := NewConnectorModel(keys.connector, networksManager)
	connector.style = styles.OverlayStyle
	profileCreator := NewProfileCreatorModel(keys.profileCreator, networksManager)
	profileCreator.style = styles.OverlayStyle
	hotspotCreator := NewHotspotCreatorModel(keys.hotspotCreator, networksManager)
	hotspotCreator.style = styles.OverlayStyle
	profileEditor := NewProfileEditorModel(keys.profileEditor, networksManager)
	profileEditor.style = styles.OverlayStyle

	available := NewAvailableNetworksModel(keys.availableNetworks, networksManager)
	available.focusedStyle = styles.BorderedFocusedStyle
	available.bluredStyle = styles.BorderedStyle
	available.SetTableStyles(styles.TableStyles, styles.DataTableStyles)
	profiles := NewNetworkProfilesModel(keys.networkProfiles, networksManager)
	profiles.focusedStyle = styles.BorderedFocusedStyle
	profiles.bluredStyle = styles.BorderedStyle
	profiles.SetTableStyles(styles.TableStyles, styles.DataTableStyles)

	networks := NewNetworksModel(available, profiles, keys.networks, networksManager, portalOpener)
	device := NewDeviceModel(keys.device, deviceManager)
	device.tableStyle = styles.BorderedStyle
	tabContentBorder := tabview.DefaultContentBorder(styles.Border)
	tabContentStyle := styles.DefaultStyle.Border(tabContentBorder)
	networks.style = tabContentStyle
	device.style = tabContentStyle

	tabs := tabview.New([]tabview.Tab{
		{Title: networks.Title(), Content: networks},
		{Title: device.Title(), Content: device},
	})
	tabs.SetStyles(styles.TabViewStyles)
	tabs.Keys = keys.tabs

	p := Popup{
		active: false,
	}

	notifStyle := lipgloss.NewStyle().Inherit(styles.NotifBorderedStyle)
	n := Notification{style: notifStyle, closeTime: mainCfg.notificationCloseTime}

	help := NewHelpModel(keys)
	help.style = styles.OverlayStyle

	return &MainModel{
		tabs:         tabs,
		popup:        p,
		notification: n,

		networks: networks,
		device:   device,

		connector:      connector,
		profileCreator: profileCreator,
		hotspotCreator: hotspotCreator,
		profileEditor:  profileEditor,

		keys:  &keys.main,
		help:  help,
		style: lipgloss.NewStyle(),
	}, nil
}

func (m *MainModel) Init() tea.Cmd {
	var cmds []tea.Cmd
	cmds = append(cmds, m.tabs.Init())

	return tea.Batch(cmds...)
}

func (m *MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.ready = true
		m.Resize(msg.Width, msg.Height)
		return m, nil
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
			m.connector.setNewNetworkCmd(string(msg)),
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
			m.profileEditor.setNewProfile(string(msg)),
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
		return m.updateOnKeyPress(msg)
	}

	var cmd tea.Cmd
	var cmds []tea.Cmd

	m.tabs, cmd = m.tabs.Update(msg)
	cmds = append(cmds, cmd)

	m.popup, cmd = m.popup.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *MainModel) updateOnKeyPress(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	if m.popup.active {
		if key.Matches(msg, m.keys.closePopup) {
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

func (m *MainModel) View() tea.View {
	if !m.ready {
		return tea.View{}
	}

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

	help := m.shortHelpView()
	view = lipgloss.JoinVertical(lipgloss.Center, view, help)
	view = m.style.Render(view)
	v := tea.NewView(view)
	v.AltScreen = true
	return v
}

func (m *MainModel) Width() int { return m.style.GetWidth() }

func (m *MainModel) Height() int { return m.style.GetHeight() }

func (m *MainModel) Resize(width, height int) {
	m.style = m.style.Width(width).Height(height)

	border := m.style.GetBorderStyle()
	width -= border.GetLeftSize() + border.GetRightSize()
	width -= border.GetBottomSize() + border.GetTopSize()

	helpHeight := lipgloss.Height(m.shortHelpView())

	m.tabs.Resize(width, height-helpHeight)
	m.help.Resize(int(float32(width)*0.8), int(float32(height)*0.8))
	m.help.help.SetWidth(width)

	notifStyle := m.notification.style.Width(width / 2)
	m.notification.style = notifStyle
}

func (m *MainModel) activeBindingsShort() []key.Binding {
	if m.popup.active {
		switch m.popup.content.(type) {
		case *ConnectorModel:
			return m.help.connectorShort()
		case *ProfileCreatorModel:
			return m.help.profileCreatorShort()
		case *HotspotCreatorModel:
			return m.help.hotspotCreatorShort()
		case *ProfileEditorModel:
			return m.help.profileEditorShort()
		}
		return m.help.mainShort()
	}

	helpKey := m.help.mainShort()

	switch m.tabs.ActiveTabIndex() {
	case 1: // Device tab
		return append(helpKey, m.help.deviceShort()...)
	default: // Networks tab: tab actions + focused window
		keys := []key.Binding{}
		keys = append(keys, m.help.networksShort()...)
		if m.networks.available.Focused() {
			keys = append(keys, m.help.availableNetworksShort()...)
		} else {
			keys = append(keys, m.help.networkProfilesShort()...)
		}
		return append(helpKey, keys...)
	}
}

func (m *MainModel) shortHelpView() string {
	return m.help.ShortViewFor(m.activeBindingsShort())
}
