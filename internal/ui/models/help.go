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
	keys     helpKeyMap
}

func NewHelpModel(keys keyMaps) *HelpModel {
	v := viewport.New(viewport.WithHeight(1), viewport.WithWidth(1))
	h := help.New()
	h.Styles = styles.HelpStyle
	h.Ellipsis = styles.SymbolEllipsis
	h.ShortSeparator = fmt.Sprintf(" %s ", styles.SymbolSeparator)

	help := HelpModel{
		viewport: v,
		keyMap:   keys,
		help:     h,
		keys:     keys.help,
	}
	help.viewport.SetContent(help.FullView())
	return &help
}

func (m *HelpModel) Resize(width, height int) {
	m.viewport.SetWidth(width - styles.BorderOffset)
	m.viewport.SetHeight(height - styles.BorderOffset)
}

func (m *HelpModel) Init() tea.Cmd {
	m.viewport.GotoTop()
	return nil
}

func (m *HelpModel) Update(msg tea.Msg) (*HelpModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.SetWidth(msg.Width - styles.BorderOffset)
		m.viewport.SetHeight(msg.Height - styles.BorderOffset)
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, m.keys.quit):
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

func (m *HelpModel) Short() []key.Binding {
	keys := []key.Binding{
		m.keyMap.main.quit,
		m.keyMap.main.closePopup,
	}
	return keys
}

func (m *HelpModel) ShortView() string {
	return m.help.ShortHelpView(m.Short())
}

func (m *HelpModel) ShortViewFor(bindings []key.Binding) string {
	return m.help.ShortHelpView(bindings)
}

func (m *HelpModel) FullView() string {
	var view string
	globalTTL := "Global"
	globalTTL = styles.AccentStyle.Render(globalTTL)
	global := m.globalFull()

	mainTTL := "Main"
	mainTTL = styles.AccentStyle.Render(mainTTL)
	main := m.mainFull()

	connectivityTTL := "Connectivity"
	connectivityTTL = styles.AccentStyle.Render(connectivityTTL)
	connectivity := m.connectivityFull()

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
		profileCreatorTTL, m.help.FullHelpView(profileCreator), "",
		hotspotCreatorTTL, m.help.FullHelpView(hotspotCreator), "",
		availableNetworksTTL, m.help.FullHelpView(availableNetworks), "",
		connectorTTL, m.help.FullHelpView(connector), "",
		networkProfilesTTL, m.help.FullHelpView(networkProfiles), "",
		profileEditorTTL, m.help.FullHelpView(profileEditor),
		connectivityTTL, m.help.FullHelpView(connectivity), "",
	)

	return view
}

func (m *HelpModel) globalFull() [][]key.Binding {
	return [][]key.Binding{{
		m.keyMap.toggle.Toggle,
	}}
}

func (m *HelpModel) mainFull() [][]key.Binding {
	return [][]key.Binding{{
		m.keyMap.main.quit,
		m.keyMap.tabs.Next,
		m.keyMap.tabs.Prev,
		m.keyMap.main.help,
	}}
}

func (m *HelpModel) connectivityFull() [][]key.Binding {
	return [][]key.Binding{{
		m.keyMap.connectivity.prev,
		m.keyMap.connectivity.next,
		m.keyMap.connectivity.rescan,
	}}
}

func (m *HelpModel) connectivityShort() []key.Binding {
	return []key.Binding{
		m.keyMap.connectivity.rescan,
	}
}

func (m *HelpModel) availableNetworksFull() [][]key.Binding {
	return [][]key.Binding{{
		m.keyMap.availableNetworks.rescan,
		m.keyMap.availableNetworks.connect,
	}}
}

func (m *HelpModel) availableNetworksShort() []key.Binding {
	return []key.Binding{
		m.keyMap.availableNetworks.rescan,
		m.keyMap.availableNetworks.connect,
	}
}

func (m *HelpModel) connectorFull() [][]key.Binding {
	return [][]key.Binding{{
		m.keyMap.connector.prev,
		m.keyMap.connector.next,
		m.keyMap.connector.togglePWVisibility,
		m.keyMap.connector.connect,
		m.keyMap.main.closePopup,
	}}
}

