package tabview

import (
	"charm.land/lipgloss/v2"
)

type Styles struct {
	ActiveTabStyle   lipgloss.Style
	InactiveTabStyle lipgloss.Style
}

func DefaultInactiveTabBorder(border lipgloss.Border) lipgloss.Border {
	border.BottomLeft = border.MiddleBottom
	border.BottomRight = border.MiddleBottom
	return border
}

func DefaultActiveTabBorder(border lipgloss.Border) lipgloss.Border {
	border.BottomRight, border.BottomLeft = border.BottomLeft, border.BottomRight
	border.Bottom = " "
	return border
}

func DefaultContentBorder(border lipgloss.Border) lipgloss.Border {
	border.Top = " "
	border.TopLeft = border.Left
	border.TopRight = border.Right
	return border
}

func GenerateStyles(style lipgloss.Style) Styles {
	border := style.GetBorderStyle()

	inactive := DefaultInactiveTabBorder(border)
	active := DefaultActiveTabBorder(border)

	return Styles{
		ActiveTabStyle:   style.Border(active).Padding(0, 1),
		InactiveTabStyle: style.Border(inactive).Padding(0, 1),
	}
}

func DefaultStyles() Styles {
	style := lipgloss.NewStyle().Border(lipgloss.RoundedBorder())
	return GenerateStyles(style)
}
