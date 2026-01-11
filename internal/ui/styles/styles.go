// Package styles provides styling for ui
package styles

import (
	"fmt"

	"github.com/alphameo/nm-tui/internal/ui/components/overlay"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

var (
	TextColor         = lipgloss.Color("#ffffff")
	AccentColor       = lipgloss.Color("99")
	BorderStyle       = lipgloss.RoundedBorder()
	DividerColor      = lipgloss.Color("240")
	TableStyle        = makeTableStyle()
	InactiveTabBorder = makeTabBorderWithBottom("┴", "─", "┴")
	ActiveTabBorder   = makeTabBorderWithBottom("┘", " ", "└")
	InactiveTabStyle  = lipgloss.NewStyle().Border(InactiveTabBorder, true).Padding(0, 1)
	ActiveTabStyle    = InactiveTabStyle.Border(ActiveTabBorder, true)
	OverlayStyle      = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				Align(lipgloss.Center, lipgloss.Center).
				Foreground(TextColor)
)

func makeTableStyle() table.Styles {
	style := table.DefaultStyles()
	style.Header = style.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(DividerColor).
		BorderBottom(true).
		Bold(false)
	style.Selected = style.Selected.
		Foreground(TextColor).
		Background(AccentColor).
		Bold(false)
	return style
}

func makeTabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}

func ConstructTabBar(
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
			border.BottomLeft = "│"
		} else if isFirst && !isActive {
			border.BottomLeft = "├"
		} else if isLast && isActive {
			border.BottomRight = "│"
		} else if isLast && !isActive {
			border.BottomRight = "┤"
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

func ApplyStyleWithTitle(view, title, keybind string, style lipgloss.Style) string {
	styledView := style.Render(view)
	inlineTitle := fmt.Sprintf("%s%s",
		lipgloss.NewStyle().Foreground(style.GetBorderTopForeground()).Render(fmt.Sprintf("[%s]%s", keybind, style.GetBorderStyle().Top)),
		lipgloss.NewStyle().Foreground(AccentColor).Render(title),
	)
	return overlay.Compose(inlineTitle, styledView, overlay.Begin, overlay.Begin, 2, 0)
}
