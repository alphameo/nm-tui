package components

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/alphameo/nm-tui/internal/infra"
	"github.com/alphameo/nm-tui/internal/ui/styles"
	"github.com/alphameo/nm-tui/internal/ui/tools/renderer"
)

type tabsKeyMap struct {
	tabNext key.Binding
	tabPrev key.Binding
}

func (k *tabsKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.tabNext, k.tabPrev}
}

func (k *tabsKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.tabNext, k.tabPrev}}
}

var tabsKeys = &tabsKeyMap{
	tabNext: key.NewBinding(
		key.WithKeys("]"),
		key.WithHelp("]", "next tab"),
	),
	tabPrev: key.NewBinding(
		key.WithKeys("["),
		key.WithHelp("[", "previous tab"),
	),
}

type Tab struct {
	title   string
	content SizedModel
}

type TabsModel struct {
	tabs      []Tab
	tabTitles []string
	activeTab int

	innerStyle *lipgloss.Style

	keys *tabsKeyMap
}

func NewTabsModel(tabs []Tab, keys *tabsKeyMap, networkManager infra.WifiManager) *TabsModel {
	tabTitles := []string{}
	for _, t := range tabs {
		tabTitles = append(tabTitles, t.title)
	}
	m := &TabsModel{
		tabs:       tabs,
		tabTitles:  tabTitles,
		activeTab:  0,
		innerStyle: &styles.TabScreenBorderStyle,
		keys:       keys,
	}
	return m
}

func (m *TabsModel) Resize(width, height int) {
	height -= styles.TabBarHeight

	style := m.innerStyle.Width(width).Height(height)
	m.innerStyle = &style

	for _, t := range m.tabs {
		t.content.Resize(width, height)
	}
}

func (m *TabsModel) Init() tea.Cmd {
	var cmds []tea.Cmd
	for _, t := range m.tabs {
		cmds = append(cmds, t.content.Init())
	}
	return tea.Batch(cmds...)
}

func (m *TabsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		return m.handleKey(msg)
	}

	upd, cmd := m.tabs[m.activeTab].content.Update(msg)
	m.tabs[m.activeTab].content = upd.(SizedModel)
	return m, cmd
}

func (m *TabsModel) handleKey(keyMsg tea.KeyPressMsg) (*TabsModel, tea.Cmd) {
	switch {
	case key.Matches(keyMsg, m.keys.tabNext):
		m.activeTab = min(m.activeTab+1, len(m.tabs)-1)
		return m, m.tabs[m.activeTab].content.Init()
	case key.Matches(keyMsg, m.keys.tabPrev):
		m.activeTab = max(m.activeTab-1, 0)
		return m, m.tabs[m.activeTab].content.Init()
	}

	upd, cmd := m.tabs[m.activeTab].content.Update(keyMsg)
	m.tabs[m.activeTab].content = upd.(SizedModel)
	return m, cmd
}

func (m *TabsModel) View() tea.View {
	tabView := m.tabs[m.activeTab].content.View().Content
	tabView = m.innerStyle.Render(tabView)
	tabBar := renderer.RenderTabBar(
		m.tabTitles,
		styles.TabTabBorderActiveStyle,
		styles.TabTabBorderInactiveStyle,
		lipgloss.Width(tabView),
		m.activeTab,
	)

	return tea.NewView(lipgloss.JoinVertical(
		lipgloss.Center,
		tabBar,
		tabView,
	))
}
