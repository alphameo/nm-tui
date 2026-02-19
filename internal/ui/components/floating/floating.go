// Package floating provides floating window model
package floating

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Model contains any tea.Model inside
type Model struct {
	Title    string    // Overlay title
	Content  tea.Model // Content of overlay
	IsActive bool      // Flag for upper composition (Default: `false`)

	Keys *FloatingKeyMap // Keycombinations for overlay

	Width   int    // Set to positive if you want specific width (Default: `0`)
	Height  int    // Set to positive if you want specific height (Default: `0`)
	XAnchor Anchor // Start position (Default: `Begin` - very top)
	YAnchor Anchor // Start position (Default: `Begin` - very left)
	XOffset int    // Counts from the `XAnchor` (Default: `0`)
	YOffset int    // Counts from the `YAnchor` (Default: `0`)
}

func NewFloatingModel(content tea.Model, title string) *Model {
	return &Model{
		Content: content,
		Title:   title,
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
	if m.Content == nil {
		return ""
	}
	return m.Content.View()
}

func (m *Model) Place(bg string, fgStyle lipgloss.Style) string {
	if m.Width > 0 {
		fgStyle = fgStyle.Width(m.Width)
	}
	if m.Height > 0 {
		fgStyle = fgStyle.Height(m.Height)
	}

	fg := fgStyle.Render(m.View())
	title := lipgloss.NewStyle().
		Foreground(fgStyle.GetBorderTopForeground()).
		Render(fmt.Sprintf("[ %s ]", m.Title))

	fg = Compose(title, fg, Center, Begin, 0, 0)
	return Compose(fg, bg, m.XAnchor, m.YAnchor, m.XOffset, m.YOffset)
}
