// Package connections provides tabbed tables with main iformation about variable connections
package connections

import (
	"fmt"
	"strings"

	"github.com/alphameo/nm-tui/internal/ui/styles"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	tabTables []ResizeableModel
	tabTitles []string
	activeTab int
}

type ResizeableModel interface {
	Resize(width, height int)
	Init() tea.Cmd
	Update(tea.Msg) (tea.Model, tea.Cmd)
	View() string
}

func New() *Model {
	wifiAvailable := NewWifiAvailable()
	wifiStored := NewWifiStored()
	tabTables := []ResizeableModel{wifiAvailable, wifiStored}
	tabTitles := &[]string{"Current", "Stored"}
	m := &Model{
		tabTables: tabTables,
		tabTitles: *tabTitles,
		activeTab: 0,
	}
	return m
}

func (m *Model) Resize(width, height int) {
	height -= styles.TabBarHeight
	for _, t := range m.tabTables {
		t.Resize(width, height)
	}
}

func (m Model) Init() tea.Cmd {
	var cmds []tea.Cmd
	for _, t := range m.tabTables {
		cmds = append(cmds, t.Init())
	}
	return tea.Batch(cmds...)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "]", "tab":
			m.activeTab = min(m.activeTab+1, len(m.tabTitles)-1)
			return m, m.tabTables[m.activeTab].Init()
		case "[", "shift+tab":
			m.activeTab = max(m.activeTab-1, 0)
			return m, m.tabTables[m.activeTab].Init()
		}
	}

	upd, cmd := m.tabTables[m.activeTab].Update(msg)
	m.tabTables[m.activeTab] = upd.(ResizeableModel)
	return m, cmd
}

func (m Model) View() string {
	view := m.tabTables[m.activeTab].View()
	tabBar := styles.ConstructTabBar(
		m.tabTitles,
		styles.ActiveTabStyle,
		styles.InactiveTabStyle,
		lipgloss.Width(view)+2,
		m.activeTab,
	)
	borderStyle := styles.BorderStyle
	borderStyle.Top = ""
	borderStyle.TopLeft = "│"
	borderStyle.TopRight = "│"
	style := lipgloss.NewStyle().Border(borderStyle)
	styledView := style.Render(view)

	sb := strings.Builder{}
	fmt.Fprintf(&sb, "%s\n%s", tabBar, styledView)
	return sb.String()
}
