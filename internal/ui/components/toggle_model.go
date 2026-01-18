package components

import (
	"github.com/alphameo/nm-tui/internal/ui/styles"
	"github.com/alphameo/nm-tui/internal/ui/tools/renderer"
	tea "github.com/charmbracelet/bubbletea"
)

type ToggleModel struct {
	value   bool
	focused bool
}

func NewToggleModel(initial bool) *ToggleModel {
	return &ToggleModel{value: initial}
}

func (t *ToggleModel) SetValue(value bool) {
	t.value = value
}

func (t *ToggleModel) Value() bool {
	return t.value
}

func (t *ToggleModel) Init() tea.Cmd {
	return nil
}

func (t *ToggleModel) Update(msg tea.Msg) (*ToggleModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == " " {
			t.value = !t.value
		}
	}
	return t, nil
}

func (t *ToggleModel) View() string {
	checkbox := renderer.RenderCheckbox(t.value)
	if t.focused {
		checkbox = styles.DefaultStyle.Foreground(styles.AccentColor).Render(checkbox)
	}
	return checkbox
}

func (t *ToggleModel) Focus() tea.Cmd {
	t.focused = true
	return nil
}

func (t *ToggleModel) Blur() {
	t.focused = false
}
