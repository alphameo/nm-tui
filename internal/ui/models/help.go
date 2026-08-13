package models

import (
	"log/slog"
	"strings"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/alphameo/nm-tui/internal/ui/styles"
)

type HelpModel struct {
	viewport viewport.Model
	help     help.Model
	keys     keyMaps
}

func NewHelpModel(keys keyMaps) *HelpModel {
	v := viewport.New(viewport.WithHeight(1), viewport.WithWidth(1))
	h := help.New()
	h.Styles = styles.HelpStyle

	help := HelpModel{
		viewport: v,
		keys:     keys,
		help:     h,
	}
	help.viewport.SetContent(help.FullView())
	return &help
}

func (m *HelpModel) Resize(width, height int) {
	m.viewport.SetWidth(width - styles.BorderOffset)
	m.viewport.SetHeight(height - styles.BorderOffset)
}

func (m *HelpModel) Init() tea.Cmd {
	return nil
}

func (m *HelpModel) Update(msg tea.Msg) (*HelpModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.SetWidth(msg.Width - styles.BorderOffset)
		m.viewport.SetHeight(msg.Height - styles.BorderOffset)
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
	view = styles.BorderedFocusedStyle.Render(view)
	return view
}

func (m *HelpModel) Short() []key.Binding {
	keys := []key.Binding{
		m.keys.main.quit,
		m.keys.main.closePopup,
	}
	return keys
}

func (m *HelpModel) ShortView() string {
	return m.help.ShortHelpView(m.Short())
}

func (m *HelpModel) FullView() string {
	var view string
	global := m.globalFull()
	main := m.mainFull()
	connectivity := m.connectivityFull()
	availableNetworks := m.availableNetworksFull()
	connector := m.connectorFull()
	hotspotCreator := m.hotspotCreatorFull()
	networkProfiles := m.networkProfilesFull()
	networks := m.networksFull()
	profileCreator := m.profileCreatorFull()
	profileEditor := m.profileEditorFull()

	view = lipgloss.JoinVertical(
		lipgloss.Left,
		view,
		"Global", m.help.FullHelpView(global),
		"",
		"Main", m.help.FullHelpView(main),
		"",
		"Connectivity", m.help.FullHelpView(connectivity),
		"",
		"Available Networks", m.help.FullHelpView(availableNetworks),
		"",
		"Connector", m.help.FullHelpView(connector),
		"",
		"Hotspot Creator", m.help.FullHelpView(hotspotCreator),
		"",
		"Network Profiles", m.help.FullHelpView(networkProfiles),
		"",
		"Networks", m.help.FullHelpView(networks),
		"",
		"Profile Creator", m.help.FullHelpView(profileCreator),
		"",
		"Profile Editor", m.help.FullHelpView(profileEditor),
	)

	slog.Error(view)
	return view
}

func (m *HelpModel) globalFull() [][]key.Binding {
	return [][]key.Binding{{
		m.keys.toggle.Toggle,
	}}
}

func (m *HelpModel) mainFull() [][]key.Binding {
	return [][]key.Binding{{
		m.keys.main.quit,
		m.keys.tabs.Next,
		m.keys.tabs.Prev,
	}}
}

func (m *HelpModel) connectivityFull() [][]key.Binding {
	return [][]key.Binding{{
		m.keys.connectivity.prev,
		m.keys.connectivity.next,
		m.keys.connectivity.rescan,
	}}
}

func (m *HelpModel) availableNetworksFull() [][]key.Binding {
	return [][]key.Binding{{
		m.keys.availableNetworks.rescan,
		m.keys.availableNetworks.connect,
	}}
}

func (m *HelpModel) connectorFull() [][]key.Binding {
	return [][]key.Binding{{
		m.keys.connector.prev,
		m.keys.connector.next,
		m.keys.connector.togglePWVisibility,
		m.keys.connector.connect,
		m.keys.main.closePopup,
	}}
}

func (m *HelpModel) hotspotCreatorFull() [][]key.Binding {
	return [][]key.Binding{{
		m.keys.hotspotCreator.prev,
		m.keys.hotspotCreator.next,
		m.keys.hotspotCreator.togglePWVisibility,
		m.keys.hotspotCreator.create,
		m.keys.main.closePopup,
	}}
}

func (m *HelpModel) networkProfilesFull() [][]key.Binding {
	return [][]key.Binding{{
		m.keys.networkProfiles.connect,
		m.keys.networkProfiles.disconnect,
		m.keys.networkProfiles.edit,
		m.keys.networkProfiles.delete,
		m.keys.networkProfiles.rescan,
	}}
}

func (m *HelpModel) networksFull() [][]key.Binding {
	return [][]key.Binding{{
		m.keys.networks.win1,
		m.keys.networks.win2,
		m.keys.networks.winPrev,
		m.keys.networks.winNext,
		m.keys.networks.createProfile,
		m.keys.networks.createHotspot,
		m.keys.networks.quickHotspot,
		m.keys.networks.openCaptivePortal,
		m.keys.networks.rescan,
	}}
}

func (m *HelpModel) profileCreatorFull() [][]key.Binding {
	return [][]key.Binding{{
		m.keys.profileCreator.prev,
		m.keys.profileCreator.next,
		m.keys.profileCreator.togglePWVisibility,
		m.keys.profileCreator.create,
		m.keys.main.closePopup,
	}}
}

func (m *HelpModel) profileEditorFull() [][]key.Binding {
	return [][]key.Binding{{
		m.keys.profileEditor.prev,
		m.keys.profileEditor.next,
		m.keys.profileEditor.togglePWVisibility,
		m.keys.profileEditor.save,
		m.keys.main.closePopup,
	}}
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
