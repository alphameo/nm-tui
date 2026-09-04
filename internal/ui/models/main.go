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
	rescanInterval        time.Duration
}

type popupKind int

const (
	popupNo popupKind = iota
	popupConfirm
	popupConnector
	popupHelp
	popupDeviceInfo
	popupHotspotCreator
	popupProfileCreator
	popupProfileEditor
)

var mainCfg = mainConfig{
	notificationCloseTime: 50 * time.Second,
	rescanInterval:        10 * time.Second,
}

type mainKeyMap struct {
	quit       key.Binding
	closePopup key.Binding
	help       key.Binding
}

type MainModel struct {
	ready bool

	tabs         tabview.Model
	notification Notification
	confirm      *ConfirmModel

	networks *NetworksModel
	device   *DeviceModel

	connector      *ConnectorModel
	deviceInfo     *DeviceInfoModel
	help           *HelpModel
	hotspotCreator *HotspotCreatorModel
	profileCreator *ProfileCreatorModel
	profileEditor  *ProfileEditorModel
	activePopup    popupKind

	keys  *mainKeyMap
	Style lipgloss.Style
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
	mainCfg.rescanInterval = time.Duration(*cfg.RescanInterval) * time.Second

	connector := NewConnectorModel(keys.connector, networksManager)
	connector.Style = styles.OverlayStyle
	profileCreator := NewProfileCreatorModel(keys.profileCreator, networksManager)
	profileCreator.Style = styles.OverlayStyle
	hotspotCreator := NewHotspotCreatorModel(keys.hotspotCreator, networksManager)
	hotspotCreator.Style = styles.OverlayStyle
	profileEditor := NewProfileEditorModel(keys.profileEditor, networksManager)
	profileEditor.Style = styles.OverlayStyle
	deviceInfo := NewDeviceInfoModel(deviceManager)
	deviceInfo.Style = styles.OverlayStyle

	available := NewAvailableNetworksModel(keys.availableNetworks, networksManager)
	available.focusedStyle = styles.BorderedFocusedStyle
	available.bluredStyle = styles.BorderedStyle
	available.SetTableStyles(styles.TableStyles, styles.DataTableStyles)

	profiles := NewNetworkProfilesModel(keys.networkProfiles, networksManager)
	profiles.focusedStyle = styles.BorderedFocusedStyle
	profiles.bluredStyle = styles.BorderedStyle
	profiles.SetTableStyles(styles.TableStyles, styles.DataTableStyles)

	networks := NewNetworksModel(available, profiles, keys.networks, networksManager, portalOpener)
	networks.IndicatorStyle = styles.DefaultStyle

	netDevices := NewNetworkDevicesModel(keys.networkDevices, deviceManager)
	netDevices.focusedStyle = styles.BorderedFocusedStyle
	netDevices.bluredStyle = styles.BorderedStyle
	netDevices.SetTableStyles(styles.TableStyles, styles.DataTableStyles)

	device := NewDeviceModel(netDevices, keys.device, deviceManager)
	device.IndicatorStyle = styles.DefaultStyle
	device.ControlsStyle = lipgloss.NewStyle().Margin(1, 0)

	tabContentBorder := tabview.DefaultContentBorder(styles.Border)
	tabContentStyle := styles.DefaultStyle.Border(tabContentBorder)
	networks.Style = tabContentStyle
	device.Style = tabContentStyle

	tabs := tabview.New([]tabview.Tab{
		{Title: networks.Title(), Content: networks},
		{Title: device.Title(), Content: device},
	})
	tabs.SetStyles(styles.TabViewStyles)
	tabs.Keys = keys.tabs

	notifStyle := lipgloss.NewStyle().Inherit(styles.NotifBorderedStyle)
	n := Notification{style: notifStyle, closeTime: mainCfg.notificationCloseTime}

	confirmStyle := styles.OverlayStyle
	confirm := NewConfirmModel(confirmStyle, keys.confirm)

	help := NewHelpModel(keys)
	help.FullStyle = styles.OverlayStyle

	return &MainModel{
		tabs:         tabs,
		notification: n,
		confirm:      confirm,

		networks: networks,
		device:   device,

		connector:      connector,
		deviceInfo:     deviceInfo,
		help:           help,
		hotspotCreator: hotspotCreator,
		profileCreator: profileCreator,
		profileEditor:  profileEditor,

		keys:  &keys.main,
		Style: lipgloss.NewStyle(),
	}, nil
}

