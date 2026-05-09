package tabview

import "charm.land/lipgloss/v2"

func RenderTabBar(
	titles []string,
	activeStyle,
	inactiveStyle lipgloss.Style,
	fullWidth int,
	active int,
) string {
	tabCount := len(titles)
	tabWidth := fullWidth / tabCount
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
		}
		if isFirst && !isActive {
			border.BottomLeft = border.MiddleLeft
		}
		if isLast && isActive {
			border.BottomRight = border.Right
		}
		if isLast && !isActive {
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
