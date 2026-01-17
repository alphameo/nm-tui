// Package styles provides styling for ui
package styles

import (
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
