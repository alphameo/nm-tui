// Package tabview provides model for tabbed view
package tabview

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/alphameo/nm-tui/internal/ui/styles"
)

type Model struct {
	tabs []Tab

	activeTab int

	cachedTabTitles   []string
	cachedTtabBarView string

	innerStyle *lipgloss.Style

	keys *KeyMap
}

func New(tabs []Tab, keys *KeyMap) *Model {
	tabTitles := []string{}
	for _, t := range tabs {
		tabTitles = append(tabTitles, t.Title)
	}
	m := &Model{
		tabs:            tabs,
		cachedTabTitles: tabTitles,
		activeTab:       0,
		innerStyle:      &styles.TabScreenBorderStyle,
		keys:            keys,
	}
	return m
}

func (m *Model) Resize(width, height int) {
	height -= styles.TabBarHeight

	style := m.innerStyle.Width(width).Height(height)
	m.innerStyle = &style

	m.cachedTtabBarView = RenderTabBar(
		m.cachedTabTitles,
		styles.TabTabBorderActiveStyle,
		styles.TabTabBorderInactiveStyle,
		width,
		m.activeTab,
	)

	width -= styles.BorderOffset
	height -= styles.BorderOffset
	for _, t := range m.tabs {
		t.Content.Resize(width, height)
	}
}

func (m *Model) Init() tea.Cmd {
	var cmds []tea.Cmd
	for _, t := range m.tabs {
		cmds = append(cmds, t.Content.Init())
	}
	return tea.Batch(cmds...)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		return m.handleKey(msg)
	}

	upd, cmd := m.tabs[m.activeTab].Content.Update(msg)
	m.tabs[m.activeTab].Content = upd.(SizedModel)
	return m, cmd
}

func (m *Model) handleKey(keyMsg tea.KeyPressMsg) (*Model, tea.Cmd) {
	switch {
	case key.Matches(keyMsg, m.keys.TabNext):
		m.activeTab = min(m.activeTab+1, len(m.tabs)-1)
		return m, m.tabs[m.activeTab].Content.Init()
	case key.Matches(keyMsg, m.keys.TabPrev):
		m.activeTab = max(m.activeTab-1, 0)
		return m, m.tabs[m.activeTab].Content.Init()
	}

	upd, cmd := m.tabs[m.activeTab].Content.Update(keyMsg)
	m.tabs[m.activeTab].Content = upd.(SizedModel)
	return m, cmd
}

func (m *Model) View() tea.View {
	tabView := m.tabs[m.activeTab].Content.View().Content
	tabView = m.innerStyle.Render(tabView)

	return tea.NewView(lipgloss.JoinVertical(
		lipgloss.Center,
		m.cachedTtabBarView,
		tabView,
	))
}
