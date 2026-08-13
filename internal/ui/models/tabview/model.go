// Package tabview provides model for tabbed view
package tabview

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type Model struct {
	activeTab int

	tabTitles   []string
	tabContents []TabModel

	cachedTabBarView string

	styles Styles

	borderOffset int
	tabBarHeight int

	Keys KeyMap
}

type Tab struct {
	Title   string
	Content TabModel
}

func New(tabs []Tab) Model {
	tabTitles := []string{}
	tabContents := []TabModel{}
	for _, t := range tabs {
		t.Content.Blur()
		tabTitles = append(tabTitles, t.Title)
		tabContents = append(tabContents, t.Content)
	}
	activeTab := 0
	tabContents[activeTab].Focus()
	return Model{
		tabTitles:   tabTitles,
		tabContents: tabContents,
		activeTab:   activeTab,
		Keys:        DefaultKeys(),
		styles:      DefaultStyles(),
	}
}

func (m *Model) SetStyles(styles Styles) {
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

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, m.Keys.Next):
			m.tabContents[m.activeTab].Blur()
			m.activeTab = min(m.activeTab+1, len(m.tabContents)-1)
			m.tabContents[m.activeTab].Focus()
			m.renderTabBar()
			return m, m.tabContents[m.activeTab].Init()
		case key.Matches(msg, m.Keys.Prev):
			m.tabContents[m.activeTab].Blur()
			m.activeTab = max(m.activeTab-1, 0)
			m.tabContents[m.activeTab].Focus()
			m.renderTabBar()
			return m, m.tabContents[m.activeTab].Init()
		}
	}

	var cmds []tea.Cmd
	for i := range m.tabContents {
		m.tabContents[i], cmd = m.tabContents[i].UpdateAsTab(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
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
