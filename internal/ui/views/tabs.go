package views

import (
	"fmt"
	"strings"

	"github.com/alphameo/nm-tui/internal/infra"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TabsModel struct {
	tabTables []ResizeableModel
	tabTitles []string
	activeTab int
	width     int
	height    int
}

type ResizeableModel interface {
	tea.Model
	Resize(width, height int)
}

func NewConnectionsModel(networkManager infra.NetworkManager) *TabsModel {
	wifi := NewWifiModel(networkManager)
	tabTables := []ResizeableModel{wifi, wifi}
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
	height -= tabBarHeight
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
	m.tabTables[m.activeTab] = upd.(ResizeableModel)
	return m, cmd
}

func (m TabsModel) View() string {
	view := m.tabTables[m.activeTab].View()
	tabBar := ConstructTabBar(
		m.tabTitles,
		ActiveTabStyle,
		InactiveTabStyle,
		lipgloss.Width(view)+2,
		m.activeTab,
	)
	borderStyle := BorderStyle
	borderStyle.Top = ""
	borderStyle.TopLeft = "│"
	borderStyle.TopRight = "│"
	style := lipgloss.NewStyle().Border(borderStyle)
	styledView := style.Render(view)

	sb := strings.Builder{}
	fmt.Fprintf(&sb, "%s\n%s", tabBar, styledView)
	return sb.String()
}
