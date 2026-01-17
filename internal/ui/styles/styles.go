// Package styles provides styling for ui
package styles

import (
	"fmt"

	"github.com/alphameo/nm-tui/internal/ui/tools/compositor"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

var (
	TextColor   = lipgloss.Color("#cbcbcb")
	AccentColor = lipgloss.Color("#865fff")
	MutedColor  = lipgloss.Color("#595959")
	RedColor    = lipgloss.Color("#ff0000")

	DefaultStyle = lipgloss.NewStyle().Foreground(TextColor)

	Border        = lipgloss.RoundedBorder()
	BorderedStyle = DefaultStyle.Border(Border)

	TableStyle = tableStyle()

	TabTabBorderInactive      = tabBorderInactive(Border)
	TabTabBorderActive        = tabBorderActive(Border)
	TabScreenBorder           = tabbedViewBorder(Border)
	TabScreenBorderStyle      = DefaultStyle.Border(TabScreenBorder)
	TabTabBorderInactiveStyle = DefaultStyle.Border(TabTabBorderInactive, true).Padding(0, 1)
	TabTabBorderActiveStyle   = TabTabBorderInactiveStyle.Border(TabTabBorderActive, true)

	OverlayStyle = DefaultStyle.
			Border(Border).
			Align(lipgloss.Center, lipgloss.Center).
			Foreground(TextColor)
)

func tableStyle() table.Styles {
	style := table.DefaultStyles()
	style.Header = style.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(MutedColor).
		BorderBottom(true).
		Bold(false)
	style.Selected = style.Selected.
		Foreground(TextColor).
		Background(AccentColor).
		Bold(false)
	return style
}

func tabBorderInactive(border lipgloss.Border) lipgloss.Border {
	border.BottomLeft = border.MiddleBottom
	border.BottomRight = border.MiddleBottom
	return border
}

func tabBorderActive(border lipgloss.Border) lipgloss.Border {
	lBot := border.BottomRight
	border.BottomRight = border.BottomLeft
	border.BottomLeft = lBot
	border.Bottom = " "
	return border
}

func tabbedViewBorder(border lipgloss.Border) lipgloss.Border {
	border.Top = " "
	border.TopLeft = border.Left
	border.TopRight = border.Right
	return border
}

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

func RenderBorderTitleWithKeybind(view, title, keybind string, border *lipgloss.Style) string {
	styledView := border.Render(view)
	keybindStyle := lipgloss.NewStyle().Foreground(border.GetBorderTopForeground())
	titleStyle := lipgloss.NewStyle().Foreground(AccentColor)
	keybind = fmt.Sprintf("[%s]%s", keybind, border.GetBorderStyle().Top)
	inlineTitle := fmt.Sprintf("%s%s",
		keybindStyle.Render(keybind),
		titleStyle.Render(title),
	)
	return compositor.Compose(inlineTitle, styledView, compositor.Begin, compositor.Begin, 2, 0)
}
