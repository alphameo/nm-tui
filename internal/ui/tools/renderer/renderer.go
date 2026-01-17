package renderer

import (
	"fmt"

	"github.com/alphameo/nm-tui/internal/ui/tools/compositor"
	"github.com/charmbracelet/lipgloss"
)

func RenderTabBar(
	titles []string,
	activeStyle,
	inactiveStyle lipgloss.Style,
	fullWidth int,
	active int,
) string {
	tabCount := len(titles)
	tabWidth := fullWidth/tabCount - 2
	tail := fullWidth % tabCount
	var renderedTabs []string
	for i, t := range titles {
		var style lipgloss.Style
		isFirst, isLast, isActive := i == 0, i == len(titles)-1, i == active
		if isActive {
			style = activeStyle
		} else {
			style = inactiveStyle
		}
		border, _, _, _, _ := style.GetBorder()
		if isFirst && isActive {
			border.BottomLeft = border.Left
		} else if isFirst && !isActive {
			border.BottomLeft = border.MiddleLeft
		} else if isLast && isActive {
			border.BottomRight = border.Right
		} else if isLast && !isActive {
			border.BottomRight = border.MiddleRight
		}
		style = style.Border(border)
		if tail > 0 {
			style = style.Width(tabWidth + 1)
			tail--
		} else {
			style = style.Width(tabWidth)
		}
		tabView := style.Render(t)
		renderedTabs = append(renderedTabs, tabView)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
}

func RenderWithTitleAndKeybind(view, title, keybind string, style *lipgloss.Style, accentColor lipgloss.TerminalColor) string {
	view = style.Render(view)
	keybind = fmt.Sprintf("[%s]", keybind)
	keybindStyle := lipgloss.NewStyle().Foreground(style.GetBorderTopForeground())
	titleStyle := lipgloss.NewStyle().Foreground(accentColor)

	title = titleStyle.Render(title)
	divider := keybindStyle.Render(style.GetBorderStyle().Top)
	keybind = keybindStyle.Render(keybind)
	extendedTitle := fmt.Sprintf("%s%s%s", keybind, divider, title)
	return compositor.PlaceTitle(view, extendedTitle)
}
