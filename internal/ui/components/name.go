package components

import (
	"charm.land/bubbles/v2/textinput"
	"charm.land/lipgloss/v2"
	"github.com/alphameo/nm-tui/internal/ui/styles"
)

type Name struct {
	*textinput.Model
}

func NewName(input *textinput.Model) Name {
	return Name{Model: input}
}

func DefaultName() Name {
	input := textinput.New()
	input.SetWidth(20)
	input.Prompt = ""
	input.Placeholder = "Name"
	return NewName(&input)
}

func (m *Name) View() string {
	view := m.Model.View()
	var style lipgloss.Style
	if m.Focused() {
		style = styles.InputFieldFocusedStyle
	} else {
		style = styles.InputFieldStyle
	}
	view = style.Render(view)
	view = lipgloss.JoinHorizontal(
		lipgloss.Center,
		view,
	)

	return view
}
