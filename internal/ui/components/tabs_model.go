package components

import (
	"github.com/alphameo/nm-tui/internal/infra"
	"github.com/alphameo/nm-tui/internal/ui/styles"
	"github.com/alphameo/nm-tui/internal/ui/tools/renderer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TabsModel struct {
	tabTables []SizedModel
	tabTitles []string
	activeTab int
	width     int
	height    int
}

func NewConnectionsModel(networkManager infra.NetworkManager) *TabsModel {
	wifi := NewWifiModel(networkManager)
	tabTables := []SizedModel{wifi, wifi}
	tabTitles := &[]string{"Wi-Fi", "VPN"}
	m := &TabsModel{
		tabTables: tabTables,
		tabTitles: *tabTitles,
		activeTab: 0,
	}
	return m
}

func (m *TabsModel) Resize(width, height int) {
	m.width = width
	m.height = height

	width -= BorderOffset
	height -= BorderOffset

	height -= TabBarHeight

	for _, t := range m.tabTables {
		t.Resize(width, height)
	}
}

func (m TabsModel) Init() tea.Cmd {
	var cmds []tea.Cmd
	for _, t := range m.tabTables {
		cmds = append(cmds, t.Init())
	}
	return tea.Batch(cmds...)
}

func (m TabsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "]":
			m.activeTab = min(m.activeTab+1, len(m.tabTitles)-1)
			return m, m.tabTables[m.activeTab].Init()
		case "[":
			m.activeTab = max(m.activeTab-1, 0)
			return m, m.tabTables[m.activeTab].Init()
		}
	}

	upd, cmd := m.tabTables[m.activeTab].Update(msg)
	m.tabTables[m.activeTab] = upd.(SizedModel)
	return m, cmd
}

func (m TabsModel) View() string {
	tabView := m.tabTables[m.activeTab].View()
	tabBar := renderer.RenderTabBar(
		m.tabTitles,
		styles.TabTabBorderActiveStyle,
		styles.TabTabBorderInactiveStyle,
		lipgloss.Width(tabView)+2,
		m.activeTab,
	)
	tabView = styles.TabScreenBorderStyle.Render(tabView)

	return lipgloss.JoinVertical(lipgloss.Center, tabBar, tabView)
}
