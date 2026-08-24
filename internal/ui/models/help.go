package models

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/alphameo/nm-tui/internal/ui/styles"
	"github.com/alphameo/nm-tui/internal/ui/tools/compositor"
	"github.com/alphameo/nm-tui/internal/ui/tools/renderer"
)

type helpKeyMap struct {
	quit key.Binding
}

type helpConfig struct {
	title string
}

var helpCfg = helpConfig{title: "Help"}

type HelpModel struct {
	viewport viewport.Model
	help     help.Model
	keyMap   keyMaps

	keys       helpKeyMap
	FullStyle  lipgloss.Style
	ShortStyle lipgloss.Style
}

func NewHelpModel(keys keyMaps) *HelpModel {
	v := viewport.New()
	h := help.New()
	h.Styles = styles.HelpStyles
	h.Ellipsis = styles.SymbolEllipsis
	h.ShortSeparator = fmt.Sprintf(" %s ", styles.SymbolSeparator)

	help := HelpModel{
		viewport:   v,
		keyMap:     keys,
		help:       h,
		keys:       keys.help,
		FullStyle:  lipgloss.NewStyle(),
		ShortStyle: lipgloss.NewStyle().MaxHeight(1),
	}
	help.viewport.SetContent(help.fullView())
	return &help
}

func (m *HelpModel) ResizeFull(width, height int) {
	m.FullStyle = m.FullStyle.Width(width).Height(height)

	border := m.FullStyle.GetBorderStyle()
	width -= border.GetLeftSize() + border.GetRightSize()
	height -= border.GetBottomSize() + border.GetTopSize()

	m.viewport.SetWidth(width)
	m.viewport.SetHeight(height)
}

func (m *HelpModel) ResizeShort(width int) {
	m.ShortStyle = m.ShortStyle.MaxWidth(width)

	border := m.ShortStyle.GetBorderStyle()
	width -= border.GetLeftSize() + border.GetRightSize()

	m.help.SetWidth(width)
}

func (m *HelpModel) Init() tea.Cmd {
	m.viewport.GotoTop()
	return nil
}

func (m *HelpModel) Update(msg tea.Msg) (*HelpModel, tea.Cmd) {
	if msg, ok := msg.(tea.KeyPressMsg); ok {
		if key.Matches(msg, m.keys.quit) {
			return m, ClosePopupCmd()
		}
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)

	return m, cmd
}

func (m *HelpModel) UpdateAsPopup(msg tea.Msg) (PopupModel, tea.Cmd) {
	return m.Update(msg)
}

func (m *HelpModel) View() string {
	view := m.viewport.View()
	title := styles.DefaultStyle.Render(renderer.RenderTitle(helpCfg.title))
	view = styles.OverlayStyle.Render(view)
	view = compositor.Compose(
		title,
		view,
		compositor.Center,
		compositor.Begin,
		0,
		0,
	)
	return view
}

func (m *HelpModel) ShortViewFor(bindings []key.Binding) string {
	return m.ShortStyle.Render(m.help.ShortHelpView(bindings))
}

