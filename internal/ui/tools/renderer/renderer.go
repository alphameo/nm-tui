// Package renderer provides methods for help to rendering app components
package renderer

import (
	"fmt"

	"charm.land/lipgloss/v2"
	"github.com/alphameo/nm-tui/internal/ui/tools/compositor"
)

func RenderWithTitleAndKeybind(view, title, keybind string, style lipgloss.Style) string {
	view = style.Render(view)
	keybind = fmt.Sprintf("[%s]", keybind)
	keybindStyle := lipgloss.NewStyle().Foreground(style.GetBorderTopForeground()).Background(style.GetBorderBottomBackground())
	titleStyle := lipgloss.NewStyle().Foreground(style.GetBorderTopForeground()).Background(style.GetBorderTopBackground())

	title = titleStyle.Render(title)
	divider := keybindStyle.Render(style.GetBorderStyle().Top)
	keybind = keybindStyle.Render(keybind)
	extendedTitle := fmt.Sprintf("%s%s%s", keybind, divider, title)
	return compositor.Compose(
		extendedTitle,
		view,
		compositor.Begin,
		compositor.Begin,
		2,
		0,
	)
}

func RenderEnabledStatus(value bool) string {
	if value {
		return "Enabled"
	}
	return "Disabled"
}

func RenderTitle(title string) string {
	return fmt.Sprintf("[ %s ]", title)
}
