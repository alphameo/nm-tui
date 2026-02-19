// Package floating provides floating window model
package floating

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Model contains any tea.Model inside
type Model struct {
	Title    string
	Content  tea.Model // Content of overlay
	IsActive bool      // Flag for upper composition (Default: `false`)

	Keys *FloatingKeyMap // Keycombinations for overlay

	ContentAlignHorizontal lipgloss.Position // Horizontal position of content
	ContentAlignVertical   lipgloss.Position // Vertical position of content

	Width   int    // Set to positive if you want specific width (Default: `0`)
	Height  int    // Set to positive if you want specific height (Default: `0`)
	XAnchor Anchor // Start position (Default: `Begin` - very top)
	YAnchor Anchor // Start position (Default: `Begin` - very left)
	XOffset int    // Counts from the `XAnchor` (Default: `0`)
	YOffset int    // Counts from the `YAnchor` (Default: `0`)
}

func New(content tea.Model) *Model {
	return &Model{
		Content: content,
		Keys:    defaultFloatingKeys,
	}
}

func (m Model) Init() tea.Cmd {
	if m.Content == nil {
		return nil
	}
	return m.Content.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, m.Keys.Quit) {
			m.IsActive = false
			return m, nil
		}
	}

	m.Content, cmd = m.Content.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	style := lipgloss.NewStyle()
	style = style.Align(m.ContentAlignHorizontal, m.ContentAlignVertical)
	if m.Width > 0 {
		style = style.Width(m.Width)
	}
	if m.Height > 0 {
		style = style.Height(m.Height)
	}
	var view string
	if m.Content != nil {
		view = m.Content.View()
	}

	return style.Render(view)
}

func (m *Model) Place(bg string, style lipgloss.Style) string {
	view := m.View()
	view = style.Render(view)
	view = Compose(m.Title, view, Center, Begin, 0, 0)
	return Compose(view, bg, m.XAnchor, m.YAnchor, m.XOffset, m.YOffset)
}