func (m *HelpModel) fullView() string {
	var view string
	globalTTL := "Global"
	globalTTL = styles.AccentStyle.Render(globalTTL)
	global := m.globalFull()

	mainTTL := "Main"
	mainTTL = styles.AccentStyle.Render(mainTTL)
	main := m.mainFull()

	deviceTTL := "Device"
	deviceTTL = styles.AccentStyle.Render(deviceTTL)
	device := m.deviceFull()

	netDevicesTTL := "Network Devices"
	netDevicesTTL = styles.AccentStyle.Render(netDevicesTTL)
	netDevices := m.netDevicesFull()

	deviceInfoTTL := "Device Info"
	deviceInfoTTL = styles.AccentStyle.Render(deviceInfoTTL)
	deviceInfo := m.deviceInfoFull()

	availableNetworksTTL := "Available Networks"
	availableNetworksTTL = styles.AccentStyle.Render(availableNetworksTTL)
	availableNetworks := m.availableNetworksFull()

	connectorTTL := "Connector"
	connectorTTL = styles.AccentStyle.Render(connectorTTL)
	connector := m.connectorFull()

	hotspotCreatorTTL := "Hotspot Creator"
	hotspotCreatorTTL = styles.AccentStyle.Render(hotspotCreatorTTL)
	hotspotCreator := m.hotspotCreatorFull()

	networkProfilesTTL := "Network Profiles"
	networkProfilesTTL = styles.AccentStyle.Render(networkProfilesTTL)
	networkProfiles := m.networkProfilesFull()

	networksTTL := "Networks"
	networksTTL = styles.AccentStyle.Render(networksTTL)
	networks := m.networksFull()

	profileCreatorTTL := "Network Creator"
	profileCreatorTTL = styles.AccentStyle.Render(profileCreatorTTL)
	profileCreator := m.profileCreatorFull()

	profileEditorTTL := "Profile Editor"
	profileEditorTTL = styles.AccentStyle.Render(profileEditorTTL)
	profileEditor := m.profileEditorFull()

	view = lipgloss.JoinVertical(
		lipgloss.Left,
		view,
		globalTTL, m.help.FullHelpView(global), "",
		mainTTL, m.help.FullHelpView(main), "",
		networksTTL, m.help.FullHelpView(networks), "",
		netDevicesTTL, m.help.FullHelpView(netDevices), "",
		deviceInfoTTL, m.help.FullHelpView(deviceInfo), "",
		profileCreatorTTL, m.help.FullHelpView(profileCreator), "",
		hotspotCreatorTTL, m.help.FullHelpView(hotspotCreator), "",
		availableNetworksTTL, m.help.FullHelpView(availableNetworks), "",
		connectorTTL, m.help.FullHelpView(connector), "",
		networkProfilesTTL, m.help.FullHelpView(networkProfiles), "",
		profileEditorTTL, m.help.FullHelpView(profileEditor), "",
		deviceTTL, m.help.FullHelpView(device), "",
	)

	return view
}

func (m *HelpModel) globalFull() [][]key.Binding {
	return [][]key.Binding{{
		m.fullKB(m.keyMap.toggle.Toggle, "Enable/Disable toggle button"),
	}}
}

func (m *HelpModel) mainFull() [][]key.Binding {
	return [][]key.Binding{{
		m.fullKB(m.keyMap.main.quit, "Exit the application"),
		m.fullKB(m.keyMap.tabs.Next, "Move to next tab"),
		m.fullKB(m.keyMap.tabs.Prev, "Move to previous tab"),
		m.fullKB(m.keyMap.main.help, "Open/Close Help menu"),
	}}
}

func (m *HelpModel) mainShort() []key.Binding {
	return []key.Binding{m.keyMap.main.help}
}

func (m *HelpModel) deviceFull() [][]key.Binding {
	return [][]key.Binding{{
		m.fullKB(m.keyMap.device.prev, "Move to previous control"),
		m.fullKB(m.keyMap.device.next, "Move to next control"),
		m.fullKB(m.keyMap.device.rescan, "Rescan device state"),
	}}
}

func (m *HelpModel) deviceShort() []key.Binding {
	k := []key.Binding{
		m.keyMap.device.rescan,
	}
	return m.shortKBs(k)
}

func (m *HelpModel) netDevicesFull() [][]key.Binding {
	return [][]key.Binding{{
		m.fullKB(m.keyMap.networkDevices.showInfo, "Show information about selected device"),
	}}
}

func (m *HelpModel) netDevicesShort() []key.Binding {
	k := []key.Binding{
		m.keyMap.networkDevices.showInfo,
	}
	return m.shortKBs(k)
}

func (m *HelpModel) deviceInfoFull() [][]key.Binding {
	return [][]key.Binding{{
		m.fullKB(m.keyMap.main.closePopup, "Close popup with information about device"),
	}}
}