func (m *HelpModel) connectorShort() []key.Binding {
	return []key.Binding{
		m.keyMap.connector.togglePWVisibility,
		m.keyMap.connector.connect,
	}
}

func (m *HelpModel) hotspotCreatorFull() [][]key.Binding {
	return [][]key.Binding{{
		m.keyMap.hotspotCreator.prev,
		m.keyMap.hotspotCreator.next,
		m.keyMap.hotspotCreator.togglePWVisibility,
		m.keyMap.hotspotCreator.create,
		m.keyMap.main.closePopup,
	}}
}

func (m *HelpModel) hotspotCreatorShort() []key.Binding {
	return []key.Binding{
		m.keyMap.hotspotCreator.togglePWVisibility,
		m.keyMap.hotspotCreator.create,
	}
}

func (m *HelpModel) networkProfilesFull() [][]key.Binding {
	return [][]key.Binding{{
		m.keyMap.networkProfiles.connect,
		m.keyMap.networkProfiles.disconnect,
		m.keyMap.networkProfiles.edit,
		m.keyMap.networkProfiles.delete,
		m.keyMap.networkProfiles.rescan,
	}}
}

func (m *HelpModel) networkProfilesShort() []key.Binding {
	return []key.Binding{
		m.keyMap.networkProfiles.connect,
		m.keyMap.networkProfiles.disconnect,
		m.keyMap.networkProfiles.edit,
		m.keyMap.networkProfiles.delete,
		m.keyMap.networkProfiles.rescan,
	}
}

func (m *HelpModel) networksFull() [][]key.Binding {
	return [][]key.Binding{{
		m.keyMap.networks.win1,
		m.keyMap.networks.win2,
		m.keyMap.networks.winPrev,
		m.keyMap.networks.winNext,
		m.keyMap.networks.createProfile,
		m.keyMap.networks.createHotspot,
		m.keyMap.networks.quickHotspot,
		m.keyMap.networks.openCaptivePortal,
		m.keyMap.networks.rescan,
	}}
}

func (m *HelpModel) networksShort() []key.Binding {
	return []key.Binding{
		m.keyMap.networks.createProfile,
		m.keyMap.networks.createHotspot,
		m.keyMap.networks.quickHotspot,
		m.keyMap.networks.openCaptivePortal,
		m.keyMap.networks.rescan,
	}
}

func (m *HelpModel) profileCreatorFull() [][]key.Binding {
	return [][]key.Binding{{
		m.keyMap.profileCreator.prev,
		m.keyMap.profileCreator.next,
		m.keyMap.profileCreator.togglePWVisibility,
		m.keyMap.profileCreator.create,
		m.keyMap.main.closePopup,
	}}
}

func (m *HelpModel) profileCreatorShort() []key.Binding {
	return []key.Binding{
		m.keyMap.profileCreator.togglePWVisibility,
		m.keyMap.profileCreator.create,
	}
}

func (m *HelpModel) profileEditorFull() [][]key.Binding {
	return [][]key.Binding{{
		m.keyMap.profileEditor.prev,
		m.keyMap.profileEditor.next,
		m.keyMap.profileEditor.togglePWVisibility,
		m.keyMap.profileEditor.save,
		m.keyMap.main.closePopup,
	}}
}

func (m *HelpModel) profileEditorShort() []key.Binding {
	return []key.Binding{
		m.keyMap.profileEditor.togglePWVisibility,
		m.keyMap.profileEditor.save,
	}
}

func HelpFromKeys(keys ...string) string {
	transformed := make([]string, len(keys))
	for i, key := range keys {
		transformed[i] = strings.ReplaceAll(key, "ctrl+", "^")
	}
	return strings.Join(transformed, "/")
}

func NewKeyBinding(b key.Binding, help string) key.Binding {
	keys := b.Keys()
	helpKeys := HelpFromKeys()
	return key.NewBinding(
		key.WithKeys(keys...),
		key.WithHelp(helpKeys, help),
	)
}
