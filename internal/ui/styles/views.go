package styles

import (
	"charm.land/bubbles/v2/textinput"
	"charm.land/lipgloss/v2"
	"github.com/alphameo/nm-tui/internal/ui/models/toggle"
)

type focusable interface {
	View() string
	Focused() bool
}

func ViewBorderedFocusable(component focusable) string {
	view := component.View()
	var style lipgloss.Style
	if component.Focused() {
		style = BorderedFocusedStyle
	} else {
		style = BorderedStyle
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
		style = BorderedFocusedStyle
	} else {
		style = BorderedStyle
	}
	view = style.Render(view)
	view = lipgloss.JoinHorizontal(
		lipgloss.Center,
		view,
		errIndicator,
	)

	return view
}
