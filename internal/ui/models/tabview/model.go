// Package tabview provides model for tabbed view
package tabview

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// TabOps is a set of operations bound to a concrete tab content. Erasing the
// concrete type at the tabview boundary lets tabview hold heterogeneous tabs
// while each tab's Update still returns its exact type.
type TabOps struct {
	init    func() tea.Cmd
	update  func(msg tea.Msg) tea.Cmd
	view    func() string
	resize  func(width, height int)
	width   func() int
	focused func() bool
	focus   func()
	blur    func()
}

// tabModel describes the operations a tab content must provide to be bound
// into a tabview. It is used only as a compile-time bound by Bind, which routes
// Update (returning the exact concrete type T) through a closure.
type tabModel[T any] interface {
	Init() tea.Cmd
	Update(msg tea.Msg) (T, tea.Cmd)
	View() string
	Resize(width, height int)
	Width() int
	Focused() bool
	Focus()
	Blur()
}

// Bind binds a concrete tab content model into a TabOps. Update is dispatched
// through the content's Update, which returns the exact concrete type T.
func Bind[T tabModel[T]](m *T) TabOps {
	return TabOps{
		init: (*m).Init,
		update: func(msg tea.Msg) tea.Cmd {
			_, cmd := (*m).Update(msg)
			return cmd
		},
		view:    (*m).View,
		resize:  (*m).Resize,
		width:   (*m).Width,
		focused: (*m).Focused,
		focus:   (*m).Focus,
		blur:    (*m).Blur,
	}
}

type Model struct {
	activeTab int

	tabTitles []string
	tabOps    []TabOps

	cachedTabBarView string

	styles Styles

	tabBarHeight int

	Keys KeyMap
}

type Tab struct {
	Title   string
	Content TabOps
}

func New(tabs []Tab) Model {
	tabTitles := []string{}
	tabOps := []TabOps{}
	for _, t := range tabs {
		t.Content.blur()
		tabTitles = append(tabTitles, t.Title)
		tabOps = append(tabOps, t.Content)
	}
	activeTab := 0
	tabOps[activeTab].focus()
	return Model{
		tabTitles: tabTitles,
		tabOps:    tabOps,
		activeTab: activeTab,
		Keys:      DefaultKeys(),
		styles:    DefaultStyles(),
	}
}

func (m *Model) SetStyles(styles Styles) {
	border := styles.ActiveTabStyle.GetBorderStyle()
	tabBarHeight := border.GetBottomSize() + border.GetTopSize() + 1
	m.styles = styles
	m.tabBarHeight = tabBarHeight
}

func (m *Model) Resize(width, height int) {
	height -= m.tabBarHeight

	for _, op := range m.tabOps {
		op.resize(width, height)
	}

	m.renderTabBar()
}

func (m *Model) Init() tea.Cmd {
	var cmds []tea.Cmd
	for _, op := range m.tabOps {
		cmds = append(cmds, op.init())
	}
	return tea.Batch(cmds...)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyPressMsg); ok {
		switch {
		case key.Matches(msg, m.Keys.Next):
			m.tabOps[m.activeTab].blur()
			m.activeTab = min(m.activeTab+1, len(m.tabOps)-1)
			m.tabOps[m.activeTab].focus()
			m.renderTabBar()
			return m, m.tabOps[m.activeTab].init()
		case key.Matches(msg, m.Keys.Prev):
			m.tabOps[m.activeTab].blur()
			m.activeTab = max(m.activeTab-1, 0)
			m.tabOps[m.activeTab].focus()
			m.renderTabBar()
			return m, m.tabOps[m.activeTab].init()
		}
	}

	var cmds []tea.Cmd
	for _, op := range m.tabOps {
		cmds = append(cmds, op.update(msg))
	}
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	tabView := m.tabOps[m.activeTab].view()

	return lipgloss.JoinVertical(
		lipgloss.Center,
		m.cachedTabBarView,
		tabView,
	)
}

func (m *Model) ActiveTabIndex() int { return m.activeTab }

func (m *Model) renderTabBar() {
	width := m.tabOps[m.activeTab].width()
	m.cachedTabBarView = RenderTabBar(
		m.tabTitles,
		m.styles.ActiveTabStyle,
		m.styles.InactiveTabStyle,
		width,
		m.activeTab,
	)
}
