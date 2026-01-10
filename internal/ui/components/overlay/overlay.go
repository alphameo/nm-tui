// Package overlay provides simple overlay windows
package overlay

import (
	"fmt"
	"slices"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Anchor int

const (
	Begin Anchor = iota
	Center
	End
)

// Model contains any tea.Model inside
type Model struct {
	Content    tea.Model // Content of overlay
	IsActive   bool      // Flag for upper composition (Default: `false`)
	Width      int       // Set to positive if you want specific width (Default: `0`)
	Height     int       // Set to positive if you want specific height (Default: `0`)
	XAnchor    Anchor    // Start position (Default: `Begin` - very top)
	YAnchor    Anchor    // Start position (Default: `Begin` - very left)
	XOffset    int       // Counts from the `XAnchor` (Default: `0`)
	YOffset    int       // Counts from the `YAnchor` (Default: `0`)
	EscapeKeys []string  // Keycombinations for close
	Title      string
}

func (m Model) Init() tea.Cmd {
	if m.Content == nil {
		return nil
	}
	return m.Content.Init()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	keyMsg, err := msg.(tea.KeyMsg)
	if err {
		if slices.Contains(m.EscapeKeys, keyMsg.String()) {
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

func New(content tea.Model, title string) *Model {
	return &Model{
		Content: content,
		Title:   title,
	}
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
