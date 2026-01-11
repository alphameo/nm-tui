package views

import (
	"fmt"
	"slices"

	"github.com/alphameo/nm-tui/internal/ui/tools/compositor"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// FloatingModel contains any tea.FloatingModel inside
type FloatingModel struct {
	Content    tea.Model      // Content of overlay
	IsActive   bool           // Flag for upper composition (Default: `false`)
	Width      int            // Set to positive if you want specific width (Default: `0`)
	Height     int            // Set to positive if you want specific height (Default: `0`)
	XAnchor    compositor.Anchor // Start position (Default: `Begin` - very top)
	YAnchor    compositor.Anchor // Start position (Default: `Begin` - very left)
	XOffset    int            // Counts from the `XAnchor` (Default: `0`)
	YOffset    int            // Counts from the `YAnchor` (Default: `0`)
	EscapeKeys []string       // Keycombinations for close
	Title      string
}

func NewFloatingModel(content tea.Model, title string) *FloatingModel {
	return &FloatingModel{
		Content: content,
		Title:   title,
	}
}

func (m FloatingModel) Init() tea.Cmd {
	if m.Content == nil {
		return nil
	}
	return m.Content.Init()
}

func (m FloatingModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m FloatingModel) View() string {
	if m.Content == nil {
		return ""
	}
	return m.Content.View()
}

func (m *FloatingModel) Place(bg string, fgStyle lipgloss.Style) string {
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

	fg = compositor.Compose(title, fg, compositor.Center, compositor.Begin, 0, 0)
	return compositor.Compose(fg, bg, m.XAnchor, m.YAnchor, m.XOffset, m.YOffset)
}
