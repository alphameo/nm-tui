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

type helpKeyMap struct {
	quit key.Binding
}

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
	view = styles.BorderedFocusedStyle.Render(view)
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
		m.keyMap.toggle.Toggle,
	}}
}

func (m *HelpModel) mainFull() [][]key.Binding {
	return [][]key.Binding{{
		m.keyMap.main.quit,
		m.keyMap.tabs.Next,
		m.keyMap.tabs.Prev,
	}}
}

func (m *HelpModel) connectivityFull() [][]key.Binding {
	return [][]key.Binding{{
		m.keyMap.connectivity.prev,
		m.keyMap.connectivity.next,
		m.keyMap.connectivity.rescan,
	}}
}

func (m *HelpModel) availableNetworksFull() [][]key.Binding {
	return [][]key.Binding{{
		m.keyMap.availableNetworks.rescan,
		m.keyMap.availableNetworks.connect,
	}}
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

func (m *HelpModel) hotspotCreatorFull() [][]key.Binding {
	return [][]key.Binding{{
		m.keyMap.hotspotCreator.prev,
		m.keyMap.hotspotCreator.next,
		m.keyMap.hotspotCreator.togglePWVisibility,
		m.keyMap.hotspotCreator.create,
		m.keyMap.main.closePopup,
	}}
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

func (m *HelpModel) profileCreatorFull() [][]key.Binding {
	return [][]key.Binding{{
		m.keyMap.profileCreator.prev,
		m.keyMap.profileCreator.next,
		m.keyMap.profileCreator.togglePWVisibility,
		m.keyMap.profileCreator.create,
		m.keyMap.main.closePopup,
	}}
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
