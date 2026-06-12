package tabview

import (
	"charm.land/lipgloss/v2"
)

type Styles struct {
	ContentStyle        lipgloss.Style
	ActiveTabBarStyle   lipgloss.Style
	InactiveTabBarStyle lipgloss.Style
}

func NewInactiveTabBarBorder(border lipgloss.Border) lipgloss.Border {
	border.BottomLeft = border.MiddleBottom
	border.BottomRight = border.MiddleBottom
	return border
}

func NewActiveTabBarBorder(border lipgloss.Border) lipgloss.Border {
	lBot := border.BottomRight
	border.BottomRight = border.BottomLeft
	border.BottomLeft = lBot
	border.Bottom = " "
	return border
}

func NewContentBorder(border lipgloss.Border) lipgloss.Border {
	border.Top = " "
	border.TopLeft = border.Left
	border.TopRight = border.Right
	return border
}

func GenerateStyles(style *lipgloss.Style) *Styles {
	if style == nil {
		return nil
	}

	border := style.GetBorderStyle()

	inactive := NewInactiveTabBarBorder(border)
	active := NewActiveTabBarBorder(border)
	content := NewContentBorder(border)

	return &Styles{
		ContentStyle:        style.Border(content),
		ActiveTabBarStyle:   style.Border(active).Padding(0, 1),
		InactiveTabBarStyle: style.Border(inactive).Padding(0, 1),
	}
}

func DefaulStyles() *Styles {
	style := lipgloss.NewStyle().Border(lipgloss.RoundedBorder())
	return GenerateStyles(&style)
}
