// Package renderer provides methods for help to rendering app components
package renderer

import (
	"fmt"
	"image/color"

	"charm.land/lipgloss/v2"
	"github.com/alphameo/nm-tui/internal/infra"
	"github.com/alphameo/nm-tui/internal/ui/tools/compositor"
)

func RenderWithTitleAndKeybind(view, title, keybind string, style lipgloss.Style, accentColor color.Color) string {
	view = style.Render(view)
	keybind = fmt.Sprintf("[%s]", keybind)
	keybindStyle := lipgloss.NewStyle().Foreground(style.GetBorderTopForeground())
	titleStyle := lipgloss.NewStyle().Foreground(accentColor)

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

func RenderNetworkModeIcon(mode infra.NetworkMode) string {
	switch mode {
	case infra.NetworkAccessPoint:
		return "󰀃"
	case infra.NetworkInfra:
		return "🖳"
	case infra.NetworkMesh:
		return ""
	case infra.NetworkAdHoc:
		return ""
	default:
		return "?"
	}
}