func (m *HelpModel) deviceInfoShort() []key.Binding {
	k := []key.Binding{
		m.keyMap.main.closePopup,
	}
	return m.shortKBs(k)
}

func (m *HelpModel) availableNetworksFull() [][]key.Binding {
	return [][]key.Binding{{
		m.fullKB(m.keyMap.availableNetworks.connect, "Open Connector for selected network"),
		m.fullKB(
			m.keyMap.availableNetworks.activate,
			"Activate the connection to the selected network by SSID."+
				"Use profile credentials if it exists, or create new one",
		),
		m.fullKB(
			m.keyMap.availableNetworks.deactivate,
			"Deactivate connection to the selected network if SSID matches profile name",
		),
	}}
}

func (m *HelpModel) availableNetworksShort() []key.Binding {
	k := []key.Binding{
		m.keyMap.availableNetworks.connect,
		m.keyMap.availableNetworks.activate,
		m.keyMap.availableNetworks.deactivate,
	}
	return m.shortKBs(k)
}

func (m *HelpModel) connectorFull() [][]key.Binding {
	return [][]key.Binding{{
		m.fullKB(m.keyMap.connector.prev, "Move to previous field"),
		m.fullKB(m.keyMap.connector.next, "Move to next field"),
		m.fullKB(m.keyMap.connector.togglePWVisibility, "Toggle password visibility"),
		m.fullKB(m.keyMap.connector.connect, "Connect with entered settings"),
		m.fullKB(m.keyMap.main.closePopup, "Close Connector"),
	}}
}

func (m *HelpModel) connectorShort() []key.Binding {
	k := []key.Binding{
		m.keyMap.connector.togglePWVisibility,
		m.keyMap.connector.connect,
		m.keyMap.main.closePopup,
	}
	return m.shortKBs(k)
}

func (m *HelpModel) hotspotCreatorFull() [][]key.Binding {
	return [][]key.Binding{{
		m.fullKB(m.keyMap.hotspotCreator.prev, "Move to previous field"),
		m.fullKB(m.keyMap.hotspotCreator.next, "Move to next field"),
		m.fullKB(m.keyMap.hotspotCreator.togglePWVisibility, "Toggle password visibility"),
		m.fullKB(m.keyMap.hotspotCreator.create, "Create hotspot profile with entered settings"),
		m.fullKB(m.keyMap.main.closePopup, "Close Hotspot Creator"),
	}}
}

func (m *HelpModel) hotspotCreatorShort() []key.Binding {
	k := []key.Binding{
		m.keyMap.hotspotCreator.togglePWVisibility,
		m.keyMap.hotspotCreator.create,
		m.keyMap.main.closePopup,
	}
	return m.shortKBs(k)
}

func (m *HelpModel) networkProfilesFull() [][]key.Binding {
	return [][]key.Binding{{
		m.fullKB(
			m.keyMap.networkProfiles.activate,
			"Activate connection to network associated with selected profile",
		),
		m.fullKB(
			m.keyMap.networkProfiles.deactivate,
			"Deactivate connection to network associated with selected profile",
		),
		m.fullKB(m.keyMap.networkProfiles.edit, "Open Profile Editor for selected profile"),
		m.fullKB(m.keyMap.networkProfiles.delete, "Delete network profile"),
	}}
}

func (m *HelpModel) networkProfilesShort() []key.Binding {
	k := []key.Binding{
		m.keyMap.networkProfiles.activate,
		m.keyMap.networkProfiles.deactivate,
		m.keyMap.networkProfiles.edit,
		m.keyMap.networkProfiles.delete,
	}
	return m.shortKBs(k)
}