func (m *MainModel) Init() tea.Cmd {
	return tea.Batch(m.tabs.Init(), IntervalRescanCmd(mainCfg.rescanInterval))
}

//nolint:funlen // main event loop
func (m *MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.ready = true
		m.Resize(msg.Width, msg.Height)
		return m, nil
	case IntervalRescanMsg:
		var cmds []tea.Cmd
		if m.networks.indicatorState == NetsDone {
			cmds = append(cmds, m.networks.rescanCmd())
		}
		if m.device.indicatorState == DeviceDone {
			cmds = append(cmds, RescanDeviceCmd())
		}
		cmds = append(cmds, IntervalRescanCmd(mainCfg.rescanInterval))
		return m, tea.Batch(cmds...)
	case NetworksRescannedMsg:
		return m, tea.Batch(
			m.networks.available.setAvailable(msg.Available, msg.ScanErr),
			m.networks.profiles.setProfiles(msg.Profiles, msg.ProfilesErr),
		)
	case OpenPopupMsg:
		m.activePopup = msg.kind
		return m, m.initPopup()
	case ClosePopupMsg:
		m.activePopup = popupNo
		return m, nil
	case ConfirmMsg:
		m.confirm.Question = msg.question
		m.confirm.Action = msg.cmd
		return m, OpenPopupCmd(popupConfirm)
	case openConnectorMsg:
		return m, tea.Batch(
			m.connector.setNewNetworkCmd(msg.ssid),
			OpenPopupCmd(popupConnector),
		)
	case openHotspotCreatorMsg:
		return m, tea.Batch(
			m.hotspotCreator.Reset(),
			OpenPopupCmd(popupHotspotCreator),
		)
	case openProfileCreatorMsg:
		return m, tea.Batch(
			m.profileCreator.Reset(),
			OpenPopupCmd(popupProfileCreator),
		)
	case openProfileEditorMsg:
		return m, tea.Batch(
			m.profileEditor.setNewProfile(msg.deviceID),
			OpenPopupCmd(popupProfileEditor),
		)
	case openDeviceInfoMsg:
		return m, tea.Batch(
			m.deviceInfo.setNewDevice(msg.deviceName),
			OpenPopupCmd(popupDeviceInfo),
		)
	case openHelpMsg:
		return m, OpenPopupCmd(popupHelp)
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

	cmd = m.updatePopup(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *MainModel) initPopup() tea.Cmd {
	switch m.activePopup {
	case popupNo:
		return nil
	case popupConfirm:
		return m.confirm.Init()
	case popupConnector:
		return m.connector.Init()
	case popupHelp:
		return m.help.Init()
	case popupDeviceInfo:
		return m.deviceInfo.Init()
	case popupHotspotCreator:
		return m.hotspotCreator.Init()
	case popupProfileCreator:
		return m.profileCreator.Init()
	case popupProfileEditor:
		return m.profileEditor.Init()
	}
	return nil
}

func (m *MainModel) updatePopup(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	switch m.activePopup {
	case popupNo:
		return nil
	case popupConfirm:
		m.confirm, cmd = m.confirm.Update(msg)
		return cmd
	case popupConnector:
		m.connector, cmd = m.connector.Update(msg)
		return cmd
	case popupHelp:
		m.help, cmd = m.help.Update(msg)
		return cmd
	case popupDeviceInfo:
		m.deviceInfo, cmd = m.deviceInfo.Update(msg)
		return cmd
	case popupHotspotCreator:
		m.hotspotCreator, cmd = m.hotspotCreator.Update(msg)
		return cmd
	case popupProfileCreator:
		m.profileCreator, cmd = m.profileCreator.Update(msg)
		return cmd
	case popupProfileEditor:
		m.profileEditor, cmd = m.profileEditor.Update(msg)
		return cmd
	}
	return nil
}

func (m *MainModel) viewPopup() string {
	switch m.activePopup {
	case popupNo:
		return ""
	case popupConfirm:
		return m.confirm.View()
	case popupConnector:
		return m.connector.View()
	case popupHelp:
		return m.help.View()
	case popupDeviceInfo:
		return m.deviceInfo.View()
	case popupHotspotCreator:
		return m.hotspotCreator.View()
	case popupProfileCreator:
		return m.profileCreator.View()
	case popupProfileEditor:
		return m.profileEditor.View()
	}
	return ""
}

func (m *MainModel) updateOnKeyPress(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	if m.activePopup != popupNo {
		if key.Matches(msg, m.keys.closePopup) {
			return m, ClosePopupCmd()
		}
		cmd = m.updatePopup(msg)
		return m, cmd
	}
	switch {
	case key.Matches(msg, m.keys.quit):
		return m, tea.Quit
	case key.Matches(msg, m.keys.help):
		return m, OpenHelpCmd()
	}
	m.tabs, cmd = m.tabs.Update(msg)
	return m, cmd
}

func (m *MainModel) View() tea.View {
	if !m.ready {
		return tea.View{}
	}

	view := m.tabs.View()

	if m.activePopup != popupNo {
		popupView := m.viewPopup()
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
	view = m.Style.Render(view)
	v := tea.NewView(view)
	v.AltScreen = true
	return v
}

func (m *MainModel) Width() int { return m.Style.GetWidth() }

func (m *MainModel) Height() int { return m.Style.GetHeight() }

func (m *MainModel) Resize(width, height int) {
	m.Style = m.Style.MaxWidth(width).MaxHeight(height)

	border := m.Style.GetBorderStyle()
	width -= border.GetLeftSize() + border.GetRightSize()
	width -= border.GetBottomSize() + border.GetTopSize()

	helpHeight := lipgloss.Height(m.shortHelpView())

	m.tabs.Resize(width, height-helpHeight)
	m.help.ResizeShort(width)
	m.help.ResizeFull(int(float32(width)*0.8), int(float32(height)*0.8))

	m.deviceInfo.Resize(int(float32(width)*0.8), int(float32(height)*0.3))

	m.notification.style = m.notification.style.Width(width / 2)
}

func (m *MainModel) activeBindingsShort() []key.Binding {
	switch m.activePopup {
	case popupConnector:
		return m.help.connectorShort()
	case popupProfileCreator:
		return m.help.profileCreatorShort()
	case popupHotspotCreator:
		return m.help.hotspotCreatorShort()
	case popupProfileEditor:
		return m.help.profileEditorShort()
	case popupDeviceInfo:
		return m.help.deviceInfoShort()
	case popupHelp:
		return m.help.helpShort()
	default:
	}

	keys := m.help.mainShort()

	switch m.tabs.ActiveTabIndex() {
	case 1: // Device tab
		keys = append(keys, m.help.deviceShort()...)
		if m.device.netDevices.Focused() {
			keys = append(keys, m.help.netDevicesShort()...)
		}
		return keys
	default: // Networks tab: tab actions + focused window
		keys = append(keys, m.help.networksShort()...)
		if m.networks.available.Focused() {
			keys = append(keys, m.help.availableNetworksShort()...)
		} else {
			keys = append(keys, m.help.networkProfilesShort()...)
		}
		return keys
	}
}

func (m *MainModel) shortHelpView() string {
	return m.help.ShortViewFor(m.activeBindingsShort())
}

// NilMsg is a fictive struct, which used to send as tea.Msg instead of nil to trigger main window re-render.
type NilMsg struct{}

// NilCmd is a function, which returns fictive Msg to trigger Model Update.
func NilCmd() tea.Cmd {
	return func() tea.Msg {
		return NilMsg{}
	}
}

type IntervalRescanMsg struct{}

func IntervalRescanCmd(interval time.Duration) tea.Cmd {
	if interval <= 0 {
		return nil
	}
	return tea.Tick(interval, func(time.Time) tea.Msg {
		return IntervalRescanMsg{}
	})
}
