// Package tabview provides model for tabbed view
package tabview

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type Model struct {
	activeTab int

	tabTitles        []string
	tabContents      []TabModel
	cachedTabBarView string

	styles       *Styles
	borderOffset int
	tabBarHeight int

	keys *KeyMap
}

type Tab struct {
	Title   string
	Content TabModel
}

func New(tabs []Tab, styles *Styles, keys *KeyMap) *Model {
	tabTitles := []string{}
	tabContents := []TabModel{}
	for _, t := range tabs {
		tabTitles = append(tabTitles, t.Title)
		tabContents = append(tabContents, t.Content)
	}
	m := &Model{
		tabTitles:   tabTitles,
		tabContents: tabContents,
		activeTab:   0,
		keys:        keys,
	}
	m.SetStyles(styles)
	return m
}

func (m *Model) SetStyles(styles *Styles) {
	if styles == nil {
		styles = DefaulStyles()
	}
	borderOffset := lipgloss.Width(styles.ContentStyle.GetBorderStyle().Left) * 2
	tabBarHeight := borderOffset + 1
	m.styles = styles
	m.borderOffset = borderOffset
	m.tabBarHeight = tabBarHeight
}

func (m *Model) Resize(width, height int) {
	height -= m.tabBarHeight

	m.styles.ContentStyle = m.styles.ContentStyle.Width(width).Height(height)

	m.renderTabBar()

	width -= m.borderOffset
	height -= m.borderOffset
	for _, t := range m.tabContents {
		t.Resize(width, height)
	}
}

func (m *Model) Init() tea.Cmd {
	var cmds []tea.Cmd
	for _, t := range m.tabContents {
		cmds = append(cmds, t.Init())
	}
	return tea.Batch(cmds...)
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		return m.handleKey(msg)
	}

	var cmd tea.Cmd
	m.tabContents[m.activeTab], cmd = m.tabContents[m.activeTab].UpdateAsTab(msg)
	return m, cmd
}

func (m *Model) handleKey(keyMsg tea.KeyPressMsg) (*Model, tea.Cmd) {
	switch {
	case key.Matches(keyMsg, m.keys.TabNext):
		m.activeTab = min(m.activeTab+1, len(m.tabContents)-1)
		m.renderTabBar()
		return m, m.tabContents[m.activeTab].Init()
	case key.Matches(keyMsg, m.keys.TabPrev):
		m.activeTab = max(m.activeTab-1, 0)
		m.renderTabBar()
		return m, m.tabContents[m.activeTab].Init()
	}

	var cmd tea.Cmd
	m.tabContents[m.activeTab], cmd = m.tabContents[m.activeTab].UpdateAsTab(keyMsg)
	return m, cmd
}

func (m *Model) View() string {
	tabView := m.tabContents[m.activeTab].View()
	tabView = m.styles.ContentStyle.Render(tabView)

	return lipgloss.JoinVertical(
		lipgloss.Center,
		m.cachedTabBarView,
		tabView,
	)
}

func (m *Model) renderTabBar() {
	width := m.styles.ContentStyle.GetWidth()
	m.cachedTabBarView = RenderTabBar(
		m.tabTitles,
		m.styles.ActiveTabBarStyle,
		m.styles.InactiveTabBarStyle,
		width,
		m.activeTab,
	)
}