func (m *HelpModel) networksFull() [][]key.Binding {
	return [][]key.Binding{{
		m.fullKB(m.keyMap.networks.win1, "Focus on 1st window"),
		m.fullKB(m.keyMap.networks.win2, "Focus on 2nd window"),
		m.fullKB(m.keyMap.networks.winPrev, "Focus on previous window"),
		m.fullKB(m.keyMap.networks.winNext, "Focus on next window"),
		m.fullKB(m.keyMap.networks.createProfile, "Open network Profile Creator"),
		m.fullKB(m.keyMap.networks.createHotspot, "Open Hotspot Creator"),
		m.fullKB(m.keyMap.networks.quickHotspot, "Enable hotspot, silently create its profile if not present"),
		m.fullKB(m.keyMap.networks.openCaptivePortal, "Open login (captive) portal in external browser"),
		m.fullKB(m.keyMap.networks.rescan, "Rescan networks"),
	}}
}

func (m *HelpModel) networksShort() []key.Binding {
	k := []key.Binding{
		m.keyMap.networks.createProfile,
		m.keyMap.networks.createHotspot,
		m.keyMap.networks.quickHotspot,
		m.keyMap.networks.openCaptivePortal,
		m.keyMap.networks.rescan,
	}
	return m.shortKBs(k)
}

func (m *HelpModel) profileCreatorFull() [][]key.Binding {
	return [][]key.Binding{{
		m.fullKB(m.keyMap.profileCreator.prev, "Move to previous field"),
		m.fullKB(m.keyMap.profileCreator.next, "Move to next field"),
		m.fullKB(m.keyMap.profileCreator.togglePWVisibility, "Toggle password visibility"),
		m.fullKB(m.keyMap.profileCreator.create, "Create network profile with entered settings"),
		m.fullKB(m.keyMap.main.closePopup, "Close Profile Creator"),
	}}
}

func (m *HelpModel) profileCreatorShort() []key.Binding {
	k := []key.Binding{
		m.keyMap.profileCreator.togglePWVisibility,
		m.keyMap.profileCreator.create,
		m.keyMap.main.closePopup,
	}
	return m.shortKBs(k)
}

func (m *HelpModel) profileEditorFull() [][]key.Binding {
	return [][]key.Binding{{
		m.fullKB(m.keyMap.profileEditor.prev, "Move to previous field"),
		m.fullKB(m.keyMap.profileEditor.next, "Move to next field"),
		m.fullKB(m.keyMap.profileEditor.togglePWVisibility, "Toggle password visibility"),
		m.fullKB(m.keyMap.profileEditor.save, "Save network profile settings"),
		m.fullKB(m.keyMap.main.closePopup, "Close Profile Editor"),
	}}
}

func (m *HelpModel) profileEditorShort() []key.Binding {
	k := []key.Binding{
		m.keyMap.profileEditor.togglePWVisibility,
		m.keyMap.profileEditor.save,
		m.keyMap.main.closePopup,
	}
	return m.shortKBs(k)
}

func (m *HelpModel) shortKB(kb key.Binding) key.Binding {
	keys := kb.Keys()
	desc := kb.Help().Desc
	helpKeys := m.compressKBKeys(keys...)
	kb.SetHelp(m.collectKeys(helpKeys...), desc)
	return kb
}

func (*HelpModel) compressKBKeys(keys ...string) []string {
	transformed := make([]string, len(keys))
	for i, key := range keys {
		transformed[i] = strings.ReplaceAll(key, "ctrl+", "^")
	}
	return transformed
}

func (*HelpModel) collectKeys(keys ...string) string {
	return strings.Join(keys, "/")
}

func (m *HelpModel) shortKBs(kbs []key.Binding) []key.Binding {
	shortKBs := make([]key.Binding, len(kbs))
	for i, v := range kbs {
		shortKBs[i] = m.shortKB(v)
	}
	return shortKBs
}

func (m *HelpModel) fullKB(kb key.Binding, desc string) key.Binding {
	helpKey := kb.Help().Key
	kb.SetHelp(helpKey, desc)
	return kb
}
