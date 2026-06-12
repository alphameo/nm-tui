package styles

import (
	"charm.land/bubbles/v2/textinput"
	"charm.land/lipgloss/v2"
	"github.com/alphameo/nm-tui/internal/ui/models/toggle"
)

func ViewInput(input *textinput.Model) string {
	view := input.View()
	var style lipgloss.Style
	if input.Focused() {
		style = InputFieldFocusedStyle
	} else {
		style = InputFieldStyle
	}
	view = style.Render(view)
	view = lipgloss.JoinHorizontal(
		lipgloss.Center,
		view,
	)

	return view
}

func ViewToggle(toggle *toggle.Model) string {
	view := toggle.View()
	var style lipgloss.Style
	if toggle.Focused() {
		style = ToggleFocusedStyle
	} else {
		style = ToggleStyle
	}
	view = style.Render(view)
	view = lipgloss.JoinHorizontal(
		lipgloss.Center,
		view,
	)

	return view
}

func ViewInputWithValidation(password *textinput.Model) string {
	view := password.View()
	var style lipgloss.Style
	errIndicator := " "
	if password.Err != nil {
		errIndicator = ErrorSymbolColored
	}
	if password.Focused() {
		style = InputFieldFocusedStyle
	} else {
		style = InputFieldStyle
	}
	view = style.Render(view)
	view = lipgloss.JoinHorizontal(
		lipgloss.Center,
		view,
		errIndicator,
	)

	return view
}
